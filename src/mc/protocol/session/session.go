package session

import (
	"net/http"
	"net/url"
)

const (
	SessionURL     = "http://session.minecraft.net/game/joinserver.jsp"
	sessionUserKey = "user"
	sessionIdKey   = "sessionId"
	sessionHashKey = "serverId"
)

type Client interface {
	JoinServer(info ServerInfo) error
}

type NullClient struct{}

func (n *NullClient) JoinServer(s ServerInfo) error {
	return nil
}

type RecorderClient struct {
	JoinRequests []ServerInfo
}

func NewRecorderClient() *RecorderClient {
	return &RecorderClient{make([]ServerInfo, 0)}
}

func (r *RecorderClient) JoinServer(s ServerInfo) error {
	r.JoinRequests = append(r.JoinRequests, s)
	return nil
}

type SessionClient struct {
	URL          string
	SessionIdKey string
	UserKey      string
	HashKey      string
	Client       *http.Client
	sessionID    string
}

type ServerInfo struct {
	Username     string
	ServerID     string
	SessionID    string
	SharedSecret []byte
	PublicKey    []byte
}

func (i *ServerInfo) serverHash() string {
	b := []byte(i.ServerID)
	b = append(b, i.SharedSecret...)
	b = append(b, i.PublicKey...)
	return sha1HexDigest(b)
}

func NewSessionClient() *SessionClient {
	return &SessionClient{
		URL:          SessionURL,
		UserKey:      sessionUserKey,
		SessionIdKey: sessionIdKey,
		HashKey:      sessionHashKey,
		Client:       &http.Client{},
	}
}

func (s *SessionClient) SetSessionID(sid string) {
	s.sessionID = sid
}

func (s *SessionClient) JoinServer(info ServerInfo) error {
	uri, err := url.Parse(s.URL)
	if err != nil {
		return err
	}
	query := uri.Query()
	query.Add(s.SessionIdKey, info.SessionID)
	query.Add(s.UserKey, info.Username)
	query.Add(s.HashKey, info.serverHash())
	uri.RawQuery = query.Encode()
	_, err = s.Client.Get(uri.String())
	return err
}
