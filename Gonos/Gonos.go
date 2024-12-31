package Gonos

import (
	"Gonos/AVTransport"
	"Gonos/AlarmClock"
	"Gonos/AudioIn"
	"Gonos/ConnectionManager"
	"Gonos/ContentDirectory"
	"Gonos/DeviceProperties"
	"Gonos/GroupManagement"
	"Gonos/GroupRenderingControl"
	"Gonos/HTControl"
	"Gonos/Helper"
	"Gonos/MusicServices"
	"Gonos/QPlay"
	"Gonos/Queue"
	"Gonos/RenderingControl"
	"Gonos/SystemProperties"
	"Gonos/VirtualLineIn"
	"Gonos/ZoneGroupTopology"
	"Gonos/lib"
	"net"
	"strings"
	"sync"
	"time"
)

type (
	ZonePlayer struct {
		// Full url address packets will be send to.
		URL string
		// GetZoneInfo call is made to confirm if the requested ZonePlayer exists opon creation, might as well store the returned data.
		ZoneInfo DeviceProperties.GetZoneInfoResponse

		// Contains helper functions.
		//
		// Prefer functions present in zp.H over the functions in Sonos SOAP Services.
		H Helper.Helper

		// Sonos SOAP Service `AlarmClock` (State: To Implement).
		//
		// Prefer functions present in zp.H over the functions in Sonos SOAP Services.
		//
		// Control the sonos alarms and times.
		AlarmClock AlarmClock.AlarmClock

		// Sonos SOAP Service `AudioIn` (State: To Implement).
		//
		// Prefer functions present in zp.H over the functions in Sonos SOAP Services.
		//
		// Control line in.
		AudioIn AudioIn.AudioIn

		// Sonos SOAP Service `AVTransport`.
		//
		// Prefer functions present in zp.H over the functions in Sonos SOAP Services.
		//
		// Service that controls stuff related to transport (play/pause/next/special URLs).
		AVTransport AVTransport.AVTransport

		// Sonos SOAP Service `ConnectionManager` (State: To Implement).
		//
		// Prefer functions present in zp.H over the functions in Sonos SOAP Services.
		//
		// Services related to connections and protocols.
		ConnectionManager ConnectionManager.ConnectionManager

		// Sonos SOAP Service `ContentDirectory`.
		//
		// Prefer functions present in zp.H over the functions in Sonos SOAP Services.
		//
		// Browse for local content.
		ContentDirectory ContentDirectory.ContentDirectory

		// Sonos SOAP Service `DeviceProperties`.
		//
		// Prefer functions present in zp.H over the functions in Sonos SOAP Services.
		//
		// Modify device properties, like LED status and stereo pairs.
		DeviceProperties DeviceProperties.DeviceProperties

		// Sonos SOAP Service `GroupManagement` (State: To Implement).
		//
		// Prefer functions present in zp.H over the functions in Sonos SOAP Services.
		//
		// Services related to groups.
		GroupManagement GroupManagement.GroupManagement

		// Sonos SOAP Service `GroupRenderingControl` (State: To Implement).
		//
		// Prefer functions present in zp.H over the functions in Sonos SOAP Services.
		//
		// Volume related controls for groups.
		GroupRenderingControl GroupRenderingControl.GroupRenderingControl

		// Sonos SOAP Service `HTControl` (State: To Implement).
		//
		// Prefer functions present in zp.H over the functions in Sonos SOAP Services.
		//
		// Service related to the TV remote control.
		HTControl HTControl.HTControl

		// Sonos SOAP Service `MusicServices` (State: To Implement).
		//
		// Prefer functions present in zp.H over the functions in Sonos SOAP Services.
		//
		// Access to external music services, like Spotify or Youtube Music.
		MusicServices MusicServices.MusicServices

		// Sonos SOAP Service `QPlay` (State: To Implement).
		//
		// Prefer functions present in zp.H over the functions in Sonos SOAP Services.
		//
		// Services related to Chinese Tencent Qplay service.
		QPlay QPlay.QPlay

		// Sonos SOAP Service `Queue` (State: To Implement).
		//
		// Prefer functions present in zp.H over the functions in Sonos SOAP Services.
		//
		// Modify and browse queues.
		Queue Queue.Queue

		// Sonos SOAP Service `RenderingControl`.
		//
		// Prefer functions present in zp.H over the functions in Sonos SOAP Services.
		//
		// Volume related controls.
		RenderingControl RenderingControl.RenderingControl

		// Sonos SOAP Service `SystemProperties` (State: To Implement).
		//
		// Prefer functions present in zp.H over the functions in Sonos SOAP Services.
		//
		// Manage system-wide settings, mainly account stuff.
		SystemProperties SystemProperties.SystemProperties

		// Sonos SOAP Service `VirtualLineIn` (State: To Implement).
		//
		// Prefer functions present in zp.H over the functions in Sonos SOAP Services.
		VirtualLineIn VirtualLineIn.VirtualLineIn

		// Sonos SOAP Service `ZoneGroupTopology` (State: To Implement).
		//
		// Prefer functions present in zp.H over the functions in Sonos SOAP Services.
		//
		// Zone config stuff, eg getting all the configured sonos zones.
		ZoneGroupTopology ZoneGroupTopology.ZoneGroupTopology
	}
)

