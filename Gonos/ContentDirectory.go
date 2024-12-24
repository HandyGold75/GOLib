package Gonos

import (
	"encoding/xml"
	"strconv"
)

type (
	Que struct {
		Count      string
		TotalCount string
		Tracks     []QueTrack
	}
	QueTrack struct {
		AlbumArtURI string
		Title       string
		Class       string
		Creator     string
		Album       string
	}
	queMetaData struct {
		XMLName     xml.Name `xml:"item"`
		Res         string   `xml:"res"`
		AlbumArtUri string   `xml:"albumArtURI"`
		Title       string   `xml:"title"`
		Class       string   `xml:"class"`
		Creator     string   `xml:"creator"`
		Album       string   `xml:"album"`
	}

	Favorites struct {
		Count      string
		TotalCount string
		Favorites  []FavoritesItem
	}
	FavoritesItem struct {
		AlbumArtURI string
		Title       string
		Description string
		Class       string
		Type        string
	}
	favoritesMetaData struct {
		XMLName     xml.Name `xml:"item"`
		Title       string   `xml:"title"`
		Class       string   `xml:"class"`
		Ordinal     string   `xml:"ordinal"`
		Res         string   `xml:"res"`
		AlbumArtUri string   `xml:"albumArtURI"`
		Type        string   `xml:"type"`
		Description string   `xml:"description"`
		ResMD       string   `xml:"resMD"`
	}

	contentDirectorResponse struct {
		XMLName        xml.Name `xml:"BrowseResponse"`
		Result         string
		NumberReturned string
		TotalMatches   string
		UpdateID       string
	}
)

func (zp *ZonePlayer) getContentDirectory(objectID string, start int, count int) (contentDirectorResponse, error) {
	res, err := zp.SendContentDirectory("Browse", "<ObjectID>"+objectID+"</ObjectID><BrowseFlag>BrowseDirectChildren</BrowseFlag><Filter>dc:title,res,dc:creator,upnp:artist,upnp:album,upnp:albumArtURI</Filter><StartingIndex>"+strconv.Itoa(start)+"</StartingIndex><RequestedCount>"+strconv.Itoa(count)+"</RequestedCount><SortCriteria></SortCriteria>", "s:Body")
	if err != nil {
		return contentDirectorResponse{}, err
	}
	data := contentDirectorResponse{}
	err = xml.Unmarshal([]byte(res), &data)
	return data, err
}

// Get information about the que. (TODO: Test)
func (zp *ZonePlayer) GetQue() (*Que, error) {
	info, err := zp.getContentDirectory(ContentTypes.QueueMain, 0, 0)
	if err != nil {
		return &Que{}, err
	}
	metadata := []queMetaData{}
	err = unmarshalMetaData(info.Result, &metadata)
	if err != nil {
		return &Que{}, err
	}
	tracks := []QueTrack{}
	for _, track := range metadata {
		tracks = append(tracks, QueTrack{
			AlbumArtURI: "http://" + zp.IpAddress.String() + ":1400" + track.AlbumArtUri,
			Title:       track.Title,
			Class:       track.Class,
			Creator:     track.Creator,
			Album:       track.Album,
		})
	}
	return &Que{
		Count:      info.NumberReturned,
		TotalCount: info.TotalMatches,
		Tracks:     tracks,
	}, nil
}

// Get information about the share. (TODO: Test)
func (zp *ZonePlayer) GetShare() (*Favorites, error) {
	info, err := zp.getContentDirectory(ContentTypes.Share, 0, 0)
	if err != nil {
		return &Favorites{}, err
	}
	metadata := []favoritesMetaData{}
	err = unmarshalMetaData(info.Result, &metadata)
	if err != nil {
		return &Favorites{}, err
	}
	favorites := []FavoritesItem{}
	for _, favorite := range metadata {
		favorites = append(favorites, FavoritesItem{
			AlbumArtURI: "http://" + zp.IpAddress.String() + ":1400" + favorite.AlbumArtUri,
			Title:       favorite.Title,
			Description: favorite.Description,
			Class:       favorite.Class,
			Type:        favorite.Type,
		})
	}
	return &Favorites{
		Count:      info.NumberReturned,
		TotalCount: info.TotalMatches,
		Favorites:  favorites,
	}, nil
}

// Get information about the favorites. (TODO: Test)
func (zp *ZonePlayer) GetFavorites() (*Favorites, error) {
	info, err := zp.getContentDirectory(ContentTypes.SonosFavorites, 0, 0)
	if err != nil {
		return &Favorites{}, err
	}
	metadata := []favoritesMetaData{}
	err = unmarshalMetaData(info.Result, &metadata)
	if err != nil {
		return &Favorites{}, err
	}
	favorites := []FavoritesItem{}
	for _, favorite := range metadata {
		favorites = append(favorites, FavoritesItem{
			AlbumArtURI: "http://" + zp.IpAddress.String() + ":1400" + favorite.AlbumArtUri,
			Title:       favorite.Title,
			Description: favorite.Description,
			Class:       favorite.Class,
			Type:        favorite.Type,
		})
	}
	return &Favorites{
		Count:      info.NumberReturned,
		TotalCount: info.TotalMatches,
		Favorites:  favorites,
	}, nil
}

// Get information about the favorites radio stations. (TODO: Test)
func (zp *ZonePlayer) GetRadioStations() (*Favorites, error) {
	info, err := zp.getContentDirectory(ContentTypes.RadioStations, 0, 0)
	if err != nil {
		return &Favorites{}, err
	}
	metadata := []favoritesMetaData{}
	err = unmarshalMetaData(info.Result, &metadata)
	if err != nil {
		return &Favorites{}, err
	}
	favorites := []FavoritesItem{}
	for _, favorite := range metadata {
		favorites = append(favorites, FavoritesItem{
			AlbumArtURI: "http://" + zp.IpAddress.String() + ":1400" + favorite.AlbumArtUri,
			Title:       favorite.Title,
			Description: favorite.Description,
			Class:       favorite.Class,
			Type:        favorite.Type,
		})
	}
	return &Favorites{
		Count:      info.NumberReturned,
		TotalCount: info.TotalMatches,
		Favorites:  favorites,
	}, nil
}

// Get information about the radio shows. (TODO: Test)
func (zp *ZonePlayer) GetRadioShows() (*Favorites, error) {
	info, err := zp.getContentDirectory(ContentTypes.RadioShows, 0, 0)
	if err != nil {
		return &Favorites{}, err
	}
	metadata := []favoritesMetaData{}
	err = unmarshalMetaData(info.Result, &metadata)
	if err != nil {
		return &Favorites{}, err
	}
	favorites := []FavoritesItem{}
	for _, favorite := range metadata {
		favorites = append(favorites, FavoritesItem{
			AlbumArtURI: "http://" + zp.IpAddress.String() + ":1400" + favorite.AlbumArtUri,
			Title:       favorite.Title,
			Description: favorite.Description,
			Class:       favorite.Class,
			Type:        favorite.Type,
		})
	}
	return &Favorites{
		Count:      info.NumberReturned,
		TotalCount: info.TotalMatches,
		Favorites:  favorites,
	}, nil
}
