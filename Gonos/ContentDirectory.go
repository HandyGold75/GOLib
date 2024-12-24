package Gonos

import (
	"encoding/xml"
	"strconv"
)

type (
	browseResponse   struct {
		XMLName               xml.Name `xml:"BrowseResponse"`
		// Encoded DIDL-Lite XML.
		Result 	string 	
		ResultParsed  struct {
			// TODO: Fill in
		}
		NumberReturned 	int 
		TotalMatches 	int 
		UpdateID 	int 
	}
	createObjectResponse   struct {
		XMLName               xml.Name `xml:"CreateObjectResponse"`
		ObjectID 	string 	 
		Result 		string
	}
	findPrefix struct {
		XMLName               xml.Name `xml:"FindPrefix "`
		StartingIndex 	int 	 
		UpdateID 	int
	}
)

// `objectID` may be one of `Gonos.ContentTypes.*` or a custom id
// 
// TODO: Test
func(zp *ZonePlayer ) Browse(objectID string,  browseFlag string,  filter string,  startingIndex int,  requestedCount int,  sortCriteria string,) (browseResponse, error ) { res, err := zp.SendContentDirectory("Browse", "<ObjectID>"+objectID+"</ObjectID><BrowseFlag>"+browseFlag+"</BrowseFlag><Filter>"+filter+"</Filter><StartingIndex>"+strconv.Itoa(startingIndex)+"</StartingIndex><RequestedCount>"+strconv.Itoa(requestedCount)+"</RequestedCount><SortCriteria>"+sortCriteria+"</SortCriteria>", ""); 
	if err != nil { return browseResponse {}, err }; data := browseResponse {}
	err = xml.Unmarshal([]byte(res), &data); if err != nil { return browseResponse {}, err }
	err = unmarshalMetaData(data.Result, &data.ResultParsed); return data, err}
 
// TODO: Test
func(zp *ZonePlayer ) CreateObject(containerID string, elements string) (createObjectResponse, error ) { res, err := zp.SendContentDirectory("CreateObject", "<ContainerID>"+containerID+"</ContainerID><Elements>"+elements+"</Elements>", ""); 
	if err != nil { return createObjectResponse {}, err }; data := createObjectResponse {}
	err = xml.Unmarshal([]byte(res), &data); return data, err}
 
// `objectID` may be one of `Gonos.ContentTypes.*` or a custom id
// 
// TODO: Test
func(zp *ZonePlayer ) DestroyObject(objectID string) ( error ) { _, err := zp.SendContentDirectory("DestroyObject", "<ObjectID>"+objectID+"</ObjectID>", ""); return err }
 
// `objectID` may be one of `Gonos.ContentTypes.*` or a custom id
// 
// TODO: Test
func(zp *ZonePlayer ) FindPrefix(objectID string, prefix string) ( findPrefix,error ) { res, err := zp.SendContentDirectory("FindPrefix", "<ObjectID>"+objectID+"</ObjectID><Prefix>"+prefix+"</Prefix>", ""); 
	if err != nil { return findPrefix {}, err }; data := findPrefix {}
	err = xml.Unmarshal([]byte(res), &data); return data, err}
 
// TODO: Test + Implement
func(zp *ZonePlayer ) GetAlbumArtistDisplayOption() ( string, error ) { return zp.SendContentDirectory("GetAlbumArtistDisplayOption", "", "AlbumArtistDisplayOption")}
 
// TODO: Test + Implement
func(zp *ZonePlayer ) GetAllPrefixLocations() ( error ) { _, err := zp.SendContentDirectory("GetAllPrefixLocations", "", ""); return err }
 
// TODO: Test + Implement
func(zp *ZonePlayer ) GetBrowseable() ( error ) { _, err := zp.SendContentDirectory("GetBrowseable", "", ""); return err }
 
// TODO: Test + Implement
func(zp *ZonePlayer ) GetLastIndexChange() ( error ) { _, err := zp.SendContentDirectory("GetLastIndexChange", "", ""); return err }
 
// TODO: Test + Implement
func(zp *ZonePlayer ) GetSearchCapabilities() ( error ) { _, err := zp.SendContentDirectory("GetSearchCapabilities", "", ""); return err }
 
// TODO: Test + Implement
func(zp *ZonePlayer ) GetShareIndexInProgress() ( error ) { _, err := zp.SendContentDirectory("GetShareIndexInProgress", "", ""); return err }
 
// TODO: Test + Implement
func(zp *ZonePlayer ) GetSortCapabilities() ( error ) { _, err := zp.SendContentDirectory("GetSortCapabilities", "", ""); return err }
 
// TODO: Test + Implement
func(zp *ZonePlayer ) GetSystemUpdateID() ( error ) { _, err := zp.SendContentDirectory("GetSystemUpdateID", "", ""); return err }
 
// TODO: Test + Implement
func(zp *ZonePlayer ) RefreshShareIndex() ( error ) { _, err := zp.SendContentDirectory("RefreshShareIndex", "", ""); return err }
 
// TODO: Test + Implement
func(zp *ZonePlayer ) RequestResort() ( error ) { _, err := zp.SendContentDirectory("RequestResort", "", ""); return err }
 
// TODO: Test + Implement
func(zp *ZonePlayer ) SetBrowseable() ( error ) { _, err := zp.SendContentDirectory("SetBrowseable", "", ""); return err }
 
// TODO: Test + Implement
func(zp *ZonePlayer ) UpdateObject() ( error ) { _, err := zp.SendContentDirectory("UpdateObject", "", ""); return err }

