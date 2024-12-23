package Gonos

// https://sonos.svrooij.io/services/device-properties

import (
	"encoding/xml"
	"strconv"
	"time"
)

// TrackInfo struct {
// 	QuePosition string
// 	Duration    string
// 	URI         string
// 	Progress    string
// 	AlbumArtURI string
// 	Title       string
// 	Class       string
// 	Creator     string
// 	Album       string
// }
// TrackInfoRaw struct {
// 	XMLName       xml.Name `xml:"GetPositionInfoResponse"`
// 	Track         string
// 	TrackDuration string
// 	TrackMetaData string
// 	TrackURI      string
// 	RelTime       string
// 	AbsTime       string
// 	RelCount      string
// 	AbsCount      string
// }

// type trackMetaData struct {
// 	XMLName       xml.Name `xml:"item"`
// 	Res           string   `xml:"res"`
// 	StreamContent string   `xml:"streamContent"`
// 	AlbumArtUri   string   `xml:"albumArtURI"`
// 	Title         string   `xml:"title"`
// 	Class         string   `xml:"class"`
// 	Creator       string   `xml:"creator"`
// 	Album         string   `xml:"album"`
// }

// func (zp *ZonePlayer) GetTrackInfo() (*TrackInfo, error) {
// 	info, err := zp.GetTrackInfoRaw()
// 	if err != nil {
// 		return &TrackInfo{}, err
// 	}
// 	metadata := trackMetaData{}
// 	err = unmarshalMetaData(info.TrackMetaData, &metadata)
// 	if err != nil {
// 		return &TrackInfo{}, err
// 	}
// 	return &TrackInfo{
// 		QuePosition: info.Track,
// 		Duration:    info.TrackDuration,
// 		URI:         info.TrackURI,
// 		Progress:    info.RelTime,
// 		AlbumArtURI: "http://" + zp.IpAddress.String() + ":1400" + metadata.AlbumArtUri,
// 		Title:       metadata.Title,
// 		Class:       metadata.Class,
// 		Creator:     metadata.Creator,
// 		Album:       metadata.Album,
// 	}, nil
// }

// func (zp *ZonePlayer) GetTrackInfoRaw() (TrackInfoRaw, error) {
// 	res, err := zp.SendAVTransport("GetPositionInfo", "<InstanceID>0</InstanceID>", "s:Body")
// 	if err != nil {
// 		return TrackInfoRaw{}, err
// 	}
// 	trackInfoRaw := TrackInfoRaw{}
// 	if err := xml.Unmarshal([]byte(res), &trackInfoRaw); err != nil {
// 		return TrackInfoRaw{}, err
// 	}
// 	return trackInfoRaw, nil
// }

// func (zp *ZonePlayer) RemoveFromQue(track int) error {
// 	_, err := zp.SendAVTransport("RemoveTrackFromQueue", "<InstanceID>0</InstanceID><ObjectID>Q:0/"+strconv.Itoa(max(1, track))+"</ObjectID><UpdateID>0</UpdateID>", "")
// 	return err
// }

