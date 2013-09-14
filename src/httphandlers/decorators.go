package httphandlers

import (
	"net/http"
)

type RequestRecorderHandler struct {
	Handler  http.Handler
	Requests []http.Request
}

func NewRequestRecorderHandler(h http.Handler) *RequestRecorderHandler {
	return &RequestRecorderHandler{h, make([]http.Request, 0)}
}

func (h *RequestRecorderHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Requests = append(h.Requests, *r)
	if h.Handler != nil {
		h.Handler.ServeHTTP(w, r)
	}
}

func (h *RequestRecorderHandler) ClearRequests() {
	h.Requests = make([]http.Request, 0)
}

func (h *RequestRecorderHandler) FilterRequests(fn func(r http.Request) bool) []http.Request {
	filtered := make([]http.Request, 0)
	for _, r := range h.Requests {
		if fn(r) {
			filtered = append(filtered, r)
		}
	}
	return filtered
}

func (h *RequestRecorderHandler) RequestsByPath(path string) []http.Request {
	return h.FilterRequests(func(r http.Request) bool {
		return r.URL.Path == path
	})
}

func (h *RequestRecorderHandler) RequestsByMethod(method string) []http.Request {
	return h.FilterRequests(func(r http.Request) bool {
		return r.Method == method
	})
}
