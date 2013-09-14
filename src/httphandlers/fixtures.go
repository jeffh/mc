package httphandlers

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

// A simple handler that returns the given fixture it has been
// constructed with. Can be configured to return a given set of
// headers and the status code.
//
// Note that all the contents are stored in memory when the
// fixture handler is created.
type FixtureHandler struct {
	Contents   []byte              // The raw bytes that holds the http body to return
	StatusCode int                 // The http status code to return
	Headers    map[string][]string // The http headers to return
}

// Creates a new fixture to emit an XML from value object.
// The value provided should support encoding/xml.Marshal.
//
// An error is returned if marshaling fails
func NewFixtureAsXML(v interface{}) (*FixtureHandler, error) {
	b, err := xml.Marshal(v)
	if err != nil {
		return nil, err
	}
	return NewFixtureFromString(xml.Header+string(b)).WithHeader("content-type", []string{"text/xml"}), nil
}

// Creates a new fixture to emit an JSON from value object.
// The value provided should support encoding/json.Marshal.
//
// An error is returned if marshaling fails
func NewFixtureAsJSON(v interface{}) (*FixtureHandler, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return NewFixture(b).WithHeader("content-type", []string{"application/json"}), nil
}

// Creates a new fixture to emit data from a given filepath.
// The file is loaded at the beginning. If the specified
// file does not exist, this function will panic
//
// An error is returned if the file cannot be read
func NewFixtureFromFile(filepath string) (*FixtureHandler, error) {
	h, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer h.Close()
	return NewFixtureFromReader(h)
}

// Creates a new fixture to emit data from a given reader.
//
// All the bytes are read at the beginning. It is up to
// the caller to close the reader.
//
// An error is returned if reading from the reader fails
func NewFixtureFromReader(r io.Reader) (*FixtureHandler, error) {
	buffer, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return NewFixture(buffer), nil
}

// Creates a new fixture that returns the given string
func NewFixtureFromString(s string) *FixtureHandler {
	return NewFixture([]byte(s))
}

// Creates a new fixture that returns the given bytes
func NewFixture(b []byte) *FixtureHandler {
	return &FixtureHandler{b, http.StatusOK, make(map[string][]string)}
}

// Alias to set the StatusCode property. Returns the same instance
// for chainability
func (f *FixtureHandler) WithStatusCode(code int) *FixtureHandler {
	f.StatusCode = code
	return f
}

// Alias to set the Header[key] property. Returns the same instance
// for chainability
func (f *FixtureHandler) WithHeader(key string, values []string) *FixtureHandler {
	f.Headers[key] = values
	return f
}

func (f *FixtureHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(f.StatusCode)
	for key, values := range f.Headers {
		w.Header()[key] = values
	}
	for written := 0; written < len(f.Contents); {
		n, err := w.Write(f.Contents[written:])
		if err != nil {
			panic(err)
		}
		written += n
	}
}