// Create new ZonePlayer for controling a Sonos speaker.
func NewZonePlayer(ipAddress string) (*ZonePlayer, error) {
	ip := net.ParseIP(ipAddress)
	if ip == nil {
		return &ZonePlayer{}, lib.ErrSonos.ErrInvalidIPAdress
	}

	zp := &ZonePlayer{URL: "http://" + ip.String() + ":1400"}
	zp.AlarmClock = AlarmClock.New(zp.SendAlarmClock)
	zp.AudioIn = AudioIn.New(zp.SendAudioIn)
	zp.AVTransport = AVTransport.New(zp.SendAVTransport)
	zp.ConnectionManager = ConnectionManager.New(zp.SendConnectionManager)
	zp.ContentDirectory = ContentDirectory.New(zp.SendContentDirectory)
	zp.DeviceProperties = DeviceProperties.New(zp.SendDeviceProperties)
	zp.GroupManagement = GroupManagement.New(zp.SendGroupManagement)
	zp.GroupRenderingControl = GroupRenderingControl.New(zp.SendGroupRenderingControl)
	zp.HTControl = HTControl.New(zp.SendHTControl)
	zp.MusicServices = MusicServices.New(zp.SendMusicServices)
	zp.QPlay = QPlay.New(zp.SendQPlay)
	zp.Queue = Queue.New(zp.SendQueue)
	zp.RenderingControl = RenderingControl.New(zp.SendRenderingControl)
	zp.SystemProperties = SystemProperties.New(zp.SendSystemProperties)
	zp.VirtualLineIn = VirtualLineIn.New(zp.SendVirtualLineIn)
	zp.ZoneGroupTopology = ZoneGroupTopology.New(zp.SendZoneGroupTopology)

	zp.H = Helper.New(
		&zp.AlarmClock,
		&zp.AudioIn,
		&zp.AVTransport,
		&zp.ConnectionManager,
		&zp.ContentDirectory,
		&zp.DeviceProperties,
		&zp.GroupManagement,
		&zp.GroupRenderingControl,
		&zp.HTControl,
		&zp.MusicServices,
		&zp.QPlay,
		&zp.Queue,
		&zp.RenderingControl,
		&zp.SystemProperties,
		&zp.VirtualLineIn,
		&zp.ZoneGroupTopology,
	)

	info, err := zp.DeviceProperties.GetZoneInfo()
	if err != nil {
		return &ZonePlayer{}, lib.ErrSonos.ErrNoZonePlayerFound
	}
	zp.ZoneInfo = info
	return zp, nil
}

// TODO: Test
//
// Create new ZonePlayer using discovery controling a Sonos speaker.
//
// `timout` of 1 second is recomended.
func DiscoverZonePlayer(timeout time.Duration) ([]*ZonePlayer, error) {
	conn, err := net.DialUDP("udp", &net.UDPAddr{Port: 1900}, &net.UDPAddr{IP: net.IPv4(239, 255, 255, 250), Port: 1900})
	if err != nil {
		return []*ZonePlayer{}, err
	}
	defer conn.Close()

	for i := 0; i < 3; i++ {
		_, _ = conn.Write([]byte("M-SEARCH * HTTP/1.1\r\nHOST: 239.255.255.250:1900\r\nMAN: \"ssdp:discover\"\r\nMX: 1\r\nST: urn:schemas-upnp-org:device:ZonePlayer:1\r\n\r\n"))
	}

	zps := []*ZonePlayer{}
	for {
		buf := make([]byte, 1024)
		conn.SetDeadline(time.Now().Add(timeout))
		_, addr, err := conn.ReadFrom(buf)
		if err.(*net.OpError).Timeout() {
			break
		} else if err != nil {
			return zps, err
		}

		zp, err := NewZonePlayer(strings.Split(addr.String(), ":")[0])
		if err != nil {
			continue
		}
		zps = append(zps, zp)
	}

	if len(zps) <= 0 {
		return zps, lib.ErrSonos.ErrNoZonePlayerFound
	}
	return zps, nil
}

