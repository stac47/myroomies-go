package data

import (
	"context"

	"github.com/stac47/myroomies/pkg/models"
)

type UserDataAccess interface {
	RetrieveUsers(ctx context.Context) []models.User
	CreateUser(ctx context.Context, newUser models.User) error
	RetrieveUser(ctx context.Context, login string) *models.User
	DeleteUser(ctx context.Context, user models.User) error
	UpdateUser(ctx context.Context, user models.User) error
}
