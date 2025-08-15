package Chapters

import (
	"encoding/json"
	"strings"

	"github.com/HandyGold75/GOLib/gotify/lib"
)

type (
	Chapters struct {
		Send   func(method lib.HttpMethod, action string, options [][2]string, body []byte) ([]byte, error)
		Market string // An ISO 3166-1 alpha-2 country code, https://en.wikipedia.org/wiki/ISO_3166-1_alpha-2
	}

	getAChapter lib.ChapterObject

	getSeveralChapters lib.Chapters
)

func New(send func(method lib.HttpMethod, action string, options [][2]string, body []byte) ([]byte, error)) Chapters {
	return Chapters{Send: send, Market: ""}
}

func (s *Chapters) GetAChapter(id string) (getAChapter, error) {
	res, err := s.Send(lib.GET, id, [][2]string{{"market", s.Market}}, []byte{})
	if err != nil {
		return getAChapter{}, err
	}
	data := getAChapter{}
	err = json.Unmarshal(res, &data)
	return data, err
}

func (s *Chapters) GetSeveralChapters(ids []string) (getSeveralChapters, error) {
	res, err := s.Send(lib.GET, "", [][2]string{{"market", s.Market}, {"ids", strings.Join(ids, ",")}}, []byte{})
	if err != nil {
		return getSeveralChapters{}, err
	}
	data := getSeveralChapters{}
	err = json.Unmarshal(res, &data)
	return data, err
}
