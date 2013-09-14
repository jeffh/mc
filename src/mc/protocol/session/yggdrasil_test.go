package session

import (
	"encoding/json"
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

func fakeSession() (*httptest.Server, *YggdrasilSession, *httphandlers.RequestRecorderHandler) {
	authHandler := fixture("fixtures/yggdrasil_authenticate.json")
	refreshHandler := fixture("fixtures/yggdrasil_refresh.json")
	validateHandler := fixture("fixtures/yggdrasil_validate.json")
	routes := httphandlers.NewURLHandler(map[string]http.Handler{
		"POST /authenticate": requireJson(authHandler),
		"POST /refresh":      requireJson(refreshHandler),
		"POST /validate":     requireJson(validateHandler),
	})
	h := httphandlers.NewRequestRecorderHandler(routes)
	server := httptest.NewServer(h)
	session := NewYggdrasilSession()
	session.URL = server.URL
	return server, session, h
}

func TestSessionServerAuthenticate(t *testing.T) {
	it := NewIt(t)
	server, session, recorder := fakeSession()
	defer server.Close()
	Must(t, session.Authenticate("username", "password"))

	it.Expects(recorder.RequestsByPath("/authenticate"), Not(ToBeEmpty))
	r := recorder.RequestsByPath("/authenticate")[0]
	it.Expects(r.Method, ToEqual, "POST")

	auth := &authRequest{}
	it.Must(json.Unmarshal(r.Body, auth))
	it.Expects(auth.Username, ToEqual, "username")
	it.Expects(auth.Password, ToEqual, "password")
	it.Expects(auth.Agent.Name, Not(ToEqual), "")
	it.Expects(auth.Agent.Version, Not(ToEqual), 0)
}

type authRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Agent    struct {
		Name    string `json:"name"`
		Version int    `json:"version"`
	} `json:"agent"`
}

/*
{
  "agent": {                             // optional
    "name": "Minecraft",                 // So far this is the only encountered value
    "version": 1                         // This number might be increased
                                         // by the vanilla client in the future
  },
  "username": "mojang account name",     // Can be an email address or player name for
                                         // unmigrated accounts
  "password": "mojang account password",
  "clientToken": "client identifier"     // optional
}
*/
