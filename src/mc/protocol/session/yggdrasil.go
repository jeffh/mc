// handles session creation to minecraft's login servers
//
// sessions are used to verify the authenticity of
// a minecraft account (like, DRM). It is used by
// servers to ban clients by their name that is
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
	YggdrasilURL           = "https://authserver.mojang.com"
	YggdrasilClientName    = "Minecraft"
	YggdrasilClientVersion = 1
)

var YggdrasilDefaultAgent = &YggdrasilAgent{
	Name:    YggdrasilClientName,
	Version: YggdrasilClientVersion,
}

type YggdrasilClient struct {
	URL    string
	Agent  YggdrasilAgent
	Client *http.Client
}

type YggdrasilSession struct {
	AccessToken string
	ClientToken string
	ProfileID   string
}

func (s *YggdrasilSession) SessionID() string {
	return fmt.Sprintf("token:%s:%s", s.AccessToken, s.ProfileID)
}

func NewYggdrasilClient() *YggdrasilClient {
	return &YggdrasilClient{
		URL: YggdrasilURL,
		Agent: YggdrasilAgent{
			Name:    YggdrasilClientName,
			Version: YggdrasilClientVersion,
		},
		Client: &http.Client{},
	}
}

func (s *YggdrasilClient) fullPath(path string) string {
	uri, err := url.Parse(s.URL + path)
	if err != nil {
		panic(err)
	}
	return uri.String()
}

func (s *YggdrasilClient) post(path string, data interface{}, respData interface{}) (*http.Response, error) {
	b := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(b)
	err := encoder.Encode(data)
	if err != nil {
		return nil, err
	}
	resp, err := s.Client.Post(s.fullPath(path), "application/json", b)
	if err == nil && respData != nil {
		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(respData)
	}
	return resp, err
}

func (s *YggdrasilClient) Authenticate(username, password string) (*YggdrasilSession, error) {
	data := &yggdrasilAuthenticateRequest{
		Agent:    *YggdrasilDefaultAgent,
		Username: username,
		Password: password,
	}
	authResponse := &yggdrasilResponse{}
	_, err := s.post("/authenticate", data, authResponse)
	if err != nil {
		return nil, err
	}

	if authResponse.IsError() {
		return nil, &YggdrasilError{
			ErrorCode:    authResponse.ErrorCode,
			ErrorMessage: authResponse.ErrorMessage,
			Cause:        authResponse.Cause,
		}
	}

	token := &YggdrasilSession{
		AccessToken: authResponse.AccessToken,
		ClientToken: authResponse.ClientToken,
		ProfileID:   authResponse.SelectedProfile.Id,
	}
	return token, err
}

func (s *YggdrasilClient) Refresh(token *YggdrasilSession) error {
	data := &yggdrasilRefreshRequest{
		AccessToken: token.AccessToken,
		ClientToken: token.ClientToken,
	}
	authResponse := &yggdrasilResponse{}
	_, err := s.post("/refresh", data, authResponse)
	if err != nil {
		return err
	}

	if authResponse.IsError() {
		return &YggdrasilError{
			ErrorCode:    authResponse.ErrorCode,
			ErrorMessage: authResponse.ErrorMessage,
			Cause:        authResponse.Cause,
		}
	}

	token.AccessToken = authResponse.AccessToken
	token.ClientToken = authResponse.ClientToken
	return err
}

func (s *YggdrasilClient) Validate(token *YggdrasilSession) error {
	data := &yggdrasilRefreshRequest{
		AccessToken: token.AccessToken,
	}
	resp, err := s.post("/validate", data, nil)
	if err != nil {
		return nil
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Auth token has expired")
	}
	return nil
}

type YggdrasilAgent struct {
	Name    string `json:"name"`
	Version int    `json:"version"`
}

type YggdrasilProfile struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type yggdrasilAuthenticateRequest struct {
	Agent       YggdrasilAgent `json:"agent",omitempty`
	Username    string         `json:"username"`
	Password    string         `json:"password"`
	ClientToken string         `json:"clientToken,omitempty"`
}

type yggdrasilRefreshRequest struct {
	AccessToken string `json:"accessToken"`
	ClientToken string `json:"clientToken,omitempty"`
}

type yggdrasilResponse struct {
	AccessToken       string             `json:"accessToken"`
	ClientToken       string             `json:"clientToken"`
	AvailableProfiles []YggdrasilProfile `json:"availableProfiles"`
	SelectedProfile   YggdrasilProfile   `json:"selectedProfile"`

	ErrorCode    string `json:"error"`
	ErrorMessage string `json:"errorMessage"`
	Cause        string `json:"cause"`
}

func (y *yggdrasilResponse) IsError() bool {
	return y.ErrorCode != ""
}

type YggdrasilError struct {
	ErrorCode    string
	ErrorMessage string
	Cause        string
}

func (y *YggdrasilError) Error() string {
	return fmt.Sprintf("%s - %s: %s", y.ErrorCode, y.ErrorMessage, y.Cause)
}
