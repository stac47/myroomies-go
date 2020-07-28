package rest

import (
	"net/http"

	"github.com/gorilla/mux"
)

func echo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	key := params["key"]
	w.Write([]byte(key))
}
