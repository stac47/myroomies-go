package rest

import (
	"net/http"

	"github.com/stac47/myroomies/internal/server/services/usermngt"

	log "github.com/sirupsen/logrus"
)

type authenticationMiddleware struct {
	// TODO: implement a way to avoid calling the Login service
}

func (amw *authenticationMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		login, password, ok := r.BasicAuth()
		if !ok {
			auth := r.Header.Get("Authorization")
			if auth == "" {
				log.Print("No HTTP Authorization header")
			} else {
				log.Printf("Authorization HTTP header [%s] is wrong", auth)
			}
		} else {
			user := usermngt.Login(r.Context(), login, password)
			if user == nil {
				log.Printf("Invalid login or password [%s, %s]", login, password)
			} else {
				r = r.WithContext(SetAuthenticatedUser(r.Context(), *user))
				log.Printf("User [%s] is authenticated", login)
			}
		}
		next.ServeHTTP(w, r)
	})
}
