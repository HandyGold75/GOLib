package Gonos

// https://sonos.svrooij.io/services/device-properties

import (
	"encoding/xml"
	"strconv"
)

type (
	resetBasicEQResponse struct {
		XMLName     xml.Name `xml:"ResetBasicEQResponse"`
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
func (zp *ZonePlayer) GetEQ(eQType string) (string, error) {
	return zp.SendRenderingControl("GetEQ", "<EQType>"+eQType+"</EQType>", "CurrentValue")
}

// TODO: Test
func (zp *ZonePlayer) GetHeadphoneConnected() (bool, error) {
	res, err := zp.SendRenderingControl("GetHeadphoneConnected", "", "CurrentHeadphoneConnected")
	return res == "1", err
}

// TODO: Test
func (zp *ZonePlayer) GetLoudness() (bool, error) {
	res, err := zp.SendRenderingControl("GetLoudness", "<Channel>"+zp.Static.RenderingControl.Channel+"</Channel>", "CurrentLoudness")
	return res == "1", err
}

// TODO: Test
func (zp *ZonePlayer) GetMute() (bool, error) {
	res, err := zp.SendRenderingControl("GetMute", "<Channel>"+zp.Static.RenderingControl.Channel+"</Channel>", "CurrentMute")
	return res == "1", err
}

// TODO: Test
func (zp *ZonePlayer) GetOutputFixed() (bool, error) {
	res, err := zp.SendRenderingControl("GetOutputFixed", "", "CurrentFixed")
	return res == "1", err
}

// TODO: Test
func (zp *ZonePlayer) GetRoomCalibrationStatus() (bool, bool, error) {
	res, err := zp.SendRenderingControl("GetRoomCalibrationStatus", "", "s:Body")
	if err != nil {
		return false, false, err
	}
	enabled, err := extractTag(res, "RoomCalibrationEnabled")
	if err != nil {
		return false, false, err
	}
	available, err := extractTag(res, "RoomCalibrationAvailable")
	if err != nil {
		return false, false, err
	}
	return enabled == "1", available == "1", err
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
	res, err := zp.SendRenderingControl("GetVolume", "<Channel>"+zp.Static.RenderingControl.Channel+"</Channel>", "CurrentVolume")
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(res)
}

// TODO: Test
func (zp *ZonePlayer) GetVolumeDB() (int, error) {
	res, err := zp.SendRenderingControl("GetVolumeDB", "<Channel>"+zp.Static.RenderingControl.Channel+"</Channel>", "CurrentVolume")
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(res)
}

// TODO: Test
func (zp *ZonePlayer) GetVolumeDBRange() (int, int, error) {
	res, err := zp.SendRenderingControl("GetVolumeDBRange", "<Channel>"+zp.Static.RenderingControl.Channel+"</Channel>", "s:Body")
	if err != nil {
		return 0, 0, err
	}
	minValue, err := extractTag(res, "MinValue")
	if err != nil {
		return 0, 0, err
	}
	maxValue, err := extractTag(res, "MaxValue")
	if err != nil {
		return 0, 0, err
	}
	minValueInt, err := strconv.Atoi(minValue)
	if err != nil {
		return 0, 0, err
	}
	maxValueInt, err := strconv.Atoi(maxValue)
	if err != nil {
		return 0, 0, err
	}
	return minValueInt, maxValueInt, err
}

// TODO: Test
func (zp *ZonePlayer) RampToVolume(rampType string, volume int, resetVolumeAfter bool, programURI string) (int, error) {
	res, err := zp.SendRenderingControl("RampToVolume", "<Channel>"+zp.Static.RenderingControl.Channel+"</Channel><RampType>"+rampType+"</RampType><DesiredVolume>"+strconv.Itoa(max(0, min(100, volume)))+"</DesiredVolume><ResetVolumeAfter>"+boolTo10(resetVolumeAfter)+"</ResetVolumeAfter><ProgramURI>"+programURI+"</ProgramURI>", "RampTime")
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(res)
}

// TODO: Test
func (zp *ZonePlayer) ResetBasicEQ() (resetBasicEQResponse, error) {
	res, err := zp.SendRenderingControl("ResetBasicEQ", "", "s:Body")
	if err != nil {
		return resetBasicEQResponse{}, err
	}
	data := resetBasicEQResponse{}
	err = xml.Unmarshal([]byte(res), &data)
	return data, err
}

// TODO: Test
func (zp *ZonePlayer) ResetExtEQ(eQType string) error {
	_, err := zp.SendRenderingControl("ResetExtEQ", "<EQType>"+eQType+"</EQType>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) RestoreVolumePriorToRamp() error {
	_, err := zp.SendRenderingControl("RestoreVolumePriorToRamp", "<Channel>"+zp.Static.RenderingControl.Channel+"</Channel>", "")
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
func (zp *ZonePlayer) SetEQ(eQType string, state string) error {
	_, err := zp.SendRenderingControl("SetEQ", "<EQType>"+eQType+"</EQType><DesiredValue>"+state+"</DesiredValue>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) SetLoudness(state bool) error {
	_, err := zp.SendRenderingControl("SetLoudness", "<Channel>"+zp.Static.RenderingControl.Channel+"</Channel><DesiredLoudness>"+boolTo10(state)+"</DesiredLoudness>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) SetMute(state bool) error {
	_, err := zp.SendRenderingControl("SetMute", "<Channel>"+zp.Static.RenderingControl.Channel+"</Channel><DesiredMute>"+boolTo10(state)+"</DesiredMute>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) SetOutputFixed(state bool) error {
	_, err := zp.SendRenderingControl("SetOutputFixed", "<DesiredFixed>"+boolTo10(state)+"</DesiredFixed>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) SetRelativeVolume(volume int) (int, error) {
	res, err := zp.SendRenderingControl("SetRelativeVolume", "<Channel>"+zp.Static.RenderingControl.Channel+"</Channel><Adjustment>"+strconv.Itoa(max(-100, min(100, volume)))+"</Adjustment>", "NewVolume")
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
	_, err := zp.SendRenderingControl("SetVolume", "<Channel>"+zp.Static.RenderingControl.Channel+"</Channel><DesiredVolume>"+strconv.Itoa(max(0, min(100, volume)))+"</DesiredVolume>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) SetVolumeDB(volume int) error {
	_, err := zp.SendRenderingControl("SetVolumeDB", "<Channel>"+zp.Static.RenderingControl.Channel+"</Channel><DesiredVolume>"+strconv.Itoa(max(0, min(100, volume)))+"</DesiredVolume>", "")
	return err
}
