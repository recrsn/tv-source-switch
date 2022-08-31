package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

var baseURL = "https://api.smartthings.com/v1/devices/%s/%s"

type SmartThingsTVClient struct {
	token    string
	deviceID string
	client   *http.Client
}

func NewSmartThingsTVClient(token string, deviceID string) *SmartThingsTVClient {
	return &SmartThingsTVClient{
		token:    token,
		deviceID: deviceID,
		client: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

func (s *SmartThingsTVClient) Close() {
	if s.client != nil {
		s.client.CloseIdleConnections()
	}
	s.client = nil
}

func (s *SmartThingsTVClient) request(method string, url string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, url, body)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.token))
	response, err := s.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "Error making request")
	}

	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, errors.Wrap(err, "Error reading response body")
	}

	if response.StatusCode != 200 {
		return nil, errors.Errorf("Error making request: %s", response.Status)
	}

	return responseBody, nil
}

func (s *SmartThingsTVClient) get(url string) ([]byte, error) {
	return s.request("GET", url, nil)
}

func (s *SmartThingsTVClient) post(url string, body []byte) ([]byte, error) {
	return s.request("POST", url, bytes.NewReader(body))
}

type DeviceStatus struct {
	DeviceID        string `json:"deviceId"`
	State           string `json:"state"`
	LastUpdatedDate string `json:"lastUpdatedDate"`
}

func (s *SmartThingsTVClient) GetStatus() (DeviceStatus, error) {
	url := fmt.Sprintf(baseURL, s.deviceID, "health")
	body, err := s.get(url)

	if err != nil {
		return DeviceStatus{}, errors.Wrap(err, "Error getting status")
	}

	var status DeviceStatus
	err = json.Unmarshal(body, &status)
	if err != nil {
		return DeviceStatus{}, errors.Wrap(err, "Error unmarshalling status")
	}
	return status, nil
}

type Command struct {
	Component  string   `json:"component"`
	Capability string   `json:"capability"`
	Command    string   `json:"command"`
	Arguments  []string `json:"arguments"`
}

type Result struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

func (s *SmartThingsTVClient) RunCommand(command Command) (Result, error) {
	url := fmt.Sprintf(baseURL, s.deviceID, "commands")

	requestBody, err := json.Marshal(map[string]interface{}{
		"commands": []Command{command},
	})

	if err != nil {
		return Result{}, errors.Wrap(err, "Error marshalling command")
	}

	body, err := s.post(url, requestBody)

	if err != nil {
		return Result{}, errors.Wrap(err, "Error executing command")
	}

	var response struct {
		Results []Result `json:"results"`
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return Result{}, errors.Wrap(err, "Error unmarshalling result")
	}

	return response.Results[0], nil
}

func (s *SmartThingsTVClient) SetPower(state string) (Result, error) {
	return s.RunCommand(Command{
		Component:  "main",
		Capability: "switch",
		Command:    state,
	})
}

func (s *SmartThingsTVClient) SetSource(source string) (Result, error) {
	return s.RunCommand(Command{
		Component:  "main",
		Capability: "mediaInputSource",
		Command:    "setInputSource",
		Arguments:  []string{source},
	})
}
