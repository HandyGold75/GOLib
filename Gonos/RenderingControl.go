package Gonos

// https://sonos.svrooij.io/services/device-properties

import (
	"encoding/xml"
	"strconv"
)

type (
	RoomCalibration struct {
		XMLName                  xml.Name `xml:"item"`
		RoomCalibrationEnabled   bool
		RoomCalibrationAvailable bool
	}

	VolumeDBRange struct {
		XMLName  xml.Name `xml:"item"`
		MinValue int
		MaxValue int
	}

	BasicEQ struct {
		XMLName     xml.Name `xml:"item"`
		Bass        int
		Treble      int
		Loudness    bool
		LeftVolume  int
		RightVolume int
	}
)

// TODO: Test
func (zp *ZonePlayer) GetBass() (int, error) {
	res, err := zp.SendRenderingControl("GetBass", "", "CurrentBass")
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(res)
}

// TODO: Test
func (zp *ZonePlayer) GetEQDialogLevel() (bool, error) {
	res, err := zp.SendRenderingControl("GetEQ", "<EQType>EQDialogLevel</EQType>", "CurrentValue")
	return res == "1", err
}

// TODO: Test
func (zp *ZonePlayer) GetEQMusicSurroundLevel() (int, error) {
	res, err := zp.SendRenderingControl("GetEQ", "<EQType>EQMusicSurroundLevel</EQType>", "CurrentValue")
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(res)
}

// TODO: Test
func (zp *ZonePlayer) GetEQNightMode() (bool, error) {
	res, err := zp.SendRenderingControl("GetEQ", "<EQType>EQNightMode</EQType>", "CurrentValue")
	return res == "1", err
}

// TODO: Test
func (zp *ZonePlayer) GetEQSubGain() (int, error) {
	res, err := zp.SendRenderingControl("GetEQ", "<EQType>EQSubGain</EQType>", "CurrentValue")
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(res)
}

// TODO: Test
func (zp *ZonePlayer) GetEQSurroundEnable() (bool, error) {
	res, err := zp.SendRenderingControl("GetEQ", "<EQType>EQSurroundEnable</EQType>", "CurrentValue")
	return res == "1", err
}

// TODO: Test
func (zp *ZonePlayer) GetEQSurroundLevel() (int, error) {
	res, err := zp.SendRenderingControl("GetEQ", "<EQType>EQSurroundLevel</EQType>", "CurrentValue")
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(res)
}

// TODO: Test
func (zp *ZonePlayer) GetEQSurroundMode() (bool, error) {
	res, err := zp.SendRenderingControl("GetEQ", "<EQType>EQSurroundMode</EQType>", "CurrentValue")
	return res == "1", err
}

// TODO: Test
func (zp *ZonePlayer) GetEQHeightChannelLevel() (int, error) {
	res, err := zp.SendRenderingControl("GetEQ", "<EQType>EQHeightChannelLevel</EQType>", "CurrentValue")
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(res)
}

// TODO: Test
func (zp *ZonePlayer) GetHeadphoneConnected() (bool, error) {
	res, err := zp.SendRenderingControl("GetHeadphoneConnected", "", "CurrentHeadphoneConnected")
	return res == "1", err
}

// TODO: Test
func (zp *ZonePlayer) GetLoudness() (bool, error) {
	res, err := zp.SendRenderingControl("GetLoudness", "<Channel>"+zp.Channel+"</Channel>", "CurrentLoudness")
	return res == "1", err
}

// TODO: Test
func (zp *ZonePlayer) GetMute() (bool, error) {
	res, err := zp.SendRenderingControl("GetMute", "<Channel>"+zp.Channel+"</Channel>", "CurrentMute")
	return res == "1", err
}

// TODO: Test
func (zp *ZonePlayer) GetOutputFixed() (bool, error) {
	res, err := zp.SendRenderingControl("GetOutputFixed", "", "CurrentFixed")
	return res == "1", err
}

// TODO: Test
func (zp *ZonePlayer) GetRoomCalibrationStatus() (RoomCalibration, error) {
	res, err := zp.SendRenderingControl("GetRoomCalibrationStatus", "", "s:Body")
	if err != nil {
		return RoomCalibration{}, err
	}
	data := RoomCalibration{}
	err = xml.Unmarshal([]byte(res), &data)
	return data, err
}

// TODO: Test
func (zp *ZonePlayer) GetSupportsOutputFixed() (bool, error) {
	res, err := zp.SendRenderingControl("GetSupportsOutputFixed", "", "CurrentSupportsFixed")
	return res == "1", err
}

// TODO: Test
func (zp *ZonePlayer) GetTreble() (int, error) {
	res, err := zp.SendRenderingControl("GetTreble", "", "CurrentTreble")
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(res)
}

