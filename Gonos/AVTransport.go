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
	return zp.SendAVTransport("GetTransportInfo", "<InstanceID>0</InstanceID>", "CurrentTransportState")
}

// Same as GetState but converts to bool based on current state (TODO: Test)
func (zp *ZonePlayer) GetPlay() (bool, error) {
	state, err := zp.GetState()
	return state == "PLAYING", err
}

// Start track. (TODO: Test)
func (zp *ZonePlayer) Play() error {
	_, err := zp.SendAVTransport("Play", "<InstanceID>0</InstanceID><Speed>1</Speed>", "")
	return err
}

// Same as GetState but converts to bool based on current state (TODO: Test)
func (zp *ZonePlayer) GetPause() (bool, error) {
	state, err := zp.GetState()
	return state == "PAUSED_PLAYBACK", err
}

// Pause track. (TODO: Test)
func (zp *ZonePlayer) Pause() error {
	_, err := zp.SendAVTransport("Pause", "<InstanceID>0</InstanceID><Speed>1</Speed>", "")
	return err
}

// Same as GetState but converts to bool based on current state (TODO: Test)
func (zp *ZonePlayer) GetStop() (bool, error) {
	state, err := zp.GetState()
	return state == "STOPPED", err
}

// Reset track progress and pause. (TODO: Test)
func (zp *ZonePlayer) Stop() error {
	_, err := zp.SendAVTransport("Stop", "<InstanceID>0</InstanceID><Speed>1</Speed>", "")
	return err
}

// Next track. (TODO: Test)
func (zp *ZonePlayer) Next() error {
	_, err := zp.SendAVTransport("Next", "<InstanceID>0</InstanceID><Speed>1</Speed>", "")
	return err
}

// Previous track. (TODO: Test)
func (zp *ZonePlayer) Previous() error {
	_, err := zp.SendAVTransport("Previous", "<InstanceID>0</InstanceID><Speed>1</Speed>", "")
	return err
}

// Set progress. (TODO: Test)
func (zp *ZonePlayer) Seek(hours int, minutes int, seconds int) error {
	_, err := zp.SendAVTransport("Seek", "<InstanceID>0</InstanceID><Unit>REL_TIME</Unit><Target>"+fmt.Sprintf("%v:%v:%v", hours, minutes, seconds)+"</Target>", "")
	return err
}

// Join player to master. (TODO: Test)
func (zp *ZonePlayer) JoinPlayer(master_uid string) error {
	_, err := zp.SendAVTransport("SetAVTransportURI", "<InstanceID>0</InstanceID><CurrentURI>x-rincon:"+master_uid+"</CurrentURI><CurrentURIMetaData></CurrentURIMetaData>", "")
	return err
}

// Unjoin player. (TODO: Test)
func (zp *ZonePlayer) UnjoinPlayer() error {
	_, err := zp.SendAVTransport("BecomeCoordinatorOfStandaloneGroup", "<InstanceID>0</InstanceID><Speed>1</Speed>", "")
	return err
}

// Get player mode. (TODO: Test)
func (zp *ZonePlayer) GetPlayMode() (shuffle bool, repeat bool, repeat_one bool, err error) {
	res, err := zp.SendAVTransport("GetTransportSettings", "<InstanceID>0</InstanceID>", "PlayMode")
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
	_, err := zp.SendAVTransport("SetPlayMode", "<InstanceID>0</InstanceID><NewPlayMode>"+mode+"</NewPlayMode>", "")
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
	_, err := zp.SendAVTransport("SetAVTransportURI", "<InstanceID>0</InstanceID><CurrentURI>x-rincon-stream:"+speaker_uid+"</CurrentURI><CurrentURIMetaData></CurrentURIMetaData>", "")
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
	res, err := zp.SendAVTransport("GetPositionInfo", "<InstanceID>0</InstanceID>", "s:Body")
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
	_, err := zp.SendAVTransport("Seek", "<InstanceID>0</InstanceID><Unit>TRACK_NR</Unit><Target>"+strconv.Itoa(max(1, track))+"</Target>", "")
	return err
}

