package goveego

import (
	//"fmt"
	"io/ioutil"
	"net/http"
	"encoding/json"
	//"os"
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

type Response struct {
	Code int `json:"code"`
	Message string `json:"message"`
	Devices []GDevice `json:"data"`
}

var (
	devUrl = "https://openapi.api.govee.com/router/api/v1/user/devices"
)

func Init(APIKey string) (*GoveeClient, error) {
	gclient := &GoveeClient{APIKey, []GDevice{}}
	
	// Create a new HTTP GET request
	req, err := http.NewRequest("GET", devUrl, nil)
	if err != nil {
		return nil, err
	}
	
	// Set the required headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Govee-API-Key", gclient.APIKey) 

	// Create an HTTP client and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read and print the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse the JSON response into the struct
	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	// Iterate over the data and get the "sku"
	for _, device := range response.Devices {
		if device.DeviceType == "devices.type.light" {
			gclient.Devices = append(gclient.Devices, device)
		}	
	}

	return gclient, nil
}



