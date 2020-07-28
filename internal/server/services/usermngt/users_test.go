package usermngt

import (
	"context"
	"fmt"
	"testing"

	"github.com/stac47/myroomies/internal/server/services"
	"github.com/stac47/myroomies/pkg/models"
)

func generateTestUsers(ctx context.Context, t *testing.T, number int) {
	for i := 0; i < number; i++ {
		err := CreateUser(ctx, models.User{
			Firstname: fmt.Sprintf("user%d_firstname", i),
			Lastname:  fmt.Sprintf("user%d_lastname", i),
			IsAdmin:   i == 0,
			Login:     fmt.Sprintf("user%d", i),
			Password:  models.PasswordType(fmt.Sprintf("password%d", i)),
		})
		if err != nil {
			t.Fatalf("Cannot create test user %d", i)
		}
	}
}

func removeTestUsers(ctx context.Context, t *testing.T) {
	users := GetUsersList(ctx)
	for _, user := range users {
		err := DeleteUser(ctx, user.Login)
		if err != nil {
			t.Errorf("Error when deleting user [%s]", user.Login)
		}
	}
}

func TestLogin(t *testing.T) {
	ctx := context.Background()
	services.Configure(nil)
	generateTestUsers(ctx, t, 4)
	defer removeTestUsers(ctx, t)
	var tests = []struct {
		name     string
		login    string
		password string
		expected bool
	}{
		{"good login, good password", "user0", "password0", true},
		{"good login, bad password", "user0", "wrong", false},
		{"bad user, bad password", "wrong", "wrong", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.expected != (Login(ctx, test.login, test.password) != nil) {
				t.Errorf("Login with [%s, %s] should have returned %t",
					test.login, test.password, test.expected)
			}
		})
	}
}

func TestGetUsersList(t *testing.T) {
	ctx := context.Background()
	services.Configure(nil)
	generateTestUsers(ctx, t, 4)
	defer removeTestUsers(ctx, t)
	users := GetUsersList(ctx)
	if number := len(users); number != 4 {
		t.Errorf("Only %d users retrieved. Expected: 4", number)
	}
}

func TestSearchUser(t *testing.T) {
	ctx := context.Background()
	services.Configure(nil)
	generateTestUsers(ctx, t, 4)
	defer removeTestUsers(ctx, t)
	user0 := SearchUser(ctx, ByLoginCriteria("user0"))
	if user0 == nil {
		t.Fatal("user0 should exist")
	}
	if !user0.IsAdmin {
		t.Fatal("user0 should be admin")
	}
}

func TestUpdateUser(t *testing.T) {
	ctx := context.Background()
	services.Configure(nil)
	generateTestUsers(ctx, t, 4)
	defer removeTestUsers(ctx, t)
	tests := []struct {
		name                   string
		authenticatedUserLogin string
		login                  string
		update                 models.User
		errorExpected          bool
	}{
		{
			name:                   "Admin changing his info",
			authenticatedUserLogin: "user0",
			login:                  "user0",
			update: models.User{
				Firstname: "changed_user0_firstname",
				Lastname:  "changed_user0_lastname",
				Password:  "changed_user0_password",
			},
			errorExpected: false,
		},
		{
			name:                   "Admin changing another info user",
			authenticatedUserLogin: "user0",
			login:                  "user1",
			update: models.User{
				Firstname: "changed_user1_firstname",
				Lastname:  "changed_user1_lastname",
				Password:  "changed_user1_password",
			},
			errorExpected: false,
		},
		{
			name:                   "Normal user changing his user info",
			authenticatedUserLogin: "user1",
			login:                  "user1",
			update: models.User{
				Firstname: "changed_user1_firstname",
				Lastname:  "changed_user1_lastname",
				Password:  "changed_user1_password",
			},
			errorExpected: false,
		},
		{
			name:                   "Normal user partially changing his user info",
			authenticatedUserLogin: "user1",
			login:                  "user1",
			update: models.User{
				Password: "changed_user1_password",
			},
			errorExpected: false,
		},
		{
			name:                   "Normal user trying to chang another user info",
			authenticatedUserLogin: "user1",
			login:                  "user0",
			update: models.User{
				Firstname: "changed_user0_firstname",
				Lastname:  "changed_user0_lastname",
				Password:  "changed_user0_password",
			},
			errorExpected: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			authenticatedUser := SearchUser(ctx, ByLoginCriteria(test.authenticatedUserLogin))
			userBackup := *SearchUser(ctx, ByLoginCriteria(test.login))
			err := UpdateUser(ctx, *authenticatedUser, test.login, test.update)
			if err != nil {
				if !test.errorExpected {
					t.Errorf("An unexpected error occured when updating a user: %s", err)
				}
			}
			if !test.errorExpected {
				modifiedUser := SearchUser(ctx, ByLoginCriteria(test.login))
				if test.update.Firstname != "" && userBackup.Firstname != test.update.Firstname {
					if modifiedUser.Firstname != test.update.Firstname {
						t.Errorf("The modified user's firstname [%s] is not the one from the request [%s]",
							modifiedUser.Firstname, test.update.Firstname)
					}
				}
				if test.update.Lastname != "" && userBackup.Lastname != test.update.Lastname {
					if modifiedUser.Lastname != test.update.Lastname {
						t.Errorf("The modified user's firstname [%s] is not the one from the request [%s]",
							modifiedUser.Lastname, test.update.Lastname)
					}
				}
				if test.update.Password != "" && userBackup.Password != test.update.Password {
					if modifiedUser.Password != test.update.Password {
						t.Errorf("The modified user's firstname [%s] is not the one from the request [%s]",
							modifiedUser.Password, test.update.Password)
					}
				}
				// Restore the modified user to its initial state
				err := UpdateUser(ctx, *authenticatedUser, test.login, userBackup)
				if err != nil {
					t.Errorf("An error occured while restoring the user data: %s", err)
				}
				// Verify if the user was correctly restored
				verificationUser := SearchUser(ctx, ByLoginCriteria(test.login))
				if *verificationUser != userBackup {
					t.Errorf("User [%s] was not restored correctly", test.login)
				}
			}

		})
	}
}
