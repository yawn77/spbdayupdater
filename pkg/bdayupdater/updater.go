package bdayupdater

import (
	"fmt"
	"math/rand"
	"net/url"
	"regexp"
	"time"

	"github.com/gookit/slog"
	"github.com/yawn77/spcontrol"
)

// URLs
const (
	urlEdit = "https://spieleplanet.eu/profile.php?do=editprofile"
)

// regex
const (
	valueRegex = "<input type=\"text\" class=\"bginput\" name=\"%s\" (id=\"tb_homepage\" )?value=\"(?P<value>.*)\" size=\"\\d+\" maxlength=\"\\d+\" dir=\"ltr\" />"
)

// errors
const (
	originalValueNotFound = Error("did not find original value")
)

type values struct {
	homepage string
	icq      string
	aim      string
	msn      string
	yahoo    string
	skype    string
}

type Error string

func (e Error) Error() string { return string(e) }

func getRandomBirthday(yearOnly bool) (string, string, string) {
	now := time.Now()
	r := rand.New(rand.NewSource(now.UnixNano()))
	if yearOnly {
		y, m, d := now.Date()
		return fmt.Sprint(y - r.Intn(120)), fmt.Sprint(int(m)), fmt.Sprint(d)
	}

	max := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Unix()
	min := time.Date(now.Year()-120, now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Unix()
	bday := time.Unix(r.Int63n(max-min)+min, 0)
	return fmt.Sprint(bday.Year()), fmt.Sprint(int(bday.Month())), fmt.Sprint(bday.Day())
}

func getValue(name string, body string) (string, error) {
	re := regexp.MustCompile(fmt.Sprintf(valueRegex, name))
	l := re.FindStringSubmatch(string(body))
	if len(l) == 0 {
		return "", originalValueNotFound
	}
	return l[2], nil
}

func getCurrentValues(client spcontrol.Client) (values, error) {
	body, err := client.Get(urlEdit)
	if err != nil {
		return values{}, err
	}
	homepage, err := getValue("homepage", body)
	if err != nil {
		return values{}, nil
	}
	icq, err := getValue("icq", body)
	if err != nil {
		return values{}, err
	}
	aim, err := getValue("aim", body)
	if err != nil {
		return values{}, err
	}
	msn, err := getValue("msn", body)
	if err != nil {
		return values{}, err
	}
	yahoo, err := getValue("yahoo", body)
	if err != nil {
		return values{}, err
	}
	skype, err := getValue("skype", body)
	if err != nil {
		return values{}, err
	}
	return values{
		homepage: homepage,
		icq:      icq,
		aim:      aim,
		msn:      msn,
		yahoo:    yahoo,
		skype:    skype,
	}, nil
}

func editBirthday(year string, month string, day string, client spcontrol.Client) error {
	cur, err := getCurrentValues(client)
	if err != nil {
		return nil
	}
	values := url.Values{
		"s":             {""},
		"securitytoken": {client.Session.SecurityToken},
		"do":            {"updateprofile"},
		"month":         {month},
		"day":           {day},
		"year":          {year},
		"showbirthday":  {"2"},
		"homepage":      {cur.homepage},
		"icq":           {cur.icq},
		"aim":           {cur.aim},
		"msn":           {cur.msn},
		"yahoo":         {cur.yahoo},
		"skype":         {cur.skype},
	}
	_, err = client.Post(urlEdit, values)
	return err
}

func Update(yearOnly bool) {
	year, month, day := getRandomBirthday(yearOnly)
	creds, err := spcontrol.GetCredentials()
	if err != nil {
		slog.Error(err)
		return
	}
	slog.Infof("update birthday to %s-%s-%s for %s", month, day, year, creds.Username)
	client, err := spcontrol.GetClient()
	if err != nil {
		slog.Error(err)
		return
	}
	err = client.Login(creds)
	if err != nil {
		slog.Error(err)
		return
	}
	err = editBirthday(year, month, day, *client)
	if err != nil {
		slog.Error(err)
		return
	}
	err = client.Logout()
	if err != nil {
		slog.Error(err)
		return
	}
	slog.Info("updated bday successfully")
}
