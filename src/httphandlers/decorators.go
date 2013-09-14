package httphandlers

import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
)

// A mostly-clone of http.Request.
//
// except Body is a byte array
type RecordedRequest struct {
	Method           string
	URL              url.URL
	Proto            string
	ProtoMajor       int
	ProtoMinor       int
	Header           http.Header
	Body             []byte
	ContentLength    int64
	TransferEncoding []string
	Host             string
	Form             url.Values
	PostForm         url.Values
	MultipartForm    *multipart.Form
	Trailer          http.Header
	RemoteAddr       string
	RequestURI       string
	TLS              *tls.ConnectionState
}

// Creates a RecordedRequest from the given http.Request.
//
// Using this consumes the http.Request object because the
// Body is consumed. Use Request() to generate a new
// http.Request from the recorded request.
func NewRecordedRequest(r *http.Request) (*RecordedRequest, error) {
	var body []byte
	if r.Body != nil {
		var err error
		body, err = ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}
	} else {
		body = []byte{}
	}

	return &RecordedRequest{
		Method:           r.Method,
		URL:              *r.URL,
		Proto:            r.Proto,
		ProtoMajor:       r.ProtoMajor,
		ProtoMinor:       r.ProtoMinor,
		Header:           r.Header,
		Body:             body,
		ContentLength:    r.ContentLength,
		TransferEncoding: r.TransferEncoding,
		Host:             r.Host,
		Form:             r.Form,
		PostForm:         r.PostForm,
		MultipartForm:    r.MultipartForm,
		Trailer:          r.Trailer,
		RemoteAddr:       r.RemoteAddr,
		RequestURI:       r.RequestURI,
		TLS:              r.TLS,
	}, nil
}

// Generates a new http request object from the previous one
func (r *RecordedRequest) Request() *http.Request {
	url := r.URL
	return &http.Request{
		Method:           r.Method,
		URL:              &url,
		Proto:            r.Proto,
		ProtoMajor:       r.ProtoMajor,
		ProtoMinor:       r.ProtoMinor,
		Header:           r.Header,
		Body:             ioutil.NopCloser(bytes.NewBuffer(r.Body)),
		ContentLength:    r.ContentLength,
		TransferEncoding: r.TransferEncoding,
		Host:             r.Host,
		Form:             r.Form,
		PostForm:         r.PostForm,
		MultipartForm:    r.MultipartForm,
		Trailer:          r.Trailer,
		RemoteAddr:       r.RemoteAddr,
		RequestURI:       r.RequestURI,
		TLS:              r.TLS,
	}
}

type RequestRecorderHandler struct {
	Handler  http.Handler
	Requests []*RecordedRequest
}

func NewRequestRecorderHandler(h http.Handler) *RequestRecorderHandler {
	return &RequestRecorderHandler{h, make([]*RecordedRequest, 0)}
}

func (h *RequestRecorderHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rr, err := NewRecordedRequest(r)
	if err != nil {
		panic(err)
	}
	h.Requests = append(h.Requests, rr)
	if h.Handler != nil {
		h.Handler.ServeHTTP(w, r)
	}
}

func (h *RequestRecorderHandler) ClearRequests() {
	h.Requests = make([]*RecordedRequest, 0)
}

func (h *RequestRecorderHandler) FilterRequests(fn func(r *RecordedRequest) bool) []*RecordedRequest {
	filtered := make([]*RecordedRequest, 0)
	for _, r := range h.Requests {
		if fn(r) {
			filtered = append(filtered, r)
		}
	}
	return filtered
}

func (h *RequestRecorderHandler) RequestsByPath(path string) []*RecordedRequest {
	return h.FilterRequests(func(r *RecordedRequest) bool {
		return r.URL.Path == path
	})
}

func (h *RequestRecorderHandler) RequestsByMethod(method string) []*RecordedRequest {
	return h.FilterRequests(func(r *RecordedRequest) bool {
		return r.Method == method
	})
}
