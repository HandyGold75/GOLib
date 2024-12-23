package Gonos

import (
	"encoding/xml"
	"strconv"
)

type queMetaData struct {
	XMLName     xml.Name `xml:"item"`
	Res         string   `xml:"res"`
	AlbumArtUri string   `xml:"albumArtURI"`
	Title       string   `xml:"title"`
	Class       string   `xml:"class"`
	Creator     string   `xml:"creator"`
	Album       string   `xml:"album"`
}

type favoritesMetaData struct {
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

type contentDirectorResponse struct {
	XMLName        xml.Name `xml:"BrowseResponse"`
	Result         string
	NumberReturned string
	TotalMatches   string
	UpdateID       string
}

func (zp *ZonePlayer) getContentDirectory(typ string, start int, count int) (contentDirectorResponse, error) {
	id, ok := ContentTypes[typ]
	if !ok {
		return contentDirectorResponse{}, ErrSonos.ErrInvalidContentType
	}
	res, err := zp.SendContentDirectory("Browse", "<ObjectID>"+id+"</ObjectID><BrowseFlag>BrowseDirectChildren</BrowseFlag><Filter>dc:title,res,dc:creator,upnp:artist,upnp:album,upnp:albumArtURI</Filter><StartingIndex>"+strconv.Itoa(start)+"</StartingIndex><RequestedCount>"+strconv.Itoa(count)+"</RequestedCount><SortCriteria></SortCriteria>", "s:Body")
	if err != nil {
		return contentDirectorResponse{}, err
	}
	data := contentDirectorResponse{}
	err = xml.Unmarshal([]byte(res), &data)
	return data, err
}

// Get information about the que. (TODO: Test)
func (zp *ZonePlayer) GetQue() (*Que, error) {
	info, err := zp.getContentDirectory("queue", 0, 0)
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
	info, err := zp.getContentDirectory("sonos share", 0, 0)
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
	info, err := zp.getContentDirectory("sonos favorites", 0, 0)
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
	info, err := zp.getContentDirectory("radio stations", 0, 0)
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
	info, err := zp.getContentDirectory("radio shows", 0, 0)
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
