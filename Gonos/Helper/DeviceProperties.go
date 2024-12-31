package Helper

func (h *Helper) GetZoneName() (string, error) {
	res, err := h.deviceProperties.GetZoneAttributes()
	return res.CurrentZoneName, err
}

func (h *Helper) GetIcon() (string, error) {
	res, err := h.deviceProperties.GetZoneAttributes()
	return res.CurrentIcon, err
}

func (h *Helper) GetConfiguration() (string, error) {
	res, err := h.deviceProperties.GetZoneAttributes()
	return res.CurrentConfiguration, err
}

func (h *Helper) GetTargetRoomName() (string, error) {
	res, err := h.deviceProperties.GetZoneAttributes()
	return res.CurrentTargetRoomName, err
}

func (h *Helper) SetZoneName(zoneName string) error {
	_, err := h.deviceProperties.Send("SetZoneAttributes", "<DesiredZoneName>"+zoneName+"</DesiredZoneName>", "")
	return err
}

func (h *Helper) SetIcon(icon string) error {
	_, err := h.deviceProperties.Send("SetZoneAttributes", "<DesiredIcon>"+icon+"</DesiredIcon>", "")
	return err
}

func (h *Helper) SetConfiguration(configuration string) error {
	_, err := h.deviceProperties.Send("SetZoneAttributes", "<DesiredConfiguration>"+configuration+"</DesiredConfiguration>", "")
	return err
}

func (h *Helper) SetTargetRoomName(targetRoomName string) error {
	_, err := h.deviceProperties.Send("SetZoneAttributes", "<DesiredTargetRoomName>"+targetRoomName+"</DesiredTargetRoomName>", "")
	return err
}

// Short for `zp.DeviceProperties.GetLEDState`.
func (h *Helper) GetLED() (bool, error) {
	return h.deviceProperties.GetLEDState()
}

// Short for `zp.DeviceProperties.SetLEDState`.
func (h *Helper) SetLED(state bool) error {
	return h.deviceProperties.SetLEDState(state)
}
