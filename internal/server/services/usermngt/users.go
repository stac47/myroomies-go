package usermngt

import (
	"context"
	"errors"
	"fmt"

	"github.com/stac47/myroomies/internal/server/services"
	"github.com/stac47/myroomies/pkg/models"

	log "github.com/sirupsen/logrus"
)

func Login(ctx context.Context, login string, password string) *models.User {
	user := SearchUser(ctx, ByLoginCriteria(login))
	if user == nil {
		log.Printf("User [%s] not found.", login)
		return nil
	}
	if user.Password.String() != password {
		log.Printf("Wrong password [%s] for user [%s]", password, login)
		return nil
	}
	return user
}

func GetUsersList(ctx context.Context) []models.User {
	return services.GetDataAccess().GetUserDataAccess().RetrieveUsers(ctx)
}

func CreateUser(ctx context.Context, user models.User) error {
	return services.GetDataAccess().GetUserDataAccess().CreateUser(ctx, user)
}

func mergeUser(userToUpdate *models.User, update models.User) {
	if update.Firstname != "" {
		userToUpdate.Firstname = update.Firstname
	}
	if update.Lastname != "" {
		userToUpdate.Lastname = update.Lastname
	}
	if update.Password != "" {
		userToUpdate.Password = update.Password
	}
}

func UpdateUser(ctx context.Context, authenticatedUser models.User, login string, update models.User) (err error) {
	if authenticatedUser.Login == login || authenticatedUser.IsAdmin {
		userToUpdate := SearchUser(ctx, ByLoginCriteria(login))
		if userToUpdate == nil {
			err = errors.New(fmt.Sprintf("User to update [%s] not found", login))
			return
		}
		mergeUser(userToUpdate, update)
		err = services.GetDataAccess().GetUserDataAccess().UpdateUser(ctx, *userToUpdate)
	} else {
		err = errors.New(fmt.Sprintf("The authenticated [%s] user is not allowed "+
			"to update user [%s]", authenticatedUser.Login, login))
		return
	}
	return
}

func DeleteUser(ctx context.Context, login string) error {
	user := SearchUser(ctx, ByLoginCriteria(login))
	if user != nil {
		return services.GetDataAccess().GetUserDataAccess().DeleteUser(ctx, *user)
	}
	return nil
}

type SearchUserCriteria func(user *models.User) bool

func ByLoginCriteria(login string) SearchUserCriteria {
	return func(user *models.User) bool {
		return user.Login == login
	}
}

func SearchUser(ctx context.Context, fn SearchUserCriteria) *models.User {
	users := services.GetDataAccess().GetUserDataAccess().RetrieveUsers(ctx)
	for _, user := range users {
		if fn(&user) {
			return &user
		}
	}
	return nil
}
