package RenderingControl

// https://sonos.svrooij.io/services/device-properties

import (
	"Gonos/lib"
	"encoding/xml"
	"strconv"
)

type (
	RenderingControl struct {
		Send func(action, body, targetTag string) (string, error)
		// Channel should be one of: `Master`, `LF` or `RF`
		Channel string
	}

	resetBasicEQResponse struct {
		XMLName     xml.Name `xml:"ResetBasicEQResponse"`
		Bass        int
		Treble      int
		Loudness    bool
		LeftVolume  int
		RightVolume int
	}
)

func New(send func(action, body, targetTag string) (string, error)) RenderingControl {
	return RenderingControl{Send: send, Channel: "Master"}
}

func (zp *RenderingControl) GetBass() (int, error) {
	res, err := zp.Send("GetBass", "", "CurrentBass")
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(res)
}

func (zp *RenderingControl) GetEQ(eQType string) (string, error) {
	return zp.Send("GetEQ", "<EQType>"+eQType+"</EQType>", "CurrentValue")
}

func (zp *RenderingControl) GetHeadphoneConnected() (bool, error) {
	res, err := zp.Send("GetHeadphoneConnected", "", "CurrentHeadphoneConnected")
	return res == "1", err
}

func (zp *RenderingControl) GetLoudness() (bool, error) {
	res, err := zp.Send("GetLoudness", "<Channel>"+zp.Channel+"</Channel>", "CurrentLoudness")
	return res == "1", err
}

func (zp *RenderingControl) GetMute() (bool, error) {
	res, err := zp.Send("GetMute", "<Channel>"+zp.Channel+"</Channel>", "CurrentMute")
	return res == "1", err
}

func (zp *RenderingControl) GetOutputFixed() (bool, error) {
	res, err := zp.Send("GetOutputFixed", "", "CurrentFixed")
	return res == "1", err
}

func (zp *RenderingControl) GetRoomCalibrationStatus() (bool, bool, error) {
	res, err := zp.Send("GetRoomCalibrationStatus", "", "s:Body")
	if err != nil {
		return false, false, err
	}
	enabled, err := lib.ExtractTag(res, "RoomCalibrationEnabled")
	if err != nil {
		return false, false, err
	}
	available, err := lib.ExtractTag(res, "RoomCalibrationAvailable")
	if err != nil {
		return false, false, err
	}
	return enabled == "1", available == "1", err
}

func (zp *RenderingControl) GetSupportsOutputFixed() (bool, error) {
	res, err := zp.Send("GetSupportsOutputFixed", "", "CurrentSupportsFixed")
	return res == "1", err
}

func (zp *RenderingControl) GetTreble() (int, error) {
	res, err := zp.Send("GetTreble", "", "CurrentTreble")
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(res)
}

func (zp *RenderingControl) GetVolume() (int, error) {
	res, err := zp.Send("GetVolume", "<Channel>"+zp.Channel+"</Channel>", "CurrentVolume")
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(res)
}

func (zp *RenderingControl) GetVolumeDB() (int, error) {
	res, err := zp.Send("GetVolumeDB", "<Channel>"+zp.Channel+"</Channel>", "CurrentVolume")
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(res)
}

