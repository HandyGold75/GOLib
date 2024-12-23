package Gonos

// https://sonos.svrooij.io/services/device-properties

import (
	"encoding/xml"
	"strconv"
)

type (
	ZoneInfo struct {
		XMLName                xml.Name `xml:"GetZoneInfoResponse"`
		SerialNumber           string
		SoftwareVersion        string
		DisplaySoftwareVersion string
		HardwareVersion        string
		IPAddress              string
		MACAddress             string
		CopyrightInfo          string
		ExtraInfo              string
		HTAudioIn              int
		Flags                  int
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
func (zp *ZonePlayer) GetAutoplayLinkedZones(source string) (bool, error) {
	res, err := zp.SendDeviceProperties("GetAutoplayLinkedZones", "<Source>"+source+"</Source>", "IncludeLinkedZones")
	return res == "1", err
}

// TODO: Test
func (zp *ZonePlayer) GetAutoplayRoomUUID(source string) (string, error) {
	res, err := zp.SendDeviceProperties("GetAutoplayRoomUUID", "<Source>"+source+"</Source>", "RoomUUID")
	return res, err
}

// TODO: Test
func (zp *ZonePlayer) GetAutoplayVolume(source string) (int, error) {
	res, err := zp.SendDeviceProperties("GetAutoplayVolume", "<Source>"+source+"</Source>", "CurrentVolume")
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
func (zp *ZonePlayer) GetUseAutoplayVolume(source string) (bool, error) {
	res, err := zp.SendDeviceProperties("GetUseAutoplayVolume", "<Source>"+source+"</Source>", "UseVolume")
	return res == "1", err
}

// TODO: Test
func (zp *ZonePlayer) GetZoneAttributeName() (string, error) {
	return zp.SendDeviceProperties("GetZoneAttributes", "", "CurrentZoneName")
}

// TODO: Test
func (zp *ZonePlayer) GetZoneAttributeIcon() (string, error) {
	return zp.SendDeviceProperties("GetZoneAttributes", "", "CurrentIcon")
}

// TODO: Test
func (zp *ZonePlayer) GetZoneAttributeConfiguration() (string, error) {
	return zp.SendDeviceProperties("GetZoneAttributes", "", "CurrentConfiguration")
}

// TODO: Test
func (zp *ZonePlayer) GetZoneAttributeTargetRoomName() (string, error) {
	return zp.SendDeviceProperties("GetZoneAttributes", "", "CurrentTargetRoomName")
}

// TODO: Test
func (zp *ZonePlayer) GetZoneInfo() (ZoneInfo, error) {
	res, err := zp.SendDeviceProperties("GetZoneInfo", "", "s:Body")
	if err != nil {
		return ZoneInfo{}, err
	}
	data := ZoneInfo{}
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
func (zp *ZonePlayer) RoomDetectionStartChirping(channel int, duration int, chirpIfPlayingSwappableAudio bool) (int, error) {
	res, err := zp.SendDeviceProperties("RoomDetectionStartChirping", "<Channel>"+strconv.Itoa(channel)+"</Channel><DurationMilliseconds>"+strconv.Itoa(duration)+"</DurationMilliseconds><ChirpIfPlayingSwappableAudio>"+boolTo10(chirpIfPlayingSwappableAudio)+"</ChirpIfPlayingSwappableAudio>", "PlayId")
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
func (zp *ZonePlayer) SetAutoplayLinkedZones(includeLinkedZones bool, source string) error {
	_, err := zp.SendDeviceProperties("SetAutoplayLinkedZones", "<IncludeLinkedZones>"+boolTo10(includeLinkedZones)+"</IncludeLinkedZones><Source>"+source+"</Source>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) SetAutoplayRoomUUID(roomUUID string, source string) error {
	_, err := zp.SendDeviceProperties("SetAutoplayRoomUUID", "<RoomUUID>"+roomUUID+"</RoomUUID><Source>"+source+"</Source>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) SetAutoplayVolume(volume int, source string) error {
	_, err := zp.SendDeviceProperties("SetAutoplayVolume", "<Volume>"+strconv.Itoa(max(0, min(100, volume)))+"</Volume><Source>"+source+"</Source>", "")
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
func (zp *ZonePlayer) SetUseAutoplayVolume(state bool, source string) error {
	_, err := zp.SendDeviceProperties("SetUseAutoplayVolume", "<UseVolume>"+boolTo10(state)+"</UseVolume><Source>"+source+"</Source>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) SetZoneAttributeZoneName(zoneName string) error {
	_, err := zp.SendDeviceProperties("SetZoneAttributes", "<DesiredZoneName>"+zoneName+"</DesiredZoneName>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) SetZoneAttributeIcon(icon string) error {
	_, err := zp.SendDeviceProperties("SetZoneAttributes", "<DesiredIcon>"+icon+"</DesiredIcon>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) SetZoneAttributeConfiguration(configuration string) error {
	_, err := zp.SendDeviceProperties("SetZoneAttributes", "<DesiredConfiguration>"+configuration+"</DesiredConfiguration>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) SetZoneAttributeTargetRoomName(targetRoomName string) error {
	_, err := zp.SendDeviceProperties("SetZoneAttributes", "<DesiredTargetRoomName>"+targetRoomName+"</DesiredTargetRoomName>", "")
	return err
}
