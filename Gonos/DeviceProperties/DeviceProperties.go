package DeviceProperties

// https://sonos.svrooij.io/services/device-properties

import (
	"Gonos/lib"
	"encoding/xml"
	"strconv"
)

type (
	DeviceProperties struct {
		Send   func(action, body, targetTag string) (string, error)
		Source string
	}

	GetZoneAttributesResponse struct {
		XMLName               xml.Name `xml:"GetZoneAttributesResponse"`
		CurrentZoneName       string
		CurrentIcon           string
		CurrentConfiguration  string
		CurrentTargetRoomName string
	}
	GetZoneInfoResponse struct {
		XMLName                xml.Name `xml:"GetZoneInfoResponse"`
		SerialNumber           string
		SoftwareVersion        string
		DisplaySoftwareVersion string
		HardwareVersion        string
		IPAddress              string
		MACAddress             string
		CopyrightInfo          string
		ExtraInfo              string
		// SPDIF input, 0 not connected / 2 stereo / 7 Dolby 2.0 / 18 dolby 5.1 / 21 not listening / 22 silence
		HTAudioIn int
		Flags     int
	}
)

func New(send func(action, body, targetTag string) (string, error)) DeviceProperties {
	return DeviceProperties{Send: send, Source: ""}
}

func (zp *DeviceProperties) AddBondedZones(channelMapSet string) error {
	_, err := zp.Send("AddBondedZones", "<ChannelMapSet>"+channelMapSet+"</ChannelMapSet>", "")
	return err
}

func (zp *DeviceProperties) AddHTSatellite(hTSatChanMapSet string) error {
	_, err := zp.Send("AddHTSatellite", "<HTSatChanMapSet>"+hTSatChanMapSet+"</HTSatChanMapSet>", "")
	return err
}

func (zp *DeviceProperties) CreateStereoPair(channelMapSet string) error {
	_, err := zp.Send("CreateStereoPair", "<ChannelMapSet>"+channelMapSet+"</ChannelMapSet>", "")
	return err
}

func (zp *DeviceProperties) EnterConfigMode(mode string, options string) (string, error) {
	res, err := zp.Send("EnterConfigMode", "<Mode>"+mode+"</Mode><Options>"+options+"</Options>", "State")
	return res, err
}

func (zp *DeviceProperties) ExitConfigMode(options string) error {
	_, err := zp.Send("ExitConfigMode", "<Options>"+options+"</Options>", "")
	return err
}

func (zp *DeviceProperties) GetAutoplayLinkedZones() (bool, error) {
	res, err := zp.Send("GetAutoplayLinkedZones", "<Source>"+zp.Source+"</Source>", "IncludeLinkedZones")
	return res == "1", err
}

func (zp *DeviceProperties) GetAutoplayRoomUUID() (string, error) {
	res, err := zp.Send("GetAutoplayRoomUUID", "<Source>"+zp.Source+"</Source>", "RoomUUID")
	return res, err
}

func (zp *DeviceProperties) GetAutoplayVolume() (int, error) {
	res, err := zp.Send("GetAutoplayVolume", "<Source>"+zp.Source+"</Source>", "CurrentVolume")
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(res)
}

func (zp *DeviceProperties) GetButtonLockState() (bool, error) {
	res, err := zp.Send("GetButtonLockState", "", "CurrentButtonLockState")
	return res == "On", err
}

func (zp *DeviceProperties) GetButtonState() (string, error) {
	res, err := zp.Send("GetButtonState", "", "State")
	return res, err
}

func (zp *DeviceProperties) GetHouseholdID() (string, error) {
	res, err := zp.Send("GetHouseholdID", "", "CurrentHouseholdID")
	return res, err
}

func (zp *DeviceProperties) GetHTForwardState() (bool, error) {
	res, err := zp.Send("GetHTForwardState", "", "IsHTForwardEnabled")
	return res == "1", err
}

func (zp *DeviceProperties) GetLEDState() (bool, error) {
	res, err := zp.Send("GetLEDState", "", "CurrentLEDState")
	return res == "On", err
}

func (zp *DeviceProperties) GetUseAutoplayVolume() (bool, error) {
	res, err := zp.Send("GetUseAutoplayVolume", "<Source>"+zp.Source+"</Source>", "UseVolume")
	return res == "1", err
}