func (zp *RenderingControl) GetVolumeDBRange() (int, int, error) {
	res, err := zp.Send("GetVolumeDBRange", "<Channel>"+zp.Channel+"</Channel>", "s:Body")
	if err != nil {
		return 0, 0, err
	}
	minValue, err := lib.ExtractTag(res, "MinValue")
	if err != nil {
		return 0, 0, err
	}
	maxValue, err := lib.ExtractTag(res, "MaxValue")
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

func (zp *RenderingControl) RampToVolume(rampType string, volume int, resetVolumeAfter bool, programURI string) (int, error) {
	res, err := zp.Send("RampToVolume", "<Channel>"+zp.Channel+"</Channel><RampType>"+rampType+"</RampType><DesiredVolume>"+strconv.Itoa(max(0, min(100, volume)))+"</DesiredVolume><ResetVolumeAfter>"+lib.BoolTo10(resetVolumeAfter)+"</ResetVolumeAfter><ProgramURI>"+programURI+"</ProgramURI>", "RampTime")
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(res)
}

func (zp *RenderingControl) ResetBasicEQ() (resetBasicEQResponse, error) {
	res, err := zp.Send("ResetBasicEQ", "", "s:Body")
	if err != nil {
		return resetBasicEQResponse{}, err
	}
	data := resetBasicEQResponse{}
	err = xml.Unmarshal([]byte(res), &data)
	return data, err
}

func (zp *RenderingControl) ResetExtEQ(eQType string) error {
	_, err := zp.Send("ResetExtEQ", "<EQType>"+eQType+"</EQType>", "")
	return err
}

func (zp *RenderingControl) RestoreVolumePriorToRamp() error {
	_, err := zp.Send("RestoreVolumePriorToRamp", "<Channel>"+zp.Channel+"</Channel>", "")
	return err
}

func (zp *RenderingControl) SetBass(volume int) error {
	_, err := zp.Send("SetBass", "<DesiredBass>"+strconv.Itoa(max(-10, min(10, volume)))+"</DesiredBass>", "")
	return err
}

func (zp *RenderingControl) SetChannelMap(channelMap string) error {
	_, err := zp.Send("SetChannelMap", "<ChannelMap>"+channelMap+"</ChannelMap>", "")
	return err
}

func (zp *RenderingControl) SetEQ(eQType string, state string) error {
	_, err := zp.Send("SetEQ", "<EQType>"+eQType+"</EQType><DesiredValue>"+state+"</DesiredValue>", "")
	return err
}

func (zp *RenderingControl) SetLoudness(state bool) error {
	_, err := zp.Send("SetLoudness", "<Channel>"+zp.Channel+"</Channel><DesiredLoudness>"+lib.BoolTo10(state)+"</DesiredLoudness>", "")
	return err
}

func (zp *RenderingControl) SetMute(state bool) error {
	_, err := zp.Send("SetMute", "<Channel>"+zp.Channel+"</Channel><DesiredMute>"+lib.BoolTo10(state)+"</DesiredMute>", "")
	return err
}

func (zp *RenderingControl) SetOutputFixed(state bool) error {
	_, err := zp.Send("SetOutputFixed", "<DesiredFixed>"+lib.BoolTo10(state)+"</DesiredFixed>", "")
	return err
}

func (zp *RenderingControl) SetRelativeVolume(volume int) (int, error) {
	res, err := zp.Send("SetRelativeVolume", "<Channel>"+zp.Channel+"</Channel><Adjustment>"+strconv.Itoa(max(-100, min(100, volume)))+"</Adjustment>", "NewVolume")
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(res)
}

func (zp *RenderingControl) SetRoomCalibrationStatus(state bool) error {
	_, err := zp.Send("SetRoomCalibrationStatus", "<RoomCalibrationEnabled>"+lib.BoolTo10(state)+"</RoomCalibrationEnabled>", "")
	return err
}

func (zp *RenderingControl) SetRoomCalibrationX(calibrationID string, coeddicients string, calibrationMode string) error {
	_, err := zp.Send("SetRoomCalibrationX", "<CalibrationID>"+calibrationID+"</CalibrationID><Coefficients>"+coeddicients+"</Coefficients><CalibrationMode>"+calibrationMode+"</CalibrationMode>", "")
	return err
}

func (zp *RenderingControl) SetTreble(volume int) error {
	_, err := zp.Send("SetTreble", "<DesiredTreble>"+strconv.Itoa(max(-10, min(10, volume)))+"</DesiredTreble>", "")
	return err
}

func (zp *RenderingControl) SetVolume(volume int) error {
	_, err := zp.Send("SetVolume", "<Channel>"+zp.Channel+"</Channel><DesiredVolume>"+strconv.Itoa(max(0, min(100, volume)))+"</DesiredVolume>", "")
	return err
}

func (zp *RenderingControl) SetVolumeDB(volume int) error {
	_, err := zp.Send("SetVolumeDB", "<Channel>"+zp.Channel+"</Channel><DesiredVolume>"+strconv.Itoa(max(0, min(100, volume)))+"</DesiredVolume>", "")
	return err
}
