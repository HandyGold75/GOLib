package Gonos

// TODO: Test
func (zp *ZonePlayer) GetZoneName() (string, error) {
	res, err := zp.GetZoneAttributes()
	return res.CurrentZoneName, err
}

// TODO: Test
func (zp *ZonePlayer) GetIcon() (string, error) {
	res, err := zp.GetZoneAttributes()
	return res.CurrentIcon, err
}

// TODO: Test
func (zp *ZonePlayer) GetConfiguration() (string, error) {
	res, err := zp.GetZoneAttributes()
	return res.CurrentConfiguration, err
}

// TODO: Test
func (zp *ZonePlayer) GetTargetRoomName() (string, error) {
	res, err := zp.GetZoneAttributes()
	return res.CurrentTargetRoomName, err
}

// TODO: Test
func (zp *ZonePlayer) SetZoneName(zoneName string) error {
	_, err := zp.SendDeviceProperties("SetZoneAttributes", "<DesiredZoneName>"+zoneName+"</DesiredZoneName>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) SetIcon(icon string) error {
	_, err := zp.SendDeviceProperties("SetZoneAttributes", "<DesiredIcon>"+icon+"</DesiredIcon>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) SetConfiguration(configuration string) error {
	_, err := zp.SendDeviceProperties("SetZoneAttributes", "<DesiredConfiguration>"+configuration+"</DesiredConfiguration>", "")
	return err
}

// TODO: Test
func (zp *ZonePlayer) SetTargetRoomName(targetRoomName string) error {
	_, err := zp.SendDeviceProperties("SetZoneAttributes", "<DesiredTargetRoomName>"+targetRoomName+"</DesiredTargetRoomName>", "")
	return err
}
