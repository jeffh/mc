package httphandlers

import (
	. "github.com/jeffh/goexpect"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestURLHandlerMapsToHandlers(t *testing.T) {
	it := NewIt(t)
	h1 := NewFixtureFromString("h1")
	h2 := NewFixtureFromString("h2")
	h3 := NewFixtureFromString("h3")
	h4 := NewFixtureFromString("h4")

	h := NewURLHandler(map[string]http.Handler{
		"GET /":   h1,
		"GET /h2": h2,
	}).WithErrorHandler(h4)

	h.Any("/h3", h3)

	r, err := http.NewRequest("GET", "http://localhost/", nil)
	r.RequestURI = "/"
	it.Must(err)
	resp := httptest.NewRecorder()
	h.ServeHTTP(resp, r)
	it.Expects(resp.Body.String(), ToEqual, "h1")

	r, err = http.NewRequest("GET", "http://localhost/h2", nil)
	r.RequestURI = "/h2"
	it.Must(err)
	resp = httptest.NewRecorder()
	h.ServeHTTP(resp, r)
	it.Expects(resp.Body.String(), ToEqual, "h2")

	r, err = http.NewRequest("GET", "http://localhost/h3", nil)
	r.RequestURI = "/h3"
	it.Must(err)
	resp = httptest.NewRecorder()
	h.ServeHTTP(resp, r)
	it.Expects(resp.Body.String(), ToEqual, "h3")

	r, err = http.NewRequest("GET", "http://localhost/404", nil)
	r.RequestURI = "/404"
	it.Must(err)

	resp = httptest.NewRecorder()
	h.ServeHTTP(resp, r)
	it.Expects(resp.Body.String(), ToEqual, "h4")
}
