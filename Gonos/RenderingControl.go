package Gonos

import "strconv"

// Get current volume. (TODO: Test)
func (zp *ZonePlayer) GetVolume() (int, error) {
	res, err := zp.sendCommand("RenderingControl", "GetVolume", "", "CurrentVolume")
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(res)
}

// Set volume. (TODO: Test)
func (zp *ZonePlayer) SetVolume(level int) error {
	_, err := zp.sendCommand("RenderingControl", "SetVolume", "<DesiredVolume>"+strconv.Itoa(max(0, min(100, level)))+"</DesiredVolume>", "")
	return err
}

// Get current mute state. (TODO: Test)
func (zp *ZonePlayer) GetMute() (bool, error) {
	res, err := zp.sendCommand("RenderingControl", "GetMute", "", "CurrentMute")
	return res == "1", err
}

// Set mute state. (TODO: Test)
func (zp *ZonePlayer) SetMute(state bool) error {
	_, err := zp.sendCommand("RenderingControl", "SetMute", "<DesiredMute>"+boolTo10(state)+"</DesiredMute>", "")
	return err
}

// Get current bass. (TODO: Test)
func (zp *ZonePlayer) GetBass() (int, error) {
	res, err := zp.sendCommand("RenderingControl", "GetBass", "", "CurrentBass")
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(res)
}

// Set bass. (TODO: Test)
func (zp *ZonePlayer) SetBass(level int) error {
	_, err := zp.sendCommand("RenderingControl", "SetBass", "<DesiredBass>"+strconv.Itoa(max(-10, min(10, level)))+"</DesiredBass>", "")
	return err
}

// Get current treble. (TODO: Test)
func (zp *ZonePlayer) GetTreble() (int, error) {
	res, err := zp.sendCommand("RenderingControl", "GetTreble", "", "CurrentTreble")
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(res)
}

// Set treble. (TODO: Test)
func (zp *ZonePlayer) SetTreble(level int) error {
	_, err := zp.sendCommand("RenderingControl", "SetTreble", "<DesiredTreble>"+strconv.Itoa(max(-10, min(10, level)))+"</DesiredTreble>", "")
	return err
}

// Get current loudness state. (TODO: Test)
func (zp *ZonePlayer) GetLoudness() (bool, error) {
	res, err := zp.sendCommand("RenderingControl", "GetLoudness", "", "CurrentLoudness")
	return res == "1", err
}

// Set loudness state. (TODO: Test)
func (zp *ZonePlayer) SetLoudness(state bool) error {
	_, err := zp.sendCommand("RenderingControl", "SetLoudness", "<DesiredLoudness>"+boolTo10(state)+"</DesiredLoudness>", "")
	return err
}
