package Player

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/HandyGold75/GOLib/gotify/lib"
)

type (
	Player struct {
		Send     func(method lib.HttpMethod, action string, options [][2]string, body []byte) ([]byte, error)
		DeviceID string
		Market   string // An ISO 3166-1 alpha-2 country code, https://en.wikipedia.org/wiki/ISO_3166-1_alpha-2
	}

	getPlaybackState struct {
		Device struct {
			ID               string `json:"id"`
			IsActive         bool   `json:"is_active"`
			IsPrivateSession bool   `json:"is_private_session"`
			IsRestricted     bool   `json:"is_restricted"`
			Name             string `json:"name"`
			Type             string `json:"type"`
			VolumePercent    int    `json:"volume_percent"`
			SupportsVolume   bool   `json:"supports_volume"`
		} `json:"device"`
		RepeatState  string `json:"repeat_state"`
		ShuffleState bool   `json:"shuffle_state"`
		Context      struct {
			Type         string `json:"type"`
			Href         string `json:"href"`
			ExternalUrls struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			URI string `json:"uri"`
		} `json:"context"`
		Timestamp  int  `json:"timestamp"`
		ProgressMs int  `json:"progress_ms"`
		IsPlaying  bool `json:"is_playing"`
		Item       struct {
			Album struct {
				AlbumType        string   `json:"album_type"`
				TotalTracks      int      `json:"total_tracks"`
				AvailableMarkets []string `json:"available_markets"`
				ExternalUrls     struct {
					Spotify string `json:"spotify"`
				} `json:"external_urls"`
				Href   string `json:"href"`
				ID     string `json:"id"`
				Images []struct {
					URL    string `json:"url"`
					Height int    `json:"height"`
					Width  int    `json:"width"`
				} `json:"images"`
				Name                 string `json:"name"`
				ReleaseDate          string `json:"release_date"`
				ReleaseDatePrecision string `json:"release_date_precision"`
				Restrictions         struct {
					Reason string `json:"reason"`
				} `json:"restrictions"`
				Type    string `json:"type"`
				URI     string `json:"uri"`
				Artists []struct {
					ExternalUrls struct {
						Spotify string `json:"spotify"`
					} `json:"external_urls"`
					Href string `json:"href"`
					ID   string `json:"id"`
					Name string `json:"name"`
					Type string `json:"type"`
					URI  string `json:"uri"`
				} `json:"artists"`
			} `json:"album"`
			Artists []struct {
				ExternalUrls struct {
					Spotify string `json:"spotify"`
				} `json:"external_urls"`
				Href string `json:"href"`
				ID   string `json:"id"`
				Name string `json:"name"`
				Type string `json:"type"`
				URI  string `json:"uri"`
			} `json:"artists"`
			AvailableMarkets []string `json:"available_markets"`
			DiscNumber       int      `json:"disc_number"`
			DurationMs       int      `json:"duration_ms"`
			Explicit         bool     `json:"explicit"`
			ExternalIds      struct {
				Isrc string `json:"isrc"`
				Ean  string `json:"ean"`
				Upc  string `json:"upc"`
			} `json:"external_ids"`
			ExternalUrls struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			Href         string   `json:"href"`
			ID           string   `json:"id"`
			IsPlayable   bool     `json:"is_playable"`
			LinkedFrom   struct{} `json:"linked_from"`
			Restrictions struct {
				Reason string `json:"reason"`
			} `json:"restrictions"`
			Name        string `json:"name"`
			Popularity  int    `json:"popularity"`
			PreviewURL  string `json:"preview_url"`
			TrackNumber int    `json:"track_number"`
			Type        string `json:"type"`
			URI         string `json:"uri"`
			IsLocal     bool   `json:"is_local"`
		} `json:"item"`
		CurrentlyPlayingType string `json:"currently_playing_type"`
		Actions              struct {
			InterruptingPlayback  bool `json:"interrupting_playback"`
			Pausing               bool `json:"pausing"`
			Resuming              bool `json:"resuming"`
			Seeking               bool `json:"seeking"`
			SkippingNext          bool `json:"skipping_next"`
			SkippingPrev          bool `json:"skipping_prev"`
			TogglingRepeatContext bool `json:"toggling_repeat_context"`
			TogglingShuffle       bool `json:"toggling_shuffle"`
			TogglingRepeatTrack   bool `json:"toggling_repeat_track"`
			TransferringPlayback  bool `json:"transferring_playback"`
		} `json:"actions"`
	}

	getAvailableDevices struct {
		Devices []struct {
			ID               string `json:"id"`
			IsActive         bool   `json:"is_active"`
			IsPrivateSession bool   `json:"is_private_session"`
			IsRestricted     bool   `json:"is_restricted"`
			Name             string `json:"name"`
			Type             string `json:"type"`
			VolumePercent    int    `json:"volume_percent"`
			SupportsVolume   bool   `json:"supports_volume"`
		} `json:"devices"`
	}

	getCurrentlyPlayingTrack struct {
		Device struct {
			ID               string `json:"id"`
			IsActive         bool   `json:"is_active"`
			IsPrivateSession bool   `json:"is_private_session"`
			IsRestricted     bool   `json:"is_restricted"`
			Name             string `json:"name"`
			Type             string `json:"type"`
			VolumePercent    int    `json:"volume_percent"`
			SupportsVolume   bool   `json:"supports_volume"`
		} `json:"device"`
		RepeatState  string `json:"repeat_state"`
		ShuffleState bool   `json:"shuffle_state"`
		Context      struct {
			Type         string `json:"type"`
			Href         string `json:"href"`
			ExternalUrls struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			URI string `json:"uri"`
		} `json:"context"`
		Timestamp  int  `json:"timestamp"`
		ProgressMs int  `json:"progress_ms"`
		IsPlaying  bool `json:"is_playing"`
		Item       struct {
			Album struct {
				AlbumType        string   `json:"album_type"`
				TotalTracks      int      `json:"total_tracks"`
				AvailableMarkets []string `json:"available_markets"`
				ExternalUrls     struct {
					Spotify string `json:"spotify"`
				} `json:"external_urls"`
				Href   string `json:"href"`
				ID     string `json:"id"`
				Images []struct {
					URL    string `json:"url"`
					Height int    `json:"height"`
					Width  int    `json:"width"`
				} `json:"images"`
				Name                 string `json:"name"`
				ReleaseDate          string `json:"release_date"`
				ReleaseDatePrecision string `json:"release_date_precision"`
				Restrictions         struct {
					Reason string `json:"reason"`
				} `json:"restrictions"`
				Type    string `json:"type"`
				URI     string `json:"uri"`
				Artists []struct {
					ExternalUrls struct {
						Spotify string `json:"spotify"`
					} `json:"external_urls"`
					Href string `json:"href"`
					ID   string `json:"id"`
					Name string `json:"name"`
					Type string `json:"type"`
					URI  string `json:"uri"`
				} `json:"artists"`
			} `json:"album"`
			Artists []struct {
				ExternalUrls struct {
					Spotify string `json:"spotify"`
				} `json:"external_urls"`
				Href string `json:"href"`
				ID   string `json:"id"`
				Name string `json:"name"`
				Type string `json:"type"`
				URI  string `json:"uri"`
			} `json:"artists"`
			AvailableMarkets []string `json:"available_markets"`
			DiscNumber       int      `json:"disc_number"`
			DurationMs       int      `json:"duration_ms"`
			Explicit         bool     `json:"explicit"`
			ExternalIds      struct {
				Isrc string `json:"isrc"`
				Ean  string `json:"ean"`
				Upc  string `json:"upc"`
			} `json:"external_ids"`
			ExternalUrls struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			Href         string   `json:"href"`
			ID           string   `json:"id"`
			IsPlayable   bool     `json:"is_playable"`
			LinkedFrom   struct{} `json:"linked_from"`
			Restrictions struct {
				Reason string `json:"reason"`
			} `json:"restrictions"`
			Name        string `json:"name"`
			Popularity  int    `json:"popularity"`
			PreviewURL  string `json:"preview_url"`
			TrackNumber int    `json:"track_number"`
			Type        string `json:"type"`
			URI         string `json:"uri"`
			IsLocal     bool   `json:"is_local"`
		} `json:"item"`
		CurrentlyPlayingType string `json:"currently_playing_type"`
		Actions              struct {
			InterruptingPlayback  bool `json:"interrupting_playback"`
			Pausing               bool `json:"pausing"`
			Resuming              bool `json:"resuming"`
			Seeking               bool `json:"seeking"`
			SkippingNext          bool `json:"skipping_next"`
			SkippingPrev          bool `json:"skipping_prev"`
			TogglingRepeatContext bool `json:"toggling_repeat_context"`
			TogglingShuffle       bool `json:"toggling_shuffle"`
			TogglingRepeatTrack   bool `json:"toggling_repeat_track"`
			TransferringPlayback  bool `json:"transferring_playback"`
		} `json:"actions"`
	}

	getRecentlyPlayedTracks struct {
		Href    string `json:"href"`
		Limit   int    `json:"limit"`
		Next    string `json:"next"`
		Cursors struct {
			After  string `json:"after"`
			Before string `json:"before"`
		} `json:"cursors"`
		Total int `json:"total"`
		Items []struct {
			Track struct {
				Album struct {
					AlbumType        string   `json:"album_type"`
					TotalTracks      int      `json:"total_tracks"`
					AvailableMarkets []string `json:"available_markets"`
					ExternalUrls     struct {
						Spotify string `json:"spotify"`
					} `json:"external_urls"`
					Href   string `json:"href"`
					ID     string `json:"id"`
					Images []struct {
						URL    string `json:"url"`
						Height int    `json:"height"`
						Width  int    `json:"width"`
					} `json:"images"`
					Name                 string `json:"name"`
					ReleaseDate          string `json:"release_date"`
					ReleaseDatePrecision string `json:"release_date_precision"`
					Restrictions         struct {
						Reason string `json:"reason"`
					} `json:"restrictions"`
					Type    string `json:"type"`
					URI     string `json:"uri"`
					Artists []struct {
						ExternalUrls struct {
							Spotify string `json:"spotify"`
						} `json:"external_urls"`
						Href string `json:"href"`
						ID   string `json:"id"`
						Name string `json:"name"`
						Type string `json:"type"`
						URI  string `json:"uri"`
					} `json:"artists"`
				} `json:"album"`
				Artists []struct {
					ExternalUrls struct {
						Spotify string `json:"spotify"`
					} `json:"external_urls"`
					Href string `json:"href"`
					ID   string `json:"id"`
					Name string `json:"name"`
					Type string `json:"type"`
					URI  string `json:"uri"`
				} `json:"artists"`
				AvailableMarkets []string `json:"available_markets"`
				DiscNumber       int      `json:"disc_number"`
				DurationMs       int      `json:"duration_ms"`
				Explicit         bool     `json:"explicit"`
				ExternalIds      struct {
					Isrc string `json:"isrc"`
					Ean  string `json:"ean"`
					Upc  string `json:"upc"`
				} `json:"external_ids"`
				ExternalUrls struct {
					Spotify string `json:"spotify"`
				} `json:"external_urls"`
				Href         string   `json:"href"`
				ID           string   `json:"id"`
				IsPlayable   bool     `json:"is_playable"`
				LinkedFrom   struct{} `json:"linked_from"`
				Restrictions struct {
					Reason string `json:"reason"`
				} `json:"restrictions"`
				Name        string `json:"name"`
				Popularity  int    `json:"popularity"`
				PreviewURL  string `json:"preview_url"`
				TrackNumber int    `json:"track_number"`
				Type        string `json:"type"`
				URI         string `json:"uri"`
				IsLocal     bool   `json:"is_local"`
			} `json:"track"`
			PlayedAt string `json:"played_at"`
			Context  struct {
				Type         string `json:"type"`
				Href         string `json:"href"`
				ExternalUrls struct {
					Spotify string `json:"spotify"`
				} `json:"external_urls"`
				URI string `json:"uri"`
			} `json:"context"`
		} `json:"items"`
	}

	getTheUsersQueue struct {
		CurrentlyPlaying struct {
			Album struct {
				AlbumType        string   `json:"album_type"`
				TotalTracks      int      `json:"total_tracks"`
				AvailableMarkets []string `json:"available_markets"`
				ExternalUrls     struct {
					Spotify string `json:"spotify"`
				} `json:"external_urls"`
				Href   string `json:"href"`
				ID     string `json:"id"`
				Images []struct {
					URL    string `json:"url"`
					Height int    `json:"height"`
					Width  int    `json:"width"`
				} `json:"images"`
				Name                 string `json:"name"`
				ReleaseDate          string `json:"release_date"`
				ReleaseDatePrecision string `json:"release_date_precision"`
				Restrictions         struct {
					Reason string `json:"reason"`
				} `json:"restrictions"`
				Type    string `json:"type"`
				URI     string `json:"uri"`
				Artists []struct {
					ExternalUrls struct {
						Spotify string `json:"spotify"`
					} `json:"external_urls"`
					Href string `json:"href"`
					ID   string `json:"id"`
					Name string `json:"name"`
					Type string `json:"type"`
					URI  string `json:"uri"`
				} `json:"artists"`
			} `json:"album"`
			Artists []struct {
				ExternalUrls struct {
					Spotify string `json:"spotify"`
				} `json:"external_urls"`
				Href string `json:"href"`
				ID   string `json:"id"`
				Name string `json:"name"`
				Type string `json:"type"`
				URI  string `json:"uri"`
			} `json:"artists"`
			AvailableMarkets []string `json:"available_markets"`
			DiscNumber       int      `json:"disc_number"`
			DurationMs       int      `json:"duration_ms"`
			Explicit         bool     `json:"explicit"`
			ExternalIds      struct {
				Isrc string `json:"isrc"`
				Ean  string `json:"ean"`
				Upc  string `json:"upc"`
			} `json:"external_ids"`
			ExternalUrls struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			Href         string   `json:"href"`
			ID           string   `json:"id"`
			IsPlayable   bool     `json:"is_playable"`
			LinkedFrom   struct{} `json:"linked_from"`
			Restrictions struct {
				Reason string `json:"reason"`
			} `json:"restrictions"`
			Name        string `json:"name"`
			Popularity  int    `json:"popularity"`
			PreviewURL  string `json:"preview_url"`
			TrackNumber int    `json:"track_number"`
			Type        string `json:"type"`
			URI         string `json:"uri"`
			IsLocal     bool   `json:"is_local"`
		} `json:"currently_playing"`
		Queue []struct {
			Album struct {
				AlbumType        string   `json:"album_type"`
				TotalTracks      int      `json:"total_tracks"`
				AvailableMarkets []string `json:"available_markets"`
				ExternalUrls     struct {
					Spotify string `json:"spotify"`
				} `json:"external_urls"`
				Href   string `json:"href"`
				ID     string `json:"id"`
				Images []struct {
					URL    string `json:"url"`
					Height int    `json:"height"`
					Width  int    `json:"width"`
				} `json:"images"`
				Name                 string `json:"name"`
				ReleaseDate          string `json:"release_date"`
				ReleaseDatePrecision string `json:"release_date_precision"`
				Restrictions         struct {
					Reason string `json:"reason"`
				} `json:"restrictions"`
				Type    string `json:"type"`
				URI     string `json:"uri"`
				Artists []struct {
					ExternalUrls struct {
						Spotify string `json:"spotify"`
					} `json:"external_urls"`
					Href string `json:"href"`
					ID   string `json:"id"`
					Name string `json:"name"`
					Type string `json:"type"`
					URI  string `json:"uri"`
				} `json:"artists"`
			} `json:"album"`
			Artists []struct {
				ExternalUrls struct {
					Spotify string `json:"spotify"`
				} `json:"external_urls"`
				Href string `json:"href"`
				ID   string `json:"id"`
				Name string `json:"name"`
				Type string `json:"type"`
				URI  string `json:"uri"`
			} `json:"artists"`
			AvailableMarkets []string `json:"available_markets"`
			DiscNumber       int      `json:"disc_number"`
			DurationMs       int      `json:"duration_ms"`
			Explicit         bool     `json:"explicit"`
			ExternalIds      struct {
				Isrc string `json:"isrc"`
				Ean  string `json:"ean"`
				Upc  string `json:"upc"`
			} `json:"external_ids"`
			ExternalUrls struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			Href         string   `json:"href"`
			ID           string   `json:"id"`
			IsPlayable   bool     `json:"is_playable"`
			LinkedFrom   struct{} `json:"linked_from"`
			Restrictions struct {
				Reason string `json:"reason"`
			} `json:"restrictions"`
			Name        string `json:"name"`
			Popularity  int    `json:"popularity"`
			PreviewURL  string `json:"preview_url"`
			TrackNumber int    `json:"track_number"`
			Type        string `json:"type"`
			URI         string `json:"uri"`
			IsLocal     bool   `json:"is_local"`
		} `json:"queue"`
	}
)

