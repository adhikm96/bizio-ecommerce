package common

import (
	"errors"
	"net/http"
)

func ValidateQueryParam(request *http.Request, name string) error {
	if request.URL.Query().Get(name) == "" {
		return errors.New(name + " is required")
	}
	return nil
}
