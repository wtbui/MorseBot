package goveego

import (
	"go.uber.org/zap"
)

type Capability string 
const (
	ON     Capability  = "on"
	OFF    Capability  = "off"
	COLOR  Capability = "color"
	TEMP   Capability  = "temp"
	EFFECT Capability  = "effect"
)

var Capabilities = map[Capability]CapabilityIdentifier{
	OFF:    CapabilityIdentifier{"devices.capabilities.on_off", "powerSwitch"}, 
	ON:     CapabilityIdentifier{"devices.capabilities.on_off", "powerSwitch"},
	COLOR:  CapabilityIdentifier{"devices.capabilities.color_setting", "colorRgb"},
	TEMP:   CapabilityIdentifier{"devices.capabilities.color_setting", "colorTemperatureK"}, 
	EFFECT: CapabilityIdentifier{"devices.capabilities.dynamic_scene", "lightScene"},
}

type CapabilityIdentifier struct {
	Type 	string
	Instance string
}

func (gclient *GoveeClient) ChangeLight(device Device, capability Capability, capInput []int) error {
	zap.S().Debugw("Setting capability", "capability", capability, "value", capInput, "device", device.DeviceName,)
	err := gclient.UpdateDevice(device, Capabilities[capability], capInput)
	return err
}

func (gclient *GoveeClient) ChangeLightAll(capability Capability, capInput []int) error {
	for _, device := range gclient.Devices {
		err := gclient.ChangeLight(device, capability, capInput)
		if err != nil {
			return err
		}
	}

	return nil
}