func New(send func(method lib.HttpMethod, action string, options [][2]string, body []byte) ([]byte, error)) Player {
	return Player{
		Send:     send,
		DeviceID: "", Market: "",
	}
}

// Scopes: `ScopeUserReadPlaybackState`
func (s *Player) GetPlaybackState() (getPlaybackState, error) {
	res, err := s.Send(lib.GET, "", [][2]string{{"market", s.Market}}, []byte{})
	if err != nil {
		return getPlaybackState{}, err
	}
	data := getPlaybackState{}
	err = json.Unmarshal(res, &data)
	return data, err
}

// Requires premium.
//
// Scopes: `ScopeUserModifyPlaybackState`
func (s *Player) TransferPlayback(deviceID string, play bool) error {
	data, err := json.Marshal(map[string]any{"device_ids": deviceID, "play": play})
	if err != nil {
		return err
	}
	_, err = s.Send(lib.PUT, "", [][2]string{}, data)
	return err
}

// Scopes: `ScopeUserReadPlaybackState`
func (s *Player) GetAvailableDevices() (getAvailableDevices, error) {
	res, err := s.Send(lib.GET, "devices", [][2]string{}, []byte{})
	if err != nil {
		return getAvailableDevices{}, err
	}
	data := getAvailableDevices{}
	err = json.Unmarshal(res, &data)
	return data, err
}

