package Gonos

import (
	"encoding/xml"
	"fmt"
	"strconv"
)

type trackMetaData struct {
	XMLName       xml.Name `xml:"item"`
	Res           string   `xml:"res"`
	StreamContent string   `xml:"streamContent"`
	AlbumArtUri   string   `xml:"albumArtURI"`
	Title         string   `xml:"title"`
	Class         string   `xml:"class"`
	Creator       string   `xml:"creator"`
	Album         string   `xml:"album"`
}

// Get current transport state. (TODO: Test)
func (zp *ZonePlayer) GetState() (string, error) {
	return zp.sendCommand("AVTransport", "GetTransportInfo", "", "CurrentTransportState")
}

// Same as GetState but converts to bool based on current state (TODO: Test)
func (zp *ZonePlayer) GetPlay() (bool, error) {
	state, err := zp.GetState()
	return state == "PLAYING", err
}

// Start track. (TODO: Test)
func (zp *ZonePlayer) Play() error {
	_, err := zp.sendCommand("AVTransport", "Play", "<Speed>1</Speed>", "")
	return err
}

// Same as GetState but converts to bool based on current state (TODO: Test)
func (zp *ZonePlayer) GetPause() (bool, error) {
	state, err := zp.GetState()
	return state == "PAUSED_PLAYBACK", err
}

// Pause track. (TODO: Test)
func (zp *ZonePlayer) Pause() error {
	_, err := zp.sendCommand("AVTransport", "Pause", "<Speed>1</Speed>", "")
	return err
}

// Same as GetState but converts to bool based on current state (TODO: Test)
func (zp *ZonePlayer) GetStop() (bool, error) {
	state, err := zp.GetState()
	return state == "STOPPED", err
}

// Reset track progress and pause. (TODO: Test)
func (zp *ZonePlayer) Stop() error {
	_, err := zp.sendCommand("AVTransport", "Stop", "<Speed>1</Speed>", "")
	return err
}

// Next track. (TODO: Test)
func (zp *ZonePlayer) Next() error {
	_, err := zp.sendCommand("AVTransport", "Next", "<Speed>1</Speed>", "")
	return err
}

// Previous track. (TODO: Test)
func (zp *ZonePlayer) Previous() error {
	_, err := zp.sendCommand("AVTransport", "Previous", "<Speed>1</Speed>", "")
	return err
}

// Set progress. (TODO: Test)
func (zp *ZonePlayer) Seek(hours int, minutes int, seconds int) error {
	_, err := zp.sendCommand("AVTransport", "Seek", "<Unit>REL_TIME</Unit><Target>"+fmt.Sprintf("%v:%v:%v", hours, minutes, seconds)+"</Target>", "")
	return err
}

// Join player to master. (TODO: Test)
func (zp *ZonePlayer) JoinPlayer(master_uid string) error {
	_, err := zp.sendCommand("AVTransport", "SetAVTransportURI", "<CurrentURI>x-rincon:"+master_uid+"</CurrentURI><CurrentURIMetaData></CurrentURIMetaData>", "")
	return err
}

// Unjoin player. (TODO: Test)
func (zp *ZonePlayer) UnjoinPlayer() error {
	_, err := zp.sendCommand("AVTransport", "BecomeCoordinatorOfStandaloneGroup", "<Speed>1</Speed>", "")
	return err
}

// Get player mode. (TODO: Test)
func (zp *ZonePlayer) GetPlayMode() (shuffle bool, repeat bool, repeat_one bool, err error) {
	res, err := zp.sendCommand("AVTransport", "GetTransportSettings", "", "PlayMode")
	if err != nil {
		return false, false, false, err
	}
	modeBools, ok := Playmodes[res]
	if !ok {
		return false, false, false, ErrSonos.ErrUnexpectedResponse
	}
	return modeBools[0], modeBools[1], modeBools[2], nil
}

// Set player mode. (TODO: Test)
func (zp *ZonePlayer) SetPlayMode(shuffle bool, repeat bool, repeat_one bool) error {
	mode, ok := PlaymodesReversed[[3]bool{shuffle, repeat, repeat_one}]
	if !ok {
		return ErrSonos.ErrInvalidPlayMode
	}
	_, err := zp.sendCommand("AVTransport", "SetPlayMode", "<NewPlayMode>"+mode+"</NewPlayMode>", "")
	return err
}

// Get shuffle mode. (TODO: Test)
func (zp *ZonePlayer) GetShuffle() (bool, error) {
	shuffle, _, _, err := zp.GetPlayMode()
	return shuffle, err
}

// Set shuffle mode. (TODO: Test)
func (zp *ZonePlayer) SetShuffle(state bool) error {
	_, repeat, repeat_one, err := zp.GetPlayMode()
	if err != nil {
		return err
	}
	return zp.SetPlayMode(state, repeat, repeat_one)
}

