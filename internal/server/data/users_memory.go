package data

import (
	"context"
	"errors"
	"fmt"

	"github.com/stac47/myroomies/pkg/models"
)

var (
	memoryUserDataAccess *MemoryUserDataAccess
)

type MemoryUserDataAccess struct {
	users []models.User
}

func GetMemoryUserDataAccess() *MemoryUserDataAccess {
	if memoryUserDataAccess == nil {
		memoryUserDataAccess = &MemoryUserDataAccess{
			users: make([]models.User, 0),
		}
	}
	return memoryUserDataAccess
}

func (dao *MemoryUserDataAccess) RetrieveUsers(ctx context.Context) []models.User {
	return dao.users
}

func (dao *MemoryUserDataAccess) CreateUser(ctx context.Context, newUser models.User) error {
	dao.users = append(dao.users, newUser)
	return nil
}

func (dao *MemoryUserDataAccess) getUserIndex(ctx context.Context, login string) (int, error) {
	for index, user := range dao.users {
		if user.Login == login {
			return index, nil
		}
	}
	return 0, errors.New(fmt.Sprintf("User [%s] not found", login))

}

func (dao *MemoryUserDataAccess) RetrieveUser(ctx context.Context, login string) *models.User {
	if index, err := dao.getUserIndex(ctx, login); err != nil {
		return &dao.users[index]
	}
	return nil
}

func (dao *MemoryUserDataAccess) DeleteUser(ctx context.Context, user models.User) error {
	for i, u := range dao.users {
		if user == u {
			dao.users[i] = dao.users[len(dao.users)-1]
			dao.users = dao.users[:len(dao.users)-1]
			return nil
		}
	}
	return errors.New(fmt.Sprintf("User [%s] cannot be found", user.Login))
}

func (dao *MemoryUserDataAccess) UpdateUser(ctx context.Context, user models.User) error {
	if index, err := dao.getUserIndex(ctx, user.Login); err != nil {
		return err
	} else {
		dao.users[index] = user
		return nil
	}
}
