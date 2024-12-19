package Gonos

import "strconv"

// Get current volume. (TODO: Test)
func (zp *ZonePlayer) GetVolume() (int, error) {
	res, err := zp.SendRenderingControl("GetVolume", "<InstanceID>0</InstanceID>", "CurrentVolume")
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(res)
}

// Set volume. (TODO: Test)
func (zp *ZonePlayer) SetVolume(level int) error {
	_, err := zp.SendRenderingControl("SetVolume", "<InstanceID>0</InstanceID><DesiredVolume>"+strconv.Itoa(max(0, min(100, level)))+"</DesiredVolume>", "")
	return err
}

// Get current mute state. (TODO: Test)
func (zp *ZonePlayer) GetMute() (bool, error) {
	res, err := zp.SendRenderingControl("GetMute", "<InstanceID>0</InstanceID>", "CurrentMute")
	return res == "1", err
}

// Set mute state. (TODO: Test)
func (zp *ZonePlayer) SetMute(state bool) error {
	_, err := zp.SendRenderingControl("SetMute", "<InstanceID>0</InstanceID><DesiredMute>"+boolTo10(state)+"</DesiredMute>", "")
	return err
}

// Get current bass. (TODO: Test)
func (zp *ZonePlayer) GetBass() (int, error) {
	res, err := zp.SendRenderingControl("GetBass", "<InstanceID>0</InstanceID>", "CurrentBass")
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(res)
}

// Set bass. (TODO: Test)
func (zp *ZonePlayer) SetBass(level int) error {
	_, err := zp.SendRenderingControl("SetBass", "<InstanceID>0</InstanceID><DesiredBass>"+strconv.Itoa(max(-10, min(10, level)))+"</DesiredBass>", "")
	return err
}

// Get current treble. (TODO: Test)
func (zp *ZonePlayer) GetTreble() (int, error) {
	res, err := zp.SendRenderingControl("GetTreble", "<InstanceID>0</InstanceID>", "CurrentTreble")
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(res)
}

// Set treble. (TODO: Test)
func (zp *ZonePlayer) SetTreble(level int) error {
	_, err := zp.SendRenderingControl("SetTreble", "<InstanceID>0</InstanceID><DesiredTreble>"+strconv.Itoa(max(-10, min(10, level)))+"</DesiredTreble>", "")
	return err
}

// Get current loudness state. (TODO: Test)
func (zp *ZonePlayer) GetLoudness() (bool, error) {
	res, err := zp.SendRenderingControl("GetLoudness", "<InstanceID>0</InstanceID>", "CurrentLoudness")
	return res == "1", err
}

// Set loudness state. (TODO: Test)
func (zp *ZonePlayer) SetLoudness(state bool) error {
	_, err := zp.SendRenderingControl("SetLoudness", "<InstanceID>0</InstanceID><DesiredLoudness>"+boolTo10(state)+"</DesiredLoudness>", "")
	return err
}
