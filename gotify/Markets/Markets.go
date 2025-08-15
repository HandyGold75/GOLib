package Markets

import (
	"encoding/json"

	"github.com/HandyGold75/GOLib/gotify/lib"
)

type Markets struct {
	Send func(method lib.HttpMethod, action string, options [][2]string, body []byte) ([]byte, error)
}

func New(send func(method lib.HttpMethod, action string, options [][2]string, body []byte) ([]byte, error)) Markets {
	return Markets{Send: send}
}

func (s *Markets) GetAvailableMarkets() ([]string, error) {
	res, err := s.Send(lib.GET, "", [][2]string{}, []byte{})
	if err != nil {
		return []string{}, err
	}
	data := struct {
		Markets []string `json:"markets"`
	}{}
	err = json.Unmarshal(res, &data)
	return data.Markets, err
}
