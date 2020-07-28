package rest

import (
	"net/http"

	"github.com/stac47/myroomies/internal/server/services"
)

func version(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(services.GetVersion()))
}
