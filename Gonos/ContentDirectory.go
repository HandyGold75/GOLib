package Gonos

import (
	"encoding/xml"
	"strconv"
	"time"
)

type (
	browseResponse struct {
		XMLName xml.Name `xml:"BrowseResponse"`
		// Encoded DIDL-Lite XML.
		//
		// Should be unmarshaled into type of `browseResponseMetaData*`
		Result         string
		NumberReturned int
		TotalMatches   int
		UpdateID       int
	}
	createObjectResponse struct {
		XMLName  xml.Name `xml:"CreateObjectResponse"`
		ObjectID string
		Result   string
	}
	findPrefix struct {
		XMLName       xml.Name `xml:"FindPrefix"`
		StartingIndex int
		UpdateID      int
	}
	getAllPrefixLocationsResponse struct {
		XMLName           xml.Name `xml:"getAllPrefixLocationsResponse"`
		TotalPrefixes     int
		PrefixAndIndexCSV string
		UpdateID          int
	}
)

// Prefer methods `zp.BrowseMusicLibrary`, `zp.BrowsePlaylist`, `zp.BrowseQue`.
//
// `objectID` may be one of `Gonos.ContentTypes.*` or a custom id
func (zp *ZonePlayer) Browse(objectID string, browseFlag string, filter string, startingIndex int, requestedCount int, sortCriteria string) (browseResponse, error) {
	res, err := zp.SendContentDirectory("Browse", "<ObjectID>"+objectID+"</ObjectID><BrowseFlag>"+browseFlag+"</BrowseFlag><Filter>"+filter+"</Filter><StartingIndex>"+strconv.Itoa(startingIndex)+"</StartingIndex><RequestedCount>"+strconv.Itoa(requestedCount)+"</RequestedCount><SortCriteria>"+sortCriteria+"</SortCriteria>", "s:Body")
	if err != nil {
		return browseResponse{}, err
	}
	data := browseResponse{}
	err = xml.Unmarshal([]byte(res), &data)
	return data, err
}

// TODO: Test
func (zp *ZonePlayer) CreateObject(containerID string, elements string) (createObjectResponse, error) {
	res, err := zp.SendContentDirectory("CreateObject", "<ContainerID>"+containerID+"</ContainerID><Elements>"+elements+"</Elements>", "s:Body")
	if err != nil {
		return createObjectResponse{}, err
	}
	data := createObjectResponse{}
	err = xml.Unmarshal([]byte(res), &data)
	return data, err
}

// `objectID` may be one of `Gonos.ContentTypes.*` or a custom id
//
// TODO: Test
func (zp *ZonePlayer) DestroyObject(objectID string) error {
	_, err := zp.SendContentDirectory("DestroyObject", "<ObjectID>"+objectID+"</ObjectID>", "")
	return err
}

// `objectID` may be one of `Gonos.ContentTypes.*` or a custom id
//
// TODO: Test
func (zp *ZonePlayer) FindPrefix(objectID string, prefix string) (findPrefix, error) {
	res, err := zp.SendContentDirectory("FindPrefix", "<ObjectID>"+objectID+"</ObjectID><Prefix>"+prefix+"</Prefix>", "s:Body")
	if err != nil {
		return findPrefix{}, err
	}
	data := findPrefix{}
	err = xml.Unmarshal([]byte(res), &data)
	return data, err
}

// TODO: Test
//
// Returns one of `WMP`, `ITUNES` or `NONE`
func (zp *ZonePlayer) GetAlbumArtistDisplayOption() (string, error) {
	return zp.SendContentDirectory("GetAlbumArtistDisplayOption", "", "AlbumArtistDisplayOption")
}

// `objectID` may be one of `Gonos.ContentTypes.*` or a custom id
//
// TODO: Test
func (zp *ZonePlayer) GetAllPrefixLocations(objectID string) (getAllPrefixLocationsResponse, error) {
	res, err := zp.SendContentDirectory("GetAllPrefixLocations", "<ObjectID>"+objectID+"</ObjectID>", "s:Body")
	if err != nil {
		return getAllPrefixLocationsResponse{}, err
	}
	data := getAllPrefixLocationsResponse{}
	err = xml.Unmarshal([]byte(res), &data)
	return data, err
}

// Check if content directory is browsable.
func (zp *ZonePlayer) GetBrowseable() (bool, error) {
	res, err := zp.SendContentDirectory("GetBrowseable", "", "IsBrowseable")
	return res == "1", err
}

// Returns when the index last changed (What is the index? It's not the que nor when any Sonos playlists/ librarys get updated? It's also not effected by `RefreshShareIndex`?).
func (zp *ZonePlayer) GetLastIndexChange() (time.Time, error) {
	res, err := zp.SendContentDirectory("GetLastIndexChange", "", "LastIndexChange")
	if err != nil {
		return time.Time{}, err
	}
	return time.Parse("S"+time.DateTime, res)
}

// TODO: Test
func (zp *ZonePlayer) GetSearchCapabilities() (string, error) {
	return zp.SendContentDirectory("GetSearchCapabilities", "", "SearchCaps")
}

// Check if share index is in progress (caused by `zp.RefreshShareIndex`).
func (zp *ZonePlayer) GetShareIndexInProgress() (bool, error) {
	res, err := zp.SendContentDirectory("GetShareIndexInProgress", "", "IsIndexing")
	return res == "1", err
}

// TODO: Test
func (zp *ZonePlayer) GetSortCapabilities() (string, error) {
	return zp.SendContentDirectory("GetSortCapabilities", "", "SortCaps")
}

// Get system update ID (What is this ID? I don't think it should be used for the UpdateID params).
func (zp *ZonePlayer) GetSystemUpdateID() (int, error) {
	res, err := zp.SendContentDirectory("GetSystemUpdateID", "", "Id")
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(res)
}

// `albumArtistDisplayOption` should be one of `zp.AlbumArtistDisplayOptionTypes`.
func (zp *ZonePlayer) RefreshShareIndex(albumArtistDisplayOption string) error {
	_, err := zp.SendContentDirectory("RefreshShareIndex", "<AlbumArtistDisplayOption>"+albumArtistDisplayOption+"</AlbumArtistDisplayOption>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) RequestResort(sortOrder string) error {
	_, err := zp.SendContentDirectory("RequestResort", "<SortOrder>sortOrder</SortOrder>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) SetBrowseable(state bool) error {
	_, err := zp.SendContentDirectory("SetBrowseable", "<Browseable>"+boolTo10(state)+"</Browseable>", "")
	return err
}

// `objectID` may be one of `Gonos.ContentTypes.*` or a custom id
//
// TODO: Test
func (zp *ZonePlayer) UpdateObject(objectID string, currentTagValue string, newTagValue string) error {
	_, err := zp.SendContentDirectory("UpdateObject", "<ObjectID>"+objectID+"</ObjectID><CurrentTagValue>"+currentTagValue+"</CurrentTagValue><NewTagValue>"+newTagValue+"</NewTagValue>", "")
	return err
}
