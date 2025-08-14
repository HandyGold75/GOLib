package gotify

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/HandyGold75/GOLib/gotify/Player"
	"github.com/HandyGold75/GOLib/gotify/Users"
	"github.com/HandyGold75/GOLib/gotify/lib"
	"golang.org/x/oauth2"
)

type (
	scope string // https://developer.spotify.com/documentation/web-api/concepts/scopes

	GotifyPlayer struct {
		URL string

		authCfg             oauth2.Config
		authUserMsgCallback func(url string)
		cl                  *http.Client

		Player Player.Player
		Users  Users.Users
	}

	errorResponse struct {
		Error struct {
			Status  int    `json:"status"`
			Message string `json:"message"`
		} `json:"error"`
	}
)

const (
	GET  = lib.GET
	PUT  = lib.PUT
	POST = lib.POST

	RepeatTrack   = lib.RepeatTrack
	RepeatContext = lib.RepeatContext
	RepeatOff     = lib.RepeatOff

	ScopeUgcImageUpload            scope = "ugc-image-upload"            // Write access to user-provided images.
	ScopeUserReadPlaybackState     scope = "user-read-playback-state"    // Read access to a user’s player state.
	ScopeUserModifyPlaybackState   scope = "user-modify-playback-state"  // Write access to a user’s playback state
	ScopeUserReadCurrentlyPlaying  scope = "user-read-currently-playing" // Read access to a user’s currently playing content.
	ScopeAppRemoteControl          scope = "app-remote-control"          // Remote control playback of Spotify. This scope is currently available to Spotify iOS and Android SDKs.
	ScopeStreaming                 scope = "streaming"                   // Control playback of a Spotify track. This scope is currently available to the Web Playback SDK. The user must have a Spotify Premium account.
	ScopePlaylistReadPrivate       scope = "playlist-read-private"       // Read access to user's private playlists.
	ScopePlaylistReadCollaborative scope = "playlist-read-collaborative" // Include collaborative playlists when requesting a user's playlists.
	ScopePlaylistModifyPrivate     scope = "playlist-modify-private"     // Write access to a user's private playlists.
	ScopePlaylistModifyPublic      scope = "playlist-modify-public"      // Write access to a user's public playlists.
	ScopeUserFollowModify          scope = "user-follow-modify"          // Write/delete access to the list of artists and other users that the user follows.
	ScopeUserFollowRead            scope = "user-follow-read"            // Read access to the list of artists and other users that the user follows.
	ScopeUserReadPlaybackPosition  scope = "user-read-playback-position" // Read access to a user’s playback position in a content.
	ScopeUserTopRead               scope = "user-top-read"               // Read access to a user's top artists and tracks.
	ScopeUserReadRecentlyPlayed    scope = "user-read-recently-played"   // Read access to a user’s recently played tracks.
	ScopeUserLibraryModify         scope = "user-library-modify"         // Write/delete access to a user's "Your Music" library.
	ScopeUserLibraryRead           scope = "user-library-read"           // Read access to a user's library.
	ScopeUserReadEmail             scope = "user-read-email"             // Read access to user’s email address.
	ScopeUserReadPrivate           scope = "user-read-private"           // Read access to user’s subscription details (type of user account).
	ScopeUserPersonalized          scope = "user-personalized"           // Get personalized content for the user.
	ScopeUserSoaLink               scope = "user-soa-link"               // Link a partner user account to a Spotify user account
	ScopeUserSoaUnlink             scope = "user-soa-unlink"             // Unlink a partner user account from a Spotify account
	ScopeSoaManageEntitlements     scope = "soa-manage-entitlements"     // Modify entitlements for linked users
	ScopeSoaManagePartner          scope = "soa-manage-partner"          // Update partner information
	ScopeSoaCreatePartner          scope = "soa-create-partner"          // Create new partners, platform partners only
)

func NewGotifyPlayer(clientID, redirectURL string, scopes ...scope) *GotifyPlayer {
	scps := []string{}
	for _, scp := range scopes {
		scps = append(scps, string(scp))
	}
	gp := &GotifyPlayer{
		URL: "https://api.spotify.com/v1",
		authCfg: oauth2.Config{
			ClientID: clientID,
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://accounts.spotify.com/authorize",
				TokenURL: "https://accounts.spotify.com/api/token",
			},
			RedirectURL: redirectURL,
			Scopes:      scps,
		},
		authUserMsgCallback: func(url string) { fmt.Print("\r\nLogin: " + url + "\r\nPaste: ") },
		cl:                  http.DefaultClient,
	}
	gp.Player = Player.New(gp.SendPlayer)
	gp.Users = Users.New(gp.SendUsers)
	return gp
}

