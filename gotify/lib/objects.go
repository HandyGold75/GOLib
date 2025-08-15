package lib

type (
	Context struct {
		Context struct {
			Type         string `json:"type"`
			Href         string `json:"href"`
			ExternalUrls struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			URI string `json:"uri"`
		} `json:"context"`
	}

	Cursors struct {
		Cursors struct {
			After  string `json:"after"`
			Before string `json:"before"`
		} `json:"cursors"`
	}

	Actions struct {
		Actions struct {
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

	device struct {
		ID               string `json:"id"`
		IsActive         bool   `json:"is_active"`
		IsPrivateSession bool   `json:"is_private_session"`
		IsRestricted     bool   `json:"is_restricted"`
		Name             string `json:"name"`
		Type             string `json:"type"`
		VolumePercent    int    `json:"volume_percent"`
		SupportsVolume   bool   `json:"supports_volume"`
	}
	Device struct {
		Device device `json:"device"`
	}
	Devices struct {
		Devices []device `json:"devices"`
	}

	artistSimple struct {
		ExternalUrls struct {
			Spotify string `json:"spotify"`
		} `json:"external_urls"`
		Href string `json:"href"`
		ID   string `json:"id"`
		Name string `json:"name"`
		Type string `json:"type"`
		URI  string `json:"uri"`
	}
	ArtistsSimple struct {
		ArtistsSimple []artistSimple `json:"artists"`
	}

	artist struct {
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
	}
	ArtistObject artist

	album struct {
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
		Type string `json:"type"`
		URI  string `json:"uri"`
		ArtistsSimple
	}
	Album struct {
		Album album `json:"album"`
	}

	track struct {
		Album
		ArtistsSimple
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
	}
	Track struct {
		Track track `json:"track"`
	}
	Tracks struct {
		Track track `json:"tracks"`
	}
	TrackObject track

	episode struct {
		AudioPreviewURL string `json:"audio_preview_url"`
		Description     string `json:"description"`
		HTMLDescription string `json:"html_description"`
		DurationMs      int    `json:"duration_ms"`
		Explicit        bool   `json:"explicit"`
		ExternalUrls    struct {
			Spotify string `json:"spotify"`
		} `json:"external_urls"`
		Href   string `json:"href"`
		ID     string `json:"id"`
		Images []struct {
			URL    string `json:"url"`
			Height int    `json:"height"`
			Width  int    `json:"width"`
		} `json:"images"`
		IsExternallyHosted   bool     `json:"is_externally_hosted"`
		IsPlayable           bool     `json:"is_playable"`
		Language             string   `json:"language"`
		Languages            []string `json:"languages"`
		Name                 string   `json:"name"`
		ReleaseDate          string   `json:"release_date"`
		ReleaseDatePrecision string   `json:"release_date_precision"`
		ResumePoint          struct {
			FullyPlayed      bool `json:"fully_played"`
			ResumePositionMs int  `json:"resume_position_ms"`
		} `json:"resume_point"`
		Type         string `json:"type"`
		URI          string `json:"uri"`
		Restrictions struct {
			Reason string `json:"reason"`
		} `json:"restrictions"`
		Show struct {
			AvailableMarkets []string `json:"available_markets"`
			Copyrights       []struct {
				Text string `json:"text"`
				Type string `json:"type"`
			} `json:"copyrights"`
			Description     string `json:"description"`
			HTMLDescription string `json:"html_description"`
			Explicit        bool   `json:"explicit"`
			ExternalUrls    struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			Href   string `json:"href"`
			ID     string `json:"id"`
			Images []struct {
				URL    string `json:"url"`
				Height int    `json:"height"`
				Width  int    `json:"width"`
			} `json:"images"`
			IsExternallyHosted bool     `json:"is_externally_hosted"`
			Languages          []string `json:"languages"`
			MediaType          string   `json:"media_type"`
			Name               string   `json:"name"`
			Publisher          string   `json:"publisher"`
			Type               string   `json:"type"`
			URI                string   `json:"uri"`
			TotalEpisodes      int      `json:"total_episodes"`
		} `json:"show"`
	}
	EpisodeObject episode

	TrackAndEpisodeObject struct {
		track
		episode
	}

	audiobook struct {
		Authors []struct {
			Name string `json:"name"`
		} `json:"authors"`
		AvailableMarkets []string `json:"available_markets"`
		Copyrights       []struct {
			Text string `json:"text"`
			Type string `json:"type"`
		} `json:"copyrights"`
		Description     string `json:"description"`
		HTMLDescription string `json:"html_description"`
		Edition         string `json:"edition"`
		Explicit        bool   `json:"explicit"`
		ExternalUrls    struct {
			Spotify string `json:"spotify"`
		} `json:"external_urls"`
		Href   string `json:"href"`
		ID     string `json:"id"`
		Images []struct {
			URL    string `json:"url"`
			Height int    `json:"height"`
			Width  int    `json:"width"`
		} `json:"images"`
		Languages []string `json:"languages"`
		MediaType string   `json:"media_type"`
		Name      string   `json:"name"`
		Narrators []struct {
			Name string `json:"name"`
		} `json:"narrators"`
		Publisher     string `json:"publisher"`
		Type          string `json:"type"`
		URI           string `json:"uri"`
		TotalChapters int    `json:"total_chapters"`
	}
	Audiobook struct {
		Audiobook audiobook `json:"audiobook"`
	}

	chapter struct {
		AudioPreviewURL  string   `json:"audio_preview_url"`
		AvailableMarkets []string `json:"available_markets"`
		ChapterNumber    int      `json:"chapter_number"`
		Description      string   `json:"description"`
		HTMLDescription  string   `json:"html_description"`
		DurationMs       int      `json:"duration_ms"`
		Explicit         bool     `json:"explicit"`
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
		IsPlayable           bool     `json:"is_playable"`
		Languages            []string `json:"languages"`
		Name                 string   `json:"name"`
		ReleaseDate          string   `json:"release_date"`
		ReleaseDatePrecision string   `json:"release_date_precision"`
		ResumePoint          struct {
			FullyPlayed      bool `json:"fully_played"`
			ResumePositionMs int  `json:"resume_position_ms"`
		} `json:"resume_point"`
		Type         string `json:"type"`
		URI          string `json:"uri"`
		Restrictions struct {
			Reason string `json:"reason"`
		} `json:"restrictions"`
		Audiobook
	}
	Chapters struct {
		Chapters []chapter `json:"chapters"`
	}
	ChapterObject chapter

	categorie struct {
		Href  string `json:"href"`
		Icons []struct {
			URL    string `json:"url"`
			Height int    `json:"height"`
			Width  int    `json:"width"`
		} `json:"icons"`
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	CategorieObject categorie
)
