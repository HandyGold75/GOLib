package Helper

import (
	"Gonos/lib"
	"fmt"
	"io"
	"strings"
)

type (
	// To implement, couldn't get a example.
	browseResponseMetaDataLibrary struct {
		Title       string `xml:"title"`
		Class       string `xml:"class"`
		Ordinal     string `xml:"ordinal"`
		Res         string `xml:"res"`
		AlbumArtUri string `xml:"albumArtURI"`
		Type        string `xml:"type"`
		Description string `xml:"description"`
		ResMD       string `xml:"resMD"`
	}
	// To implement, couldn't get a example.
	libraryInfo struct {
		Count      int
		TotalCount int
		Librarys   []libraryInfoItem
	}
	// To implement, couldn't get a example.
	libraryInfoItem struct {
		AlbumArtURI string
		Title       string
		Description string
		Ordinal     string
		Class       string
		Type        string
	}

	browseResponseMetaDataQuePlaylist struct {
		Title       string `xml:"title"`
		Class       string `xml:"class"`
		Ordinal     string `xml:"ordinal"`
		Res         string `xml:"res"`
		AlbumArtUri string `xml:"albumArtURI"`
		Type        string `xml:"type"`
		Description string `xml:"description"`
		ResMD       string `xml:"resMD"`
	}
	playlistInfo struct {
		Count      int
		TotalCount int
		Playlists  []playlistInfoItem
	}
	playlistInfoItem struct {
		AlbumArtURI string
		Title       string
		Description string
		Ordinal     string
		Class       string
		Type        string
	}

	browseResponseMetaDataQue struct {
		Res         string `xml:"res"`
		AlbumArtUri string `xml:"albumArtURI"`
		Title       string `xml:"title"`
		Class       string `xml:"class"`
		Creator     string `xml:"creator"`
		Album       string `xml:"album"`
	}
	queInfo struct {
		Count      int
		TotalCount int
		Tracks     []queInfoItem
	}
	queInfoItem struct {
		AlbumArtURI string
		Title       string
		Class       string
		Creator     string
		Album       string
	}
)

// Prefer methods `h.LibraryArtist`, `h.LibraryAlbumArtist`, `h.LibraryAlbum`, `h.LibraryGenre`, `h.LibraryComposer`, `h.LibraryTracks`, `h.LibraryPlaylists`.
//
// `objectID` may be one of `Gonos.lib.ContentTypes.*` or a custom id
func (h *Helper) BrowseMusicLibrary(objectID string) (libraryInfo, error) {
	info, err := h.contentDirectory.Browse(objectID, "BrowseDirectChildren", "dc:title,res,dc:creator,upnp:artist,upnp:album,upnp:albumArtURI", 0, 0, "")
	if err != nil {
		return libraryInfo{}, err
	}
	metadata := []browseResponseMetaDataLibrary{}
	fmt.Println(strings.ReplaceAll(info.Result, "id=", "\r\nid="))
	err = lib.UnmarshalMetaData(info.Result, &metadata)
	if err == io.EOF {
		return libraryInfo{}, nil
	} else if err != nil {
		return libraryInfo{}, err
	}
	librarys := []libraryInfoItem{}
	for _, library := range metadata {
		librarys = append(librarys, libraryInfoItem{
			AlbumArtURI: library.AlbumArtUri,
			Title:       library.Title,
			Description: library.Description,
			Class:       library.Class,
			Type:        library.Type,
		})
	}
	return libraryInfo{Count: info.NumberReturned, TotalCount: info.TotalMatches, Librarys: librarys}, nil
}

func (h *Helper) GetLibraryArtist() (libraryInfo, error) {
	return h.BrowseMusicLibrary(lib.ContentTypes.Artist)
}

func (h *Helper) GetLibraryAlbumArtist() (libraryInfo, error) {
	return h.BrowseMusicLibrary(lib.ContentTypes.AlbumArtist)
}

func (h *Helper) GetLibraryAlbum() (libraryInfo, error) {
	return h.BrowseMusicLibrary(lib.ContentTypes.Album)
}

func (h *Helper) GetLibraryGenre() (libraryInfo, error) {
	return h.BrowseMusicLibrary(lib.ContentTypes.Genre)
}

