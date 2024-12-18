package Gonos

// Get current led state. (TODO: Test)
func (zp *ZonePlayer) GetLedState() (bool, error) {
	res, err := zp.sendCommand("DeviceProperties", "GetLEDState", "", "CurrentLEDState")
	return res == "On", err
}

// Set led state. (TODO: Test)
func (zp *ZonePlayer) SetLedState(state bool) error {
	_, err := zp.sendCommand("DeviceProperties", "SetLEDState", "<DesiredLEDState>"+boolToOnOff(state)+"</DesiredLEDState>", "")
	return err
}

// Get player name. (TODO: Test)
func (zp *ZonePlayer) GetPlayerName() (string, error) {
	return zp.sendCommand("DeviceProperties", "GetZoneAttributes", "", "CurrentZoneName")
}

// Set player name. (TODO: Test)
func (zp *ZonePlayer) SetPlayerName(name string) error {
	_, err := zp.sendCommand("DeviceProperties", "SetZoneAttributes", "<DesiredZoneName>"+name+"</DesiredZoneName>", "")
	return err
}

// Set player icon. (TODO: Test)
func (zp *ZonePlayer) SetPlayerIcon(icon string) error {
	_, err := zp.sendCommand("DeviceProperties", "SetZoneAttributes", "<DesiredIcon>"+icon+"</DesiredIcon/>", "")
	return err
}

// Set player config. (TODO: Test)
func (zp *ZonePlayer) SetPlayeConfig(config string) error {
	_, err := zp.sendCommand("DeviceProperties", "SetZoneAttributes", "<DesiredConfiguration>"+config+"</DesiredConfiguration>", "")
	return err
}
