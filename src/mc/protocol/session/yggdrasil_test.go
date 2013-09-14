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

func fakeSession() (*httptest.Server, *YggdrasilClient, *httphandlers.RequestRecorderHandler) {
	authHandler := fixture("fixtures/yggdrasil_authenticate.json")
	refreshHandler := fixture("fixtures/yggdrasil_refresh.json")
	validateHandler := httphandlers.NewFixtureFromString("")
	routes := httphandlers.NewURLHandler(map[string]http.Handler{
		"POST /authenticate": requireJson(authHandler),
		"POST /refresh":      requireJson(refreshHandler),
		"POST /validate":     requireJson(validateHandler),
	})
	h := httphandlers.NewRequestRecorderHandler(routes)
	server := httptest.NewServer(h)
	session := NewYggdrasilClient()
	session.URL = server.URL
	return server, session, h
}

func TestSessionServerAuthenticate(t *testing.T) {
	it := NewIt(t)
	server, session, recorder := fakeSession()
	defer server.Close()
	token, err := session.Authenticate("username", "password")
	it.Must(err)

	it.Expects(token, ToEqual, &YggdrasilSession{
		AccessToken: "deadbeef",
		ClientToken: "clientID",
	})

	requests := recorder.RequestsByPath("/authenticate")
	it.Expects(requests, Not(ToBeEmpty))
	r := requests[0]
	it.Expects(r.Method, ToEqual, "POST")

	auth := &authRequest{}
	it.Must(json.Unmarshal(r.Body, auth))
	it.Expects(auth.Username, ToEqual, "username")
	it.Expects(auth.Password, ToEqual, "password")
	it.Expects(auth.ClientToken, ToEqual, "")
	it.Expects(auth.Agent.Name, Not(ToEqual), "")
	it.Expects(auth.Agent.Version, Not(ToEqual), 0)
}

type authRequest struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	ClientToken string `json:"clientToken,omitempty"`
	Agent       struct {
		Name    string `json:"name"`
		Version int    `json:"version"`
	} `json:"agent"`
}

func TestSessionServiceRefresh(t *testing.T) {
	it := NewIt(t)
	server, session, recorder := fakeSession()
	defer server.Close()
	token := &YggdrasilSession{
		AccessToken: "accessToken",
		ClientToken: "clientToken",
	}
	it.Must(session.Refresh(token))

	it.Expects(token, ToEqual, &YggdrasilSession{
		AccessToken: "deadbeef",
		ClientToken: "clientID",
	})

	requests := recorder.RequestsByPath("/refresh")
	it.Expects(requests, Not(ToBeEmpty))
	r := requests[0]
	it.Expects(r.Method, ToEqual, "POST")

	auth := &refreshRequest{}
	it.Must(json.Unmarshal(r.Body, auth))
	it.Expects(auth.AccessToken, ToEqual, "accessToken")
	it.Expects(auth.ClientToken, ToEqual, "clientToken")
}

type refreshRequest struct {
	AccessToken string `json:"accessToken"`
	ClientToken string `json:"clientToken"`
}

func TestSessionServiceValidate(t *testing.T) {
	it := NewIt(t)
	server, session, recorder := fakeSession()
	defer server.Close()
	token := &YggdrasilSession{
		AccessToken: "accessToken",
		ClientToken: "clientToken",
	}
	it.Must(session.Validate(token))

	requests := recorder.RequestsByPath("/validate")
	it.Expects(requests, Not(ToBeEmpty))
	r := requests[0]
	it.Expects(r.Method, ToEqual, "POST")

	auth := &refreshRequest{}
	it.Must(json.Unmarshal(r.Body, auth))
	it.Expects(auth.AccessToken, ToEqual, "accessToken")
	it.Expects(auth.ClientToken, ToEqual, "")
}
