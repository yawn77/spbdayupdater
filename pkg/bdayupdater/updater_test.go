package bdayupdater

import (
	"flag"
	"fmt"
	"testing"
	"time"

	"github.com/yawn77/spcontrol"
)

// credentials for command-line tests
var username = flag.String("username", "", "Username")
var password = flag.String("password", "", "User password")

// get credentials either from env or from cmdline if env does not exists
func getCredentials() spcontrol.Credentials {
	cred, err := spcontrol.GetCredentials()
	if err != nil {
		cred = spcontrol.Credentials{
			Username:    *username,
			Password:    *password,
			PasswordMd5: spcontrol.GetMD5Hash(*password),
		}
	}
	return cred
}

func TestGetRandomBirthday(t *testing.T) {
	// arrange
	ey, em, ed := time.Now().Date()
	cnt := 10
	diff := false

	// act
	for i := 0; i < cnt; i++ {
		year, month, day := getRandomBirthday(false)
		if year != fmt.Sprint(ey) && month != fmt.Sprint(em) && day != fmt.Sprint(ed) {
			diff = true
		}
	}

	// assert
	if !diff {
		t.Fatal("no entire differnt date was generated")
	}
}

func TestGetRandomBirthdayYearOnly(t *testing.T) {
	// arrange
	ey, em, ed := time.Now().Date()

	// act
	year, month, day := getRandomBirthday(true)

	// assert
	if year == fmt.Sprint(ey) || month != fmt.Sprint(int(em)) || day != fmt.Sprint(ed) {
		t.Fatalf("expect only year to be changed: got %s-%s-%s vs %d-%d-%d", month, day, year, em, ed, ey)
	}
}

func TestUpdate(t *testing.T) {
	// arrange
	year, month, day := getRandomBirthday(false)
	creds := getCredentials()
	client, err := spcontrol.GetClient()
	if err != nil {
		t.Fatalf("could not create client: %v", err)
	}
	err = client.Login(creds)
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}

	// act
	err = editBirthday(year, month, day, *client)

	// assert
	if err != nil {
		t.Fatalf("update failed: %v", err)
	}

	// cleanup
	err = client.Logout()
	if err != nil {
		t.Fatalf("lgout failed: %v", err)
	}
}
