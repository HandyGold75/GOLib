package Gonos

// https://sonos.svrooij.io/services/device-properties

import (
	"encoding/xml"
	"strconv"
	"time"
)

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
		CurrentURIMetaData       string
		CurrentURIMetaDataParsed struct {
			// TODO: Fill in
		}

		NextURI string
		// Embedded XML
		NextURIMetaData string
		// Possible values: `NONE` / `NETWORK`
		NextURIMetaDataParsed struct {
			// TODO: Fill in
		}

		PlayMedium string
		// Possible values: `NONE`
		RecordMedium string
		WriteStatus  string
	}
	getPositionInfoResponse struct {
		XMLName       xml.Name `xml:"GetPositionInfoResponse"`
		Track         int
		TrackDuration string
		// Embedded XML
		TrackMetaData       string
		TrackMetaDataParsed struct {
			XMLName       xml.Name `xml:"item"`
			Res           string   `xml:"res"`
			StreamContent string   `xml:"streamContent"`
			AlbumArtUri   string   `xml:"albumArtURI"`
			Title         string   `xml:"title"`
			Class         string   `xml:"class"`
			Creator       string   `xml:"creator"`
			Album         string   `xml:"album"`
		}
		TrackURI string
		RelTime  string
		AbsTime  string
		RelCount int
		AbsCount int
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
	getTransportInfoResponse struct {
		XMLName xml.Name `xml:"GetTransportInfoResponse"`
		// Possible values: `STOPPED` / `PLAYING` / `PAUSED_PLAYBACK` / `TRANSITIONING`
		CurrentTransportState  string
		CurrentTransportStatus string
		// Possible values: `1`
		CurrentSpeed string
	}
	getTransportSettingsResponse struct {
		XMLName xml.Name `xml:"GetTransportSettingsResponse"`
		// Possible values: `NORMAL` / `REPEAT_ALL` / `REPEAT_ONE` / `SHUFFLE_NOREPEAT` / `SHUFFLE` / `SHUFFLE_REPEAT_ONE`
		PlayMode       string
		RecQualityMode string
	}
	reorderTracksInSavedQueueResponse struct {
		XMLName           xml.Name `xml:"ReorderTracksInSavedQueueResponse"`
		QueueLengthChange int
		NewQueueLength    int
		NewUpdateID       int
	}
)

