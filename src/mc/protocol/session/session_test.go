package session

import (
	. "github.com/jeffh/goexpect"
	"httphandlers"
	"net/http/httptest"
	"testing"
)

func clientAndServerPair() (*SessionClient, *httptest.Server, *httphandlers.RequestRecorderHandler) {
	client := NewSessionClient()
	responseHandler := httphandlers.NewFixtureFromString("OK")
	h := httphandlers.NewRequestRecorderHandler(responseHandler)
	server := httptest.NewServer(h)
	client.URL = server.URL
	return client, server, h
}

func TestSessionClientCanJoinServer(t *testing.T) {
	it := NewIt(t)
	client, server, recorder := clientAndServerPair()
	defer server.Close()

	it.Must(client.JoinServer(ServerInfo{
		Username:     "John",
		SessionID:    "sessionID",
		ServerID:     "myServer",
		SharedSecret: []byte("secret"),
		PublicKey:    []byte("publicKey"),
	}))

	requests := recorder.RequestsByPath("/")
	it.Expects(requests, ToBeLengthOf, 1)
	r := requests[0]
	it.Expects(r.Method, ToEqual, "GET")
	it.Expects(r.URL.Query().Get("user"), ToEqual, "John")
	it.Expects(r.URL.Query().Get("sessionId"), ToEqual, "sessionID")
	it.Expects(r.URL.Query().Get("serverId"), ToEqual, "-f6217b3fe196685c9cfef5eea9a02125855af37")
}
