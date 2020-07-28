package models

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestMarshallingUser(t *testing.T) {
	const password PasswordType = "superpassword"

	user := User{
		Firstname: "user_firstname",
		Lastname:  "user_lastname",
		IsAdmin:   true,
		Login:     "user",
		Password:  password,
	}
	b, err := json.Marshal(&user)
	if err != nil {
		t.Fatal("Error when marshalling")
	}
	marshalledStr := string(b)
	if strings.Contains(marshalledStr, password.String()) {
		t.Errorf("Marshalled user should not contain %s: %s", password, marshalledStr)
	}
}