func (h *Helper) GetLibraryComposer() (libraryInfo, error) {
	return h.BrowseMusicLibrary(lib.ContentTypes.Composer)
}

func (h *Helper) GetLibraryTracks() (libraryInfo, error) {
	return h.BrowseMusicLibrary(lib.ContentTypes.Tracks)
}

func (h *Helper) GetLibraryPlaylists() (libraryInfo, error) {
	return h.BrowseMusicLibrary(lib.ContentTypes.Playlists)
}

// Prefer methods `h.GetShare`, `h.GetSonosPlaylists`, `h.GetSonosFavorites`, `h.GetRadioStations` or `h.GetRadioShows`.
//
// `objectID` may be one of `Gonos.lib.ContentTypes.*` or a custom id
func (h *Helper) BrowsePlaylist(objectID string) (playlistInfo, error) {
	info, err := h.contentDirectory.Browse(objectID, "BrowseDirectChildren", "dc:title,res,dc:creator,upnp:artist,upnp:album,upnp:albumArtURI", 0, 0, "")
	if err != nil {
		return playlistInfo{}, err
	}
	metadata := []browseResponseMetaDataQuePlaylist{}
	err = lib.UnmarshalMetaData(info.Result, &metadata)
	if err == io.EOF {
		return playlistInfo{}, nil
	} else if err != nil {
		return playlistInfo{}, err
	}
	playlists := []playlistInfoItem{}
	for _, playlist := range metadata {
		playlists = append(playlists, playlistInfoItem{
			AlbumArtURI: playlist.AlbumArtUri,
			Title:       playlist.Title,
			Description: playlist.Description,
			Class:       playlist.Class,
			Type:        playlist.Type,
		})
	}
	return playlistInfo{Count: info.NumberReturned, TotalCount: info.TotalMatches, Playlists: playlists}, nil
}

func (h *Helper) GetShare() (playlistInfo, error) {
	return h.BrowsePlaylist(lib.ContentTypes.Share)
}

func (h *Helper) GetSonosPlaylists() (playlistInfo, error) {
	return h.BrowsePlaylist(lib.ContentTypes.SonosPlaylists)
}

// Get Sonos playlists, in case no sonos playlists are present a empty playlist will be returned
func (h *Helper) GetSonosFavorites() (playlistInfo, error) {
	return h.BrowsePlaylist(lib.ContentTypes.SonosFavorites)
}

func (h *Helper) GetRadioStations() (playlistInfo, error) {
	return h.BrowsePlaylist(lib.ContentTypes.RadioStations)
}

func (h *Helper) GetRadioShows() (playlistInfo, error) {
	return h.BrowsePlaylist(lib.ContentTypes.RadioShows)
}

// Prefer methods `h.GetQue` or `h.GetQueSecond`.
func (h *Helper) BrowseQue(objectID string) (queInfo, error) {
	info, err := h.contentDirectory.Browse(objectID, "BrowseDirectChildren", "dc:title,res,dc:creator,upnp:artist,upnp:album,upnp:albumArtURI", 0, 0, "")
	if err != nil {
		return queInfo{}, err
	}
	metadata := []browseResponseMetaDataQue{}
	err = lib.UnmarshalMetaData(info.Result, &metadata)
	if err == io.EOF {
		return queInfo{}, nil
	} else if err != nil {
		return queInfo{}, err
	}
	fmt.Println(len(metadata))
	tracks := []queInfoItem{}
	for _, track := range metadata {
		tracks = append(tracks, queInfoItem{
			AlbumArtURI: track.AlbumArtUri,
			Title:       track.Title,
			Class:       track.Class,
			Creator:     track.Creator,
			Album:       track.Album,
		})
	}
	return queInfo{Count: info.NumberReturned, TotalCount: info.TotalMatches, Tracks: tracks}, nil
}

// Get que, in case no que is active a empty que will be returned.
//
// Will return incorrect data if a third party application is controling playback.
func (h *Helper) GetQue() (queInfo, error) {
	return h.BrowseQue(lib.ContentTypes.QueueMain)
}

// Get secondairy que, in case no que is active a empty que will be returned.
//
// Will return incorrect data if a third party application is controling playback.
func (h *Helper) GetQueSecond() (queInfo, error) {
	return h.BrowseQue(lib.ContentTypes.QueueSecond)
}