// TODO: Test
func (zp *ZonePlayer) GetVolume() (int, error) {
	res, err := zp.SendRenderingControl("GetVolume", "<Channel>"+zp.Channel+"</Channel>", "CurrentVolume")
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(res)
}

// TODO: Test
func (zp *ZonePlayer) GetVolumeDB() (int, error) {
	res, err := zp.SendRenderingControl("GetVolumeDB", "<Channel>"+zp.Channel+"</Channel>", "CurrentVolume")
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(res)
}

// TODO: Test
func (zp *ZonePlayer) GetVolumeDBRange() (VolumeDBRange, error) {
	res, err := zp.SendRenderingControl("GetVolumeDBRange", "<Channel>"+zp.Channel+"</Channel>", "s:Body")
	if err != nil {
		return VolumeDBRange{}, err
	}
	data := VolumeDBRange{}
	err = xml.Unmarshal([]byte(res), &data)
	return data, err
}

// TODO: Test
func (zp *ZonePlayer) RampToVolumeSleepTimer(volume int, resetVolumeAfter bool, ProgramURI string) (int, error) {
	res, err := zp.SendRenderingControl("RampToVolume", "<Channel>"+zp.Channel+"</Channel><RampType>SleepTimer</RampType><DesiredVolume>"+strconv.Itoa(max(0, min(100, volume)))+"</DesiredVolume><ResetVolumeAfter>"+boolTo10(resetVolumeAfter)+"</ResetVolumeAfter><ProgramURI>"+ProgramURI+"</ProgramURI>", "RampTime")
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(res)
}

// TODO: Test
func (zp *ZonePlayer) RampToVolumeAlarm(volume int, resetVolumeAfter bool, ProgramURI string) (int, error) {
	res, err := zp.SendRenderingControl("RampToVolume", "<Channel>"+zp.Channel+"</Channel><RampType>Alarm</RampType><DesiredVolume>"+strconv.Itoa(max(0, min(100, volume)))+"</DesiredVolume><ResetVolumeAfter>"+boolTo10(resetVolumeAfter)+"</ResetVolumeAfter><ProgramURI>"+ProgramURI+"</ProgramURI>", "RampTime")
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(res)
}

// TODO: Test
func (zp *ZonePlayer) RampToVolumeAutoPlay(volume int, resetVolumeAfter bool, ProgramURI string) (int, error) {
	res, err := zp.SendRenderingControl("RampToVolume", "<Channel>"+zp.Channel+"</Channel><RampType>AutoPlay</RampType><DesiredVolume>"+strconv.Itoa(max(0, min(100, volume)))+"</DesiredVolume><ResetVolumeAfter>"+boolTo10(resetVolumeAfter)+"</ResetVolumeAfter><ProgramURI>"+ProgramURI+"</ProgramURI>", "RampTime")
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(res)
}

// TODO: Test
func (zp *ZonePlayer) ResetBasicEQ() (BasicEQ, error) {
	res, err := zp.SendRenderingControl("ResetBasicEQ", "", "s:Body")
	if err != nil {
		return BasicEQ{}, err
	}
	data := BasicEQ{}
	err = xml.Unmarshal([]byte(res), &data)
	return data, err
}

