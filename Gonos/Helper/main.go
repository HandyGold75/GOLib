package Helper

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
	"Gonos/MusicServices"
	"Gonos/QPlay"
	"Gonos/Queue"
	"Gonos/RenderingControl"
	"Gonos/SystemProperties"
	"Gonos/VirtualLineIn"
	"Gonos/ZoneGroupTopology"
)

type Helper struct {
	alarmClock            *AlarmClock.AlarmClock
	audioIn               *AudioIn.AudioIn
	aVTransport           *AVTransport.AVTransport
	connectionManager     *ConnectionManager.ConnectionManager
	contentDirectory      *ContentDirectory.ContentDirectory
	deviceProperties      *DeviceProperties.DeviceProperties
	groupManagement       *GroupManagement.GroupManagement
	groupRenderingControl *GroupRenderingControl.GroupRenderingControl
	hTControl             *HTControl.HTControl
	musicServices         *MusicServices.MusicServices
	qPlay                 *QPlay.QPlay
	queue                 *Queue.Queue
	renderingControl      *RenderingControl.RenderingControl
	systemProperties      *SystemProperties.SystemProperties
	virtualLineIn         *VirtualLineIn.VirtualLineIn
	zoneGroupTopology     *ZoneGroupTopology.ZoneGroupTopology
}

func New(
	AlarmClock *AlarmClock.AlarmClock,
	AudioIn *AudioIn.AudioIn,
	AVTransport *AVTransport.AVTransport,
	ConnectionManager *ConnectionManager.ConnectionManager,
	ContentDirectory *ContentDirectory.ContentDirectory,
	DeviceProperties *DeviceProperties.DeviceProperties,
	GroupManagement *GroupManagement.GroupManagement,
	GroupRenderingControl *GroupRenderingControl.GroupRenderingControl,
	HTControl *HTControl.HTControl,
	MusicServices *MusicServices.MusicServices,
	QPlay *QPlay.QPlay,
	Queue *Queue.Queue,
	RenderingControl *RenderingControl.RenderingControl,
	SystemProperties *SystemProperties.SystemProperties,
	VirtualLineIn *VirtualLineIn.VirtualLineIn,
	ZoneGroupTopology *ZoneGroupTopology.ZoneGroupTopology,
) Helper {
	return Helper{
		alarmClock:            AlarmClock,
		audioIn:               AudioIn,
		aVTransport:           AVTransport,
		connectionManager:     ConnectionManager,
		contentDirectory:      ContentDirectory,
		deviceProperties:      DeviceProperties,
		groupManagement:       GroupManagement,
		groupRenderingControl: GroupRenderingControl,
		hTControl:             HTControl,
		musicServices:         MusicServices,
		qPlay:                 QPlay,
		queue:                 Queue,
		renderingControl:      RenderingControl,
		systemProperties:      SystemProperties,
		virtualLineIn:         VirtualLineIn,
		zoneGroupTopology:     ZoneGroupTopology,
	}
}
