package gapo

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type (
	Tapo struct {
		ip              net.IP
		email, password string
		authHash        []byte

		httpClient    *http.Client
		handshakeData *handshakeData

		// Total delay to wait after handshakes.
		// Higher takes longer, lower is more unstable, a value of around 100 millisecond usually works.
		HandshakeDelay time.Duration
	}

	requestParams struct {
		DeviceOn bool `json:"device_on"`
	}
	request struct {
		Method          string         `json:"method"`
		RequestTimeMils int            `json:"requestTimeMils"`
		Params          *requestParams `json:"params,omitempty"`
	}

	DeviceInfo struct {
		DeviceID            string
		FwVer, HwVer        string
		Type, Model         string
		Mac                 string
		HwID, FwID          string
		OemID               string
		IP                  string
		TimeDiff            int
		Ssid                string
		Rssi                int
		SignalLevel         int
		AutoOffStatus       string
		AutoOffRemainTime   int
		Latitude, Longitude int
		Lang                string
		Avatar              string
		Region              string
		Specs               string
		Nickname            string
		HasSetLocationInfo  bool
		DeviceOn            bool
		OnTime              int
		DefaultStates       *struct {
			Type  string
			State *struct{}
		}
		Overheated                               bool
		PowerProtectionStatus, OvercurrentStatus string
	}
	EnergyUsage struct {
		TodayRuntime, MonthRuntime int
		TodayEnergy, MonthEnergy   int
		LocalTime                  string
		ElectricityCharge          []int
		CurrentPower               int
	}
	response struct {
		Result *struct {
			// DeviceInfo
			DeviceID           string `json:"device_id,omitempty"`
			FwVer              string `json:"fw_ver,omitempty"`
			HwVer              string `json:"hw_ver,omitempty"`
			Type               string `json:"type,omitempty"`
			Model              string `json:"model,omitempty"`
			Mac                string `json:"mac,omitempty"`
			HwID               string `json:"hw_id,omitempty"`
			FwID               string `json:"fw_id,omitempty"`
			OemID              string `json:"oem_id,omitempty"`
			IP                 string `json:"ip,omitempty"`
			TimeDiff           int    `json:"time_diff,omitempty"`
			Ssid               string `json:"ssid,omitempty"`
			Rssi               int    `json:"rssi,omitempty"`
			SignalLevel        int    `json:"signal_level,omitempty"`
			AutoOffStatus      string `json:"auto_off_status,omitempty"`
			AutoOffRemainTime  int    `json:"auto_off_remain_time,omitempty"`
			Latitude           int    `json:"latitude,omitempty"`
			Longitude          int    `json:"longitude,omitempty"`
			Lang               string `json:"lang,omitempty"`
			Avatar             string `json:"avatar,omitempty"`
			Region             string `json:"region,omitempty"`
			Specs              string `json:"specs,omitempty"`
			Nickname           string `json:"nickname,omitempty"`
			HasSetLocationInfo bool   `json:"has_set_location_info,omitempty"`
			DeviceOn           bool   `json:"device_on,omitempty"`
			OnTime             int    `json:"on_time,omitempty"`
			DefaultStates      *struct {
				Type  string    `json:"type,omitempty"`
				State *struct{} `json:"state,omitempty"`
			} `json:"default_states,omitempty"`
			Overheated            bool   `json:"overheated,omitempty"`
			PowerProtectionStatus string `json:"power_protection_status,omitempty"`
			OvercurrentStatus     string `json:"overcurrent_status,omitempty"`

			// EnergyUsage
			TodayRuntime      int    `json:"today_runtime,omitempty"`
			MonthRuntime      int    `json:"month_runtime,omitempty"`
			TodayEnergy       int    `json:"today_energy,omitempty"`
			MonthEnergy       int    `json:"month_energy,omitempty"`
			LocalTime         string `json:"local_time,omitempty"`
			ElectricityCharge []int  `json:"electricity_charge,omitempty"`
			CurrentPower      int    `json:"current_power,omitempty"`
		} `json:"result,omitempty"`
		ErrorCode int `json:"error_code"`
	}
)

// Create a new tapo session using email and password.
func NewTapo(ip, email, password string) (*Tapo, error) {
	t := &Tapo{
		ip:    net.ParseIP(ip),
		email: email, password: password,
		authHash: []byte{},

		httpClient:    &http.Client{Timeout: time.Second * 2},
		handshakeData: nil,

		HandshakeDelay: time.Millisecond * 100,
	}
	if err := t.handshake(); err != nil {
		return &Tapo{}, err
	}
	return t, nil
}

// Create a new tapo session using a auth hash.
//
// Auth hash: sha256(sha1(username)sha1(password))
func NewTapoHash(ip, authHash string) (*Tapo, error) {
	authHashBytes, err := hex.DecodeString(authHash)
	if err != nil {
		return &Tapo{}, err
	}
	t := &Tapo{
		ip:    net.ParseIP(ip),
		email: "", password: "",
		authHash: authHashBytes,

		httpClient:    &http.Client{Timeout: time.Second * 2},
		handshakeData: nil,

		HandshakeDelay: time.Millisecond * 100,
	}
	if err := t.handshake(); err != nil {
		return &Tapo{}, err
	}
	return t, nil
}

