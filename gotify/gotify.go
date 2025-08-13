package gotify

import (
	"io"
	"net/http"
	"slices"
	"strings"

	"github.com/HandyGold75/GOLib/gotify/Player"
	"github.com/HandyGold75/GOLib/gotify/lib"
)

type (
	GotifyPlayer struct {
		URL string

		cl *http.Client

		Player Player.Player
	}
)

const (
	GET  = lib.GET
	PUT  = lib.PUT
	POST = lib.POST
)

func NewGotifyPlayer(cl *http.Client) (*GotifyPlayer, error) {
	gp := &GotifyPlayer{URL: "https://api.spotify.com/v1/me", cl: cl}
	gp.Player = Player.New(gp.SendPlayer)
	return gp, nil
}

func (gp *GotifyPlayer) SendPlayer(method lib.HttpMethod, action string, options [][2]string) (string, error) {
	return Send(method, gp.URL+"/player/"+action, options)
}

func Send(method lib.HttpMethod, url string, options [][2]string) (string, error) {
	opts := ""
	for _, opt := range slices.DeleteFunc(options, func(o [2]string) bool { return o[0] == "" || o[1] == "" }) {
		if opts != "" {
			opts += "&"
		}
		opts += opt[0] + "=" + opt[1]
	}
	if opts != "" {
		opts = "?" + opts
	}

	req, err := http.NewRequest(string(method), url+opts, strings.NewReader(""))
	if err != nil {
		return "", err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(res[:]), nil
}
