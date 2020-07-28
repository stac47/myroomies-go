package rest

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/stac47/myroomies/pkg/models"
)

func encodeJson(w http.ResponseWriter, obj interface{}) error {
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(obj); err != nil {
		http.Error(w, models.NewRestError("Error during JSON encoding: "+err.Error()).ToJSON(),
			http.StatusBadRequest)
		return err
	}
	return nil
}

func decodeJson(w http.ResponseWriter, r io.Reader, obj interface{}) error {
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&obj); err != nil {
		http.Error(w, models.NewRestError("Invalid JSON: "+err.Error()).ToJSON(), http.StatusBadRequest)
		return err
	}
	return nil
}
