package httphelpers

import (
	"io"
	"net/http"
	"os"
	"reflect"
)

type FixtureHandler struct {
	Contents   []byte
	StatusCode int
	Headers    map[string][]string
}

func NewFixtureFromFile(filepath string) *FixtureHandler {
	h, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer h.Close()
	return NewFixtureFromReader(h)
}

func NewFixtureFromReader(r io.Reader) *FixtureHandler {
	bytes := make([]byte, 0)
	buf := make([]byte, 1024)
	for _, err := r.Read(buf); err != nil; {
		if err != nil && err != io.EOF {
			panic(err)
		}
		for _, b := range buf {
			bytes = append(bytes, b)
		}
	}
	return NewFixture(bytes)
}

func NewFixtureFromString(s string) *FixtureHandler {
	return NewFixture([]byte(s))
}

func NewFixture(b []byte) *FixtureHandler {
	return &FixtureHandler{b, 200, make(map[string][]string)}
}

func (f *FixtureHandler) WithStatusCode(code int) *FixtureHandler {
	f.StatusCode = code
	return f
}

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

type RequiredHeadersHandler struct {
	RequiredHeaders map[string][]string
	Handler         http.Handler
	FailedHandler   http.Handler
}

func NewRequiredHeadersHandler(headers map[string][]string, handler http.Handler) *RequiredHeadersHandler {
	h := NewFixtureFromString("").WithStatusCode(http.StatusBadRequest)
	return &RequiredHeadersHandler{headers, handler, h}
}

func (d *RequiredHeadersHandler) WithFailedHandler(h http.Handler) *RequiredHeadersHandler {
	d.FailedHandler = h
	return d
}

func (d *RequiredHeadersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	valid := true
	for key, values := range d.RequiredHeaders {
		actualValues, ok := r.Header[key]
		valid = ok && (values == nil || reflect.DeepEqual(values, actualValues))
		if !valid {
			break
		}
	}

	if valid {
		d.Handler.ServeHTTP(w, r)
	} else {
		d.FailedHandler.ServeHTTP(w, r)
	}
}

type URLHandler struct {
	Handlers     map[string]http.Handler
	ErrorHandler http.Handler
}

func NewURLHandler(handlers map[string]http.Handler) *URLHandler {
	h := NewFixtureFromString("Not Found").WithStatusCode(http.StatusNotFound)
	return &URLHandler{handlers, h}
}

func (u *URLHandler) WithErrorHandler(h http.Handler) *URLHandler {
	u.ErrorHandler = h
	return u
}

func (h *URLHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler, ok := h.Handlers[r.RequestURI]
	if ok {
		handler.ServeHTTP(w, r)
	} else {
		h.ErrorHandler.ServeHTTP(w, r)
	}
}
