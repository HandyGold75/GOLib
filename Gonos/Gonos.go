package Gonos

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

type (
	command    struct{ Endpoint, Action, Body, ExpectedResponse, TargetTag string }
	errSonos   struct{ ErrUnexpectedResponse, ErrInvalidIPAdress, ErrNoZonePlayerFound, ErrInvalidEndpoint, ErrInvalidContentType, ErrInvalidPlayMode error }
	ZonePlayer struct{ IpAddress net.IP }
)

var ErrSonos = errSonos{
	ErrUnexpectedResponse: errors.New("unexpected response"),
	ErrInvalidIPAdress:    errors.New("unable to discover zone player"),
	ErrNoZonePlayerFound:  errors.New("unable to find zone player"),
	ErrInvalidEndpoint:    errors.New("invalid endpoint"),
	ErrInvalidContentType: errors.New("invalid content type"),
	ErrInvalidPlayMode:    errors.New("invalid play mode"),
}

var Endpoints = map[string]string{
	"AVTransport":      "/MediaRenderer/AVTransport/Control",
	"RenderingControl": "/MediaRenderer/RenderingControl/Control",
	"DeviceProperties": "/DeviceProperties/Control",
	"ContentDirectory": "/MediaServer/ContentDirectory/Control",
}
var EndpointsBodyPrefix = map[string]string{
	"AVTransport":      "<InstanceID>0</InstanceID>",
	"RenderingControl": "<InstanceID>0</InstanceID><Channel>Master</Channel>",
	"DeviceProperties": "<InstanceID>0</InstanceID><Channel>Master</Channel>",
	"ContentDirectory": "",
}

var ContentTypes = map[string]string{
	"artist":          "A:ARTIST",
	"albumartist":     "A:ALBUMARTIST",
	"album":           "A:ALBUM",
	"genre":           "A:GENRE",
	"composer":        "A:COMPOSER",
	"track":           "A:TRACKS",
	"playlists":       "A:PLAYLISTS",
	"music library":   "A",
	"share":           "S:",
	"sonos playlists": "SQ:",
	"sonos favorites": "FV:2",
	"radio stations":  "R:0/0",
	"radio shows":     "R:0/1",
	"queue":           "Q:", // Maybe Q:0 ??
}

var Playmodes = map[string][3]bool{
	// "MODE": [2]bool{shuffle, repeat, repeat_one}
	"NORMAL":             {false, false, false},
	"SHUFFLE_NOREPEAT":   {true, false, false},
	"SHUFFLE":            {true, true, false},
	"REPEAT_ALL":         {false, true, false},
	"SHUFFLE_REPEAT_ONE": {true, false, true},
	"REPEAT_ONE":         {false, false, true},
}
var PlaymodesReversed = func() map[[3]bool]string {
	PMS := map[[3]bool]string{}
	for k, v := range Playmodes {
		PMS[v] = k
	}
	return PMS
}()

func unmarshalMetaData[T any](data string, v T) error {
	data = strings.ReplaceAll(data, "&apos;", "'")
	data = strings.ReplaceAll(data, "&quot;", "\"")
	data = strings.ReplaceAll(data, "&gt;", ">")
	data = strings.ReplaceAll(data, "&lt;", "<")
	data, err := extractTag(data, "DIDL-Lite")
	if err != nil {
		return err
	}
	if err := xml.Unmarshal([]byte(data), v); err != nil {
		return err
	}
	return nil
}

func extractTag(data, tag string) (string, error) {
	if start, end := strings.Index(data, "<"+tag+">"), strings.Index(data, "</"+tag+">"); start != -1 && end != -1 {
		return data[start+len(tag)+2 : end], nil
	}
	if start, end := strings.Index(data, "<"+tag+" "), strings.Index(data, "</"+tag+">"); start != -1 && end != -1 {
		data = data[start+len(tag)+2 : end]
		if mid := strings.Index(data, ">"); mid != -1 {
			return data[mid+1:], nil
		}
	}
	return data, ErrSonos.ErrUnexpectedResponse
}

func boolTo10(b bool) string {
	if b {
		return "1"
	}
	return "0"
}
func boolToOnOff(b bool) string {
	if b {
		return "On"
	}
	return "Off"
}

// Create new ZonePlayer for controling a Sonos speaker.
func NewZonePlayer(ipAddress string) (*ZonePlayer, error) {
	ip := net.ParseIP(ipAddress)
	if ip == nil {
		return &ZonePlayer{}, ErrSonos.ErrInvalidIPAdress
	}
	return &ZonePlayer{IpAddress: ip}, nil
}

