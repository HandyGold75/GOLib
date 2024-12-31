package AVTransport

// https://sonos.svrooij.io/services/device-properties

import (
	"Gonos/lib"
	"encoding/xml"
	"strconv"
	"time"
)

type (
	AVTransport struct {
		Send func(action, body, targetTag string) (string, error)
		// Play speed usually `1`, can be a fraction of 1 Allowed values: `1`
		Speed int
		// Should always be `0`
		UpdateID int
	}

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
		CurrentURIMetaDataParsed struct{}

		NextURI string
		// Embedded XML
		NextURIMetaData string
		// Possible values: `NONE` / `NETWORK`
		NextURIMetaDataParsed struct{}

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

func New(send func(action, body, targetTag string) (string, error)) AVTransport {
	return AVTransport{Send: send, Speed: 1, UpdateID: 0}
}

func (zp *AVTransport) AddMultipleURIsToQueue(numberOfURIs int, enqueuedURIs string, enqueuedURIsMetaData string, containerURI string, containerMetaData string, desiredFirstTrackNumberEnqueued int, enqueueAsNext bool) (addMultipleURIsToQueueResponse, error) {
	res, err := zp.Send("AddMultipleURIsToQueue", "<UpdateID>"+strconv.Itoa(zp.UpdateID)+"</UpdateID><NumberOfURIs>"+strconv.Itoa(numberOfURIs)+"</NumberOfURIs><EnqueuedURIs>"+enqueuedURIs+"</EnqueuedURIs><EnqueuedURIsMetaData>"+enqueuedURIsMetaData+"</EnqueuedURIsMetaData><ContainerURI>"+containerURI+"</ContainerURI><ContainerMetaData>"+containerMetaData+"</ContainerMetaData><DesiredFirstTrackNumberEnqueued>"+strconv.Itoa(desiredFirstTrackNumberEnqueued)+"</DesiredFirstTrackNumberEnqueued><EnqueueAsNext>"+lib.BoolTo10(enqueueAsNext)+"</EnqueueAsNext>", "s:Body")
	if err != nil {
		return addMultipleURIsToQueueResponse{}, err
	}
	data := addMultipleURIsToQueueResponse{}
	err = xml.Unmarshal([]byte(res), &data)
	return data, err
}

func (zp *AVTransport) AddURIToQueue(enqueuedURI string, enqueuedURIMetaData string, desiredFirstTrackNumberEnqueued int, enqueueAsNext bool) (addURIToQueueResponse, error) {
	res, err := zp.Send("AddURIToQueue", "<EnqueuedURI>"+enqueuedURI+"</EnqueuedURI><EnqueuedURIMetaData>"+enqueuedURIMetaData+"</EnqueuedURIMetaData><DesiredFirstTrackNumberEnqueued>"+strconv.Itoa(desiredFirstTrackNumberEnqueued)+"</DesiredFirstTrackNumberEnqueued><EnqueueAsNext>"+lib.BoolTo10(enqueueAsNext)+"</EnqueueAsNext>", "s:Body")
	if err != nil {
		return addURIToQueueResponse{}, err
	}
	data := addURIToQueueResponse{}
	err = xml.Unmarshal([]byte(res), &data)
	return data, err
}

// `contentType` should be one of `Gonos.ContentTypes.*`
func (zp *AVTransport) AddURIToSavedQueue(contentType string, enqueuedURI string, enqueuedURIMetaData string, addAtIndex int) (addURIToSavedQueueResponse, error) {
	res, err := zp.Send("AddURIToSavedQueue", "<ObjectID>"+contentType+"</ObjectID><UpdateID>"+strconv.Itoa(zp.UpdateID)+"</UpdateID><EnqueuedURI>"+enqueuedURI+"</EnqueuedURI><EnqueuedURIMetaData>"+enqueuedURIMetaData+"</EnqueuedURIMetaData><AddAtIndex>"+strconv.Itoa(addAtIndex)+"</AddAtIndex>", "s:Body")
	if err != nil {
		return addURIToSavedQueueResponse{}, err
	}
	data := addURIToSavedQueueResponse{}
	err = xml.Unmarshal([]byte(res), &data)
	return data, err
}

func (zp *AVTransport) BackupQueue() error {
	_, err := zp.Send("BackupQueue", "", "")
	return err
}

func (zp *AVTransport) BecomeCoordinatorOfStandaloneGroup() (becomeCoordinatorOfStandaloneGroupResponse, error) {
	res, err := zp.Send("BecomeCoordinatorOfStandaloneGroup", "", "s:Body")
	if err != nil {
		return becomeCoordinatorOfStandaloneGroupResponse{}, err
	}
	data := becomeCoordinatorOfStandaloneGroupResponse{}
	err = xml.Unmarshal([]byte(res), &data)
	return data, err
}

func (zp *AVTransport) BecomeGroupCoordinator(currentCoordinator string, currentGroupID string, otherMembers string, transportSettings string, currentURI string, currentURIMetaData string, sleepTimerState string, alarmState string, streamRestartState string, currentQueueTrackList string, currentVLIState string) error {
	_, err := zp.Send("BecomeGroupCoordinator", "<CurrentCoordinator>"+currentCoordinator+"</CurrentCoordinator><CurrentGroupID>"+currentGroupID+"</CurrentGroupID><OtherMembers>"+otherMembers+"</OtherMembers><TransportSettings>"+transportSettings+"</TransportSettings><CurrentURI>"+currentURI+"</CurrentURI><CurrentURIMetaData>"+currentURIMetaData+"</CurrentURIMetaData><SleepTimerState>"+sleepTimerState+"</SleepTimerState><AlarmState>"+alarmState+"</AlarmState><StreamRestartState>"+streamRestartState+"</StreamRestartState><CurrentQueueTrackList>"+currentQueueTrackList+"</CurrentQueueTrackList><CurrentVLIState>"+currentVLIState+"</CurrentVLIState>", "")
	return err
}

func (zp *AVTransport) BecomeGroupCoordinatorAndSource(currentCoordinator string, currentGroupID string, otherMembers string, currentURI string, currentURIMetaData string, sleepTimerState string, alarmState string, streamRestartState string, currentAVTTrackList string, currentQueueTrackList string, currentSourceState string, resumePlayback bool) error {
	_, err := zp.Send("BecomeGroupCoordinatorAndSource", "<CurrentCoordinator>"+currentCoordinator+"</CurrentCoordinator><CurrentGroupID>"+currentGroupID+"</CurrentGroupID><OtherMembers>"+otherMembers+"</OtherMembers><CurrentURI>"+currentURI+"</CurrentURI><CurrentURIMetaData>"+currentURIMetaData+"</CurrentURIMetaData><SleepTimerState>"+sleepTimerState+"</SleepTimerState><AlarmState>"+alarmState+"</AlarmState><StreamRestartState>"+streamRestartState+"</StreamRestartState><CurrentAVTTrackList>"+currentAVTTrackList+"</CurrentAVTTrackList><CurrentQueueTrackList>"+currentQueueTrackList+"</CurrentQueueTrackList><CurrentSourceState>"+currentSourceState+"</CurrentSourceState><ResumePlayback>"+lib.BoolTo10(resumePlayback)+"</ResumePlayback>", "")
	return err
}

func (zp *AVTransport) ChangeCoordinator(currentCoordinator string, newCoordinator string, newTransportSettings string, currentAVTransportURI string) error {
	_, err := zp.Send("ChangeCoordinator", "<CurrentCoordinator>"+currentCoordinator+"</CurrentCoordinator><NewCoordinator>"+newCoordinator+"</NewCoordinator><NewTransportSettings>"+newTransportSettings+"</NewTransportSettings><CurrentAVTransportURI>"+currentAVTransportURI+"</CurrentAVTransportURI>", "")
	return err
}

func (zp *AVTransport) ChangeTransportSettings(newTransportSettings string, currentAVTransportURI string) error {
	_, err := zp.Send("ChangeTransportSettings", "<NewTransportSettings>"+newTransportSettings+"</NewTransportSettings><CurrentAVTransportURI>"+currentAVTransportURI+"</CurrentAVTransportURI>", "")
	return err
}

func (zp *AVTransport) ConfigureSleepTimer(seconds int) error {
	_, err := zp.Send("ConfigureSleepTimer", "<NewSleepTimerDuration>"+time.Time.Add(time.Time{}, time.Second*time.Duration(seconds)).Format("15:04:05")+"</NewSleepTimerDuration>", "")
	return err
}

func (zp *AVTransport) CreateSavedQueue(title string, enqueuedURI string, enqueuedURIMetaData string) (createSavedQueueResponse, error) {
	res, err := zp.Send("CreateSavedQueue", "<Title>title</Title><EnqueuedURI>enqueuedURI</EnqueuedURI><EnqueuedURIMetaData>enqueuedURIMetaData</EnqueuedURIMetaData>", "s:Body")
	if err != nil {
		return createSavedQueueResponse{}, err
	}
	data := createSavedQueueResponse{}
	err = xml.Unmarshal([]byte(res), &data)
	return data, err
}

func (zp *AVTransport) DelegateGroupCoordinationTo(newCoordinator string, rejoinGroup bool) error {
	_, err := zp.Send("DelegateGroupCoordinationTo", "<NewCoordinator>"+newCoordinator+"</NewCoordinator><RejoinGroup>"+lib.BoolTo10(rejoinGroup)+"</RejoinGroup>", "")
	return err
}

func (zp *AVTransport) EndDirectControlSession() error {
	_, err := zp.Send("EndDirectControlSession", "", "")
	return err
}

func (zp *AVTransport) GetCrossfadeMode() (bool, error) {
	res, err := zp.Send("GetCrossfadeMode", "", "CrossfadeMode")
	return res == "1", err
}

func (zp *AVTransport) GetCurrentTransportActions() (string, error) {
	return zp.Send("GetCurrentTransportActions", "", "Actions")
}

func (zp *AVTransport) GetDeviceCapabilities() (getDeviceCapabilitiesResponse, error) {
	res, err := zp.Send("GetDeviceCapabilities", "", "s:Body")
	if err != nil {
		return getDeviceCapabilitiesResponse{}, err
	}
	data := getDeviceCapabilitiesResponse{}
	err = xml.Unmarshal([]byte(res), &data)
	return data, err
}

func (zp *AVTransport) GetMediaInfo() (getMediaInfoResponse, error) {
	res, err := zp.Send("GetMediaInfo", "", "s:Body")
	if err != nil {
		return getMediaInfoResponse{}, err
	}
	data := getMediaInfoResponse{}
	err = xml.Unmarshal([]byte(res), &data)
	if err != nil {
		return getMediaInfoResponse{}, err
	}
	err = lib.UnmarshalMetaData(data.CurrentURIMetaData, &data.CurrentURIMetaDataParsed)
	if err != nil {
		return getMediaInfoResponse{}, err
	}
	err = lib.UnmarshalMetaData(data.NextURIMetaData, &data.NextURIMetaDataParsed)
	return data, err
}

func (zp *AVTransport) GetPositionInfo() (getPositionInfoResponse, error) {
	res, err := zp.Send("GetPositionInfo", "", "s:Body")
	if err != nil {
		return getPositionInfoResponse{}, err
	}
	data := getPositionInfoResponse{}
	err = xml.Unmarshal([]byte(res), &data)
	if err != nil {
		return getPositionInfoResponse{}, err
	}
	err = lib.UnmarshalMetaData(data.TrackMetaData, &data.TrackMetaDataParsed)
	return data, err
}

func (zp *AVTransport) GetRemainingSleepTimerDuration() (getRemainingSleepTimerDurationResponse, error) {
	res, err := zp.Send("GetRemainingSleepTimerDuration", "", "s:Body")
	if err != nil {
		return getRemainingSleepTimerDurationResponse{}, err
	}
	data := getRemainingSleepTimerDurationResponse{}
	err = xml.Unmarshal([]byte(res), &data)
	return data, err
}

func (zp *AVTransport) GetRunningAlarmProperties() (getRunningAlarmPropertiesResponse, error) {
	res, err := zp.Send("GetRunningAlarmProperties", "", "s:Body")
	if err != nil {
		return getRunningAlarmPropertiesResponse{}, err
	}
	data := getRunningAlarmPropertiesResponse{}
	err = xml.Unmarshal([]byte(res), &data)
	return data, err
}

func (zp *AVTransport) GetTransportInfo() (getTransportInfoResponse, error) {
	res, err := zp.Send("GetTransportInfo", "", "s:Body")
	if err != nil {
		return getTransportInfoResponse{}, err
	}
	data := getTransportInfoResponse{}
	err = xml.Unmarshal([]byte(res), &data)
	return data, err
}

func (zp *AVTransport) GetTransportSettings() (getTransportSettingsResponse, error) {
	res, err := zp.Send("GetTransportSettings", "", "s:Body")
	if err != nil {
		return getTransportSettingsResponse{}, err
	}
	data := getTransportSettingsResponse{}
	err = xml.Unmarshal([]byte(res), &data)
	return data, err
}

func (zp *AVTransport) Next() error {
	_, err := zp.Send("Next", "", "")
	return err
}

func (zp *AVTransport) NotifyDeletedURI(deletedURI string) error {
	_, err := zp.Send("NotifyDeletedURI", "<DeletedURI>"+deletedURI+"</DeletedURI>", "")
	return err
}

func (zp *AVTransport) Pause() error {
	_, err := zp.Send("Pause", "", "")
	return err
}

func (zp *AVTransport) Play() error {
	_, err := zp.Send("Play", "<Speed>"+strconv.Itoa(zp.Speed)+"</Speed>", "")
	return err
}

func (zp *AVTransport) Previous() error {
	_, err := zp.Send("Previous", "", "")
	return err
}

func (zp *AVTransport) RemoveAllTracksFromQueue() error {
	_, err := zp.Send("RemoveAllTracksFromQueue", "", "")
	return err
}

// `contentType` should be one of `Gonos.ContentTypes.*`
func (zp *AVTransport) RemoveTrackFromQueue(contentType string, track int) error {
	_, err := zp.Send("RemoveTrackFromQueue", "<ObjectID>"+contentType+"/"+strconv.Itoa(max(1, track))+"</ObjectID><UpdateID>"+strconv.Itoa(zp.UpdateID)+"</UpdateID>", "")
	return err
}

func (zp *AVTransport) RemoveTrackRangeFromQueue(start int, count int) (int, error) {
	res, err := zp.Send("RemoveTrackRangeFromQueue", "<UpdateID>"+strconv.Itoa(zp.UpdateID)+"</UpdateID><StartingIndex>"+strconv.Itoa(start)+"</StartingIndex><NumberOfTracks>"+strconv.Itoa(count)+"</NumberOfTracks>", "NewUpdateID")
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(res)
}

func (zp *AVTransport) ReorderTracksInQueue(start int, count int, insertBefore int) error {
	_, err := zp.Send("ReorderTracksInQueue", "<StartingIndex>"+strconv.Itoa(start)+"</StartingIndex><NumberOfTracks>"+strconv.Itoa(count)+"</NumberOfTracks><InsertBefore>"+strconv.Itoa(insertBefore)+"</InsertBefore><UpdateID>"+strconv.Itoa(zp.UpdateID)+"</UpdateID>", "")
	return err
}

// `contentType` should be one of `Gonos.ContentTypes.*`
func (zp *AVTransport) ReorderTracksInSavedQueue(contentType string, trackList string, newPositionList string) (reorderTracksInSavedQueueResponse, error) {
	res, err := zp.Send("ReorderTracksInSavedQueue", "<ObjectID>"+contentType+"</ObjectID><UpdateID>"+strconv.Itoa(zp.UpdateID)+"</UpdateID><TrackList>"+trackList+"</TrackList><NewPositionList>"+newPositionList+"</NewPositionList>", "")
	if err != nil {
		return reorderTracksInSavedQueueResponse{}, err
	}
	data := reorderTracksInSavedQueueResponse{}
	err = xml.Unmarshal([]byte(res), &data)
	return data, err
}

// `playMode` should be one of `Gonos.PlayModes.*`
func (zp *AVTransport) RunAlarm(alarmID int, loggedStartTime string, duration string, programURI string, programMetaData string, playMode string, volume int, includeLinkedZones bool) error {
	_, err := zp.Send("RunAlarm", "<AlarmID>"+strconv.Itoa(alarmID)+"</AlarmID><LoggedStartTime>"+loggedStartTime+"</LoggedStartTime><Duration>"+duration+"</Duration><ProgramURI>"+programURI+"</ProgramURI><ProgramMetaData>"+programMetaData+"</ProgramMetaData><PlayMode>"+playMode+"</PlayMode><Volume>"+strconv.Itoa(max(0, min(100, volume)))+"</Volume><IncludeLinkedZones>"+lib.BoolTo10(includeLinkedZones)+"</IncludeLinkedZones>", "")
	return err
}

// `contentType` should be one of `Gonos.ContentTypes.*`
//
// Returns the objectID of the new que.
func (zp *AVTransport) SaveQueue(title string) (string, error) {
	return zp.Send("SaveQueue", "<Title>"+title+"</Title><ObjectID></ObjectID>", "AssignedObjectID")
}

// Prefer methods `zp.SeekTrack`, `zp.SeekTime` or `zp.SeekTimeDelta`.
//
// `unit` should be one of `Gonos.SeekModes.*`.
func (zp *AVTransport) Seek(unit string, target string) error {
	_, err := zp.Send("Seek", "<Unit>"+unit+"</Unit><Target>"+target+"</Target>", "")
	return err
}

func (zp *AVTransport) SetAVTransportURI(currentURI string, currentURIMetaData string) error {
	_, err := zp.Send("SetAVTransportURI", "<CurrentURI>"+currentURI+"</CurrentURI><CurrentURIMetaData>"+currentURIMetaData+"</CurrentURIMetaData>", "")
	return err
}

func (zp *AVTransport) SetCrossfadeMode(state bool) error {
	_, err := zp.Send("SetCrossfadeMode", "<CrossfadeMode>"+lib.BoolTo10(state)+"</CrossfadeMode>", "")
	return err
}

func (zp *AVTransport) SetNextAVTransportURI(nextURI string, nextURIMetaData string) error {
	_, err := zp.Send("SetNextAVTransportURI", "<NextURI>"+nextURI+"</NextURI><NextURIMetaData>"+nextURIMetaData+"</NextURIMetaData>", "")
	return err
}

func (zp *AVTransport) SetPlayMode(shuffle bool, repeat bool, repeatOne bool) error {
	mode, ok := lib.PlayModeMapReversed[[3]bool{shuffle, repeat, repeatOne}]
	if !ok {
		return lib.ErrSonos.ErrInvalidPlayMode
	}
	_, err := zp.Send("SetPlayMode", "<NewPlayMode>"+mode+"</NewPlayMode>", "")
	return err
}

func (zp *AVTransport) SnoozeAlarm(seconds int) error {
	_, err := zp.Send("SnoozeAlarm", "<Duration>"+time.Time.Add(time.Time{}, time.Second*time.Duration(max(0, seconds))).Format("15:04:05")+"</Duration>", "")
	return err
}

func (zp *AVTransport) StartAutoplay(programURI string, programMetaData string, volume int, includeLinkedZones bool, resetVolumeAfter bool) error {
	_, err := zp.Send("StartAutoplay", "<ProgramURI>"+programURI+"</ProgramURI><ProgramMetaData>"+programMetaData+"</ProgramMetaData><Volume>"+strconv.Itoa(volume)+"</Volume><IncludeLinkedZones>"+lib.BoolTo10(includeLinkedZones)+"</IncludeLinkedZones><ResetVolumeAfter>"+lib.BoolTo10(resetVolumeAfter)+"</ResetVolumeAfter>", "")
	return err
}

func (zp *AVTransport) Stop() error {
	_, err := zp.Send("Stop", "", "")
	return err
}
