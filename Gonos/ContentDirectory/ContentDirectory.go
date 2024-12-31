package ContentDirectory

import (
	"Gonos/lib"
	"encoding/xml"
	"fmt"
	"strconv"
	"time"
)

type (
	ContentDirectory struct {
		Send func(action, body, targetTag string) (string, error)
	}

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

func New(send func(action, body, targetTag string) (string, error)) ContentDirectory {
	return ContentDirectory{Send: send}
}

// Prefer methods `zp.BrowseMusicLibrary`, `zp.BrowsePlaylist`, `zp.BrowseQue`.
//
// `objectID` may be one of `Gonos.ContentTypes.*` or a custom id
func (zp *ContentDirectory) Browse(objectID string, browseFlag string, filter string, startingIndex int, requestedCount int, sortCriteria string) (browseResponse, error) {
	res, err := zp.Send("Browse", "<ObjectID>"+objectID+"</ObjectID><BrowseFlag>"+browseFlag+"</BrowseFlag><Filter>"+filter+"</Filter><StartingIndex>"+strconv.Itoa(startingIndex)+"</StartingIndex><RequestedCount>"+strconv.Itoa(requestedCount)+"</RequestedCount><SortCriteria>"+sortCriteria+"</SortCriteria>", "s:Body")
	if err != nil {
		return browseResponse{}, err
	}
	data := browseResponse{}
	err = xml.Unmarshal([]byte(res), &data)
	return data, err
}

func (zp *ContentDirectory) CreateObject(containerID string, elements string) (createObjectResponse, error) {
	res, err := zp.Send("CreateObject", "<ContainerID>"+containerID+"</ContainerID><Elements>"+elements+"</Elements>", "s:Body")
	if err != nil {
		return createObjectResponse{}, err
	}
	data := createObjectResponse{}
	err = xml.Unmarshal([]byte(res), &data)
	return data, err
}

// `objectID` may be one of `Gonos.ContentTypes.*` or a custom id
func (zp *ContentDirectory) DestroyObject(objectID string) error {
	_, err := zp.Send("DestroyObject", "<ObjectID>"+objectID+"</ObjectID>", "")
	return err
}

// `objectID` may be one of `Gonos.ContentTypes.*` or a custom id
func (zp *ContentDirectory) FindPrefix(objectID string, prefix string) (findPrefix, error) {
	res, err := zp.Send("FindPrefix", "<ObjectID>"+objectID+"</ObjectID><Prefix>"+prefix+"</Prefix>", "s:Body")
	if err != nil {
		return findPrefix{}, err
	}
	data := findPrefix{}
	err = xml.Unmarshal([]byte(res), &data)
	return data, err
}

// Returns one of `WMP`, `ITUNES` or `NONE`
func (zp *ContentDirectory) GetAlbumArtistDisplayOption() (string, error) {
	return zp.Send("GetAlbumArtistDisplayOption", "", "AlbumArtistDisplayOption")
}

// `objectID` may be one of `Gonos.ContentTypes.*` or a custom id
func (zp *ContentDirectory) GetAllPrefixLocations(objectID string) (getAllPrefixLocationsResponse, error) {
	res, err := zp.Send("GetAllPrefixLocations", "<ObjectID>"+objectID+"</ObjectID>", "s:Body")
	if err != nil {
		return getAllPrefixLocationsResponse{}, err
	}
	data := getAllPrefixLocationsResponse{}
	err = xml.Unmarshal([]byte(res), &data)
	return data, err
}

func (zp *ContentDirectory) GetBrowseable() (bool, error) {
	res, err := zp.Send("GetBrowseable", "", "IsBrowseable")
	return res == "1", err
}

func (zp *ContentDirectory) GetLastIndexChange() (time.Time, error) {
	res, err := zp.Send("GetLastIndexChange", "", "LastIndexChange")
	if err != nil {
		return time.Time{}, err
	}
	return time.Parse("S"+time.DateTime, res)
}

func (zp *ContentDirectory) GetSearchCapabilities() (string, error) {
	return zp.Send("GetSearchCapabilities", "", "SearchCaps")
}

func (zp *ContentDirectory) GetShareIndexInProgress() (bool, error) {
	res, err := zp.Send("GetShareIndexInProgress", "", "IsIndexing")
	return res == "1", err
}

func (zp *ContentDirectory) GetSortCapabilities() (string, error) {
	return zp.Send("GetSortCapabilities", "", "SortCaps")
}

func (zp *ContentDirectory) GetSystemUpdateID() (int, error) {
	res, err := zp.Send("GetSystemUpdateID", "", "Id")
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(res)
}

// `albumArtistDisplayOption` should be one of `zp.AlbumArtistDisplayOptionModes.*`.
func (zp *ContentDirectory) RefreshShareIndex(albumArtistDisplayOption string) error {
	_, err := zp.Send("RefreshShareIndex", "<AlbumArtistDisplayOption>"+albumArtistDisplayOption+"</AlbumArtistDisplayOption>", "")
	return err
}

func (zp *ContentDirectory) RequestResort(sortOrder string) error {
	_, err := zp.Send("RequestResort", "<SortOrder>sortOrder</SortOrder>", "")
	return err
}

func (zp *ContentDirectory) SetBrowseable(state bool) error {
	_, err := zp.Send("SetBrowseable", "<Browseable>"+lib.BoolTo10(state)+"</Browseable>", "")
	return err
}

// `objectID` may be one of `Gonos.ContentTypes.*` or a custom id
func (zp *ContentDirectory) UpdateObject(objectID string, currentTagValue string, newTagValue string) error {
	out, err := zp.Send("UpdateObject", "<ObjectID>"+objectID+"</ObjectID><CurrentTagValue>"+currentTagValue+"</CurrentTagValue><NewTagValue>"+newTagValue+"</NewTagValue>", "")
	fmt.Println(out)
	return err
}