// Remove from que. (TODO: Test)
func (zp *ZonePlayer) RemoveFromQue(track int) error {
	_, err := zp.SendAVTransport("RemoveTrackFromQueue", "<InstanceID>0</InstanceID><ObjectID>Q:0/"+strconv.Itoa(max(1, track))+"</ObjectID><UpdateID>0</UpdateID>", "")
	return err
}

// Add URI to que. (TODO: Test)
func (zp *ZonePlayer) AddToQue(uri string, index string, next bool) error {
	_, err := zp.SendAVTransport("AddURIToQueue", "<InstanceID>0</InstanceID><EnqueuedURI>"+uri+"</EnqueuedURI><EnqueuedURIMetaData></EnqueuedURIMetaData><DesiredFirstTrackNumberEnqueued>"+index+"</DesiredFirstTrackNumberEnqueued><EnqueueAsNext>"+boolTo10(next)+"</EnqueueAsNext>", "")
	return err
}

// Clear que. (TODO: Test)
func (zp *ZonePlayer) ClearQue() error {
	_, err := zp.SendAVTransport("RemoveAllTracksFromQueue", "<InstanceID>0</InstanceID>", "")
	return err
}

// Set URI. (TODO: Test)
func (zp *ZonePlayer) PlayUri(uri string, meta string) error {
	_, err := zp.SendAVTransport("SetAVTransportURI", "<InstanceID>0</InstanceID><CurrentURI>"+uri+"</CurrentURI><CurrentURIMetaData>"+meta+"</CurrentURIMetaData>", "")
	return err
}

