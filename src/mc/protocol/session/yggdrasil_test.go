package session

import (
	. "github.com/jeffh/goexpect"
	"httphandlers"
	"net/http"
	"net/http/httptest"
	"testing"
)

func requireJson(h http.Handler) http.Handler {
	return httphandlers.NewRequiredHeadersHandler(http.Header{
		"content-type": []string{"application/json"},
	}, h)
}

func fixture(path string) *httphandlers.FixtureHandler {
	handler := httphandlers.Required(httphandlers.NewFixtureFromFile(path))
	return handler.(*httphandlers.FixtureHandler)
}

func fakeSession() (*httptest.Server, *YggdrasilSession) {
	authHandler := fixture("fixtures/yggdrasil_authenticate.json")
	refreshHandler := fixture("fixtures/yggdrasil_refresh.json")
	validateHandler := fixture("fixtures/yggdrasil_validate.json")
	h := httphandlers.NewURLHandler(map[string]http.Handler{
		"POST /authenticate": requireJson(authHandler),
		"POST /refresh":      requireJson(refreshHandler),
		"POST /validate":     requireJson(validateHandler),
	})
	server := httptest.NewTLSServer(h)
	session := NewYggdrasilSession()
	session.URL = server.URL
	return server, session
}

func TestSessionServerAuthenticate(t *testing.T) {
	server, session := fakeSession()
	defer server.Close()
	Must(t, session.Authenticate("username", "password"))
}
