package Gonos

import (
	"strconv"
	"time"
)

type (
	trackInfo struct {
		QuePosition int
		Duration    string
		URI         string
		Progress    string
		AlbumArtURI string
		Title       string
		Class       string
		Creator     string
		Album       string
	}
)

// Get simplified info about currently playing track.
func (zp *ZonePlayer) GetTrackInfo() (*trackInfo, error) {
	info, err := zp.GetPositionInfo()
	if err != nil {
		return &trackInfo{}, err
	}
	return &trackInfo{
		QuePosition: info.Track,
		Duration:    info.TrackDuration,
		URI:         info.TrackURI,
		Progress:    info.RelTime,
		AlbumArtURI: "http://" + zp.IpAddress.String() + ":1400" + info.TrackMetaDataParsed.AlbumArtUri,
		Title:       info.TrackMetaDataParsed.Title,
		Class:       info.TrackMetaDataParsed.Class,
		Creator:     info.TrackMetaDataParsed.Creator,
		Album:       info.TrackMetaDataParsed.Album,
	}, nil
}

// Get current transport state, this can be one of `STOPPED`, `PLAYING`, `PAUSED_PLAYBACK`, `TRANSITIONING`.
func (zp *ZonePlayer) GetCurrentTransportState() (string, error) {
	res, err := zp.GetTransportInfo()
	return res.CurrentTransportState, err
}

// Short for `zp.GetTrackInfo() == "STOPPED"`.
func (zp *ZonePlayer) GetStop() (bool, error) {
	state, err := zp.GetCurrentTransportState()
	return state == "STOPPED", err
}

// Short for `zp.GetTrackInfo() == "PLAYING"`.
func (zp *ZonePlayer) GetPlay() (bool, error) {
	state, err := zp.GetCurrentTransportState()
	return state == "PLAYING", err
}

// Short for `zp.GetTrackInfo() == "PAUSED_PLAYBACK"`.
func (zp *ZonePlayer) GetPause() (bool, error) {
	state, err := zp.GetCurrentTransportState()
	return state == "PAUSED_PLAYBACK", err
}

// Short for `zp.GetTrackInfo() == "TRANSITIONING"`.
func (zp *ZonePlayer) GetTransitioning() (bool, error) {
	state, err := zp.GetCurrentTransportState()
	return state == "TRANSITIONING", err
}

// Get current transport status.
func (zp *ZonePlayer) GetCurrentTransportStatus() (string, error) {
	res, err := zp.GetTransportInfo()
	return res.CurrentTransportStatus, err
}

// Get current speed.
func (zp *ZonePlayer) GetCurrentSpeed() (string, error) {
	res, err := zp.GetTransportInfo()
	return res.CurrentSpeed, err
}

// Will always return false for all if a third party application is controling playback.
func (zp *ZonePlayer) GetPlayMode() (shuffle bool, repeat bool, repeatOne bool, err error) {
	res, err := zp.GetTransportSettings()
	if err != nil {
		return false, false, false, err
	}
	modeBools, ok := PlayModeMap[res.PlayMode]
	if !ok {
		return false, false, false, ErrSonos.ErrUnexpectedResponse
	}
	return modeBools[0], modeBools[1], modeBools[2], nil
}

// Will always return false if a third party application is controling playback.
func (zp *ZonePlayer) GetShuffle() (bool, error) {
	shuffle, _, _, err := zp.GetPlayMode()
	return shuffle, err
}

// Will always disable other playmodes if a third party application is controling playback, as we can not determine the actual state.
func (zp *ZonePlayer) SetShuffle(state bool) error {
	_, repeat, repeatOne, err := zp.GetPlayMode()
	if err != nil {
		return err
	}
	return zp.SetPlayMode(state, repeat, repeatOne)
}

// Will always return false if a third party application is controling playback.
func (zp *ZonePlayer) GetRepeat() (bool, error) {
	_, repeat, _, err := zp.GetPlayMode()
	return repeat, err
}

// If enabled then repeat one will be disabled.
//
// Will always disable other playmodes if a third party application is controling playback, as we can not determine the actual state.
func (zp *ZonePlayer) SetRepeat(state bool) error {
	shuffle, _, repeatOne, err := zp.GetPlayMode()
	if err != nil {
		return err
	}
	return zp.SetPlayMode(shuffle, state, repeatOne && !state)
}

// Will always return false if a third party application is controling playback.
func (zp *ZonePlayer) GetRepeatOne() (bool, error) {
	_, _, repeatOne, err := zp.GetPlayMode()
	return repeatOne, err
}

// If enabled then repeat will be disabled.
//
// Will always disable other playmodes if a third party application is controling playback, as we can not determine the actual state.
func (zp *ZonePlayer) SetRepeatOne(state bool) error {
	shuffle, repeat, _, err := zp.GetPlayMode()
	if err != nil {
		return err
	}
	return zp.SetPlayMode(shuffle, repeat && !state, state)
}

// Returns `NOT_IMPLEMENTED`.
func (zp *ZonePlayer) GetRecQualityMode() (string, error) {
	res, err := zp.GetTransportSettings()
	return res.RecQualityMode, err
}

// Go to track by index (index starts at 1).
//
// Will always fail if a third party application is controling playback.
func (zp *ZonePlayer) SeekTrack(track int) error {
	return zp.Seek(SeekModes.Track, strconv.Itoa(max(1, track)))
}

// Go to track time (Absolute).
func (zp *ZonePlayer) SeekTime(seconds int) error {
	return zp.Seek(SeekModes.Relative, time.Time.Add(time.Time{}, time.Second*time.Duration(max(0, seconds))).Format("15:04:05"))
}

// Go to track time (Relative).
func (zp *ZonePlayer) SeekTimeDelta(seconds int) error {
	prefix := "+"
	if seconds < 0 {
		seconds = -seconds
		prefix = "-"
	}
	return zp.Seek(SeekModes.Absolute, prefix+time.Time.Add(time.Time{}, time.Second*time.Duration(seconds)).Format("15:04:05"))
}
