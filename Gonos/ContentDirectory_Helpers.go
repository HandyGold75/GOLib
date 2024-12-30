package Gonos

import (
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

// Prefer methods `zp.LibraryArtist`, `zp.LibraryAlbumArtist`, `zp.LibraryAlbum`, `zp.LibraryGenre`, `zp.LibraryComposer`, `zp.LibraryTracks`, `zp.LibraryPlaylists`.
//
// `objectID` may be one of `Gonos.ContentTypes.*` or a custom id
func (zp *ZonePlayer) BrowseMusicLibrary(objectID string) (libraryInfo, error) {
	info, err := zp.Browse(objectID, "BrowseDirectChildren", "dc:title,res,dc:creator,upnp:artist,upnp:album,upnp:albumArtURI", 0, 0, "")
	if err != nil {
		return libraryInfo{}, err
	}
	metadata := []browseResponseMetaDataLibrary{}
	fmt.Println(strings.ReplaceAll(info.Result, "id=", "\r\nid="))
	err = unmarshalMetaData(info.Result, &metadata)
	if err == io.EOF {
		return libraryInfo{}, nil
	} else if err != nil {
		return libraryInfo{}, err
	}
	librarys := []libraryInfoItem{}
	for _, library := range metadata {
		librarys = append(librarys, libraryInfoItem{
			AlbumArtURI: "http://" + zp.IpAddress.String() + ":1400" + library.AlbumArtUri,
			Title:       library.Title,
			Description: library.Description,
			Class:       library.Class,
			Type:        library.Type,
		})
	}
	return libraryInfo{Count: info.NumberReturned, TotalCount: info.TotalMatches, Librarys: librarys}, nil
}

// TODO: Test
func (zp *ZonePlayer) GetLibraryArtist() (libraryInfo, error) {
	return zp.BrowseMusicLibrary(ContentTypes.Artist)
}

// TODO: Test
func (zp *ZonePlayer) GetLibraryAlbumArtist() (libraryInfo, error) {
	return zp.BrowseMusicLibrary(ContentTypes.AlbumArtist)
}

// TODO: Test
func (zp *ZonePlayer) GetLibraryAlbum() (libraryInfo, error) {
	return zp.BrowseMusicLibrary(ContentTypes.Album)
}

// TODO: Test
func (zp *ZonePlayer) GetLibraryGenre() (libraryInfo, error) {
	return zp.BrowseMusicLibrary(ContentTypes.Genre)
}

// TODO: Test
func (zp *ZonePlayer) GetLibraryComposer() (libraryInfo, error) {
	return zp.BrowseMusicLibrary(ContentTypes.Composer)
}

// TODO: Test
func (zp *ZonePlayer) GetLibraryTracks() (libraryInfo, error) {
	return zp.BrowseMusicLibrary(ContentTypes.Tracks)
}

// TODO: Test
func (zp *ZonePlayer) GetLibraryPlaylists() (libraryInfo, error) {
	return zp.BrowseMusicLibrary(ContentTypes.Playlists)
}

// Prefer methods `zp.GetShare`, `zp.GetSonosPlaylists`, `zp.GetSonosFavorites`, `zp.GetRadioStations` or `zp.GetRadioShows`.
//
// `objectID` may be one of `Gonos.ContentTypes.*` or a custom id
func (zp *ZonePlayer) BrowsePlaylist(objectID string) (playlistInfo, error) {
	info, err := zp.Browse(objectID, "BrowseDirectChildren", "dc:title,res,dc:creator,upnp:artist,upnp:album,upnp:albumArtURI", 0, 0, "")
	if err != nil {
		return playlistInfo{}, err
	}
	metadata := []browseResponseMetaDataQuePlaylist{}
	err = unmarshalMetaData(info.Result, &metadata)
	if err == io.EOF {
		return playlistInfo{}, nil
	} else if err != nil {
		return playlistInfo{}, err
	}
	playlists := []playlistInfoItem{}
	for _, playlist := range metadata {
		playlists = append(playlists, playlistInfoItem{
			AlbumArtURI: "http://" + zp.IpAddress.String() + ":1400" + playlist.AlbumArtUri,
			Title:       playlist.Title,
			Description: playlist.Description,
			Class:       playlist.Class,
			Type:        playlist.Type,
		})
	}
	return playlistInfo{Count: info.NumberReturned, TotalCount: info.TotalMatches, Playlists: playlists}, nil
}

// TODO: Test
func (zp *ZonePlayer) GetShare() (playlistInfo, error) {
	return zp.BrowsePlaylist(ContentTypes.Share)
}

// TODO: Test
func (zp *ZonePlayer) GetSonosPlaylists() (playlistInfo, error) {
	return zp.BrowsePlaylist(ContentTypes.SonosPlaylists)
}

// Get Sonos playlists, in case no sonos playlists are present a empty playlist will be returned
func (zp *ZonePlayer) GetSonosFavorites() (playlistInfo, error) {
	return zp.BrowsePlaylist(ContentTypes.SonosFavorites)
}

// TODO: Test
func (zp *ZonePlayer) GetRadioStations() (playlistInfo, error) {
	return zp.BrowsePlaylist(ContentTypes.RadioStations)
}

// TODO: Test
func (zp *ZonePlayer) GetRadioShows() (playlistInfo, error) {
	return zp.BrowsePlaylist(ContentTypes.RadioShows)
}

// Prefer methods `zp.GetQue` or `zp.GetQueSecond`.
func (zp *ZonePlayer) BrowseQue(objectID string) (queInfo, error) {
	info, err := zp.Browse(objectID, "BrowseDirectChildren", "dc:title,res,dc:creator,upnp:artist,upnp:album,upnp:albumArtURI", 0, 0, "")
	if err != nil {
		return queInfo{}, err
	}
	metadata := []browseResponseMetaDataQue{}
	err = unmarshalMetaData(info.Result, &metadata)
	if err == io.EOF {
		return queInfo{}, nil
	} else if err != nil {
		return queInfo{}, err
	}
	fmt.Println(len(metadata))
	tracks := []queInfoItem{}
	for _, track := range metadata {
		tracks = append(tracks, queInfoItem{
			AlbumArtURI: "http://" + zp.IpAddress.String() + ":1400" + track.AlbumArtUri,
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
func (zp *ZonePlayer) GetQue() (queInfo, error) {
	return zp.BrowseQue(ContentTypes.QueueMain)
}

// Get secondairy que, in case no que is active a empty que will be returned.
//
// Will return incorrect data if a third party application is controling playback.
func (zp *ZonePlayer) GetQueSecond() (queInfo, error) {
	return zp.BrowseQue(ContentTypes.QueueSecond)
}
