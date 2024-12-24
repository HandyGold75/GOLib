package Gonos

// https://sonos.svrooij.io/services/device-properties

import (
	"encoding/xml"
	"strconv"
)

type (
	getZoneAttributesResponse struct {
		XMLName               xml.Name `xml:"GetZoneAttributesResponse"`
		CurrentZoneName       string
		CurrentIcon           string
		CurrentConfiguration  string
		CurrentTargetRoomName string
	}
	getZoneInfoResponse struct {
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

// TODO: Test
func (zp *ZonePlayer) AddBondedZones(channelMapSet string) error {
	_, err := zp.SendDeviceProperties("AddBondedZones", "<ChannelMapSet>"+channelMapSet+"</ChannelMapSet>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) AddHTSatellite(hTSatChanMapSet string) error {
	_, err := zp.SendDeviceProperties("AddHTSatellite", "<HTSatChanMapSet>"+hTSatChanMapSet+"</HTSatChanMapSet>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) CreateStereoPair(channelMapSet string) error {
	_, err := zp.SendDeviceProperties("CreateStereoPair", "<ChannelMapSet>"+channelMapSet+"</ChannelMapSet>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) EnterConfigMode(mode string, options string) (string, error) {
	res, err := zp.SendDeviceProperties("EnterConfigMode", "<Mode>"+mode+"</Mode><Options>"+options+"</Options>", "State")
	return res, err
}

// TODO: Test
func (zp *ZonePlayer) ExitConfigMode(options string) error {
	_, err := zp.SendDeviceProperties("ExitConfigMode", "<Options>"+options+"</Options>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) GetAutoplayLinkedZones() (bool, error) {
	res, err := zp.SendDeviceProperties("GetAutoplayLinkedZones", "<Source>"+zp.Static.DeviceProperties.Source+"</Source>", "IncludeLinkedZones")
	return res == "1", err
}

// TODO: Test
func (zp *ZonePlayer) GetAutoplayRoomUUID() (string, error) {
	res, err := zp.SendDeviceProperties("GetAutoplayRoomUUID", "<Source>"+zp.Static.DeviceProperties.Source+"</Source>", "RoomUUID")
	return res, err
}

// TODO: Test
func (zp *ZonePlayer) GetAutoplayVolume() (int, error) {
	res, err := zp.SendDeviceProperties("GetAutoplayVolume", "<Source>"+zp.Static.DeviceProperties.Source+"</Source>", "CurrentVolume")
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(res)
}

// TODO: Test
func (zp *ZonePlayer) GetButtonLockState() (bool, error) {
	res, err := zp.SendDeviceProperties("GetButtonLockState", "", "CurrentButtonLockState")
	return res == "On", err
}

// TODO: Test
func (zp *ZonePlayer) GetButtonState() (string, error) {
	res, err := zp.SendDeviceProperties("GetButtonState", "", "State")
	return res, err
}

// TODO: Test
func (zp *ZonePlayer) GetHouseholdID() (string, error) {
	res, err := zp.SendDeviceProperties("GetHouseholdID", "", "CurrentHouseholdID")
	return res, err
}

// TODO: Test
func (zp *ZonePlayer) GetHTForwardState() (bool, error) {
	res, err := zp.SendDeviceProperties("GetHTForwardState", "", "IsHTForwardEnabled")
	return res == "1", err
}

// TODO: Test
func (zp *ZonePlayer) GetLEDState() (bool, error) {
	res, err := zp.SendDeviceProperties("GetLEDState", "", "CurrentLEDState")
	return res == "On", err
}

// TODO: Test
func (zp *ZonePlayer) GetUseAutoplayVolume() (bool, error) {
	res, err := zp.SendDeviceProperties("GetUseAutoplayVolume", "<Source>"+zp.Static.DeviceProperties.Source+"</Source>", "UseVolume")
	return res == "1", err
}

// TODO: Test
func (zp *ZonePlayer) GetZoneAttributes() (getZoneAttributesResponse, error) {
	res, err := zp.SendDeviceProperties("GetZoneAttributes", "", "s:Body")
	if err != nil {
		return getZoneAttributesResponse{}, err
	}
	data := getZoneAttributesResponse{}
	err = xml.Unmarshal([]byte(res), &data)
	return data, err
}

// TODO: Test
func (zp *ZonePlayer) GetZoneInfo() (getZoneInfoResponse, error) {
	res, err := zp.SendDeviceProperties("GetZoneInfo", "", "s:Body")
	if err != nil {
		return getZoneInfoResponse{}, err
	}
	data := getZoneInfoResponse{}
	err = xml.Unmarshal([]byte(res), &data)
	return data, err
}

// TODO: Test
func (zp *ZonePlayer) RemoveBondedZones(channelMapSet string, keepGrouped bool) error {
	_, err := zp.SendDeviceProperties("RemoveBondedZones", "<ChannelMapSet>"+channelMapSet+"</ChannelMapSet><KeepGrouped>"+boolTo10(keepGrouped)+"</KeepGrouped>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) RemoveHTSatellite(satRoomUUID string) error {
	_, err := zp.SendDeviceProperties("RemoveHTSatellite", "<SatRoomUUID>"+satRoomUUID+"</SatRoomUUID>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) RoomDetectionStartChirping(channel int, milliseconds int, chirpIfPlayingSwappableAudio bool) (int, error) {
	res, err := zp.SendDeviceProperties("RoomDetectionStartChirping", "<Channel>"+strconv.Itoa(channel)+"</Channel><DurationMilliseconds>"+strconv.Itoa(milliseconds)+"</DurationMilliseconds><ChirpIfPlayingSwappableAudio>"+boolTo10(chirpIfPlayingSwappableAudio)+"</ChirpIfPlayingSwappableAudio>", "PlayId")
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(res)
}

// TODO: Test
func (zp *ZonePlayer) RoomDetectionStopChirping(playId int) error {
	_, err := zp.SendDeviceProperties("RoomDetectionStopChirping", "<PlayId>"+strconv.Itoa(playId)+"</PlayId>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) SeparateStereoPair(channelMapSet string) error {
	_, err := zp.SendDeviceProperties("SeparateStereoPair", "<ChannelMapSet>"+channelMapSet+"</ChannelMapSet>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) SetAutoplayLinkedZones(includeLinkedZones bool) error {
	_, err := zp.SendDeviceProperties("SetAutoplayLinkedZones", "<IncludeLinkedZones>"+boolTo10(includeLinkedZones)+"</IncludeLinkedZones><Source>"+zp.Static.DeviceProperties.Source+"</Source>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) SetAutoplayRoomUUID(roomUUID string) error {
	_, err := zp.SendDeviceProperties("SetAutoplayRoomUUID", "<RoomUUID>"+roomUUID+"</RoomUUID><Source>"+zp.Static.DeviceProperties.Source+"</Source>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) SetAutoplayVolume(volume int) error {
	_, err := zp.SendDeviceProperties("SetAutoplayVolume", "<Volume>"+strconv.Itoa(max(0, min(100, volume)))+"</Volume><Source>"+zp.Static.DeviceProperties.Source+"</Source>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) SetButtonLockState(state bool) error {
	_, err := zp.SendDeviceProperties("SetButtonLockState", "<DesiredButtonLockState>"+boolToOnOff(state)+"</DesiredButtonLockState>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) SetLEDState(state bool) error {
	_, err := zp.SendDeviceProperties("SetLEDState", "<DesiredLEDState>"+boolToOnOff(state)+"</DesiredLEDState>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) SetUseAutoplayVolume(state bool) error {
	_, err := zp.SendDeviceProperties("SetUseAutoplayVolume", "<UseVolume>"+boolTo10(state)+"</UseVolume><Source>"+zp.Static.DeviceProperties.Source+"</Source>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) SetZoneAttributes(zoneName string, icon string, configuration string, targetRoomName string) error {
	_, err := zp.SendDeviceProperties("SetZoneAttributes", "<DesiredZoneName>"+zoneName+"</DesiredZoneName><DesiredIcon>"+icon+"</DesiredIcon><DesiredConfiguration>"+configuration+"</DesiredConfiguration><DesiredTargetRoomName>"+targetRoomName+"</DesiredTargetRoomName>", "")
	return err
}