// func (zp *ZonePlayer) AddMultipleURIsToQueue() ( error ) { _, err := zp.SendAVTransport("AddMultipleURIsToQueue", "", ""); return err }
// func (zp *ZonePlayer) AddURIToQueue() ( error ) { _, err := zp.SendAVTransport("AddURIToQueue", "", ""); return err }
// func (zp *ZonePlayer) AddURIToSavedQueue() ( error ) { _, err := zp.SendAVTransport("AddURIToSavedQueue", "", ""); return err }
// func (zp *ZonePlayer) BackupQueue() ( error ) { _, err := zp.SendAVTransport("BackupQueue", "", ""); return err }
// func (zp *ZonePlayer) BecomeCoordinatorOfStandaloneGroup() ( error ) { _, err := zp.SendAVTransport("BecomeCoordinatorOfStandaloneGroup", "", ""); return err }
// func (zp *ZonePlayer) BecomeGroupCoordinator() ( error ) { _, err := zp.SendAVTransport("BecomeGroupCoordinator", "", ""); return err }
// func (zp *ZonePlayer) BecomeGroupCoordinatorAndSource() ( error ) { _, err := zp.SendAVTransport("BecomeGroupCoordinatorAndSource", "", ""); return err }
// func (zp *ZonePlayer) ChangeCoordinator() ( error ) { _, err := zp.SendAVTransport("ChangeCoordinator", "", ""); return err }
// func (zp *ZonePlayer) ChangeTransportSettings() ( error ) { _, err := zp.SendAVTransport("ChangeTransportSettings", "", ""); return err }
// func (zp *ZonePlayer) ConfigureSleepTimer() ( error ) { _, err := zp.SendAVTransport("ConfigureSleepTimer", "", ""); return err }
// func (zp *ZonePlayer) CreateSavedQueue() ( error ) { _, err := zp.SendAVTransport("CreateSavedQueue", "", ""); return err }
// func (zp *ZonePlayer) DelegateGroupCoordinationTo() ( error ) { _, err := zp.SendAVTransport("DelegateGroupCoordinationTo", "", ""); return err }
// func (zp *ZonePlayer) EndDirectControlSession() ( error ) { _, err := zp.SendAVTransport("EndDirectControlSession", "", ""); return err }
// func (zp *ZonePlayer) GetCrossfadeMode() ( error ) { _, err := zp.SendAVTransport("GetCrossfadeMode", "", ""); return err }
// func (zp *ZonePlayer) GetCurrentTransportActions() ( error ) { _, err := zp.SendAVTransport("GetCurrentTransportActions", "", ""); return err }
// func (zp *ZonePlayer) GetDeviceCapabilities() ( error ) { _, err := zp.SendAVTransport("GetDeviceCapabilities", "", ""); return err }
// func (zp *ZonePlayer) GetMediaInfo() ( error ) { _, err := zp.SendAVTransport("GetMediaInfo", "", ""); return err }
// func (zp *ZonePlayer) GetPositionInfo() ( error ) { _, err := zp.SendAVTransport("GetPositionInfo", "", ""); return err }
// func (zp *ZonePlayer) GetRemainingSleepTimerDuration() ( error ) { _, err := zp.SendAVTransport("GetRemainingSleepTimerDuration", "", ""); return err }
// func (zp *ZonePlayer) GetRunningAlarmProperties() ( error ) { _, err := zp.SendAVTransport("GetRunningAlarmProperties", "", ""); return err }
// func (zp *ZonePlayer) GetTransportInfo() ( error ) { _, err := zp.SendAVTransport("GetTransportInfo", "", ""); return err }
// func (zp *ZonePlayer) GetTransportSettings() ( error ) { _, err := zp.SendAVTransport("GetTransportSettings", "", ""); return err }
// func (zp *ZonePlayer) Next() ( error ) { _, err := zp.SendAVTransport("Next", "", ""); return err }
// func (zp *ZonePlayer) NotifyDeletedURI() ( error ) { _, err := zp.SendAVTransport("NotifyDeletedURI", "", ""); return err }
// func (zp *ZonePlayer) Pause() ( error ) { _, err := zp.SendAVTransport("Pause", "", ""); return err }
// func (zp *ZonePlayer) Play() ( error ) { _, err := zp.SendAVTransport("Play", "", ""); return err }
// func (zp *ZonePlayer) Previous() ( error ) { _, err := zp.SendAVTransport("Previous", "", ""); return err }
// func (zp *ZonePlayer) RemoveAllTracksFromQueue() ( error ) { _, err := zp.SendAVTransport("RemoveAllTracksFromQueue", "", ""); return err }
// func (zp *ZonePlayer) RemoveTrackFromQueue() ( error ) { _, err := zp.SendAVTransport("RemoveTrackFromQueue", "", ""); return err }
// func (zp *ZonePlayer) RemoveTrackRangeFromQueue() ( error ) { _, err := zp.SendAVTransport("RemoveTrackRangeFromQueue", "", ""); return err }
// func (zp *ZonePlayer) ReorderTracksInQueue() ( error ) { _, err := zp.SendAVTransport("ReorderTracksInQueue", "", ""); return err }
// func (zp *ZonePlayer) ReorderTracksInSavedQueue() ( error ) { _, err := zp.SendAVTransport("ReorderTracksInSavedQueue", "", ""); return err }
// func (zp *ZonePlayer) RunAlarm() ( error ) { _, err := zp.SendAVTransport("RunAlarm", "", ""); return err }
// func (zp *ZonePlayer) SaveQueue() ( error ) { _, err := zp.SendAVTransport("SaveQueue", "", ""); return err }
// func (zp *ZonePlayer) Seek() ( error ) { _, err := zp.SendAVTransport("Seek", "", ""); return err }
// func (zp *ZonePlayer) SetAVTransportURI() ( error ) { _, err := zp.SendAVTransport("SetAVTransportURI", "", ""); return err }
// func (zp *ZonePlayer) SetCrossfadeMode() ( error ) { _, err := zp.SendAVTransport("SetCrossfadeMode", "", ""); return err }
// func (zp *ZonePlayer) SetNextAVTransportURI() ( error ) { _, err := zp.SendAVTransport("SetNextAVTransportURI", "", ""); return err }
// func (zp *ZonePlayer) SetPlayMode() ( error ) { _, err := zp.SendAVTransport("SetPlayMode", "", ""); return err }
// func (zp *ZonePlayer) SnoozeAlarm() ( error ) { _, err := zp.SendAVTransport("SnoozeAlarm", "", ""); return err }
// func (zp *ZonePlayer) StartAutoplay() ( error ) { _, err := zp.SendAVTransport("StartAutoplay", "", ""); return err }
// func (zp *ZonePlayer) Stop() ( error ) { _, err := zp.SendAVTransport("Stop", "", ""); return err }
