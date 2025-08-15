package Users

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/HandyGold75/GOLib/gotify/lib"
)

type (
	Users struct {
		Send     func(method lib.HttpMethod, action string, options [][2]string, body []byte) ([]byte, error)
		DeviceID string
	}

	getCurrentUsersProfile struct {
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

	getUsersTopArtists struct {
		Href     string             `json:"href"`
		Limit    int                `json:"limit"`
		Next     string             `json:"next"`
		Offset   int                `json:"offset"`
		Previous string             `json:"previous"`
		Total    int                `json:"total"`
		Items    []lib.ArtistObject `json:"items"`
	}

	getUsersTopTracks struct {
		Href     string            `json:"href"`
		Limit    int               `json:"limit"`
		Next     string            `json:"next"`
		Offset   int               `json:"offset"`
		Previous string            `json:"previous"`
		Total    int               `json:"total"`
		Items    []lib.TrackObject `json:"items"`
	}

	getUsersProfile struct {
		DisplayName  string `json:"display_name"`
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
		Type string `json:"type"`
		URI  string `json:"uri"`
	}

	getFollowedArtists struct {
		Artists struct {
			Href    string `json:"href"`
			Limit   int    `json:"limit"`
			Next    string `json:"next"`
			Cursors struct {
				After  string `json:"after"`
				Before string `json:"before"`
			} `json:"cursors"`
			Total int `json:"total"`
			Items []struct {
				ExternalUrls struct {
					Spotify string `json:"spotify"`
				} `json:"external_urls"`
				Followers struct {
					Href  string `json:"href"`
					Total int    `json:"total"`
				} `json:"followers"`
				Genres []string `json:"genres"`
				Href   string   `json:"href"`
				ID     string   `json:"id"`
				Images []struct {
					URL    string `json:"url"`
					Height int    `json:"height"`
					Width  int    `json:"width"`
				} `json:"images"`
				Name       string `json:"name"`
				Popularity int    `json:"popularity"`
				Type       string `json:"type"`
				URI        string `json:"uri"`
			} `json:"items"`
		} `json:"artists"`
	}
)

func New(send func(method lib.HttpMethod, action string, options [][2]string, body []byte) ([]byte, error)) Users {
	return Users{Send: send, DeviceID: ""}
}

// Scopes: `ScopeUserReadPrivate`, `ScopeUserReadEmail`
func (s *Users) GetCurrentUsersProfile() (getCurrentUsersProfile, error) {
	res, err := s.Send(lib.GET, "me", [][2]string{}, []byte{})
	if err != nil {
		return getCurrentUsersProfile{}, err
	}
	data := getCurrentUsersProfile{}
	err = json.Unmarshal(res, &data)
	return data, err
}

// Scopes: `ScopeUserTopRead`
func (s *Users) GetUsersTopArtists(time lib.TimeRange, limit, offset int) (getUsersTopArtists, error) {
	res, err := s.Send(lib.GET, "me/top/artists", [][2]string{{"time_range", string(time)}, {"limit", strconv.Itoa(max(1, min(50, limit)))}, {"offset", strconv.Itoa(max(0, offset))}}, []byte{})
	if err != nil {
		return getUsersTopArtists{}, err
	}
	data := getUsersTopArtists{}
	err = json.Unmarshal(res, &data)
	return data, err
}

// Scopes: `ScopeUserTopRead`
func (s *Users) GetUsersTopTracks(time lib.TimeRange, limit, offset int) (getUsersTopTracks, error) {
	res, err := s.Send(lib.GET, "me/top/tracks", [][2]string{{"time_range", string(time)}, {"limit", strconv.Itoa(max(1, min(50, limit)))}, {"offset", strconv.Itoa(max(0, offset))}}, []byte{})
	if err != nil {
		return getUsersTopTracks{}, err
	}
	data := getUsersTopTracks{}
	err = json.Unmarshal(res, &data)
	return data, err
}

// Scopes:
func (s *Users) GetUsersProfile(id string) (getUsersProfile, error) {
	res, err := s.Send(lib.GET, "users/"+id, [][2]string{}, []byte{})
	if err != nil {
		return getUsersProfile{}, err
	}
	data := getUsersProfile{}
	err = json.Unmarshal(res, &data)
	return data, err
}

// Scopes: `ScopePlaylistModifyPublic`, `ScopePlaylistModifyPrivate`
func (s *Users) FollowPlaylist(id string, public bool) error {
	data, err := json.Marshal(map[string]any{"public": public})
	if err != nil {
		return err
	}
	_, err = s.Send(lib.PUT, "playlists/"+id+"/followers", [][2]string{}, data)
	return err
}

// Scopes: `ScopePlaylistModifyPublic`, `ScopePlaylistModifyPrivate`
func (s *Users) UnfollowPlaylist(id string) error {
	_, err := s.Send(lib.DELETE, "playlists/"+id+"/followers", [][2]string{}, []byte{})
	return err
}

// Scopes: `ScopeUserFollowRead`
func (s *Users) GetFollowedArtists(after string, limit int) (getFollowedArtists, error) {
	res, err := s.Send(lib.GET, "me/following", [][2]string{{"type", "artists"}, {"after", after}, {"limit", strconv.Itoa(max(1, min(50, limit)))}}, []byte{})
	if err != nil {
		return getFollowedArtists{}, err
	}
	data := getFollowedArtists{}
	err = json.Unmarshal(res, &data)
	return data, err
}

// Scopes: `ScopeUserFollowModify`
func (s *Users) FollowArtists(ids []string) error {
	_, err := s.Send(lib.PUT, "me/following", [][2]string{{"type", "artist"}, {"ids", strings.Join(ids, ",")}}, []byte{})
	return err
}

// Scopes: `ScopeUserFollowModify`
func (s *Users) FollowUsers(ids []string) error {
	_, err := s.Send(lib.PUT, "me/following", [][2]string{{"type", "user"}, {"ids", strings.Join(ids, ",")}}, []byte{})
	return err
}

// Scopes: `ScopeUserFollowModify`
func (s *Users) UnfollowArtists(ids []string) error {
	_, err := s.Send(lib.DELETE, "me/following", [][2]string{{"type", "artist"}, {"ids", strings.Join(ids, ",")}}, []byte{})
	return err
}

// Scopes: `ScopeUserFollowModify`
func (s *Users) UnfollowUsers(ids []string) error {
	_, err := s.Send(lib.DELETE, "me/following", [][2]string{{"type", "user"}, {"ids", strings.Join(ids, ",")}}, []byte{})
	return err
}

// Scopes: `ScopeUserFollowRead`
func (s *Users) CheckIfUserFollowsArtists(ids []string) ([]bool, error) {
	res, err := s.Send(lib.GET, "me/following/contains", [][2]string{{"type", "artist"}, {"ids", strings.Join(ids, ",")}}, []byte{})
	if err != nil {
		return []bool{}, err
	}
	data := []bool{}
	err = json.Unmarshal(res, &data)
	return data, err
}

// Scopes: `ScopeUserFollowRead`
func (s *Users) CheckIfUserFollowsUsers(ids []string) ([]bool, error) {
	res, err := s.Send(lib.GET, "me/following/contains", [][2]string{{"type", "user"}, {"ids", strings.Join(ids, ",")}}, []byte{})
	if err != nil {
		return []bool{}, err
	}
	data := []bool{}
	err = json.Unmarshal(res, &data)
	return data, err
}

// Scopes:
func (s *Users) CheckIfCurrentUserFollowsPlaylist(id string) (bool, error) {
	res, err := s.Send(lib.GET, "playlists/"+id+"/followers/contains", [][2]string{}, []byte{})
	if err != nil {
		return false, err
	}
	data := []bool{}
	err = json.Unmarshal(res, &data)
	if len(data) < 1 {
		return false, lib.Errors.UnexpectedResponse
	}
	return data[0], err
}
