package goveego

import (
	"go.uber.org/zap"
	"context"
	"golang.org/x/sync/errgroup"
)

type Capability string 
const (
	ON     Capability  = "on"
	OFF    Capability  = "off"
	COLOR  Capability  = "color"
	TEMP   Capability  = "temp"
	EFFECT Capability  = "effect"
	BRIGHT Capability  = "bright" 
)

var Capabilities = map[Capability]CapabilityIdentifier{
	OFF:    CapabilityIdentifier{"devices.capabilities.on_off", "powerSwitch"}, 
	ON:     CapabilityIdentifier{"devices.capabilities.on_off", "powerSwitch"},
	COLOR:  CapabilityIdentifier{"devices.capabilities.color_setting", "colorRgb"},
	TEMP:   CapabilityIdentifier{"devices.capabilities.color_setting", "colorTemperatureK"}, 
	EFFECT: CapabilityIdentifier{"devices.capabilities.dynamic_scene", "lightScene"},
	BRIGHT: CapabilityIdentifier{"devices.capabilities.range", "brightness"},
}

type CapabilityIdentifier struct {
	Type     string
	Instance string
}

func (gclient *GoveeClient) ChangeLight(ctx context.Context, device Device, capability Capability, capInput []int) error {
	zap.S().Debugw("Setting capability", "capability", capability, "value", capInput, "device", device.DeviceName,)
	err := gclient.UpdateDevice(device, Capabilities[capability], capInput)
	return err
}

func (gclient *GoveeClient) ChangeLightAll(ctx context.Context, capability Capability, capInput []int) error {
	group, ctx := errgroup.WithContext(ctx) 
	for _, device := range gclient.Devices {
		device := device
		group.Go(func() error {
			return gclient.ChangeLight(ctx, device, capability, capInput)
		})
	}
	
	err := group.Wait()
	return err
}
