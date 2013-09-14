// handles session creation to minecraft.net
//
// sessions are used to verify the authenticity of
// a minecraft account (aka, DRM). It is used by
// servers to ban clients, which have be the name
// registered with minecraft.net.
package session

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	YggdrasilURL     = "https://authserver.mojang.com"
	YggdrasilClient  = "Minecraft"
	YggdrasilVersion = 1
)

type YggdrasilSession struct {
	URL    string
	Agent  YggdrasilAgent
	Client *http.Client
}

func NewYggdrasilSession() *YggdrasilSession {
	return &YggdrasilSession{
		URL: YggdrasilURL,
		Agent: YggdrasilAgent{
			Name:    YggdrasilClient,
			Version: YggdrasilVersion,
		},
		Client: &http.Client{},
	}
}

func (s *YggdrasilSession) fullPath(path string) string {
	uri, err := url.Parse(s.URL + path)
	if err != nil {
		panic(err)
	}
	return uri.String()
}

func (s *YggdrasilSession) post(path string, data interface{}) (*http.Response, error) {
	b := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(b)
	err := encoder.Encode(data)
	if err != nil {
		return nil, err
	}
	return s.Client.Post(s.fullPath(path), "application/json", b)
}

func (s *YggdrasilSession) Authenticate(username, password string) error {
	data := &yggdrasilAuthenticateRequest{
		Agent:    *YggdrasilDefaultAgent,
		Username: username,
		Password: password,
	}
	_, err := s.post("/authenticate", data)
	return err
}

func (s *YggdrasilSession) Refresh(accessToken, clientToken string) error {
	return nil
}

func (s *YggdrasilSession) Validate(accessToken string) bool {
	return false
}

type YggdrasilAgent struct {
	Name    string `json:"name"`
	Version int    `json:"version"`
}

var YggdrasilDefaultAgent = &YggdrasilAgent{
	Name:    YggdrasilClient,
	Version: YggdrasilVersion,
}

type YggdrasilProfile struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type yggdrasilAuthenticateRequest struct {
	Agent       YggdrasilAgent `json:"agent"`
	Username    string         `json:"username"`
	Password    string         `json:"password"`
	ClientToken string         `json:"clientToken,omitempty"`
}

type yggdrasilAuthenticateResponse struct {
	AccessToken       string
	ClientToken       string
	AvailableProfiles []YggdrasilProfile
	SelectedProfile   YggdrasilProfile
}

type YggdrasilError struct {
	ErrorCode    string `json:"error"`
	ErrorMessage string `json:"errorMessage"`
	Cause        string `json:"cause"`
}

func (e *YggdrasilError) Error() string {
	return fmt.Sprintf("%s - %s: %s", e.ErrorCode, e.ErrorMessage, e.Cause)
}
