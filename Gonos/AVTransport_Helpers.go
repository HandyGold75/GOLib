package Gonos

import (
	"strconv"
	"time"
)

type (
	TrackInfo struct {
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

// TODO: Test
func (zp *ZonePlayer) GetTrackInfo() (*TrackInfo, error) {
	info, err := zp.GetPositionInfo()
	if err != nil {
		return &TrackInfo{}, err
	}
	return &TrackInfo{
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

// TODO: Test
func (zp *ZonePlayer) GetCurrentTransportState() (string, error) {
	res, err := zp.GetTransportInfo()
	return res.CurrentTransportState, err
}

// TODO: Test
func (zp *ZonePlayer) GetStop() (bool, error) {
	state, err := zp.GetCurrentTransportState()
	return state == "STOPPED", err
}

// TODO: Test
func (zp *ZonePlayer) GetPlay() (bool, error) {
	state, err := zp.GetCurrentTransportState()
	return state == "PLAYING", err
}

// TODO: Test
func (zp *ZonePlayer) GetPause() (bool, error) {
	state, err := zp.GetCurrentTransportState()
	return state == "PAUSED_PLAYBACK", err
}

// TODO: Test
func (zp *ZonePlayer) GetTransitioning() (bool, error) {
	state, err := zp.GetCurrentTransportState()
	return state == "TRANSITIONING", err
}

// TODO: Test
func (zp *ZonePlayer) GetCurrentTransportStatus() (string, error) {
	res, err := zp.GetTransportInfo()
	return res.CurrentTransportStatus, err
}

// TODO: Test
func (zp *ZonePlayer) GetCurrentSpeed() (string, error) {
	res, err := zp.GetTransportInfo()
	return res.CurrentSpeed, err
}

// TODO: Test
func (zp *ZonePlayer) GetPlayMode() (shuffle bool, repeat bool, repeatOne bool, err error) {
	res, err := zp.GetTransportSettings()
	if err != nil {
		return false, false, false, err
	}
	modeBools, ok := PlaymodesMap[res.PlayMode]
	if !ok {
		return false, false, false, ErrSonos.ErrUnexpectedResponse
	}
	return modeBools[0], modeBools[1], modeBools[2], nil
}

// TODO: test
func (zp *ZonePlayer) GetShuffle() (bool, error) {
	shuffle, _, _, err := zp.GetPlayMode()
	return shuffle, err
}

// TODO: test
func (zp *ZonePlayer) SetShuffle(state bool) error {
	_, repeat, repeatOne, err := zp.GetPlayMode()
	if err != nil {
		return err
	}
	return zp.SetPlayMode(state, repeat, repeatOne)
}

// TODO: test
func (zp *ZonePlayer) GetRepeat() (bool, error) {
	_, repeat, _, err := zp.GetPlayMode()
	return repeat, err
}

// TODO: test
func (zp *ZonePlayer) SetRepeat(state bool) error {
	shuffle, _, repeatOne, err := zp.GetPlayMode()
	if err != nil {
		return err
	}
	return zp.SetPlayMode(shuffle, state, repeatOne && !state)
}

// TODO: test
func (zp *ZonePlayer) GetRepeatOne() (bool, error) {
	_, _, repeatOne, err := zp.GetPlayMode()
	return repeatOne, err
}

// TODO: test
func (zp *ZonePlayer) SetRepeatOne(state bool) error {
	shuffle, repeat, _, err := zp.GetPlayMode()
	if err != nil {
		return err
	}
	return zp.SetPlayMode(shuffle, repeat && !state, state)
}

// TODO: Test
func (zp *ZonePlayer) GetRecQualityMode() (string, error) {
	res, err := zp.GetTransportSettings()
	return res.RecQualityMode, err
}

// TODO: Test
func (zp *ZonePlayer) SeekTrack(track int) error {
	return zp.Seek("TRACK_NR", strconv.Itoa(max(1, track)))
}

// TODO: Test
func (zp *ZonePlayer) SeekTime(seconds int) error {
	return zp.Seek("REL_TIME", time.Time.Add(time.Time{}, time.Second*time.Duration(max(0, seconds))).Format("15:04:05"))
}

// TODO: Test
func (zp *ZonePlayer) SeekTimeDelta(seconds int) error {
	prefix := "+"
	if seconds < 0 {
		seconds = -seconds
		prefix = "-"
	}
	return zp.Seek("TIME_DELTA", prefix+time.Time.Add(time.Time{}, time.Second*time.Duration(seconds)).Format("15:04:05"))
}
