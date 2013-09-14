// A collection of simple HTTP handlers.
//
// Composing these handlers can be used to build a simple stub
// HTTP server. Along with net/http/httptest package, a stubbed
// service can be built.
package httphandlers

import (
	"fmt"
	"net/http"
)

// A simple router of requests. Can be composed with other handlers
// to build a stub web service.
//
// If a request is received that is not in the router, the
// ErrorHandler is called instead
type URLHandler struct {
	Handlers     map[string]http.Handler
	ErrorHandler http.Handler
}

// Creates a new URLHandler. It accepts the map of strings to
// handlers. The string is in:
//
//   "<METHOD> <PATH>"
//
// Such as "GET /" for the index page.
//
// Use WithErrorHandler() to quickly specify a handler to be
// called if an incoming request does not map.
//
// The default error handler will simply return 404.
func NewURLHandler(routes map[string]http.Handler) *URLHandler {
	if routes == nil {
		routes = make(map[string]http.Handler)
	}
	h := NewFixtureFromString("Not Found").WithStatusCode(http.StatusNotFound)
	return &URLHandler{routes, h}
}

// Alias to setting a route.
//
// Handles the "<METHOD> <URL>" format when setting the route.
func (u *URLHandler) Set(method, url string, h http.Handler) {
	u.Handlers[fmt.Sprintf("%s %s", method, url)] = h
}

func (u *URLHandler) Any(url string, h http.Handler) {
	u.Set("ANY", url, h)
}

// Shorthand to setting a GET route
func (u *URLHandler) Get(url string, h http.Handler) {
	u.Set("GET", url, h)
}

// Shorthand to setting a POST route
func (u *URLHandler) Post(url string, h http.Handler) {
	u.Set("POST", url, h)
}

// Shorthand to setting a PUT route
func (u *URLHandler) Put(url string, h http.Handler) {
	u.Set("PUT", url, h)
}

// Shorthand to setting a DELETE route
func (u *URLHandler) Delete(url string, h http.Handler) {
	u.Set("DELETE", url, h)
}

// Shorthand to setting a HEAD route
func (u *URLHandler) Head(url string, h http.Handler) {
	u.Set("HEAD", url, h)
}

// Shorthand to setting an OPTIONS route
func (u *URLHandler) Options(url string, h http.Handler) {
	u.Set("HEAD", url, h)
}

// Builder-syntax to specify an ErrorHandler to resolve any
// request that are not handling by the routes.
//
// Returns the same instance
func (u *URLHandler) WithErrorHandler(h http.Handler) *URLHandler {
	u.ErrorHandler = h
	return u
}

func (u *URLHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler, ok := u.Handlers[fmt.Sprintf("%s %s", r.Method, r.RequestURI)]
	if ok {
		handler.ServeHTTP(w, r)
		return
	}

	handler, ok = u.Handlers[fmt.Sprintf("ANY %s", r.RequestURI)]
	if ok {
		handler.ServeHTTP(w, r)
	} else {
		u.ErrorHandler.ServeHTTP(w, r)
	}
}
