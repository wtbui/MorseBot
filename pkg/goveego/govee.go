package goveego

import (
	"io"
	"net/http"
	"encoding/json"
	"bytes"
	"strconv"
	"go.uber.org/zap"
)

type GoveeClient struct {
	APIKey string
	Devices []GDevice
}

type GDevice struct {
	SKU string `json:"sku"`
	DeviceAddr string `json:"device"`
	DeviceName string `json:"deviceName"`
	DeviceType string `json:"type"`
}

type GDevResponse struct {
	Code int `json:"code"`
	Message string `json:"message"`
	Devices []GDevice `json:"data"`
}

type GConResponse struct {
	SKU string `json:"sku"`
}

type GConBody struct {
	RequestId string `json:"requestId"`
	Payload GConBodyPayload `json:"payload"`
}

type GConBodyPayload struct {
	SKU string `json:"sku"`
	Device string `json:"device"`
	Capability GConBodyCap `json:"capability"`
}

type GConBodyCapabilitySingle struct {
	Type string `json:"type"`
	Instance string `json:"instance"`
	Value int `json:"value"` 
}

type GConBodyCapabilityEffect struct {
	Type string `json:"type"`
	Instance string `json:"instance"`
	Value GConCapEffectValue `json:"value"`
}

type GConCapEffectValue struct {
	Id int `json:"id"`
	ParamId int `json:"paramId"`
}

type GConBodyCap interface {
	GetType() string
}

func (gcap GConBodyCapabilitySingle) GetType() string {
	return gcap.Type
}

func (gcap GConBodyCapabilityEffect) GetType() string {
	return gcap.Type
}

type CapData struct {
	Type string
	Value []int
}

var (
	devUrl = "https://openapi.api.govee.com/router/api/v1/user/devices"
	controlUrl = "https://openapi.api.govee.com/router/api/v1/device/control"
)

func Init(APIKey string) (*GoveeClient, error) {
	gclient := &GoveeClient{APIKey, []GDevice{}}
	
	zap.S().Debug("Establishing connection to GoveeAPI and getting devices with api-key: " + gclient.APIKey)
	clientInfo, err := makeRequest(gclient, devUrl, "GET", nil)
	if err != nil {
		return nil, err
	}

	zap.S().Debug("Request from api recieved, grabbing response")
	// Parse the JSON response into the struct
	var response GDevResponse
	err = json.Unmarshal(clientInfo, &response)
	if err != nil {
		return nil, err
	}

	// Iterate over the data and get the "sku"
	for _, device := range response.Devices {
		zap.S().Debug("Found device " + device.DeviceName + " with device type " + device.DeviceType)
		if device.DeviceType == "devices.types.light" {
			gclient.Devices = append(gclient.Devices, device)
		}	
	}
	zap.S().Debug("Devices recieved: " + strconv.Itoa(len(gclient.Devices)))

	return gclient, nil
}

func makeRequest(gclient *GoveeClient, url string, reqType string, reqBody []byte) ([]byte, error) {
	// Create a new HTTP request
	zap.S().Debug("Generating request: " + reqType + " " + url)
	
	req, err := http.NewRequest(reqType, url, bytes.NewBuffer(reqBody))
	if reqBody == nil {
		req, err = http.NewRequest(reqType, url, nil)
	}

	if err != nil {
		return nil, err
	}
	
	// Set the required headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Govee-API-Key", gclient.APIKey) 

	// Create an HTTP client and send the request
	zap.S().Debug("Making request...")
	zap.S().Debug(req)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	zap.S().Debug("Reading response...")
	respInfo, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	zap.S().Debug("Recieved Response: " + string(respInfo))

	return respInfo, err
}

func (gclient GoveeClient) UpdateDevice(device GDevice, capType string, capInst string, capData CapData) error { 
	var cap GConBodyCap

	if capData.Type == "single" {
		cap = GConBodyCapabilitySingle{
			Type: capType,
			Instance: capInst,
			Value: capData.Value[0],
		}
	} else if capData.Type == "effect" {
		cap = GConBodyCapabilityEffect{
			Type: capType,
			Instance: capInst,
			Value: GConCapEffectValue{
				Id: capData.Value[1],
				ParamId: capData.Value[0],
			},
		}
	}

	reqBody := GConBody{
		RequestId: "uuid",
		Payload: GConBodyPayload{
			SKU: device.SKU,
			Device: device.DeviceAddr,
			Capability: cap,
		},
	}
	
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	clientInfo, err := makeRequest(&gclient, controlUrl, "POST", jsonData)
	if err != nil {
		return err
	}

	var response GConResponse
	err = json.Unmarshal(clientInfo, &response)
	if err != nil {
		return err
	}

	return nil
}

func (gclient GoveeClient) TurnOnOff(device GDevice, value int) error {
	zap.S().Debug("Setting light powerswitch value to " + strconv.Itoa(value) + " for device " + device.DeviceName)
	err := gclient.UpdateDevice(device, "devices.capabilities.on_off", "powerSwitch", CapData{"single", []int{value}})
	return err
}

func (gclient GoveeClient) TurnOnOffAll(value int) error {
	for _, device := range gclient.Devices {
		err := gclient.TurnOnOff(device, value)
		if err != nil {
			return err
		}
	}

	return nil
}

func (gclient GoveeClient) ChangeColor(device GDevice, value int) error {
	zap.S().Debug("Setting light colorId value to " + strconv.Itoa(value) + " for device " + device.DeviceName)
	err := gclient.UpdateDevice(device, "devices.capabilities.color_setting", "colorRgb", CapData{"single", []int{value}})
	return err
}

func (gclient GoveeClient) ChangeColorAll(value int) error {
	for _, device := range gclient.Devices {
		err := gclient.ChangeColor(device, value)
		if err != nil {
			return err
		}
	}

	return nil
}

func (gclient GoveeClient) ChangeTemp(device GDevice, value int) error {
	zap.S().Debug("Setting light colorTemperatureK value to " + strconv.Itoa(value) + " for device " + device.DeviceName)
	err := gclient.UpdateDevice(device, "devices.capabilities.color_setting", "colorTemperatureK", CapData{"single", []int{value}})
	return err
}

func (gclient GoveeClient) ChangeTempAll(value int) error {
	for _, device := range gclient.Devices {
		err := gclient.ChangeTemp(device, value)
		if err != nil {
			return err
		}
	}

	return nil
}

func (gclient GoveeClient) ChangeEffect(device GDevice, paramId int, id int) error {
	zap.S().Debug("Setting light lightscene value to " + strconv.Itoa(id) + " for device " + device.DeviceName)
	err := gclient.UpdateDevice(device, "devices.capabilities.dynamic_scene", "lightScene", CapData{"effect", []int{paramId, id}})
	return err
}

func (gclient GoveeClient) ChangeEffectAll(paramId int, id int) error {
	for _, device := range gclient.Devices {
		err := gclient.ChangeEffect(device, paramId, id)
		if err != nil {
			return err
		}
	}

	return nil
}



