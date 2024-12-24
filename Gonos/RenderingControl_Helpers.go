package Gonos

import "strconv"

// TODO: Test
func (zp *ZonePlayer) GetEQDialogLevel() (bool, error) {
	res, err := zp.GetEQ("EQDialogLevel")
	return res == "1", err
}

// TODO: Test
func (zp *ZonePlayer) GetEQMusicSurroundLevel() (int, error) {
	res, err := zp.GetEQ("EQMusicSurroundLevel")
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(res)
}

// TODO: Test
func (zp *ZonePlayer) GetEQNightMode() (bool, error) {
	res, err := zp.GetEQ("EQNightMode")
	return res == "1", err
}

// TODO: Test
func (zp *ZonePlayer) GetEQSubGain() (int, error) {
	res, err := zp.GetEQ("EQSubGain")
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(res)
}

// TODO: Test
func (zp *ZonePlayer) GetEQSurroundEnable() (bool, error) {
	res, err := zp.GetEQ("EQSurroundEnable")
	return res == "1", err
}

// TODO: Test
func (zp *ZonePlayer) GetEQSurroundLevel() (int, error) {
	res, err := zp.GetEQ("EQSurroundLevel")
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(res)
}

// TODO: Test
func (zp *ZonePlayer) GetEQSurroundMode() (bool, error) {
	res, err := zp.GetEQ("EQSurroundMode")
	return res == "1", err
}

// TODO: Test
func (zp *ZonePlayer) GetEQHeightChannelLevel() (int, error) {
	res, err := zp.GetEQ("EQHeightChannelLevel")
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(res)
}

// TODO: Test
func (zp *ZonePlayer) RampToVolumeSleepTimer(volume int, resetVolumeAfter bool, programURI string) (int, error) {
	return zp.RampToVolume("SleepTimer", volume, resetVolumeAfter, programURI)
}

// TODO: Test
func (zp *ZonePlayer) RampToVolumeAlarm(volume int, resetVolumeAfter bool, programURI string) (int, error) {
	return zp.RampToVolume("Alarm", volume, resetVolumeAfter, programURI)
}

// TODO: Test
func (zp *ZonePlayer) RampToVolumeAutoPlay(volume int, resetVolumeAfter bool, programURI string) (int, error) {
	return zp.RampToVolume("AutoPlay", volume, resetVolumeAfter, programURI)
}

// TODO: Test
func (zp *ZonePlayer) ResetExtEQDialogLevel() error {
	return zp.ResetExtEQ("EQDialogLevel")
}

// TODO: Test
func (zp *ZonePlayer) ResetExtEQMusicSurroundLevel() error {
	return zp.ResetExtEQ("EQMusicSurroundLevel")
}

// TODO: Test
func (zp *ZonePlayer) ResetExtEQNightMode() error {
	return zp.ResetExtEQ("EQNightMode")
}

// TODO: Test
func (zp *ZonePlayer) ResetExtEQSubGain() error {
	return zp.ResetExtEQ("EQSubGain")
}

// TODO: Test
func (zp *ZonePlayer) ResetExtEQSurroundEnable() error {
	return zp.ResetExtEQ("EQSurroundEnable")
}

// TODO: Test
func (zp *ZonePlayer) ResetExtEQSurroundLevel() error {
	return zp.ResetExtEQ("EQSurroundLevel")
}

// TODO: Test
func (zp *ZonePlayer) ResetExtEQSurroundMode() error {
	return zp.ResetExtEQ("EQSurroundMode")
}

// TODO: Test
func (zp *ZonePlayer) ResetExtEQHeightChannelLevel() error {
	return zp.ResetExtEQ("EQHeightChannelLevel")
}

// TODO: Test
func (zp *ZonePlayer) SetEQDialogLevel(state bool) error {
	return zp.SetEQ("DialogLevel", boolTo10(state))
}

// TODO: Test
func (zp *ZonePlayer) SetEQMusicSurroundLevel(volume int) error {
	return zp.SetEQ("MusicSurroundLevel", strconv.Itoa(max(-15, min(15, volume))))
}

// TODO: Test
func (zp *ZonePlayer) SetEQNightMode(state bool) error {
	return zp.SetEQ("NightMode", boolTo10(state))
}

// TODO: Test
func (zp *ZonePlayer) SetEQSubGain(volume int) error {
	return zp.SetEQ("SubGain", strconv.Itoa(max(-10, min(10, volume))))
}

// TODO: Test
func (zp *ZonePlayer) SetEQSurroundEnable(state bool) error {
	return zp.SetEQ("SurroundEnable", boolTo10(state))
}

// TODO: Test
func (zp *ZonePlayer) SetEQSurroundLevel(volume int) error {
	return zp.SetEQ("SurroundLevel", strconv.Itoa(max(-15, min(15, volume))))
}

// TODO: Test
func (zp *ZonePlayer) SetEQSurroundMode(full bool) error {
	return zp.SetEQ("SurroundMode", boolTo10(full))
}

// TODO: Test
func (zp *ZonePlayer) SetEQHeightChannelLevel(volume int) error {
	return zp.SetEQ("HeightChannelLevel", strconv.Itoa(max(-10, min(10, volume))))
}