// Authenticate using stdin.
func (gp *GotifyPlayer) AuthenticateStdin() error {
	verifier, state, ch := oauth2.GenerateVerifier(), oauth2.GenerateVerifier(), make(chan string)
	go func() {
		defer close(ch)
		msg := ""
		for msg == "" {
			m, err := bufio.NewReader(os.Stdin).ReadString('\n')
			if err != nil {
				return
			}
			msgSplit := strings.Split(strings.TrimSuffix(m, "\n"), "?")
			msg = msgSplit[len(msgSplit)-1]
		}
		ch <- msg
	}()

	gp.authUserMsgCallback(gp.authCfg.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.S256ChallengeOption(verifier)))
	msg, ok := <-ch
	if !ok {
		return errors.New("failed authentication")
	}
	code, actualState := "", ""
	for pair := range strings.SplitSeq(msg, "&") {
		if strings.HasPrefix(pair, "code=") {
			code = strings.Replace(pair, "code=", "", 1)
		} else if strings.HasPrefix(pair, "state=") {
			actualState = strings.Replace(pair, "state=", "", 1)
		}
	}
	if code == "" || actualState != state {
		return errors.New("failed authentication")
	}
	token, err := gp.authCfg.Exchange(context.Background(), code, oauth2.VerifierOption(verifier))
	if err != nil {
		return err
	}
	gp.cl = gp.authCfg.Client(context.Background(), token)
	return nil
}

// Authenticate using local http server.
func (gp *GotifyPlayer) AuthenticateHTTP(port uint16) error {
	verifier, state, ch := oauth2.GenerateVerifier(), oauth2.GenerateVerifier(), make(chan string)
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		defer close(ch)
		values := r.URL.Query()
		if e := values.Get("error"); e != "" || values.Get("state") != string(state) || r.FormValue("state") != string(state) {
			return
		}
		ch <- values.Get("code")
	})
	server := &http.Server{Addr: ":" + strconv.FormatUint(uint64(port), 10), Handler: nil}
	go func() { _ = server.ListenAndServe() }()
	defer server.Close()

	gp.authUserMsgCallback(gp.authCfg.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.S256ChallengeOption(verifier)))
	code, ok := <-ch
	if !ok {
		return errors.New("failed authentication")
	}
	token, err := gp.authCfg.Exchange(context.Background(), code)
	if err != nil {
		return err
	}
	gp.cl = gp.authCfg.Client(context.Background(), token)
	return nil
}

// Authenticate using a token.
func (gp *GotifyPlayer) AuthenticateToken(token *oauth2.Token) error {
	token.Expiry.Add(-(time.Hour * 2))
	token, err := gp.authCfg.TokenSource(context.Background(), token).Token()
	if err != nil {
		return err
	}
	gp.cl = gp.authCfg.Client(context.Background(), token)
	return nil
}

// Token get current active token.
func (gp *GotifyPlayer) Token() (*oauth2.Token, error) {
	transport, ok := gp.cl.Transport.(*oauth2.Transport)
	if !ok {
		return nil, errors.New("client not backed by oauth2 transport")
	}
	return transport.Source.Token()
}

func (gp *GotifyPlayer) SendPlayer(method lib.HttpMethod, action string, options [][2]string, body []byte) ([]byte, error) {
	return gp.Send(method, gp.URL+"/me/player/"+action, options, body)
}

func (gp *GotifyPlayer) SendUsers(method lib.HttpMethod, action string, options [][2]string, body []byte) ([]byte, error) {
	return gp.Send(method, gp.URL+"/"+action, options, body)
}

func (gp *GotifyPlayer) Send(method lib.HttpMethod, url string, options [][2]string, body []byte) ([]byte, error) {
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

	req, err := http.NewRequest(string(method), url+opts, strings.NewReader(string(body[:])))
	if err != nil {
		return []byte{}, err
	}
	if len(body) > 0 {
		req.Header.Add("Content-Type", "application/json")
	}
	resp, err := gp.cl.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer func() { _ = resp.Body.Close() }()

	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}
	data := errorResponse{}
	if err := json.Unmarshal(res, &data); err == nil && data.Error.Status != 0 {
		return []byte{}, errors.New(strconv.Itoa(data.Error.Status) + ": " + data.Error.Message)
	}
	return res, nil
}
