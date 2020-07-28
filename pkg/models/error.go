package models

import (
	"encoding/json"
	"strings"
)

type Error struct {
	Message string
}

func NewRestError(msg string) Error {
	return Error{
		Message: msg,
	}
}

func (e Error) ToJSON() string {
	builder := strings.Builder{}
	encoder := json.NewEncoder(&builder)
	err := encoder.Encode(e)
	if err != nil {
		return "An error occured during the encoding of this error"
	}
	return builder.String()
}