type (
	addMultipleURIsToQueueResponse struct {
		XMLName                  xml.Name `xml:"AddMultipleURIsToQueueResponse"`
		FirstTrackNumberEnqueued int
		NumTracksAdded           int
		NewQueueLength           int
		NewUpdateID              int
	}
	addURIToQueueResponse struct {
		XMLName                  xml.Name `xml:"AddURIToQueueResponse"`
		FirstTrackNumberEnqueued int
		NumTracksAdded           int
		NewQueueLength           int
	}
	addURIToSavedQueueResponse struct {
		XMLName        xml.Name `xml:"AddURIToSavedQueueResponse"`
		NumTracksAdded int
		NewQueueLength int
		NewUpdateID    int
	}
	becomeCoordinatorOfStandaloneGroupResponse struct {
		XMLName                     xml.Name `xml:"BecomeCoordinatorOfStandaloneGroupResponse"`
		DelegatedGroupCoordinatorID string
		NewGroupID                  string
	}
	createSavedQueueResponse struct {
		XMLName          xml.Name `xml:"CreateSavedQueueResponse"`
		NewQueueLength   int
		AssignedObjectID string
		NewUpdateID      int
	}
	getDeviceCapabilitiesResponse struct {
		XMLName         xml.Name `xml:"GetDeviceCapabilitiesResponse"`
		PlayMedia       string
		RecMedia        string
		RecQualityModes string
	}
	getMediaInfoResponse struct {
		XMLName       xml.Name `xml:"GetMediaInfoResponse"`
		NrTracks      int
		MediaDuration string
		CurrentURI    string
		// Embedded XML
		CurrentURIMetaData string
		NextURI            string
		// Embedded XML
		NextURIMetaData string
		// Possible values: NONE / NETWORK
		PlayMedium string
		// Possible values: NONE
		RecordMedium string
		WriteStatus  string
	}
	getPositionInfoResponse struct {
		XMLName       xml.Name `xml:"GetPositionInfoResponse"`
		Track         int
		TrackDuration string
		// Embedded XML
		TrackMetaData string
		TrackURI      string
		RelTime       string
		AbsTime       string
		RelCount      int
		AbsCount      int
	}
	getRemainingSleepTimerDurationResponse struct {
		XMLName xml.Name `xml:"GetRemainingSleepTimerDurationResponse"`
		// Format hh:mm:ss or empty string if not set
		RemainingSleepTimerDuration string
		CurrentSleepTimerGeneration int
	}
	getRunningAlarmPropertiesResponse struct {
		XMLName         xml.Name `xml:"GetRunningAlarmPropertiesResponse"`
		AlarmID         int
		GroupID         string
		LoggedStartTime string
	}
	reorderTracksInSavedQueueResponse struct {
		XMLName           xml.Name `xml:"ReorderTracksInSavedQueueResponse"`
		QueueLengthChange int
		NewQueueLength    int
		NewUpdateID       int
	}
)

var Playmodes = map[string][3]bool{
	// "MODE": [2]bool{shuffle, repeat, repeat_one}
	"NORMAL":             {false, false, false},
	"REPEAT_ALL":         {false, true, false},
	"REPEAT_ONE":         {false, false, true},
	"SHUFFLE_NOREPEAT":   {true, false, false},
	"SHUFFLE":            {true, true, false},
	"SHUFFLE_REPEAT_ONE": {true, false, true},
}

var PlaymodesReversed = func() map[[3]bool]string {
	PMS := map[[3]bool]string{}
	for k, v := range Playmodes {
		PMS[v] = k
	}
	return PMS
}()

// TODO: test
func (zp *ZonePlayer) AddMultipleURIsToQueue(updateID int, numberOfURIs int, enqueuedURIs string, enqueuedURIsMetaData string, containerURI string, containerMetaData string, desiredFirstTrackNumberEnqueued int, enqueueAsNext bool) (addMultipleURIsToQueueResponse, error) {
	res, err := zp.SendAVTransport("AddMultipleURIsToQueue", "<UpdateID>"+strconv.Itoa(updateID)+"</UpdateID><NumberOfURIs>"+strconv.Itoa(numberOfURIs)+"</NumberOfURIs><EnqueuedURIs>"+enqueuedURIs+"</EnqueuedURIs><EnqueuedURIsMetaData>"+enqueuedURIsMetaData+"</EnqueuedURIsMetaData><ContainerURI>"+containerURI+"</ContainerURI><ContainerMetaData>"+containerMetaData+"</ContainerMetaData><DesiredFirstTrackNumberEnqueued>"+strconv.Itoa(desiredFirstTrackNumberEnqueued)+"</DesiredFirstTrackNumberEnqueued><EnqueueAsNext>"+boolTo10(enqueueAsNext)+"</EnqueueAsNext>", "s:Body")
	if err != nil {
		return addMultipleURIsToQueueResponse{}, err
	}
	data := addMultipleURIsToQueueResponse{}
	err = xml.Unmarshal([]byte(res), &data)
	return data, err
}

