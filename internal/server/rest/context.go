package rest

import (
	"context"

	"github.com/stac47/myroomies/pkg/models"
)

type authenticatedUserKeyType int

const (
	authenticatedUserKey authenticatedUserKeyType = 1
)

func SetAuthenticatedUser(ctx context.Context, user models.User) context.Context {
	return context.WithValue(ctx, authenticatedUserKey, user)
}

func GetAuthenticatedUser(ctx context.Context) (user models.User, ok bool) {
	user, ok = ctx.Value(authenticatedUserKey).(models.User)
	return
}
