package goveego

import (
	"io"
	"net/http"
	"encoding/json"
	"bytes"
	"strconv"
	"go.uber.org/zap"
	"fmt"
	"github.com/google/uuid"
)

type GoveeClient struct {
	APIKey string
	Devices []Device
}

type Device struct {
	SKU        string          `json:"sku"`
	DeviceAddr string          `json:"device"`
	DeviceName string          `json:"deviceName"`
	DeviceType string          `json:"type"`
}

type DeviceResponse struct {
	Code       int             `json:"code"`
	Message    string          `json:"message"`
	Devices    []Device        `json:"data"`
}

type ControlResponse struct {
	SKU        string          `json:"sku"`
}

type ControlReq struct {
	RequestId  string          `json:"requestId"`
	Payload    json.RawMessage `json:"payload"`
}

type CapabilityBody struct {
	Type       string          `json:"type"`
	Instance   string          `json:"instance"`
	Value      json.RawMessage `json:"value"`
}

const (
	deviceURL = "https://openapi.api.govee.com/router/api/v1/user/devices"
	controlURL = "https://openapi.api.govee.com/router/api/v1/device/control"
)

func NewClient(APIKey string) (*GoveeClient, error) {
	gclient := &GoveeClient{APIKey, []Device{}}
	
	zap.S().Debug("Establishing connection to GoveeAPI and getting devices with api-key: " + gclient.APIKey)
	clientInfo, err := makeRequest(gclient, deviceURL, "GET", nil)
	if err != nil {
		return nil, err
	}

	zap.S().Debug("Request from api recieved, grabbing response")
	var resp DeviceResponse
	err = json.Unmarshal(clientInfo, &resp)
	if err != nil {
		return nil, err
	}

	// Iterate over the data and get the "sku"
	for _, device := range resp.Devices {
		zap.S().Debugw("Found device", "name", device.DeviceName, "type", device.DeviceType)
		if device.DeviceType == "devices.types.light" {
			gclient.Devices = append(gclient.Devices, device)
		}	
	}
	zap.S().Debug("Devices recieved: " + strconv.Itoa(len(gclient.Devices)))

	return gclient, nil
}

func makeRequest(gclient *GoveeClient, url string, reqType string, reqBody []byte) ([]byte, error) {
	zap.S().Debug("Generating request: " + reqType + " " + url)
	var body io.Reader
	if reqBody != nil {
		body = bytes.NewBuffer(reqBody)
	}
	req, err := http.NewRequest(reqType, url, body)
	if err != nil {
		return nil, err
	}
	
	// Set the required headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Govee-API-Key", gclient.APIKey) 

	// Create an HTTP client and send the request
	zap.S().Debug("Making request...")
	zap.S().Debugw("http request", "method", req.Method, "url", req.URL.String())
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	zap.S().Debug("Reading response...")
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("http error %d: %s", resp.StatusCode, respBody)
	}
	zap.S().Debug("Recieved Response: " + string(respBody))

	return respBody, err
}

func (gclient *GoveeClient) UpdateDevice(device Device, capIdentity CapabilityIdentifier, capInput []int) error { 
	reqBody := CapabilityBody{
		Type: capIdentity.Type,
		Instance: capIdentity.Instance,
	}

	if len(capInput) == 1 {
		valueBytes, err := json.Marshal(capInput[0])
		if err != nil { 
			return err 
		}
		reqBody.Value = valueBytes
	} else if len(capInput) == 2 {
		valueBytes, err := json.Marshal(
			&struct{
				Id 		int `json:"id"`
				ParamId int `json:"paramId"`
			}{
				Id: capInput[0],
				ParamId: capInput[1],
			},
		)
		if err != nil { 
			return err 
		}
		reqBody.Value = valueBytes
	} else {
		return fmt.Errorf("Invalid capability input length: %d", len(capInput))
	}

	payloadBytes, err := json.Marshal(&struct {
		SKU        string          `json:"sku"`
		Device     string          `json:"device"`
		Capability CapabilityBody  `json:"capability"`
	}{
		SKU:        device.SKU,
		Device:     device.DeviceAddr,
		Capability: reqBody,
	})
	if err != nil { 
		return err 
	}

	reqJson := ControlReq{
		RequestId: uuid.NewString(),
		Payload: payloadBytes, 
	}
	
	reqBytes, err := json.Marshal(reqJson)
	if err != nil {
		return err
	}

	respBytes, err := makeRequest(gclient, controlURL, "POST", reqBytes)
	if err != nil {
		return err
	}

	var respJson ControlResponse
	err = json.Unmarshal(respBytes, &respJson)
	if err != nil {
		return err
	}

	return nil
}
