package Users

import (
	"encoding/json"

	"github.com/HandyGold75/GOLib/gotify/lib"
)

type (
	Users struct {
		Send     func(method lib.HttpMethod, action string, options [][2]string, body []byte) ([]byte, error)
		DeviceID string
	}

	currentUserResponse struct {
		Country         string `json:"country"`
		DisplayName     string `json:"display_name"`
		Email           string `json:"email"`
		ExplicitContent struct {
			FilterEnabled bool `json:"filter_enabled"`
			FilterLocked  bool `json:"filter_locked"`
		} `json:"explicit_content"`
		ExternalUrls struct {
			Spotify string `json:"spotify"`
		} `json:"external_urls"`
		Followers struct {
			Href  string `json:"href"`
			Total int    `json:"total"`
		} `json:"followers"`
		Href   string `json:"href"`
		ID     string `json:"id"`
		Images []struct {
			URL    string `json:"url"`
			Height int    `json:"height"`
			Width  int    `json:"width"`
		} `json:"images"`
		Product string `json:"product"`
		Type    string `json:"type"`
		URI     string `json:"uri"`
	}
)

func New(send func(method lib.HttpMethod, action string, options [][2]string, body []byte) ([]byte, error)) Users {
	return Users{Send: send, DeviceID: ""}
}

// Scopes: `ScopeUserReadPrivate` (optional), `ScopeUserReadEmail` (optional)
func (s *Users) GetCurrentUser() (currentUserResponse, error) {
	res, err := s.Send(lib.GET, "me", [][2]string{}, []byte{})
	if err != nil {
		return currentUserResponse{}, err
	}
	data := currentUserResponse{}
	err = json.Unmarshal(res, &data)
	return data, err
}