// Turn device on.
func (t *Tapo) On() (response, error) {
	res, err := t.doReq(&request{
		Method: "set_device_info", RequestTimeMils: int(time.Now().Unix()),
		Params: &requestParams{DeviceOn: true},
	})
	if err != nil {
		return response{}, err
	}
	return res, nil
}

// Turn device off.
func (t *Tapo) Off() (response, error) {
	res, err := t.doReq(&request{
		Method: "set_device_info", RequestTimeMils: int(time.Now().Unix()),
		Params: &requestParams{DeviceOn: false},
	})
	if err != nil {
		return response{}, err
	}
	return res, nil
}

// Get device info.
func (t *Tapo) GetDeviceInfo() (DeviceInfo, error) {
	res, err := t.doReq(&request{
		Method: "get_device_info", RequestTimeMils: int(time.Now().Unix()),
	})
	if err != nil {
		return DeviceInfo{}, err
	}
	return DeviceInfo{
		DeviceID: res.Result.DeviceID,
		FwVer:    res.Result.FwVer, HwVer: res.Result.HwVer,
		Type: res.Result.Type, Model: res.Result.Model,
		Mac:  res.Result.Mac,
		HwID: res.Result.HwID, FwID: res.Result.FwID,
		OemID:             res.Result.OemID,
		IP:                res.Result.IP,
		TimeDiff:          res.Result.TimeDiff,
		Ssid:              res.Result.Ssid,
		Rssi:              res.Result.Rssi,
		SignalLevel:       res.Result.SignalLevel,
		AutoOffStatus:     res.Result.AutoOffStatus,
		AutoOffRemainTime: res.Result.AutoOffRemainTime,
		Latitude:          res.Result.Latitude, Longitude: res.Result.Longitude,
		Lang:               res.Result.Lang,
		Avatar:             res.Result.Avatar,
		Region:             res.Result.Region,
		Specs:              res.Result.Specs,
		Nickname:           res.Result.Nickname,
		HasSetLocationInfo: res.Result.HasSetLocationInfo,
		DeviceOn:           res.Result.DeviceOn,
		OnTime:             res.Result.OnTime,
		DefaultStates: &struct {
			Type  string
			State *struct{}
		}{
			Type:  res.Result.DefaultStates.Type,
			State: res.Result.DefaultStates.State,
		},
		Overheated:            res.Result.Overheated,
		PowerProtectionStatus: res.Result.PowerProtectionStatus, OvercurrentStatus: res.Result.OvercurrentStatus,
	}, nil
}

// Get device info.
func (t *Tapo) GetEnergyUsage() (EnergyUsage, error) {
	res, err := t.doReq(&request{
		Method: "get_energy_usage", RequestTimeMils: int(time.Now().Unix()),
	})
	if err != nil {
		return EnergyUsage{}, err
	}
	return EnergyUsage{
		TodayRuntime: res.Result.TodayRuntime, MonthRuntime: res.Result.MonthRuntime,
		TodayEnergy: res.Result.TodayEnergy, MonthEnergy: res.Result.MonthEnergy,
		LocalTime:         res.Result.LocalTime,
		ElectricityCharge: res.Result.ElectricityCharge,
		CurrentPower:      res.Result.CurrentPower,
	}, nil
}

// Make a request to the device.
// If no session is active, tries to start a session.
func (t *Tapo) doReq(req *request) (response, error) {
	if t.handshakeData == nil {
		if err := t.handshake(); err != nil {
			return response{}, err
		}
	}

	dataJSON, err := json.Marshal(req)
	if err != nil {
		return response{}, err
	}
	dataEncrypted, seq, err := t.handshakeData.klapSession.encrypt(string(dataJSON))
	if err != nil {
		return response{}, err
	}
	u, err := url.Parse(fmt.Sprintf("http://%s/app/request?seq=%d", t.ip, seq))
	if err != nil {
		return response{}, err
	}
	reqHTTP, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(dataEncrypted))
	if err != nil {
		return response{}, err
	}
	reqHTTP.Header.Set("Content-Type", "application/json")
	for _, cookie := range t.handshakeData.Cookies {
		reqHTTP.AddCookie(cookie)
	}

	resHTTP, err := t.httpClient.Do(reqHTTP)
	if err != nil {
		return response{}, err
	}
	defer resHTTP.Body.Close()
	if resHTTP.StatusCode != 200 {
		return response{}, errors.New("unexpected status code: " + strconv.Itoa(resHTTP.StatusCode))
	}

	resHTTPBody, err := io.ReadAll(resHTTP.Body)
	if err != nil {
		return response{}, err
	}
	ret := response{}
	resHTTPBodyDecrypted, err := t.handshakeData.klapSession.decrypt(resHTTPBody)
	if err != nil {
		return response{}, err
	}
	if err = json.Unmarshal([]byte(resHTTPBodyDecrypted), &ret); err != nil {
		return response{}, err
	}
	return ret, nil
}