// Create new ZonePlayer using discovery controling a Sonos speaker. (TODO: Broken?)
func DiscoverZonePlayer() (*ZonePlayer, error) {
	conn, err := net.DialUDP("udp", &net.UDPAddr{Port: 1900}, &net.UDPAddr{IP: net.IPv4(239, 255, 255, 250), Port: 1900})
	if err != nil {
		return &ZonePlayer{}, err
	}
	defer conn.Close()

	chOut := make(chan string)
	go func() {
		buf := make([]byte, 1024)
		n, addr, err := conn.ReadFrom(buf)
		if err != nil {
			chOut <- ""
			return
		}
		fmt.Println("---")
		fmt.Println(addr)
		fmt.Println(buf[:n])
		chOut <- strings.Split(addr.String(), ":")[0]
	}()

	for i := 0; i < 3; i++ {
		_, _ = conn.Write([]byte("M-SEARCH * HTTP/1.1\r\nHOST: 239.255.255.250:1900\r\nMAN: \"ssdp:discover\"\r\nMX: 1\r\nST: urn:schemas-upnp-org:device:ZonePlayer:1\r\n\r\n"))
	}

	select {
	case addr := <-chOut:
		if addr == "" {
			return &ZonePlayer{}, ErrSonos.ErrNoZonePlayerFound
		}
		return NewZonePlayer(addr)
	case <-time.After(time.Second):
		return &ZonePlayer{}, ErrSonos.ErrNoZonePlayerFound
	}
}

// Create new ZonePlayer using network scanning controling a Sonos speaker.
func ScanZonePlayer(cidr string) ([]*ZonePlayer, error) {
	var incIP = func(ip net.IP) {
		for j := len(ip) - 1; j >= 0; j-- {
			ip[j]++
			if ip[j] > 0 {
				break
			}
		}
	}

	ip, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return []*ZonePlayer{}, err
	}

	wg, zps := sync.WaitGroup{}, []*ZonePlayer{}
	for ip := ip.Mask(ipNet.Mask); ipNet.Contains(ip); incIP(ip) {
		wg.Add(1)
		go func(ip string) {
			defer wg.Done()

			conn, err := net.DialTimeout("tcp", ip+":"+"1400", time.Second)
			if err != nil {
				return
			}
			defer conn.Close()

			zp, err := NewZonePlayer(ip)
			if err != nil {
				return
			}
			if _, err = zp.GetState(); err == nil {
				zps = append(zps, zp)
			}
		}(ip.String())
	}
	wg.Wait()

	if len(zps) <= 0 {
		return zps, ErrSonos.ErrNoZonePlayerFound
	}

	return zps, nil
}

type (
	TrackInfo struct {
		QuePosition string
		Duration    string
		URI         string
		Progress    string
		AlbumArtURI string
		Title       string
		Class       string
		Creator     string
		Album       string
	}

	TrackInfoRaw struct {
		XMLName       xml.Name `xml:"GetPositionInfoResponse"`
		Track         string
		TrackDuration string
		TrackMetaData string
		TrackURI      string
		RelTime       string
		AbsTime       string
		RelCount      string
		AbsCount      string
	}
)

type (
	Que struct {
		Count      string
		TotalCount string
		Tracks     []QueTrack
	}

	QueTrack struct {
		AlbumArtURI string
		Title       string
		Class       string
		Creator     string
		Album       string
	}
)

type (
	Favorites struct {
		Count      string
		TotalCount string
		Favorites  []FavoritesItem
	}

	FavoritesItem struct {
		AlbumArtURI string
		Title       string
		Description string
		Class       string
		Type        string
	}
)

func (zp *ZonePlayer) sendCommand(endpoint string, action string, body string, targetTag string) (string, error) {
	endpointUri, ok := Endpoints[endpoint]
	if !ok {
		return "", ErrSonos.ErrInvalidEndpoint
	}
	endpointBodyPrefix, ok := EndpointsBodyPrefix[endpoint]
	if !ok {
		return "", ErrSonos.ErrInvalidEndpoint
	}

	req, err := http.NewRequest(
		"POST",
		"http://"+zp.IpAddress.String()+":1400"+endpointUri,
		strings.NewReader(`<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/" s:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/"><s:Body><u:`+action+` xmlns:u="urn:schemas-upnp-org:service:`+endpoint+`:1">`+endpointBodyPrefix+body+`</u:`+action+`></s:Body></s:Envelope>`),
	)
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "text/xml")
	req.Header.Add("SOAPACTION", "urn:schemas-upnp-org:service:"+endpoint+":1#"+action)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	resultStr := string(result[:])

	if targetTag != "" {
		return extractTag(resultStr, targetTag)
	}

	if resultStr != `<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/" s:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/"><s:Body><u:`+action+`Response xmlns:u="urn:schemas-upnp-org:service:`+endpoint+`:1"></u:`+action+`Response></s:Body></s:Envelope>` {
		fmt.Print("\r\n" + resultStr)
		fmt.Print("\r\n" + `<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/" s:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/"><s:Body><u:` + action + `Response xmlns:u="urn:schemas-upnp-org:service:` + endpoint + `:1"></u:` + action + `Response></s:Body></s:Envelope>`)
		fmt.Print("\r\n")
		fmt.Print("\r\n")
		return resultStr, ErrSonos.ErrUnexpectedResponse
	}

	return resultStr, nil
}