// TODO: Test
func (zp *ZonePlayer) ResetExtEQDialogLevel() error {
	_, err := zp.SendRenderingControl("ResetExtEQ", "<EQType>EQDialogLevel</EQType>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) ResetExtEQMusicSurroundLevel() error {
	_, err := zp.SendRenderingControl("ResetExtEQ", "<EQType>EQMusicSurroundLevel</EQType>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) ResetExtEQNightMode() error {
	_, err := zp.SendRenderingControl("ResetExtEQ", "<EQType>EQNightMode</EQType>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) ResetExtEQSubGain() error {
	_, err := zp.SendRenderingControl("ResetExtEQ", "<EQType>EQSubGain</EQType>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) ResetExtEQSurroundEnable() error {
	_, err := zp.SendRenderingControl("ResetExtEQ", "<EQType>EQSurroundEnable</EQType>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) ResetExtEQSurroundLevel() error {
	_, err := zp.SendRenderingControl("ResetExtEQ", "<EQType>EQSurroundLevel</EQType>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) ResetExtEQSurroundMode() error {
	_, err := zp.SendRenderingControl("ResetExtEQ", "<EQType>EQSurroundMode</EQType>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) ResetExtEQHeightChannelLevel() error {
	_, err := zp.SendRenderingControl("ResetExtEQ", "<EQType>EQHeightChannelLevel</EQType>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) RestoreVolumePriorToRamp() error {
	_, err := zp.SendRenderingControl("RestoreVolumePriorToRamp", "<Channel>"+zp.Channel+"</Channel>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) SetBass(volume int) error {
	_, err := zp.SendRenderingControl("SetBass", "<DesiredBass>"+strconv.Itoa(max(-10, min(10, volume)))+"</DesiredBass>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) SetChannelMap(channelMap string) error {
	_, err := zp.SendRenderingControl("SetChannelMap", "<ChannelMap>"+channelMap+"</ChannelMap>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) SetEQDialogLevel(state bool) error {
	_, err := zp.SendRenderingControl("SetEQ", "<EQType>DialogLevel</EQType><DesiredValue>"+boolTo10(state)+"</DesiredValue>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) SetEQMusicSurroundLevel(volume int) error {
	_, err := zp.SendRenderingControl("SetEQ", "<EQType>MusicSurroundLevel</EQType><DesiredValue>"+strconv.Itoa(max(-15, min(15, volume)))+"</DesiredValue>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) SetEQNightMode(state bool) error {
	_, err := zp.SendRenderingControl("SetEQ", "<EQType>NightMode</EQType><DesiredValue>"+boolTo10(state)+"</DesiredValue>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) SetEQSubGain(volume int) error {
	_, err := zp.SendRenderingControl("SetEQ", "<EQType>SubGain</EQType><DesiredValue>"+strconv.Itoa(max(-10, min(10, volume)))+"</DesiredValue>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) SetEQSurroundEnable(state bool) error {
	_, err := zp.SendRenderingControl("SetEQ", "<EQType>SurroundEnable</EQType><DesiredValue>"+boolTo10(state)+"</DesiredValue>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) SetEQSurroundLevel(volume int) error {
	_, err := zp.SendRenderingControl("SetEQ", "<EQType>SurroundLevel</EQType><DesiredValue>"+strconv.Itoa(max(-15, min(15, volume)))+"</DesiredValue>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) SetEQSurroundMode(full bool) error {
	_, err := zp.SendRenderingControl("SetEQ", "<EQType>SurroundMode</EQType><DesiredValue>"+boolTo10(full)+"</DesiredValue>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) SetEQHeightChannelLevel(volume int) error {
	_, err := zp.SendRenderingControl("SetEQ", "<EQType>HeightChannelLevel</EQType><DesiredValue>"+strconv.Itoa(max(-10, min(10, volume)))+"</DesiredValue>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) SetLoudness(state bool) error {
	_, err := zp.SendRenderingControl("SetLoudness", "<Channel>"+zp.Channel+"</Channel><DesiredLoudness>"+boolTo10(state)+"</DesiredLoudness>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) SetMute(state bool) error {
	_, err := zp.SendRenderingControl("SetMute", "<Channel>"+zp.Channel+"</Channel><DesiredMute>"+boolTo10(state)+"</DesiredMute>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) SetOutputFixed(state bool) error {
	_, err := zp.SendRenderingControl("SetOutputFixed", "<DesiredFixed>"+boolTo10(state)+"</DesiredFixed>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) SetRelativeVolume(volume int) (int, error) {
	res, err := zp.SendRenderingControl("SetRelativeVolume", "<Channel>"+zp.Channel+"</Channel><Adjustment>"+strconv.Itoa(max(-100, min(100, volume)))+"</Adjustment>", "NewVolume")
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(res)
}

// TODO: Test
func (zp *ZonePlayer) SetRoomCalibrationStatus(state bool) error {
	_, err := zp.SendRenderingControl("SetRoomCalibrationStatus", "<RoomCalibrationEnabled>"+boolTo10(state)+"</RoomCalibrationEnabled>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) SetRoomCalibrationX(calibrationID string, coeddicients string, calibrationMode string) error {
	_, err := zp.SendRenderingControl("SetRoomCalibrationX", "<CalibrationID>"+calibrationID+"</CalibrationID><Coefficients>"+coeddicients+"</Coefficients><CalibrationMode>"+calibrationMode+"</CalibrationMode>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) SetTreble(volume int) error {
	_, err := zp.SendRenderingControl("SetTreble", "<DesiredTreble>"+strconv.Itoa(max(-10, min(10, volume)))+"</DesiredTreble>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) SetVolume(volume int) error {
	_, err := zp.SendRenderingControl("SetVolume", "<Channel>"+zp.Channel+"</Channel><DesiredVolume>"+strconv.Itoa(max(0, min(100, volume)))+"</DesiredVolume>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) SetVolumeDB(volume int) error {
	_, err := zp.SendRenderingControl("SetVolumeDB", "<Channel>"+zp.Channel+"</Channel><DesiredVolume>"+strconv.Itoa(max(0, min(100, volume)))+"</DesiredVolume>", "")
	return err
}