// Scopes: `ScopeUserReadCurrentlyPlaying`
func (s *Player) GetCurrentlyPlayingTrack() (getCurrentlyPlayingTrack, error) {
	res, err := s.Send(lib.GET, "currently-playing", [][2]string{{"market", s.Market}}, []byte{})
	if err != nil {
		return getCurrentlyPlayingTrack{}, err
	}
	data := getCurrentlyPlayingTrack{}
	err = json.Unmarshal(res, &data)
	return data, err
}

// Requires premium.
//
// Scopes: `ScopeUserModifyPlaybackState`
func (s *Player) StartResumePlayback(position time.Duration) error {
	data, err := json.Marshal(map[string]any{"position_ms": strconv.Itoa(int(position.Milliseconds()))})
	if err != nil {
		return err
	}
	_, err = s.Send(lib.PUT, "play", [][2]string{{"device_id", s.DeviceID}}, data)
	return err
}

// Requires premium.
//
// Scopes: `ScopeUserModifyPlaybackState`
func (s *Player) StartResumePlaybackSimple() error {
	_, err := s.Send(lib.PUT, "play", [][2]string{{"device_id", s.DeviceID}}, []byte{})
	return err
}

// Requires premium.
//
// Scopes: `ScopeUserModifyPlaybackState`
//
// Body:
//
//	{
//	    "context_uri": "spotify:album:5ht7ItJgpBH7W6vJ5BqpPr",
//	    "uris": ["spotify:track:4iV5W9uYEdYUVa79Axb7Rh", "spotify:track:1301WleyT98MSxVHPZCA6M"],
//	    "offset": {
//	        "position": 5,
//	        "uri": "spotify:track:1301WleyT98MSxVHPZCA6M"
//	    },
//	    "position_ms": 0
//	}
func (s *Player) StartResumePlaybackRaw(body map[string]any) error {
	data, err := json.Marshal(body)
	if err != nil {
		return err
	}
	_, err = s.Send(lib.PUT, "play", [][2]string{{"device_id", s.DeviceID}}, data)
	return err
}

