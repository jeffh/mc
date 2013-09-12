package httphelpers

import (
	. "github.com/jeffh/goexpect"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFixtureHandler(t *testing.T) {
	it := NewIt(t)
	resp := httptest.NewRecorder()
	f := NewFixtureFromString("hello world")
	f.ServeHTTP(resp, nil)
	it.Expects(resp.Code, ToEqual, http.StatusOK)
	it.Expects(resp.Body.String(), ToEqual, "hello world")
	it.Expects(resp.Header(), ToBeEmpty)
}

func TestFail(t *testing.T) {
}