// TODO: Test
func (zp *ZonePlayer) AddURIToQueue(enqueuedURI string, enqueuedURIMetaData string, desiredFirstTrackNumberEnqueued int, enqueueAsNext bool) (addURIToQueueResponse, error) {
	res, err := zp.SendAVTransport("AddURIToQueue", "<EnqueuedURI>"+enqueuedURI+"</EnqueuedURI><EnqueuedURIMetaData>"+enqueuedURIMetaData+"</EnqueuedURIMetaData><DesiredFirstTrackNumberEnqueued>"+strconv.Itoa(desiredFirstTrackNumberEnqueued)+"</DesiredFirstTrackNumberEnqueued><EnqueueAsNext>"+boolTo10(enqueueAsNext)+"</EnqueueAsNext>", "s:Body")
	if err != nil {
		return addURIToQueueResponse{}, err
	}
	data := addURIToQueueResponse{}
	err = xml.Unmarshal([]byte(res), &data)
	return data, err
}

// TODO: Test
func (zp *ZonePlayer) AddURIToSavedQueue(objectID string, updateID int, enqueuedURI string, enqueuedURIMetaData string, addAtIndex int) (addURIToSavedQueueResponse, error) {
	res, err := zp.SendAVTransport("AddURIToSavedQueue", "<ObjectID>"+objectID+"</ObjectID><UpdateID>"+strconv.Itoa(updateID)+"</UpdateID><EnqueuedURI>"+enqueuedURI+"</EnqueuedURI><EnqueuedURIMetaData>"+enqueuedURIMetaData+"</EnqueuedURIMetaData><AddAtIndex>"+strconv.Itoa(addAtIndex)+"</AddAtIndex>", "s:Body")
	if err != nil {
		return addURIToSavedQueueResponse{}, err
	}
	data := addURIToSavedQueueResponse{}
	err = xml.Unmarshal([]byte(res), &data)
	return data, err
}

