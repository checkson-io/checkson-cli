package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stefan-hudelmaier/checkson-cli/output"
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
