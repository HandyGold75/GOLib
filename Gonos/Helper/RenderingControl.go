package Helper

import (
	"Gonos/lib"
	"strconv"
)

func (h *Helper) GetEQDialogLevel() (bool, error) {
	res, err := h.renderingControl.GetEQ("EQDialogLevel")
	return res == "1", err
}

func (h *Helper) GetEQMusicSurroundLevel() (int, error) {
	res, err := h.renderingControl.GetEQ("EQMusicSurroundLevel")
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(res)
}

func (h *Helper) GetEQNightMode() (bool, error) {
	res, err := h.renderingControl.GetEQ("EQNightMode")
	return res == "1", err
}

func (h *Helper) GetEQSubGain() (int, error) {
	res, err := h.renderingControl.GetEQ("EQSubGain")
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(res)
}

func (h *Helper) GetEQSurroundEnable() (bool, error) {
	res, err := h.renderingControl.GetEQ("EQSurroundEnable")
	return res == "1", err
}

func (h *Helper) GetEQSurroundLevel() (int, error) {
	res, err := h.renderingControl.GetEQ("EQSurroundLevel")
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(res)
}

func (h *Helper) GetEQSurroundMode() (bool, error) {
	res, err := h.renderingControl.GetEQ("EQSurroundMode")
	return res == "1", err
}

func (h *Helper) GetEQHeightChannelLevel() (int, error) {
	res, err := h.renderingControl.GetEQ("EQHeightChannelLevel")
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(res)
}

func (h *Helper) RampToVolumeSleepTimer(volume int, resetVolumeAfter bool, programURI string) (int, error) {
	return h.renderingControl.RampToVolume("SleepTimer", volume, resetVolumeAfter, programURI)
}

func (h *Helper) RampToVolumeAlarm(volume int, resetVolumeAfter bool, programURI string) (int, error) {
	return h.renderingControl.RampToVolume("Alarm", volume, resetVolumeAfter, programURI)
}

func (h *Helper) RampToVolumeAutoPlay(volume int, resetVolumeAfter bool, programURI string) (int, error) {
	return h.renderingControl.RampToVolume("AutoPlay", volume, resetVolumeAfter, programURI)
}

func (h *Helper) ResetExtEQDialogLevel() error {
	return h.renderingControl.ResetExtEQ("EQDialogLevel")
}

func (h *Helper) ResetExtEQMusicSurroundLevel() error {
	return h.renderingControl.ResetExtEQ("EQMusicSurroundLevel")
}

func (h *Helper) ResetExtEQNightMode() error {
	return h.renderingControl.ResetExtEQ("EQNightMode")
}

func (h *Helper) ResetExtEQSubGain() error {
	return h.renderingControl.ResetExtEQ("EQSubGain")
}

func (h *Helper) ResetExtEQSurroundEnable() error {
	return h.renderingControl.ResetExtEQ("EQSurroundEnable")
}

func (h *Helper) ResetExtEQSurroundLevel() error {
	return h.renderingControl.ResetExtEQ("EQSurroundLevel")
}

func (h *Helper) ResetExtEQSurroundMode() error {
	return h.renderingControl.ResetExtEQ("EQSurroundMode")
}

func (h *Helper) ResetExtEQHeightChannelLevel() error {
	return h.renderingControl.ResetExtEQ("EQHeightChannelLevel")
}

func (h *Helper) SetEQDialogLevel(state bool) error {
	return h.renderingControl.SetEQ("DialogLevel", lib.BoolTo10(state))
}

func (h *Helper) SetEQMusicSurroundLevel(volume int) error {
	return h.renderingControl.SetEQ("MusicSurroundLevel", strconv.Itoa(max(-15, min(15, volume))))
}

func (h *Helper) SetEQNightMode(state bool) error {
	return h.renderingControl.SetEQ("NightMode", lib.BoolTo10(state))
}

func (h *Helper) SetEQSubGain(volume int) error {
	return h.renderingControl.SetEQ("SubGain", strconv.Itoa(max(-10, min(10, volume))))
}

func (h *Helper) SetEQSurroundEnable(state bool) error {
	return h.renderingControl.SetEQ("SurroundEnable", lib.BoolTo10(state))
}

func (h *Helper) SetEQSurroundLevel(volume int) error {
	return h.renderingControl.SetEQ("SurroundLevel", strconv.Itoa(max(-15, min(15, volume))))
}

func (h *Helper) SetEQSurroundMode(full bool) error {
	return h.renderingControl.SetEQ("SurroundMode", lib.BoolTo10(full))
}

func (h *Helper) SetEQHeightChannelLevel(volume int) error {
	return h.renderingControl.SetEQ("HeightChannelLevel", strconv.Itoa(max(-10, min(10, volume))))
}

// Short for `zp.RenderingControl.SetLoudness`.
func (h *Helper) SetLoudness(volume bool) error {
	return h.renderingControl.SetLoudness(volume)
}

// Short for `zp.RenderingControl.GetLoudness`.
func (h *Helper) GetLoudness() (bool, error) {
	return h.renderingControl.GetLoudness()
}

// Short for `zp.RenderingControl.SetMute`.
func (h *Helper) SetMute(volume bool) error {
	return h.renderingControl.SetMute(volume)
}

// Short for `zp.RenderingControl.GetMute`.
func (h *Helper) GetMute() (bool, error) {
	return h.renderingControl.GetMute()
}

// Short for `zp.RenderingControl.SetBass`.
func (h *Helper) SetBass(volume int) error {
	return h.renderingControl.SetBass(volume)
}

// Short for `zp.RenderingControl.GetBass`.
func (h *Helper) GetBass() (int, error) {
	return h.renderingControl.GetBass()
}

// Short for `zp.RenderingControl.SetTreble`.
func (h *Helper) SetTreble(volume int) error {
	return h.renderingControl.SetTreble(volume)
}

// Short for `zp.RenderingControl.GetTreble`.
func (h *Helper) GetTreble() (int, error) {
	return h.renderingControl.GetTreble()
}

// Short for `zp.RenderingControl.SetVolume`.
func (h *Helper) SetVolume(volume int) error {
	return h.renderingControl.SetVolume(volume)
}

// Short for `zp.RenderingControl.SetRelativeVolume`.
func (h *Helper) SetVolumeDelta(volume int) (int, error) {
	return h.renderingControl.SetRelativeVolume(volume)
}

// Short for `zp.RenderingControl.GetVolume`.
func (h *Helper) GetVolume() (int, error) {
	return h.renderingControl.GetVolume()
}

// Short for `zp.RenderingControl.SetVolumeDB`.
func (h *Helper) SetVolumeDB(volume int) error {
	return h.renderingControl.SetVolumeDB(volume)
}

// Short for `zp.RenderingControl.GetVolumeDB`.
func (h *Helper) GetVolumeDB() (int, error) {
	return h.renderingControl.GetVolumeDB()
}