// TODO: Test
func (zp *ZonePlayer) BackupQueue() error {
	_, err := zp.SendAVTransport("BackupQueue", "", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) BecomeCoordinatorOfStandaloneGroup() (becomeCoordinatorOfStandaloneGroupResponse, error) {
	res, err := zp.SendAVTransport("BecomeCoordinatorOfStandaloneGroup", "", "s:Body")
	if err != nil {
		return becomeCoordinatorOfStandaloneGroupResponse{}, err
	}
	data := becomeCoordinatorOfStandaloneGroupResponse{}
	err = xml.Unmarshal([]byte(res), &data)
	return data, err
}

// TODO: Test
func (zp *ZonePlayer) BecomeGroupCoordinator(currentCoordinator string, currentGroupID string, otherMembers string, transportSettings string, currentURI string, currentURIMetaData string, sleepTimerState string, alarmState string, streamRestartState string, currentQueueTrackList string, currentVLIState string) error {
	_, err := zp.SendAVTransport("BecomeGroupCoordinator", "<CurrentCoordinator>"+currentCoordinator+"</CurrentCoordinator><CurrentGroupID>"+currentGroupID+"</CurrentGroupID><OtherMembers>"+otherMembers+"</OtherMembers><TransportSettings>"+transportSettings+"</TransportSettings><CurrentURI>"+currentURI+"</CurrentURI><CurrentURIMetaData>"+currentURIMetaData+"</CurrentURIMetaData><SleepTimerState>"+sleepTimerState+"</SleepTimerState><AlarmState>"+alarmState+"</AlarmState><StreamRestartState>"+streamRestartState+"</StreamRestartState><CurrentQueueTrackList>"+currentQueueTrackList+"</CurrentQueueTrackList><CurrentVLIState>"+currentVLIState+"</CurrentVLIState>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) BecomeGroupCoordinatorAndSource(currentCoordinator string, currentGroupID string, otherMembers string, currentURI string, currentURIMetaData string, sleepTimerState string, alarmState string, streamRestartState string, currentAVTTrackList string, currentQueueTrackList string, currentSourceState string, resumePlayback bool) error {
	_, err := zp.SendAVTransport("BecomeGroupCoordinatorAndSource", "<CurrentCoordinator>"+currentCoordinator+"</CurrentCoordinator><CurrentGroupID>"+currentGroupID+"</CurrentGroupID><OtherMembers>"+otherMembers+"</OtherMembers><CurrentURI>"+currentURI+"</CurrentURI><CurrentURIMetaData>"+currentURIMetaData+"</CurrentURIMetaData><SleepTimerState>"+sleepTimerState+"</SleepTimerState><AlarmState>"+alarmState+"</AlarmState><StreamRestartState>"+streamRestartState+"</StreamRestartState><CurrentAVTTrackList>"+currentAVTTrackList+"</CurrentAVTTrackList><CurrentQueueTrackList>"+currentQueueTrackList+"</CurrentQueueTrackList><CurrentSourceState>"+currentSourceState+"</CurrentSourceState><ResumePlayback>"+boolTo10(resumePlayback)+"</ResumePlayback>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) ChangeCoordinator(currentCoordinator string, newCoordinator string, newTransportSettings string, currentAVTransportURI string) error {
	_, err := zp.SendAVTransport("ChangeCoordinator", "<CurrentCoordinator>"+currentCoordinator+"</CurrentCoordinator><NewCoordinator>"+newCoordinator+"</NewCoordinator><NewTransportSettings>"+newTransportSettings+"</NewTransportSettings><CurrentAVTransportURI>"+currentAVTransportURI+"</CurrentAVTransportURI>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) ChangeTransportSettings(newTransportSettings string, currentAVTransportURI string) error {
	_, err := zp.SendAVTransport("ChangeTransportSettings", "<NewTransportSettings>"+newTransportSettings+"</NewTransportSettings><CurrentAVTransportURI>"+currentAVTransportURI+"</CurrentAVTransportURI>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) ConfigureSleepTimer(seconds int) error {
	_, err := zp.SendAVTransport("ConfigureSleepTimer", "<NewSleepTimerDuration>"+time.Time.Add(time.Time{}, time.Second*time.Duration(seconds)).Format("15:04:05")+"</NewSleepTimerDuration>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) CreateSavedQueue(title string, enqueuedURI string, enqueuedURIMetaData string) (createSavedQueueResponse, error) {
	res, err := zp.SendAVTransport("CreateSavedQueue", "<Title>title</Title><EnqueuedURI>enqueuedURI</EnqueuedURI><EnqueuedURIMetaData>enqueuedURIMetaData</EnqueuedURIMetaData>", "s:Body")
	if err != nil {
		return createSavedQueueResponse{}, err
	}
	data := createSavedQueueResponse{}
	err = xml.Unmarshal([]byte(res), &data)
	return data, err
}

// TODO: Test
func (zp *ZonePlayer) DelegateGroupCoordinationTo(newCoordinator string, rejoinGroup bool) error {
	_, err := zp.SendAVTransport("DelegateGroupCoordinationTo", "<NewCoordinator>"+newCoordinator+"</NewCoordinator><RejoinGroup>"+boolTo10(rejoinGroup)+"</RejoinGroup>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) EndDirectControlSession() error {
	_, err := zp.SendAVTransport("EndDirectControlSession", "", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) GetCrossfadeMode() (bool, error) {
	res, err := zp.SendAVTransport("GetCrossfadeMode", "", "CrossfadeMode")
	return res == "1", err
}

// TODO: Test
func (zp *ZonePlayer) GetCurrentTransportActions() (string, error) {
	return zp.SendAVTransport("GetCurrentTransportActions", "", "Actions")
}

// TODO: Test
func (zp *ZonePlayer) GetDeviceCapabilities() (getDeviceCapabilitiesResponse, error) {
	res, err := zp.SendAVTransport("GetDeviceCapabilities", "", "s:Body")
	if err != nil {
		return getDeviceCapabilitiesResponse{}, err
	}
	data := getDeviceCapabilitiesResponse{}
	err = xml.Unmarshal([]byte(res), &data)
	return data, err
}

// TODO: Test
func (zp *ZonePlayer) GetMediaInfo() (getMediaInfoResponse, error) {
	res, err := zp.SendAVTransport("GetMediaInfo", "", "s:Body")
	if err != nil {
		return getMediaInfoResponse{}, err
	}
	data := getMediaInfoResponse{}
	err = xml.Unmarshal([]byte(res), &data)
	return data, err
}

// TODO: Test
func (zp *ZonePlayer) GetPositionInfo() (getPositionInfoResponse, error) {
	res, err := zp.SendAVTransport("GetPositionInfo", "", "s:Body")
	if err != nil {
		return getPositionInfoResponse{}, err
	}
	data := getPositionInfoResponse{}
	err = xml.Unmarshal([]byte(res), &data)
	return data, err
}

// TODO: Test
func (zp *ZonePlayer) GetRemainingSleepTimerDuration() (getRemainingSleepTimerDurationResponse, error) {
	res, err := zp.SendAVTransport("GetRemainingSleepTimerDuration", "", "s:Body")
	if err != nil {
		return getRemainingSleepTimerDurationResponse{}, err
	}
	data := getRemainingSleepTimerDurationResponse{}
	err = xml.Unmarshal([]byte(res), &data)
	return data, err
}

// TODO: Test
func (zp *ZonePlayer) GetRunningAlarmProperties() (getRunningAlarmPropertiesResponse, error) {
	res, err := zp.SendAVTransport("GetRunningAlarmProperties", "", "s:Body")
	if err != nil {
		return getRunningAlarmPropertiesResponse{}, err
	}
	data := getRunningAlarmPropertiesResponse{}
	err = xml.Unmarshal([]byte(res), &data)
	return data, err
}

// TODO: Test
func (zp *ZonePlayer) GetCurrentTransportStatus() (string, error) {
	return zp.SendAVTransport("GetTransportInfo", "", "CurrentTransportStatus")
}

// TODO: Test
func (zp *ZonePlayer) GetCurrentTransportState() (string, error) {
	return zp.SendAVTransport("GetTransportInfo", "", "CurrentTransportState")
}

// TODO: Test
func (zp *ZonePlayer) GetStop() (bool, error) {
	state, err := zp.GetCurrentTransportState()
	return state == "STOPPED", err
}

// TODO: Test
func (zp *ZonePlayer) GetPlay() (bool, error) {
	state, err := zp.GetCurrentTransportState()
	return state == "PLAYING", err
}

// TODO: Test
func (zp *ZonePlayer) GetPause() (bool, error) {
	state, err := zp.GetCurrentTransportState()
	return state == "PAUSED_PLAYBACK", err
}

// TODO: Test
func (zp *ZonePlayer) GetTransitioning() (bool, error) {
	state, err := zp.GetCurrentTransportState()
	return state == "TRANSITIONING", err
}

// TODO: Test
func (zp *ZonePlayer) GetRecQualityMode() (string, error) {
	return zp.SendAVTransport("GetTransportSettings", "", "RecQualityMode")
}

// TODO: Test
func (zp *ZonePlayer) GetPlayMode() (shuffle bool, repeat bool, repeat_one bool, err error) {
	res, err := zp.SendAVTransport("GetTransportSettings", "", "PlayMode")
	if err != nil {
		return false, false, false, err
	}
	modeBools, ok := Playmodes[res]
	if !ok {
		return false, false, false, ErrSonos.ErrUnexpectedResponse
	}
	return modeBools[0], modeBools[1], modeBools[2], nil
}

// TODO: Test
func (zp *ZonePlayer) Next() error {
	_, err := zp.SendAVTransport("Next", "", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) NotifyDeletedURI(deletedURI string) error {
	_, err := zp.SendAVTransport("NotifyDeletedURI", "<DeletedURI>"+deletedURI+"</DeletedURI>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) Pause() error {
	_, err := zp.SendAVTransport("Pause", "", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) Play() error {
	_, err := zp.SendAVTransport("Play", "<Speed>"+strconv.Itoa(zp.Speed)+"</Speed>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) Previous() error {
	_, err := zp.SendAVTransport("Previous", "", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) RemoveAllTracksFromQueue() error {
	_, err := zp.SendAVTransport("RemoveAllTracksFromQueue", "", "")
	return err
}

// TODO: Check with old implementation
// TODO: Test
func (zp *ZonePlayer) RemoveTrackFromQueue(objectID string, updateID int) error {
	_, err := zp.SendAVTransport("RemoveTrackFromQueue", "<ObjectID>"+objectID+"</ObjectID><UpdateID>"+strconv.Itoa(updateID)+"</UpdateID>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) RemoveTrackRangeFromQueue(updateID int, startingIndex int, numberOfTracks int) (int, error) {
	res, err := zp.SendAVTransport("RemoveTrackRangeFromQueue", "<UpdateID>"+strconv.Itoa(updateID)+"</UpdateID><StartingIndex>"+strconv.Itoa(startingIndex)+"</StartingIndex><NumberOfTracks>strconv.Itoa(numberOfTracks)</NumberOfTracks>", "NewUpdateID")
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(res)
}

// TODO: Test
func (zp *ZonePlayer) ReorderTracksInQueue(startingIndex int, numberOfTracks int, insertBefore int, updateID int) error {
	_, err := zp.SendAVTransport("ReorderTracksInQueue", "<StartingIndex>"+strconv.Itoa(startingIndex)+"</StartingIndex><NumberOfTracks>"+strconv.Itoa(numberOfTracks)+"</NumberOfTracks><InsertBefore>"+strconv.Itoa(insertBefore)+"</InsertBefore><UpdateID>"+strconv.Itoa(updateID)+"</UpdateID>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) ReorderTracksInSavedQueue(objectID string, updateID int, trackList string, newPositionList string) (reorderTracksInSavedQueueResponse, error) {
	res, err := zp.SendAVTransport("ReorderTracksInSavedQueue", "<ObjectID>"+objectID+"</ObjectID><UpdateID>"+strconv.Itoa(updateID)+"</UpdateID><TrackList>"+trackList+"</TrackList><NewPositionList>"+newPositionList+"</NewPositionList>", "")
	if err != nil {
		return reorderTracksInSavedQueueResponse{}, err
	}
	data := reorderTracksInSavedQueueResponse{}
	err = xml.Unmarshal([]byte(res), &data)
	return data, err
}

// TODO: Test
func (zp *ZonePlayer) RunAlarm(alarmID int, loggedStartTime string, duration string, programURI string, programMetaData string, playMode string, volume int, includeLinkedZones bool) error {
	_, err := zp.SendAVTransport("RunAlarm", "<AlarmID>"+strconv.Itoa(alarmID)+"</AlarmID><LoggedStartTime>"+loggedStartTime+"</LoggedStartTime><Duration>"+duration+"</Duration><ProgramURI>"+programURI+"</ProgramURI><ProgramMetaData>"+programMetaData+"</ProgramMetaData><PlayMode>"+playMode+"</PlayMode><Volume>"+strconv.Itoa(max(0, min(100, volume)))+"</Volume><IncludeLinkedZones>"+boolTo10(includeLinkedZones)+"</IncludeLinkedZones>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) SaveQueue(title string, objectID string) (string, error) {
	return zp.SendAVTransport("SaveQueue", "<Title>"+title+"</Title><ObjectID>"+objectID+"</ObjectID>", "AssignedObjectID")
}

// TODO: Test
func (zp *ZonePlayer) SeekTrack(track int) error {
	_, err := zp.SendAVTransport("Seek", "<Unit>TRACK_NR</Unit><Target>"+strconv.Itoa(max(1, track))+"</Target>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) SeekTime(seconds int) error {
	_, err := zp.SendAVTransport("Seek", "<Unit>REL_TIME</Unit><Target>"+time.Time.Add(time.Time{}, time.Second*time.Duration(max(0, seconds))).Format("15:04:05")+"</Target>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) SeekTimeDelta(seconds int) error {
	factor := "+"
	if seconds < 0 {
		seconds = -seconds
		factor = "-"
	}
	_, err := zp.SendAVTransport("Seek", "<Unit>TIME_DELTA</Unit><Target>"+factor+time.Time.Add(time.Time{}, time.Second*time.Duration(seconds)).Format("15:04:05")+"</Target>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) SetAVTransportURI(currentURI string, currentURIMetaData string) error {
	_, err := zp.SendAVTransport("SetAVTransportURI", "<CurrentURI>"+currentURI+"</CurrentURI><CurrentURIMetaData>"+currentURIMetaData+"</CurrentURIMetaData>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) SetCrossfadeMode(crossfadeMode bool) error {
	_, err := zp.SendAVTransport("SetCrossfadeMode", "<CrossfadeMode>"+boolTo10(crossfadeMode)+"</CrossfadeMode>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) SetNextAVTransportURI(nextURI string, nextURIMetaData string) error {
	_, err := zp.SendAVTransport("SetNextAVTransportURI", "<NextURI>"+nextURI+"</NextURI><NextURIMetaData>"+nextURIMetaData+"</NextURIMetaData>", "")
	return err
}

// TODO: test
func (zp *ZonePlayer) SetPlayMode(shuffle bool, repeat bool, repeat_one bool) error {
	mode, ok := PlaymodesReversed[[3]bool{shuffle, repeat, repeat_one}]
	if !ok {
		return ErrSonos.ErrInvalidPlayMode
	}
	_, err := zp.SendAVTransport("SetPlayMode", "<NewPlayMode>"+mode+"</NewPlayMode>", "")
	return err
}

// TODO: test
func (zp *ZonePlayer) SnoozeAlarm(seconds int) error {
	_, err := zp.SendAVTransport("SnoozeAlarm", "<Duration>"+time.Time.Add(time.Time{}, time.Second*time.Duration(max(0, seconds))).Format("15:04:05")+"</Duration>", "")
	return err
}

// TODO: test
func (zp *ZonePlayer) StartAutoplay(programURI string, programMetaData string, volume int, includeLinkedZones bool, resetVolumeAfter bool) error {
	_, err := zp.SendAVTransport("StartAutoplay", "<ProgramURI>"+programURI+"</ProgramURI><ProgramMetaData>"+programMetaData+"</ProgramMetaData><Volume>"+strconv.Itoa(volume)+"</Volume><IncludeLinkedZones>"+boolTo10(includeLinkedZones)+"</IncludeLinkedZones><ResetVolumeAfter>"+boolTo10(resetVolumeAfter)+"</ResetVolumeAfter>", "")
	return err
}

// TODO: test
func (zp *ZonePlayer) Stop() error {
	_, err := zp.SendAVTransport("Stop", "", "")
	return err
}

// TODO: test
func (zp *ZonePlayer) GetShuffle() (bool, error) {
	shuffle, _, _, err := zp.GetPlayMode()
	return shuffle, err
}

// TODO: test
func (zp *ZonePlayer) SetShuffle(state bool) error {
	_, repeat, repeat_one, err := zp.GetPlayMode()
	if err != nil {
		return err
	}
	return zp.SetPlayMode(state, repeat, repeat_one)
}

// TODO: test
func (zp *ZonePlayer) GetRepeat() (bool, error) {
	_, repeat, _, err := zp.GetPlayMode()
	return repeat, err
}

// TODO: test
func (zp *ZonePlayer) SetRepeat(state bool) error {
	shuffle, _, repeat_one, err := zp.GetPlayMode()
	if err != nil {
		return err
	}
	return zp.SetPlayMode(shuffle, state, repeat_one && !state)
}

// TODO: test
func (zp *ZonePlayer) GetRepeatOne() (bool, error) {
	_, _, repeat_one, err := zp.GetPlayMode()
	return repeat_one, err
}

// TODO: test
func (zp *ZonePlayer) SetRepeatOne(state bool) error {
	shuffle, repeat, _, err := zp.GetPlayMode()
	if err != nil {
		return err
	}
	return zp.SetPlayMode(shuffle, repeat && !state, state)
}
