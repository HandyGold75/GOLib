package Gonos

import (
	"encoding/xml"
)

type (
	ZoneInfo struct {
		XMLName                xml.Name `xml:"item"`
		SerialNumber           string
		SoftwareVersion        string
		DisplaySoftwareVersion string
		HardwareVersion        string
		IPAddress              string
		MACAddress             string
		CopyrightInfo          string
		ExtraInfo              string
		HTAudioIn              string
		Flags                  string
	}
)

// TODO: Test
func (zp *ZonePlayer) AddBondedZones() error {
	_, err := zp.SendDeviceProperties("AddBondedZones", "<ChannelMapSet>string</ChannelMapSet>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) AddHTSatellite() error {
	_, err := zp.SendDeviceProperties("AddHTSatellite", "<HTSatChanMapSet>string</HTSatChanMapSet>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) CreateStereoPair() error {
	_, err := zp.SendDeviceProperties("CreateStereoPair", "<ChannelMapSet>string</ChannelMapSet>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) EnterConfigMode() (string, error) {
	res, err := zp.SendDeviceProperties("EnterConfigMode", "<Mode>string</Mode><Options>string</Options>", "State")
	return res, err
}

// TODO: Test
func (zp *ZonePlayer) ExitConfigMode() error {
	_, err := zp.SendDeviceProperties("ExitConfigMode", "<Options>string</Options>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) GetAutoplayLinkedZones() (bool, error) {
	res, err := zp.SendDeviceProperties("GetAutoplayLinkedZones", "<Source>string</Source>", "IncludeLinkedZones")
	return res == "1", err
}

// TODO: Test
func (zp *ZonePlayer) GetAutoplayRoomUUID() (string, error) {
	res, err := zp.SendDeviceProperties("GetAutoplayRoomUUID", "<Source>string</Source>", "RoomUUID")
	return res, err
}

// TODO: Test
func (zp *ZonePlayer) GetAutoplayVolume() (string, error) {
	res, err := zp.SendDeviceProperties("GetAutoplayVolume", "<Source>string</Source>", "CurrentVolume")
	return res, err
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
	res, err := zp.SendDeviceProperties("GetUseAutoplayVolume", "<Source>string</Source>", "UseVolume")
	return res == "1", err
}

// TODO: Test
func (zp *ZonePlayer) GetZoneName() (string, error) {
	return zp.SendDeviceProperties("GetZoneAttributes", "", "CurrentZoneName")
}

// TODO: Test
func (zp *ZonePlayer) GetIcon() (string, error) {
	return zp.SendDeviceProperties("GetZoneAttributes", "", "CurrentIcon")
}

// TODO: Test
func (zp *ZonePlayer) GetConfiguration() (string, error) {
	return zp.SendDeviceProperties("GetZoneAttributes", "", "CurrentConfiguration")
}

// TODO: Test
func (zp *ZonePlayer) GetTargetRoomName() (string, error) {
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
func (zp *ZonePlayer) RemoveBondedZones() error {
	_, err := zp.SendDeviceProperties("RemoveBondedZones", "<ChannelMapSet>string</ChannelMapSet><KeepGrouped>boolean</KeepGrouped>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) RemoveHTSatellite() error {
	_, err := zp.SendDeviceProperties("RemoveHTSatellite", "<SatRoomUUID>string</SatRoomUUID>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) RoomDetectionStartChirping() (string, error) {
	return zp.SendDeviceProperties("RoomDetectionStartChirping", "<Channel>ui2</Channel><DurationMilliseconds>ui4</DurationMilliseconds><ChirpIfPlayingSwappableAudio>boolean</ChirpIfPlayingSwappableAudio>", "PlayId")
}

// TODO: Test
func (zp *ZonePlayer) RoomDetectionStopChirping() (string, error) {
	return zp.SendDeviceProperties("RoomDetectionStopChirping", "<PlayId>ui4</PlayId>", "PlayId")
}

// TODO: Test
func (zp *ZonePlayer) SeparateStereoPair() error {
	_, err := zp.SendDeviceProperties("SeparateStereoPair", "<ChannelMapSet>string</ChannelMapSet>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) SetAutoplayLinkedZones() error {
	_, err := zp.SendDeviceProperties("SetAutoplayLinkedZones", "<IncludeLinkedZones>boolean</IncludeLinkedZones><Source>string</Source>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) SetAutoplayRoomUUID() error {
	_, err := zp.SendDeviceProperties("SetAutoplayRoomUUID", "<RoomUUID>string</RoomUUID><Source>string</Source>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) SetAutoplayVolume() error {
	_, err := zp.SendDeviceProperties("SetAutoplayVolume", "<Volume>ui2</Volume><Source>string</Source>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) SetButtonLockState() error {
	_, err := zp.SendDeviceProperties("SetButtonLockState", "<DesiredButtonLockState>string</DesiredButtonLockState>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) SetLEDState() error {
	_, err := zp.SendDeviceProperties("SetLEDState", "<DesiredLEDState>string</DesiredLEDState>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) SetUseAutoplayVolume() error {
	_, err := zp.SendDeviceProperties("SetUseAutoplayVolume", "<UseVolume>boolean</UseVolume><Source>string</Source>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) SetZoneName() error {
	_, err := zp.SendDeviceProperties("SetZoneAttributes", "<DesiredZoneName>string</DesiredZoneName>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) SetIcon() error {
	_, err := zp.SendDeviceProperties("SetZoneAttributes", "<DesiredIcon>string</DesiredIcon>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) SetConfiguration() error {
	_, err := zp.SendDeviceProperties("SetZoneAttributes", "<DesiredConfiguration>string</DesiredConfiguration>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) SetTargetRoomName() error {
	_, err := zp.SendDeviceProperties("SetZoneAttributes", "<DesiredTargetRoomName>string</DesiredTargetRoomName>", "")
	return err
}
