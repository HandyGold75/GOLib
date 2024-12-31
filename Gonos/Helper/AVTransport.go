package Helper

import (
	"Gonos/lib"
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
func (h *Helper) GetTrackInfo() (*trackInfo, error) {
	info, err := h.aVTransport.GetPositionInfo()
	if err != nil {
		return &trackInfo{}, err
	}
	return &trackInfo{
		QuePosition: info.Track,
		Duration:    info.TrackDuration,
		URI:         info.TrackURI,
		Progress:    info.RelTime,
		AlbumArtURI: info.TrackMetaDataParsed.AlbumArtUri,
		Title:       info.TrackMetaDataParsed.Title,
		Class:       info.TrackMetaDataParsed.Class,
		Creator:     info.TrackMetaDataParsed.Creator,
		Album:       info.TrackMetaDataParsed.Album,
	}, nil
}

// Get current transport state, this can be one of `STOPPED`, `PLAYING`, `PAUSED_PLAYBACK`, `TRANSITIONING`.
func (h *Helper) GetCurrentTransportState() (string, error) {
	res, err := h.aVTransport.GetTransportInfo()
	return res.CurrentTransportState, err
}

// Short for `zp.AVTransport.Stop`
func (h *Helper) Stop() error {
	return h.aVTransport.Stop()
}

// Short for `h.GetCurrentTransportState() == "STOPPED"`.
func (h *Helper) GetStop() (bool, error) {
	state, err := h.GetCurrentTransportState()
	return state == "STOPPED", err
}

// Short for `zp.AVTransport.Play`
func (h *Helper) Play() error {
	return h.aVTransport.Play()
}

// Short for `h.GetCurrentTransportState() == "PLAYING"`.
func (h *Helper) GetPlay() (bool, error) {
	state, err := h.GetCurrentTransportState()
	return state == "PLAYING", err
}

// Short for `zp.AVTransport.Pause`
func (h *Helper) Pause() error {
	return h.aVTransport.Pause()
}

// Short for `h.GetCurrentTransportState() == "PAUSED_PLAYBACK"`.
func (h *Helper) GetPause() (bool, error) {
	state, err := h.GetCurrentTransportState()
	return state == "PAUSED_PLAYBACK", err
}

// Short for `h.GetCurrentTransportState() == "TRANSITIONING"`.
func (h *Helper) GetTransitioning() (bool, error) {
	state, err := h.GetCurrentTransportState()
	return state == "TRANSITIONING", err
}

// Short for `zp.AVTransport.Next`.
func (h *Helper) Next() error {
	return h.aVTransport.Next()
}

// Short for `zp.AVTransport.`.
func (h *Helper) Previous() error {
	return h.aVTransport.Previous()
}

// Get current transport status.
func (h *Helper) GetCurrentTransportStatus() (string, error) {
	res, err := h.aVTransport.GetTransportInfo()
	return res.CurrentTransportStatus, err
}

// Get current speed.
func (h *Helper) GetCurrentSpeed() (string, error) {
	res, err := h.aVTransport.GetTransportInfo()
	return res.CurrentSpeed, err
}

// Will always return false for all if a third party application is controling playback.
func (h *Helper) GetPlayMode() (shuffle bool, repeat bool, repeatOne bool, err error) {
	res, err := h.aVTransport.GetTransportSettings()
	if err != nil {
		return false, false, false, err
	}
	modeBools, ok := lib.PlayModeMap[res.PlayMode]
	if !ok {
		return false, false, false, lib.ErrSonos.ErrUnexpectedResponse
	}
	return modeBools[0], modeBools[1], modeBools[2], nil
}

// Will always return false if a third party application is controling playback.
func (h *Helper) GetShuffle() (bool, error) {
	shuffle, _, _, err := h.GetPlayMode()
	return shuffle, err
}

// Will always disable other playmodes if a third party application is controling playback, as we can not determine the actual state.
func (h *Helper) SetShuffle(state bool) error {
	_, repeat, repeatOne, err := h.GetPlayMode()
	if err != nil {
		return err
	}
	return h.aVTransport.SetPlayMode(state, repeat, repeatOne)
}

// Will always return false if a third party application is controling playback.
func (h *Helper) GetRepeat() (bool, error) {
	_, repeat, _, err := h.GetPlayMode()
	return repeat, err
}

// If enabled then repeat one will be disabled.
//
// Will always disable other playmodes if a third party application is controling playback, as we can not determine the actual state.
func (h *Helper) SetRepeat(state bool) error {
	shuffle, _, repeatOne, err := h.GetPlayMode()
	if err != nil {
		return err
	}
	return h.aVTransport.SetPlayMode(shuffle, state, repeatOne && !state)
}

// Will always return false if a third party application is controling playback.
func (h *Helper) GetRepeatOne() (bool, error) {
	_, _, repeatOne, err := h.GetPlayMode()
	return repeatOne, err
}

// If enabled then repeat will be disabled.
//
// Will always disable other playmodes if a third party application is controling playback, as we can not determine the actual state.
func (h *Helper) SetRepeatOne(state bool) error {
	shuffle, repeat, _, err := h.GetPlayMode()
	if err != nil {
		return err
	}
	return h.aVTransport.SetPlayMode(shuffle, repeat && !state, state)
}

// Returns `NOT_IMPLEMENTED`.
func (h *Helper) GetRecQualityMode() (string, error) {
	res, err := h.aVTransport.GetTransportSettings()
	return res.RecQualityMode, err
}

// Go to track by index (index starts at 1).
//
// Will always fail if a third party application is controling playback.
func (h *Helper) SeekTrack(track int) error {
	return h.aVTransport.Seek(lib.SeekModes.Track, strconv.Itoa(max(1, track)))
}

// Go to track time (Absolute).
func (h *Helper) SeekTime(seconds int) error {
	return h.aVTransport.Seek(lib.SeekModes.Relative, time.Time.Add(time.Time{}, time.Second*time.Duration(max(0, seconds))).Format("15:04:05"))
}

// Go to track time (Relative).
func (h *Helper) SeekTimeDelta(seconds int) error {
	prefix := "+"
	if seconds < 0 {
		seconds = -seconds
		prefix = "-"
	}
	return h.aVTransport.Seek(lib.SeekModes.Absolute, prefix+time.Time.Add(time.Time{}, time.Second*time.Duration(seconds)).Format("15:04:05"))
}

// Short for `zp.AVTransport.RemoveTrackFromQueue`.
func (h *Helper) QueRemove(track int) error {
	return h.aVTransport.RemoveTrackFromQueue(lib.ContentTypes.QueueMain, track)
}

// Short for `zp.AVTransport.RemoveTrackFromQueue`.
func (h *Helper) QueSecondRemove(track int) error {
	return h.aVTransport.RemoveTrackFromQueue(lib.ContentTypes.QueueMain, track)
}

// Short for `zp.AVTransport.RemoveAllTracksFromQueue`.
func (h *Helper) QueClear() error {
	return h.aVTransport.RemoveAllTracksFromQueue()
}
