package Player

import (
	"strconv"
	"time"

	"github.com/HandyGold75/GOLib/gotify/lib"
)

type Player struct {
	Send     func(method lib.HttpMethod, action string, options [][2]string) (string, error)
	DeviceID string
}

func New(send func(method lib.HttpMethod, action string, options [][2]string) (string, error)) Player {
	return Player{Send: send, DeviceID: ""}
}

func (s *Player) SkipToNext() error {
	_, err := s.Send(lib.POST, "next", [][2]string{{"device_id", s.DeviceID}})
	return err
}

func (s *Player) SkipToPrevious() error {
	_, err := s.Send(lib.POST, "next", [][2]string{{"device_id", s.DeviceID}})
	return err
}

func (s *Player) SeekToPosition(pos time.Duration) error {
	_, err := s.Send(lib.PUT, "seek", [][2]string{{"device_id", s.DeviceID}, {"position_ms", strconv.Itoa(int(pos.Milliseconds()))}})
	return err
}
