package Gonos

import "encoding/xml"

type (
	browseResponseMetaDataQue struct {
		XMLName     xml.Name `xml:"item"`
		Res         string   `xml:"res"`
		AlbumArtUri string   `xml:"albumArtURI"`
		Title       string   `xml:"title"`
		Class       string   `xml:"class"`
		Creator     string   `xml:"creator"`
		Album       string   `xml:"album"`
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

	browseResponseMetaDataQueFavorites struct {
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
	favoritesInfo struct {
		Count      int
		TotalCount int
		Favorites  []favoritesInfoItem
	}
	favoritesInfoItem struct {
		AlbumArtURI string
		Title       string
		Description string
		Class       string
		Type        string
	}
)

// TODO: Test
func (zp *ZonePlayer) GetShare() (favoritesInfo, error) {
	info, err := zp.Browse(ContentTypes.Share, "BrowseDirectChildren", "dc:title,res,dc:creator,upnp:artist,upnp:album,upnp:albumArtURI", 0, 0, "")
	if err != nil {
		return favoritesInfo{}, err
	}
	metadata := []browseResponseMetaDataQueFavorites{}
	err = unmarshalMetaData(info.Result, &metadata)
	if err != nil {
		return favoritesInfo{}, err
	}
	favorites := []favoritesInfoItem{}
	for _, favorite := range metadata {
		favorites = append(favorites, favoritesInfoItem{
			AlbumArtURI: "http://" + zp.IpAddress.String() + ":1400" + favorite.AlbumArtUri,
			Title:       favorite.Title,
			Description: favorite.Description,
			Class:       favorite.Class,
			Type:        favorite.Type,
		})
	}
	return favoritesInfo{Count: info.NumberReturned, TotalCount: info.TotalMatches, Favorites: favorites}, nil
}

// TODO: Test
func (zp *ZonePlayer) GetPlaylists() (favoritesInfo, error) {
	info, err := zp.Browse(ContentTypes.SonosPlaylists, "BrowseDirectChildren", "dc:title,res,dc:creator,upnp:artist,upnp:album,upnp:albumArtURI", 0, 0, "")
	if err != nil {
		return favoritesInfo{}, err
	}
	metadata := []browseResponseMetaDataQueFavorites{}
	err = unmarshalMetaData(info.Result, &metadata)
	if err != nil {
		return favoritesInfo{}, err
	}
	favorites := []favoritesInfoItem{}
	for _, favorite := range metadata {
		favorites = append(favorites, favoritesInfoItem{
			AlbumArtURI: "http://" + zp.IpAddress.String() + ":1400" + favorite.AlbumArtUri,
			Title:       favorite.Title,
			Description: favorite.Description,
			Class:       favorite.Class,
			Type:        favorite.Type,
		})
	}
	return favoritesInfo{Count: info.NumberReturned, TotalCount: info.TotalMatches, Favorites: favorites}, nil
}

// TODO: Test
func (zp *ZonePlayer) GetFavorites() (favoritesInfo, error) {
	info, err := zp.Browse(ContentTypes.SonosFavorites, "BrowseDirectChildren", "dc:title,res,dc:creator,upnp:artist,upnp:album,upnp:albumArtURI", 0, 0, "")
	if err != nil {
		return favoritesInfo{}, err
	}
	metadata := []browseResponseMetaDataQueFavorites{}
	err = unmarshalMetaData(info.Result, &metadata)
	if err != nil {
		return favoritesInfo{}, err
	}
	favorites := []favoritesInfoItem{}
	for _, favorite := range metadata {
		favorites = append(favorites, favoritesInfoItem{
			AlbumArtURI: "http://" + zp.IpAddress.String() + ":1400" + favorite.AlbumArtUri,
			Title:       favorite.Title,
			Description: favorite.Description,
			Class:       favorite.Class,
			Type:        favorite.Type,
		})
	}
	return favoritesInfo{Count: info.NumberReturned, TotalCount: info.TotalMatches, Favorites: favorites}, nil
}

// TODO: Test
func (zp *ZonePlayer) GetRadioStations() (favoritesInfo, error) {
	info, err := zp.Browse(ContentTypes.RadioStations, "BrowseDirectChildren", "dc:title,res,dc:creator,upnp:artist,upnp:album,upnp:albumArtURI", 0, 0, "")
	if err != nil {
		return favoritesInfo{}, err
	}
	metadata := []browseResponseMetaDataQueFavorites{}
	err = unmarshalMetaData(info.Result, &metadata)
	if err != nil {
		return favoritesInfo{}, err
	}
	favorites := []favoritesInfoItem{}
	for _, favorite := range metadata {
		favorites = append(favorites, favoritesInfoItem{
			AlbumArtURI: "http://" + zp.IpAddress.String() + ":1400" + favorite.AlbumArtUri,
			Title:       favorite.Title,
			Description: favorite.Description,
			Class:       favorite.Class,
			Type:        favorite.Type,
		})
	}
	return favoritesInfo{Count: info.NumberReturned, TotalCount: info.TotalMatches, Favorites: favorites}, nil
}

// TODO: Test
func (zp *ZonePlayer) GetRadioShows() (favoritesInfo, error) {
	info, err := zp.Browse(ContentTypes.RadioShows, "BrowseDirectChildren", "dc:title,res,dc:creator,upnp:artist,upnp:album,upnp:albumArtURI", 0, 0, "")
	if err != nil {
		return favoritesInfo{}, err
	}
	metadata := []browseResponseMetaDataQueFavorites{}
	err = unmarshalMetaData(info.Result, &metadata)
	if err != nil {
		return favoritesInfo{}, err
	}
	favorites := []favoritesInfoItem{}
	for _, favorite := range metadata {
		favorites = append(favorites, favoritesInfoItem{
			AlbumArtURI: "http://" + zp.IpAddress.String() + ":1400" + favorite.AlbumArtUri,
			Title:       favorite.Title,
			Description: favorite.Description,
			Class:       favorite.Class,
			Type:        favorite.Type,
		})
	}
	return favoritesInfo{Count: info.NumberReturned, TotalCount: info.TotalMatches, Favorites: favorites}, nil
}

// TODO: Test
func (zp *ZonePlayer) GetQue() (queInfo, error) {
	info, err := zp.Browse(ContentTypes.QueueMain, "BrowseDirectChildren", "dc:title,res,dc:creator,upnp:artist,upnp:album,upnp:albumArtURI", 0, 0, "")
	if err != nil {
		return queInfo{}, err
	}
	metadata := []browseResponseMetaDataQue{}
	err = unmarshalMetaData(info.Result, &metadata)
	if err != nil {
		return queInfo{}, err
	}
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
