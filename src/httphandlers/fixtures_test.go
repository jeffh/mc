package httphandlers

import (
	"bytes"
	"encoding/xml"
	. "github.com/jeffh/goexpect"
	"net/http"
	"net/http/httptest"
	"testing"
)

type FixtureData struct {
	XMLName xml.Name `xml:"root" json:"-"`
	Name    string   `xml:"name" json:"name"`
	Body    string   `xml:"body" json:"body"`
}

func TestFixtureHandlerAsXML(t *testing.T) {
	it := NewIt(t)
	resp := httptest.NewRecorder()
	f, err := NewFixtureAsXML(&FixtureData{Name: "First", Body: "Post"})
	it.Must(err)
	f.ServeHTTP(resp, nil)
	it.Expects(resp.Code, ToEqual, http.StatusOK)
	it.Expects(resp.Body.String(), ToEqual, xml.Header+`<root><name>First</name><body>Post</body></root>`)
	it.Expects(resp.Header(), ToEqual, http.Header{
		"content-type": []string{"text/xml"},
	})
}

func TestFixtureHandlerAsJSON(t *testing.T) {
	it := NewIt(t)
	resp := httptest.NewRecorder()
	f, err := NewFixtureAsJSON(&FixtureData{Name: "First", Body: "Post"})
	it.Must(err)
	f.ServeHTTP(resp, nil)
	it.Expects(resp.Code, ToEqual, http.StatusOK)
	it.Expects(resp.Body.String(), ToEqual, `{"name":"First","body":"Post"}`)
	it.Expects(resp.Header(), ToEqual, http.Header{
		"content-type": []string{"application/json"},
	})
}

func TestFixtureHandlerFromStringCanHandleRequest(t *testing.T) {
	it := NewIt(t)
	resp := httptest.NewRecorder()
	f := NewFixtureFromString("hello world")
	f.ServeHTTP(resp, nil)
	it.Expects(resp.Code, ToEqual, http.StatusOK)
	it.Expects(resp.Body.String(), ToEqual, "hello world")
	it.Expects(resp.Header(), ToBeEmpty)
}

func TestFixtureHandlerCanHandleRequest(t *testing.T) {
	it := NewIt(t)
	resp := httptest.NewRecorder()
	f := NewFixture([]byte("hello world")).WithStatusCode(http.StatusBadRequest)
	f.ServeHTTP(resp, nil)
	it.Expects(resp.Code, ToEqual, http.StatusBadRequest)
	it.Expects(resp.Body.String(), ToEqual, "hello world")
	it.Expects(resp.Header(), ToBeEmpty)
}

func TestFixtureHandlerFromReaderCanHandleMultipleRequest(t *testing.T) {
	it := NewIt(t)
	resp := httptest.NewRecorder()
	f, err := NewFixtureFromReader(bytes.NewBufferString("Hello world"))
	it.Must(err)
	f.StatusCode = http.StatusNotFound
	f.Headers["foo"] = []string{"Bar"}
	f.ServeHTTP(resp, nil)
	it.Expects(resp.Code, ToEqual, http.StatusNotFound)
	it.Expects(resp.Body.String(), ToEqual, "Hello world")
	it.Expects(resp.Header(), ToEqual, http.Header{
		"foo": []string{"Bar"},
	})

	resp = httptest.NewRecorder()
	f.ServeHTTP(resp, nil)
	it.Expects(resp.Body.String(), ToEqual, "Hello world")
}
