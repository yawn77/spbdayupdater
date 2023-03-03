package updater

import (
	"os"
	"testing"
)

func TestGetCredentialsSucess(t *testing.T) {
	os.Setenv("SP_USERNAME", "username")
	defer os.Unsetenv("SP_USERNAME")
	os.Setenv("SP_PASSWORD", "password")
	defer os.Unsetenv("SP_PASSWORD")

	creds, err := GetCredentials()

	if err != nil {
		t.Fatalf("err got %v, want nil", err)
	}
	expectedCredentials := credentials{
		username:    "username",
		password:    "password",
		passwordMd5: "5f4dcc3b5aa765d61d8327deb882cf99",
	}
	if creds != expectedCredentials {
		t.Fatalf("creds got %v, want %v", creds, expectedCredentials)
	}
}

func TestGetCredentialsWithoutUsernameOrPassword(t *testing.T) {
	_, err := GetCredentials()

	if err != UserNameOrPasswordNotFound {
		t.Fatalf("err got %v, want nil", err)
	}
}