// Requires premium.
//
// Scopes: `ScopeUserModifyPlaybackState`
func (s *Player) PausePlayback() error {
	_, err := s.Send(lib.PUT, "pause", [][2]string{{"device_id", s.DeviceID}}, []byte{})
	return err
}

// Requires premium.
//
// Scopes: `ScopeUserModifyPlaybackState`
func (s *Player) SkipToNext() error {
	_, err := s.Send(lib.POST, "next", [][2]string{{"device_id", s.DeviceID}}, []byte{})
	return err
}

// Requires premium.
//
// Scopes: `ScopeUserModifyPlaybackState`
func (s *Player) SkipToPrevious() error {
	_, err := s.Send(lib.POST, "previous", [][2]string{{"device_id", s.DeviceID}}, []byte{})
	return err
}

// Requires premium.
//
// Scopes: `ScopeUserModifyPlaybackState`
func (s *Player) SeekToPosition(position time.Duration) error {
	_, err := s.Send(lib.PUT, "seek", [][2]string{{"device_id", s.DeviceID}, {"position_ms", strconv.Itoa(int(position.Milliseconds()))}}, []byte{})
	return err
}

// Requires premium.
//
// Scopes: `ScopeUserModifyPlaybackState`
func (s *Player) SetRepeatMode(state lib.RepeatMode) error {
	_, err := s.Send(lib.PUT, "repeat", [][2]string{{"device_id", s.DeviceID}, {"state", string(state)}}, []byte{})
	return err
}