// Get repeat mode. (TODO: Test)
func (zp *ZonePlayer) GetRepeat() (bool, error) {
	_, repeat, _, err := zp.GetPlayMode()
	return repeat, err
}

// Set repeat mode. (TODO: Test)
func (zp *ZonePlayer) SetRepeat(state bool) error {
	shuffle, _, repeat_one, err := zp.GetPlayMode()
	if err != nil {
		return err
	}
	return zp.SetPlayMode(shuffle, state, repeat_one && !state)
}

// Get repeat one mode. (TODO: Test)
func (zp *ZonePlayer) GetRepeatOne() (bool, error) {
	_, _, repeat_one, err := zp.GetPlayMode()
	return repeat_one, err
}

// Set repeat one mode. (TODO: Test)
func (zp *ZonePlayer) SetRepeatOne(state bool) error {
	shuffle, repeat, _, err := zp.GetPlayMode()
	if err != nil {
		return err
	}
	return zp.SetPlayMode(shuffle, repeat && !state, state)
}

// Set line in. (TODO: Test)
func (zp *ZonePlayer) SetLineIn(speaker_uid string) error {
	_, err := zp.sendCommand("AVTransport", "SetAVTransportURI", "<CurrentURI>x-rincon-stream:"+speaker_uid+"</CurrentURI><CurrentURIMetaData></CurrentURIMetaData>", "")
	return err
}

// Get information about the current track. (TODO: Test)
func (zp *ZonePlayer) GetTrackInfo() (*TrackInfo, error) {
	info, err := zp.GetTrackInfoRaw()
	if err != nil {
		return &TrackInfo{}, err
	}
	metadata := trackMetaData{}
	err = unmarshalMetaData(info.TrackMetaData, &metadata)
	if err != nil {
		return &TrackInfo{}, err
	}
	return &TrackInfo{
		QuePosition: info.Track,
		Duration:    info.TrackDuration,
		URI:         info.TrackURI,
		Progress:    info.RelTime,
		AlbumArtURI: "http://" + zp.IpAddress.String() + ":1400" + metadata.AlbumArtUri,
		Title:       metadata.Title,
		Class:       metadata.Class,
		Creator:     metadata.Creator,
		Album:       metadata.Album,
	}, nil
}

// Same as GetTrackInfo but won't parse the information as much. (TODO: Test)
func (zp *ZonePlayer) GetTrackInfoRaw() (TrackInfoRaw, error) {
	res, err := zp.sendCommand("AVTransport", "GetPositionInfo", "", "s:Body")
	if err != nil {
		return TrackInfoRaw{}, err
	}
	trackInfoRaw := TrackInfoRaw{}
	if err := xml.Unmarshal([]byte(res), &trackInfoRaw); err != nil {
		return TrackInfoRaw{}, err
	}
	return trackInfoRaw, nil
}

// Play from que. (TODO: Test)
func (zp *ZonePlayer) PlayFromQue(track int) error {
	_, err := zp.sendCommand("AVTransport", "Seek", "<Unit>TRACK_NR</Unit><Target>"+strconv.Itoa(max(1, track))+"</Target>", "")
	return err
}

// Remove from que. (TODO: Test)
func (zp *ZonePlayer) RemoveFromQue(track int) error {
	_, err := zp.sendCommand("AVTransport", "RemoveTrackFromQueue", "<ObjectID>Q:0/"+strconv.Itoa(max(1, track))+"</ObjectID><UpdateID>0</UpdateID>", "")
	return err
}

// Add URI to que. (TODO: Test)
func (zp *ZonePlayer) AddToQue(uri string, index string, next bool) error {
	_, err := zp.sendCommand("AVTransport", "AddURIToQueue", "<EnqueuedURI>"+uri+"</EnqueuedURI><EnqueuedURIMetaData></EnqueuedURIMetaData><DesiredFirstTrackNumberEnqueued>"+index+"</DesiredFirstTrackNumberEnqueued><EnqueueAsNext>"+boolTo10(next)+"</EnqueueAsNext>", "")
	return err
}

// Clear que. (TODO: Test)
func (zp *ZonePlayer) ClearQue() error {
	_, err := zp.sendCommand("AVTransport", "RemoveAllTracksFromQueue", "", "")
	return err
}

// Set URI. (TODO: Test)
func (zp *ZonePlayer) PlayUri(uri string, meta string) error {
	_, err := zp.sendCommand("AVTransport", "SetAVTransportURI", "<CurrentURI>"+uri+"</CurrentURI><CurrentURIMetaData>"+meta+"</CurrentURIMetaData>", "")
	return err
}
