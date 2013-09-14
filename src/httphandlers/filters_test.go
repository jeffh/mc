package httphandlers

import (
	. "github.com/jeffh/goexpect"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequiredHeadersHandlerCallsBadHandlerWithoutTheRequiredHeaders(t *testing.T) {
	it := NewIt(t)
	resp := httptest.NewRecorder()
	okHandler := NewFixtureFromString("ok")
	failHandler := NewFixtureFromString("fail")
	h := NewRequiredHeadersHandler(http.Header{
		"X-Token": []string{"Foobar"},
	}, okHandler).WithFailedHandler(failHandler)

	r, err := http.NewRequest("GET", "/", nil)
	it.Must(err)

	h.ServeHTTP(resp, r)
	it.Expects(resp.Body.String(), ToEqual, "fail")
}

func TestRequiredHeadersHandlerCallsGoodHandlerWithTheRequiredHeaderValues(t *testing.T) {
	it := NewIt(t)
	resp := httptest.NewRecorder()
	okHandler := NewFixtureFromString("ok")
	failHandler := NewFixtureFromString("fail")
	h := NewRequiredHeadersHandler(http.Header{
		"X-Token": []string{"Foobar"},
	}, okHandler).WithFailedHandler(failHandler)

	r, err := http.NewRequest("GET", "/", nil)
	r.Header.Set("X-Token", "Foobar")
	it.Must(err)

	h.ServeHTTP(resp, r)
	it.Expects(resp.Body.String(), ToEqual, "ok")
}

func TestRequiredHeadersHandlerCallsGoodHandlerWithTheRequiredHeaders(t *testing.T) {
	it := NewIt(t)
	resp := httptest.NewRecorder()
	okHandler := NewFixtureFromString("ok")
	failHandler := NewFixtureFromString("fail")
	h := NewRequiredHeadersHandler(http.Header{
		"X-Token": nil,
	}, okHandler).WithFailedHandler(failHandler)

	r, err := http.NewRequest("GET", "/", nil)
	r.Header.Set("X-Token", "Cakes")
	it.Must(err)

	h.ServeHTTP(resp, r)
	it.Expects(resp.Body.String(), ToEqual, "ok")
}

func TestRequiredHeadersHandlerCallsGoodHandlerWithTheRequiredHeaderInMultipleHeaders(t *testing.T) {
	it := NewIt(t)
	resp := httptest.NewRecorder()
	okHandler := NewFixtureFromString("ok")
	failHandler := NewFixtureFromString("fail")
	h := NewRequiredHeadersHandler(http.Header{
		"X-Token": []string{"Cakes"},
	}, okHandler).WithFailedHandler(failHandler)

	r, err := http.NewRequest("GET", "/", nil)
	r.Header.Add("X-Token", "Cakes")
	r.Header.Add("X-Token", "Cream")
	it.Must(err)

	h.ServeHTTP(resp, r)
	it.Expects(resp.Body.String(), ToEqual, "ok")
}

func TestRequiredHeadersHandlerCallsBadHandlerWithTheRequiredHeaderIsNotInMultipleHeaders(t *testing.T) {
	it := NewIt(t)
	resp := httptest.NewRecorder()
	okHandler := NewFixtureFromString("ok")
	failHandler := NewFixtureFromString("fail")
	h := NewRequiredHeadersHandler(http.Header{
		"X-Token": []string{"Cakes"},
	}, okHandler).WithFailedHandler(failHandler)

	r, err := http.NewRequest("GET", "/", nil)
	r.Header.Add("X-Token", "Crab")
	r.Header.Add("X-Token", "Lobster")
	it.Must(err)

	h.ServeHTTP(resp, r)
	it.Expects(resp.Body.String(), ToEqual, "fail")
}
