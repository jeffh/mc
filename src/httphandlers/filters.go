package httphandlers

import (
	"net/http"
)

// Dispatches to either of two handlers if the specified headers are
// present in the request.
type RequiredHeadersHandler struct {
	RequiredHeaders http.Header
	Handler         http.Handler
	FailedHandler   http.Handler
}

// Creates a new header filter handlers. If the given headers are
// satisfied, then the passed handler is passed.
//
// The default failure handler returns a http.StatusBadRequest status code
func NewRequiredHeadersHandler(headers http.Header, handler http.Handler) *RequiredHeadersHandler {
	h := NewFixtureFromString("failed required headers").WithStatusCode(http.StatusBadRequest)
	return &RequiredHeadersHandler{headers, handler, h}
}

// Alias to setting the FailedHandler together.
func (d *RequiredHeadersHandler) WithFailedHandler(h http.Handler) *RequiredHeadersHandler {
	d.FailedHandler = h
	return d
}

func canonicalizeHeaders(h http.Header) http.Header {
	newHeaders := make(http.Header, 0)
	for key, value := range h {
		newHeaders[http.CanonicalHeaderKey(key)] = value
	}
	return newHeaders
}

func (d *RequiredHeadersHandler) hasRequiredHeaders(h http.Header) bool {
	h = canonicalizeHeaders(h)
	for key, values := range canonicalizeHeaders(d.RequiredHeaders) {
		actualValues, ok := h[key]
		if !ok {
			return false
		}
		for _, value := range values {
			match := false
			for _, actualValue := range actualValues {
				if value == actualValue {
					match = true
					break
				}
			}
			if !match {
				return false
			}
		}
	}
	return true
}

func (d *RequiredHeadersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if d.hasRequiredHeaders(r.Header) {
		d.Handler.ServeHTTP(w, r)
	} else {
		d.FailedHandler.ServeHTTP(w, r)
	}
}
