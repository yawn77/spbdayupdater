package updater

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gookit/slog"
)

const (
	UrlEdit   = "https://spieleplanet.eu/profile.php?do=editprofile"
	UrlUpdate = "https://spieleplanet.eu/profile.php?do=updateprofile"
	UrlLogin  = "https://www.spieleplanet.eu/login.php?do=login"
	UrlRoot   = "https://www.spieleplanet.eu/"
)

type Error string

func (e Error) Error() string { return string(e) }

const (
	LoginFailed                = Error("login failed, please check username and password")
	LogoutFailed               = Error("logout failed")
	OriginalValueNotFound      = Error("did not find original value")
	SecurityTokenNotFound      = Error("failed to find security token")
	UpdateFailed               = Error("update of birthday failed")
	UserNameOrPasswordNotFound = Error("environment variable SP_USERNAME or SP_PASSWORD not set")
)

type credentials struct {
	username    string
	password    string
	passwordMd5 string
}

type session struct {
	securityToken string
	logoutUrl     string
}

func getMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func GetCredentials() (credentials, error) {
	username := os.Getenv("SP_USERNAME")
	password := os.Getenv("SP_PASSWORD")
	if username == "" || password == "" {
		return credentials{}, UserNameOrPasswordNotFound
	}
	passwordMd5 := getMD5Hash(password)

	return credentials{
		username:    username,
		password:    password,
		passwordMd5: passwordMd5,
	}, nil
}

func login(credentials credentials) (*http.Client, error) {
	options := cookiejar.Options{}
	jar, err := cookiejar.New(&options)
	if err != nil {
		return nil, err
	}
	client := http.Client{Jar: jar}
	resp, err := client.PostForm(UrlLogin, url.Values{
		"vb_login_username":        {credentials.username},
		"vb_login_password":        {credentials.password},
		"s":                        {""},
		"securitytoken":            {"guest"},
		"do":                       {"login"},
		"vb_login_md5password":     {credentials.passwordMd5},
		"vb_login_md5password_utf": {credentials.passwordMd5},
	})
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}
	if !strings.Contains(strings.ToUpper(string(body)), fmt.Sprintf("DANKE FÜR DEINE ANMELDUNG, %s.", strings.ToUpper(credentials.username))) {
		return nil, LoginFailed
	}

	return &client, nil
}

func getSessionInformation(client *http.Client) (session, error) {
	resp, err := client.Get(UrlRoot)
	if err != nil {
		return session{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return session{}, err
	}

	securityTokenRegex := "<input type=\"hidden\" name=\"securitytoken\" value=\"(?P<token>.*)\" \\/>"
	re := regexp.MustCompile(securityTokenRegex)
	l := re.FindStringSubmatch(string(body))
	if len(l) == 0 {
		return session{}, SecurityTokenNotFound
	}
	token := l[1]

	return session{
		securityToken: token,
		logoutUrl:     fmt.Sprintf("%slogin.php?do=logout&amp;logouthash=%s", UrlRoot, token),
	}, nil
}

func getBirthday() (string, string, string) {
	rand.Seed(time.Now().UnixNano())
	year, month, day := time.Now().Date()
	return fmt.Sprint(year - rand.Intn(120)), fmt.Sprint(int(month)), fmt.Sprint(day)
}

func getValue(name string, body string) (string, error) {
	regex := fmt.Sprintf("<input type=\"text\" class=\"bginput\" name=\"%s\" (id=\"tb_homepage\" )?value=\"(?P<value>.*)\" size=\"\\d+\" maxlength=\"\\d+\" dir=\"ltr\" />", name)
	re := regexp.MustCompile(regex)
	l := re.FindStringSubmatch(string(body))
	if len(l) == 0 {
		return "", OriginalValueNotFound
	}
	return l[2], nil
}

func editBirthday(year string, month string, day string, session session, client *http.Client) error {
	resp, err := client.Get(UrlEdit)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	bodyStr := string(body)

	homepage, err := getValue("homepage", bodyStr)
	if err != nil {
		return err
	}
	icq, err := getValue("icq", bodyStr)
	if err != nil {
		return err
	}
	aim, err := getValue("aim", bodyStr)
	if err != nil {
		return err
	}
	msn, err := getValue("msn", bodyStr)
	if err != nil {
		return err
	}
	yahoo, err := getValue("yahoo", bodyStr)
	if err != nil {
		return err
	}
	skype, err := getValue("skype", bodyStr)
	if err != nil {
		return err
	}

	resp, err = client.PostForm(UrlEdit, url.Values{
		"s":             {""},
		"securitytoken": {session.securityToken},
		"do":            {"updateprofile"},
		"month":         {month},
		"day":           {day},
		"year":          {year},
		"showbirthday":  {"2"},
		"homepage":      {homepage},
		"icq":           {icq},
		"aim":           {aim},
		"msn":           {msn},
		"yahoo":         {yahoo},
		"skype":         {skype},
	})
	if err != nil {
		return err
	}

	_, err = io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return err
	}

	return nil
}

func Update() {
	year, month, day := getBirthday()
	slog.Infof("update birthday to %s-%s-%s", month, day, year)
	creds, err := GetCredentials()
	if err != nil {
		slog.Error(err)
		return
	}
	client, err := login(creds)
	if err != nil {
		slog.Error(err)
		return
	}
	session, err := getSessionInformation(client)
	if err != nil {
		slog.Error(err)
		return
	}
	err = editBirthday(year, month, day, session, client)
	if err != nil {
		slog.Error(err)
		return
	}
	slog.Info("updated bday successfully")
}