// Requires premium.
//
// Scopes: `ScopeUserModifyPlaybackState`
func (s *Player) SetPlaybackVolume(volume int) error {
	_, err := s.Send(lib.PUT, "volume", [][2]string{{"device_id", s.DeviceID}, {"volume_percent", strconv.Itoa(max(0, min(100, volume)))}}, []byte{})
	return err
}

// Requires premium.
//
// Scopes: `ScopeUserModifyPlaybackState`
func (s *Player) TogglePlaybackShuffle(state bool) error {
	_, err := s.Send(lib.PUT, "shuffle", [][2]string{{"device_id", s.DeviceID}, {"state", strconv.FormatBool(state)}}, []byte{})
	return err
}

// Scopes: `ScopeUserReadRecentlyPlayed`
//
// Return items after stamp if after is true, otherwise returns items before time.
// Use `time.Time{}` to disable this filter.
func (s *Player) GetRecentlyPlayedTracks(limit int, stamp time.Time, after bool) (getRecentlyPlayedTracks, error) {
	key, value := "before", strconv.Itoa(int(stamp.Unix()))
	if stamp.Unix() == (time.Time{}.Unix()) {
		value = ""
	} else if after {
		key = "after"
	}
	res, err := s.Send(lib.GET, "recently-played", [][2]string{{"limit", strconv.Itoa(max(1, min(50, limit)))}, {key, value}}, []byte{})
	if err != nil {
		return getRecentlyPlayedTracks{}, err
	}
	data := getRecentlyPlayedTracks{}
	err = json.Unmarshal(res, &data)
	return data, err
}

// Scopes: `ScopeUserReadCurrentlyPlaying`, `ScopeUserReadPlaybackState`
func (s *Player) GetTheUsersQueue() (getTheUsersQueue, error) {
	res, err := s.Send(lib.GET, "queue", [][2]string{}, []byte{})
	if err != nil {
		return getTheUsersQueue{}, err
	}
	data := getTheUsersQueue{}
	err = json.Unmarshal(res, &data)
	return data, err
}

// Requires premium.
//
// Scopes: `ScopeUserModifyPlaybackState`
func (s *Player) AddItemToPlaybackQueue(uri string) error {
	_, err := s.Send(lib.POST, "", [][2]string{{"device_id", s.DeviceID}, {"uri", uri}}, []byte{})
	return err
}