func (zp *DeviceProperties) GetZoneAttributes() (GetZoneAttributesResponse, error) {
	res, err := zp.Send("GetZoneAttributes", "", "s:Body")
	if err != nil {
		return GetZoneAttributesResponse{}, err
	}
	data := GetZoneAttributesResponse{}
	err = xml.Unmarshal([]byte(res), &data)
	return data, err
}

func (zp *DeviceProperties) GetZoneInfo() (GetZoneInfoResponse, error) {
	res, err := zp.Send("GetZoneInfo", "", "s:Body")
	if err != nil {
		return GetZoneInfoResponse{}, err
	}
	data := GetZoneInfoResponse{}
	err = xml.Unmarshal([]byte(res), &data)
	return data, err
}

func (zp *DeviceProperties) RemoveBondedZones(channelMapSet string, keepGrouped bool) error {
	_, err := zp.Send("RemoveBondedZones", "<ChannelMapSet>"+channelMapSet+"</ChannelMapSet><KeepGrouped>"+lib.BoolTo10(keepGrouped)+"</KeepGrouped>", "")
	return err
}

func (zp *DeviceProperties) RemoveHTSatellite(satRoomUUID string) error {
	_, err := zp.Send("RemoveHTSatellite", "<SatRoomUUID>"+satRoomUUID+"</SatRoomUUID>", "")
	return err
}

func (zp *DeviceProperties) RoomDetectionStartChirping(channel int, milliseconds int, chirpIfPlayingSwappableAudio bool) (int, error) {
	res, err := zp.Send("RoomDetectionStartChirping", "<Channel>"+strconv.Itoa(channel)+"</Channel><DurationMilliseconds>"+strconv.Itoa(milliseconds)+"</DurationMilliseconds><ChirpIfPlayingSwappableAudio>"+lib.BoolTo10(chirpIfPlayingSwappableAudio)+"</ChirpIfPlayingSwappableAudio>", "PlayId")
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(res)
}

func (zp *DeviceProperties) RoomDetectionStopChirping(playId int) error {
	_, err := zp.Send("RoomDetectionStopChirping", "<PlayId>"+strconv.Itoa(playId)+"</PlayId>", "")
	return err
}

func (zp *DeviceProperties) SeparateStereoPair(channelMapSet string) error {
	_, err := zp.Send("SeparateStereoPair", "<ChannelMapSet>"+channelMapSet+"</ChannelMapSet>", "")
	return err
}

func (zp *DeviceProperties) SetAutoplayLinkedZones(includeLinkedZones bool) error {
	_, err := zp.Send("SetAutoplayLinkedZones", "<IncludeLinkedZones>"+lib.BoolTo10(includeLinkedZones)+"</IncludeLinkedZones><Source>"+zp.Source+"</Source>", "")
	return err
}

func (zp *DeviceProperties) SetAutoplayRoomUUID(roomUUID string) error {
	_, err := zp.Send("SetAutoplayRoomUUID", "<RoomUUID>"+roomUUID+"</RoomUUID><Source>"+zp.Source+"</Source>", "")
	return err
}

func (zp *DeviceProperties) SetAutoplayVolume(volume int) error {
	_, err := zp.Send("SetAutoplayVolume", "<Volume>"+strconv.Itoa(max(0, min(100, volume)))+"</Volume><Source>"+zp.Source+"</Source>", "")
	return err
}

func (zp *DeviceProperties) SetButtonLockState(state bool) error {
	_, err := zp.Send("SetButtonLockState", "<DesiredButtonLockState>"+lib.BoolToOnOff(state)+"</DesiredButtonLockState>", "")
	return err
}

func (zp *DeviceProperties) SetLEDState(state bool) error {
	_, err := zp.Send("SetLEDState", "<DesiredLEDState>"+lib.BoolToOnOff(state)+"</DesiredLEDState>", "")
	return err
}

func (zp *DeviceProperties) SetUseAutoplayVolume(state bool) error {
	_, err := zp.Send("SetUseAutoplayVolume", "<UseVolume>"+lib.BoolTo10(state)+"</UseVolume><Source>"+zp.Source+"</Source>", "")
	return err
}

func (zp *DeviceProperties) SetZoneAttributes(zoneName string, icon string, configuration string, targetRoomName string) error {
	_, err := zp.Send("SetZoneAttributes", "<DesiredZoneName>"+zoneName+"</DesiredZoneName><DesiredIcon>"+icon+"</DesiredIcon><DesiredConfiguration>"+configuration+"</DesiredConfiguration><DesiredTargetRoomName>"+targetRoomName+"</DesiredTargetRoomName>", "")
	return err
}
