package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/stac47/myroomies/internal/server/services/usermngt"
	"github.com/stac47/myroomies/pkg/models"
)

func retrieveUsers(w http.ResponseWriter, r *http.Request) {
	users := usermngt.GetUsersList(r.Context())
	encodeJson(w, users)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var newUser models.User
	if err := decodeJson(w, r.Body, &newUser); err != nil {
		return
	}
	err := usermngt.CreateUser(r.Context(), newUser)
	if err != nil {
		http.Error(w, models.NewRestError("Error while adding a new user").ToJSON(),
			http.StatusInternalServerError)
		return
	}
	w.Header().Add("Location", fmt.Sprintf("/users/%s", newUser.Login))
	w.WriteHeader(http.StatusCreated)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	login := params["login"]
	var newUser models.User

	if err := decodeJson(w, r.Body, &newUser); err != nil {
		return
	}
	authenticatedUser, _ := GetAuthenticatedUser(r.Context())
	if err := usermngt.UpdateUser(r.Context(), authenticatedUser, login, newUser); err != nil {
		http.Error(w, models.NewRestError("Error while updating user").ToJSON(),
			http.StatusInternalServerError)
		return
	}
}

func getUserInfo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	login := params["login"]
	foundUser := usermngt.SearchUser(r.Context(), usermngt.ByLoginCriteria(login))
	if foundUser == nil {
		msg := fmt.Sprintf("The user %s does not exist", login)
		http.Error(w, models.NewRestError(msg).ToJSON(),
			http.StatusInternalServerError)
	} else {
		encodeJson(w, foundUser)
	}
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	login := params["login"]
	err := usermngt.DeleteUser(r.Context(), login)
	if err != nil {
		http.Error(w, models.NewRestError("Error while deleting  user").ToJSON(),
			http.StatusInternalServerError)
	}
}
