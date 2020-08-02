package rest

import (
	"net/http"

	"github.com/stac47/myroomies/internal/server/services"
)

func getGlobalInfo(w http.ResponseWriter, r *http.Request) {
	info := services.GetGlobalInfo()
	encodeJson(w, info)
}
