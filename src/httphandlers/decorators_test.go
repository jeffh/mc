package httphandlers

import (
	. "github.com/jeffh/goexpect"
	"net/http"
	"net/http/httptest"
	"testing"
)

func doRequest(it *It, h http.Handler, method, url string) (*http.Request, *httptest.ResponseRecorder) {
	resp := httptest.NewRecorder()
	r, err := http.NewRequest(method, url, nil)
	it.Must(err)
	h.ServeHTTP(resp, r)
	return r, resp
}

func TestRequestRecorderCanRecordIncomingRequestsWithoutAHandler(t *testing.T) {
	it := NewIt(t)
	h := NewRequestRecorderHandler(nil)
	r1, _ := doRequest(it, h, "GET", "http://localhost/")
	r2, _ := doRequest(it, h, "POST", "http://localhost/path1")

	it.Expects(h.Requests, ToBeLengthOf, 2)
	it.Expects(h.Requests[0].URL.String(), ToEqual, r1.URL.String())
	it.Expects(h.Requests[1].URL.String(), ToEqual, r2.URL.String())

	it.Expects(h.RequestsByPath("/"), ToBeLengthOf, 1)
	it.Expects(h.RequestsByPath("/")[0].URL.String(), ToEqual, r1.URL.String())
	it.Expects(h.RequestsByMethod("POST"), ToBeLengthOf, 1)
	it.Expects(h.RequestsByMethod("POST")[0].URL.String(), ToEqual, r2.URL.String())
}

func TestRequestRecorderCanClearRecordedRequests(t *testing.T) {
	it := NewIt(t)
	resp := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "http://localhost/", nil)
	it.Must(err)
	f := NewFixtureFromString("ok")
	h := NewRequestRecorderHandler(f)
	h.ServeHTTP(resp, r)
	h.ClearRequests()
	it.Expects(h.Requests, ToBeEmpty)
}

func TestRequestRecorderCanRecordIncomingRequestsBy(t *testing.T) {
	it := NewIt(t)
	resp := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "http://localhost/", nil)
	it.Must(err)
	f := NewFixtureFromString("ok")
	h := NewRequestRecorderHandler(f)
	h.ServeHTTP(resp, r)
	it.Expects(resp.Body.String(), ToEqual, "ok")

	it.Expects(h.Requests, ToBeLengthOf, 1)
	it.Expects(h.Requests[0].URL.String(), ToEqual, r.URL.String())
}
