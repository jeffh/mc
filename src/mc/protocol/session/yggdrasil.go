// handles session creation to minecraft.net
//
// sessions are used to verify the authenticity of
// a minecraft account (aka, DRM). It is used by
// servers to ban clients, which have be the name
// registered with minecraft.net.
package session

import (
	"fmt"
)

const (
	YggdrasilUrl = "https://authserver.mojang.com"
)

type YggdrasilSession struct {
	Url   string
	Agent YggdrasilAgent
}

func (s *YggdrasilSession) Authenticate(username, password string) {
}

func (s *YggdrasilSession) Refresh(accessToken, clientToken string) {
}

func (s *YggdrasilSession) Validate(accessToken string) bool {
}

type yggdrasilAgent struct {
	Name    string `json:"name"`
	Version int    `json:"version"`
}

type yggdrasilProfile struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type yggdrasilAuthenticateRequest struct {
	Agent       yggdrasilAgent `json:"agent"`
	Username    string         `json:"username"`
	Password    string         `json:"password"`
	ClientToken string         `json:"clientToken"`
}

type yggdrasilAuthenticateResponse struct {
	AccessToken       string
	ClientToken       string
	AvailableProfiles []yggdrasilProfile
	SelectedProfile   yggdrasilProfile
}

type YggdrasilError struct {
	Error        string `json:"error"`
	ErrorMessage string `json:"errorMessage"`
	Cause        string `json:"cause"`
}

func (e *YggdrasilError) Error() string {
	return fmt.Sprintf("%s - %s: %s", e.Error, e.ErrorMessage, e.Cause)
}