// Create new ZonePlayer using network scanning controling a Sonos speaker.
//
// `timout` of 1 second is recomended.
func ScanZonePlayer(cidr string, timeout time.Duration) ([]*ZonePlayer, error) {
	incIP := func(ip net.IP) {
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
			zps = append(zps, zp)
		}(ip.String())
	}
	wg.Wait()

	if len(zps) <= 0 {
		return zps, lib.ErrSonos.ErrNoZonePlayerFound
	}
	return zps, nil
}

func (zp *ZonePlayer) SendAVTransport(action, body, targetTag string) (string, error) {
	return lib.SendAndVerify(zp.URL+"/MediaRenderer/AVTransport/Control", "AVTransport", action, "<InstanceID>0</InstanceID>"+body, targetTag)
}

func (zp *ZonePlayer) SendAlarmClock(action, body, targetTag string) (string, error) {
	return lib.SendAndVerify(zp.URL+"/AlarmClock/Control", "AlarmClock", action, body, targetTag)
}

func (zp *ZonePlayer) SendAudioIn(action, body, targetTag string) (string, error) {
	return lib.SendAndVerify(zp.URL+"/AudioIn/Control", "AudioIn", action, body, targetTag)
}

func (zp *ZonePlayer) SendConnectionManager(action, body, targetTag string) (string, error) {
	return lib.SendAndVerify(zp.URL+"/MediaRenderer/ConnectionManager/Control", "ConnectionManager", action, body, targetTag)
}

func (zp *ZonePlayer) SendContentDirectory(action, body, targetTag string) (string, error) {
	return lib.SendAndVerify(zp.URL+"/MediaServer/ContentDirectory/Control", "ContentDirectory", action, body, targetTag)
}

func (zp *ZonePlayer) SendDeviceProperties(action, body, targetTag string) (string, error) {
	return lib.SendAndVerify(zp.URL+"/DeviceProperties/Control", "DeviceProperties", action, body, targetTag)
}

func (zp *ZonePlayer) SendGroupManagement(action, body, targetTag string) (string, error) {
	return lib.SendAndVerify(zp.URL+"/GroupManagement/Control", "GroupManagement", action, body, targetTag)
}

func (zp *ZonePlayer) SendGroupRenderingControl(action, body, targetTag string) (string, error) {
	return lib.SendAndVerify(zp.URL+"/MediaRenderer/GroupRenderingControl/Control", "GroupRenderingControl", action, body, targetTag)
}

func (zp *ZonePlayer) SendHTControl(action, body, targetTag string) (string, error) {
	return lib.SendAndVerify(zp.URL+"/HTControl/Control", "HTControl", action, body, targetTag)
}

func (zp *ZonePlayer) SendMusicServices(action, body, targetTag string) (string, error) {
	return lib.SendAndVerify(zp.URL+"/MusicServices/Control", "MusicServices", action, body, targetTag)
}

func (zp *ZonePlayer) SendQPlay(action, body, targetTag string) (string, error) {
	return lib.SendAndVerify(zp.URL+"/QPlay/Control", "QPlay", action, body, targetTag)
}

func (zp *ZonePlayer) SendQueue(action, body, targetTag string) (string, error) {
	return lib.SendAndVerify(zp.URL+"/MediaRenderer/Queue/Control", "Queue", action, body, targetTag)
}

func (zp *ZonePlayer) SendRenderingControl(action, body, targetTag string) (string, error) {
	return lib.SendAndVerify(zp.URL+"/MediaRenderer/RenderingControl/Control", "RenderingControl", action, "<InstanceID>0</InstanceID>"+body, targetTag)
}

func (zp *ZonePlayer) SendSystemProperties(action, body, targetTag string) (string, error) {
	return lib.SendAndVerify(zp.URL+"/SystemProperties/Control", "SystemProperties", action, body, targetTag)
}

func (zp *ZonePlayer) SendVirtualLineIn(action, body, targetTag string) (string, error) {
	return lib.SendAndVerify(zp.URL+"/MediaRenderer/VirtualLineIn/Control", "VirtualLineIn", action, body, targetTag)
}

func (zp *ZonePlayer) SendZoneGroupTopology(action, body, targetTag string) (string, error) {
	return lib.SendAndVerify(zp.URL+"/ZoneGroupTopology/Control", "ZoneGroupTopology", action, body, targetTag)
}