// TODO: test
func (zp *ZonePlayer) AddMultipleURIsToQueue(numberOfURIs int, enqueuedURIs string, enqueuedURIsMetaData string, containerURI string, containerMetaData string, desiredFirstTrackNumberEnqueued int, enqueueAsNext bool) (addMultipleURIsToQueueResponse, error) {
	res, err := zp.SendAVTransport("AddMultipleURIsToQueue", "<UpdateID>"+strconv.Itoa(zp.Static.AVTransport.UpdateID)+"</UpdateID><NumberOfURIs>"+strconv.Itoa(numberOfURIs)+"</NumberOfURIs><EnqueuedURIs>"+enqueuedURIs+"</EnqueuedURIs><EnqueuedURIsMetaData>"+enqueuedURIsMetaData+"</EnqueuedURIsMetaData><ContainerURI>"+containerURI+"</ContainerURI><ContainerMetaData>"+containerMetaData+"</ContainerMetaData><DesiredFirstTrackNumberEnqueued>"+strconv.Itoa(desiredFirstTrackNumberEnqueued)+"</DesiredFirstTrackNumberEnqueued><EnqueueAsNext>"+boolTo10(enqueueAsNext)+"</EnqueueAsNext>", "s:Body")
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

// `contentType` should be one of `Gonos.ContentTypes.*`
//
// TODO: Test
func (zp *ZonePlayer) AddURIToSavedQueue(contentType string, enqueuedURI string, enqueuedURIMetaData string, addAtIndex int) (addURIToSavedQueueResponse, error) {
	res, err := zp.SendAVTransport("AddURIToSavedQueue", "<ObjectID>"+contentType+"</ObjectID><UpdateID>"+strconv.Itoa(zp.Static.AVTransport.UpdateID)+"</UpdateID><EnqueuedURI>"+enqueuedURI+"</EnqueuedURI><EnqueuedURIMetaData>"+enqueuedURIMetaData+"</EnqueuedURIMetaData><AddAtIndex>"+strconv.Itoa(addAtIndex)+"</AddAtIndex>", "s:Body")
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
	if err != nil {
		return getMediaInfoResponse{}, err
	}
	err = unmarshalMetaData(data.CurrentURIMetaData, &data.CurrentURIMetaDataParsed)
	if err != nil {
		return getMediaInfoResponse{}, err
	}
	err = unmarshalMetaData(data.NextURIMetaData, &data.NextURIMetaDataParsed)
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
	if err != nil {
		return getPositionInfoResponse{}, err
	}
	err = unmarshalMetaData(data.TrackMetaData, &data.TrackMetaDataParsed)
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
func (zp *ZonePlayer) GetTransportInfo() (getTransportInfoResponse, error) {
	res, err := zp.SendAVTransport("GetTransportInfo", "", "s:Body")
	if err != nil {
		return getTransportInfoResponse{}, err
	}
	data := getTransportInfoResponse{}
	err = xml.Unmarshal([]byte(res), &data)
	return data, err
}

// TODO: Test
func (zp *ZonePlayer) GetTransportSettings() (getTransportSettingsResponse, error) {
	res, err := zp.SendAVTransport("GetTransportSettings", "", "s:Body")
	if err != nil {
		return getTransportSettingsResponse{}, err
	}
	data := getTransportSettingsResponse{}
	err = xml.Unmarshal([]byte(res), &data)
	return data, err
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
	_, err := zp.SendAVTransport("Play", "<Speed>"+strconv.Itoa(zp.Static.AVTransport.Speed)+"</Speed>", "")
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

// `contentType` should be one of `Gonos.ContentTypes.*`
//
// TODO: Test
func (zp *ZonePlayer) RemoveTrackFromQueue(contentType string, track int) error {
	_, err := zp.SendAVTransport("RemoveTrackFromQueue", "<ObjectID>"+contentType+"/"+strconv.Itoa(max(1, track))+"</ObjectID><UpdateID>"+strconv.Itoa(zp.Static.AVTransport.UpdateID)+"</UpdateID>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) RemoveTrackRangeFromQueue(start int, count int) (int, error) {
	res, err := zp.SendAVTransport("RemoveTrackRangeFromQueue", "<UpdateID>"+strconv.Itoa(zp.Static.AVTransport.UpdateID)+"</UpdateID><StartingIndex>"+strconv.Itoa(start)+"</StartingIndex><NumberOfTracks>"+strconv.Itoa(count)+"</NumberOfTracks>", "NewUpdateID")
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(res)
}

// TODO: Test
func (zp *ZonePlayer) ReorderTracksInQueue(start int, count int, insertBefore int) error {
	_, err := zp.SendAVTransport("ReorderTracksInQueue", "<StartingIndex>"+strconv.Itoa(start)+"</StartingIndex><NumberOfTracks>"+strconv.Itoa(count)+"</NumberOfTracks><InsertBefore>"+strconv.Itoa(insertBefore)+"</InsertBefore><UpdateID>"+strconv.Itoa(zp.Static.AVTransport.UpdateID)+"</UpdateID>", "")
	return err
}

// `contentType` should be one of `Gonos.ContentTypes.*`
//
// TODO: Test
func (zp *ZonePlayer) ReorderTracksInSavedQueue(contentType string, trackList string, newPositionList string) (reorderTracksInSavedQueueResponse, error) {
	res, err := zp.SendAVTransport("ReorderTracksInSavedQueue", "<ObjectID>"+contentType+"</ObjectID><UpdateID>"+strconv.Itoa(zp.Static.AVTransport.UpdateID)+"</UpdateID><TrackList>"+trackList+"</TrackList><NewPositionList>"+newPositionList+"</NewPositionList>", "")
	if err != nil {
		return reorderTracksInSavedQueueResponse{}, err
	}
	data := reorderTracksInSavedQueueResponse{}
	err = xml.Unmarshal([]byte(res), &data)
	return data, err
}

// `playMode` should be one of `Gonos.PlayModes.*`
//
// TODO: Test
func (zp *ZonePlayer) RunAlarm(alarmID int, loggedStartTime string, duration string, programURI string, programMetaData string, playMode string, volume int, includeLinkedZones bool) error {
	_, err := zp.SendAVTransport("RunAlarm", "<AlarmID>"+strconv.Itoa(alarmID)+"</AlarmID><LoggedStartTime>"+loggedStartTime+"</LoggedStartTime><Duration>"+duration+"</Duration><ProgramURI>"+programURI+"</ProgramURI><ProgramMetaData>"+programMetaData+"</ProgramMetaData><PlayMode>"+playMode+"</PlayMode><Volume>"+strconv.Itoa(max(0, min(100, volume)))+"</Volume><IncludeLinkedZones>"+boolTo10(includeLinkedZones)+"</IncludeLinkedZones>", "")
	return err
}

// `contentType` should be one of `Gonos.ContentTypes.*`
//
// Returns the objectID of the new que.
func (zp *ZonePlayer) SaveQueue(title string) (string, error) {
	return zp.SendAVTransport("SaveQueue", "<Title>"+title+"</Title><ObjectID></ObjectID>", "AssignedObjectID")
}

// Prefer methods `zp.SeekTrack`, `zp.SeekTime` or `zp.SeekTimeDelta`.
//
// `unit` should be one of `Gonos.SeekModes.*`.
func (zp *ZonePlayer) Seek(unit string, target string) error {
	_, err := zp.SendAVTransport("Seek", "<Unit>"+unit+"</Unit><Target>"+target+"</Target>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) SetAVTransportURI(currentURI string, currentURIMetaData string) error {
	_, err := zp.SendAVTransport("SetAVTransportURI", "<CurrentURI>"+currentURI+"</CurrentURI><CurrentURIMetaData>"+currentURIMetaData+"</CurrentURIMetaData>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) SetCrossfadeMode(state bool) error {
	_, err := zp.SendAVTransport("SetCrossfadeMode", "<CrossfadeMode>"+boolTo10(state)+"</CrossfadeMode>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) SetNextAVTransportURI(nextURI string, nextURIMetaData string) error {
	_, err := zp.SendAVTransport("SetNextAVTransportURI", "<NextURI>"+nextURI+"</NextURI><NextURIMetaData>"+nextURIMetaData+"</NextURIMetaData>", "")
	return err
}

// TODO: test
func (zp *ZonePlayer) SetPlayMode(shuffle bool, repeat bool, repeatOne bool) error {
	mode, ok := PlayModeMapReversed[[3]bool{shuffle, repeat, repeatOne}]
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
