package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stefan-hudelmaier/checkson-cli/output"
	"io/ioutil"
	"net/http"
)

func CreateChannel(channel NotificationChannel, authToken string, devMode bool) error {

	client := &http.Client{}

	jsonBytes, jsonErr := json.Marshal(channel)
	output.Debugf("Sending:", string(jsonBytes))
	if jsonErr != nil {
		return errors.New("cannot serialize channel")
	}

	url := getApiUrl(devMode, "api/notification-channels/")
	req, err := http.NewRequest("PUT", url+channel.Name, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return fmt.Errorf("problem preparing request: %w", err)
	}
	addHeaders(req, authToken)

	resp, err1 := client.Do(req)
	if err1 != nil {
		return fmt.Errorf("problem performing request: %w", err1)
	}
	defer resp.Body.Close()

	return handleRestResponse("Notification channel", resp)
}

func GetChannel(channelName string, authToken string, devMode bool) (*NotificationChannel, error) {

	client := &http.Client{}

	url := getApiUrl(devMode, "api/notification-channels/")
	req, err := http.NewRequest("GET", url+channelName, nil)
	if err != nil {
		return nil, fmt.Errorf("problem preparing request: %w", err)
	}
	addHeaders(req, authToken)

	resp, err1 := client.Do(req)
	if err1 != nil {
		return nil, fmt.Errorf("problem performing request: %w", err1)
	}
	defer resp.Body.Close()

	err2 := handleRestResponse("Notification channel", resp)
	if err2 != nil {
		return nil, err2
	}

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return nil, readErr
	}

	var channel NotificationChannel
	jsonErr := json.Unmarshal(body, &channel)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return &channel, nil
}

func ListChannels(authToken string, devMode bool) ([]NotificationChannel, error) {

	client := &http.Client{}

	url := getApiUrl(devMode, "api/notification-channels")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("problem preparing request: %w", err)
	}
	addHeaders(req, authToken)

	resp, err1 := client.Do(req)
	if err1 != nil {
		return nil, fmt.Errorf("problem performing request: %w", err1)
	}
	defer resp.Body.Close()

	err2 := handleRestResponse("Notification channel", resp)
	if err2 != nil {
		return nil, err2
	}

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return nil, readErr
	}

	var channels []NotificationChannel
	jsonErr := json.Unmarshal(body, &channels)
	if jsonErr != nil {
		return nil, fmt.Errorf("unable to parse value: %q, error: %s", string(body), jsonErr.Error())
	}

	return channels, nil
}

func DeleteChannel(channelName string, authToken string, devMode bool) error {

	client := &http.Client{}

	url := getApiUrl(devMode, "api/notification-channels/")
	req, err := http.NewRequest("DELETE", url+channelName, nil)
	if err != nil {
		return fmt.Errorf("problem preparing request: %w", err)
	}
	addHeaders(req, authToken)

	resp, err1 := client.Do(req)
	if err1 != nil {
		return fmt.Errorf("problem performing request: %w", err1)
	}
	defer resp.Body.Close()

	return handleRestResponse("Notification channel", resp)
}
