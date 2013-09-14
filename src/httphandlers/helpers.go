package httphandlers

import (
	"net/http"
)

// Ensures the creation of a handler, or panics otherwise.
func Required(h http.Handler, err error) http.Handler {
	if err != nil {
		panic(err)
	}
	return h
}
