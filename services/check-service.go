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

func CreateCheck(check Check, authToken string, devMode bool) error {

	client := &http.Client{}

	jsonBytes, jsonErr := json.Marshal(check)
	output.Debugf("Sending:", string(jsonBytes))
	if jsonErr != nil {
		return errors.New("Cannot serialize check")
	}

	url := getApiUrl(devMode, "api/checks/")
	req, err := http.NewRequest("PUT", url+check.Name, bytes.NewBuffer(jsonBytes))

	if err != nil {
		return fmt.Errorf("problem preparing request: %w", err)
	}
	addHeaders(req, authToken)

	resp, err1 := client.Do(req)
	if err1 != nil {
		return fmt.Errorf("problem performing request: %w", err1)
	}
	defer resp.Body.Close()

	return handleRestResponse("Check", resp)
}

func DeleteCheck(checkName string, authToken string, devMode bool) error {

	client := &http.Client{}

	req, err := http.NewRequest("DELETE", getApiUrl(devMode, "api/checks/")+checkName, nil)
	if err != nil {
		return fmt.Errorf("problem preparing request: %w", err)
	}
	addHeaders(req, authToken)

	resp, err1 := client.Do(req)
	if err1 != nil {
		return fmt.Errorf("problem performing request: %w", err1)
	}
	defer resp.Body.Close()

	return handleRestResponse("Check", resp)
}

func ListChecks(authToken string, devMode bool) ([]Check, error) {

	client := &http.Client{}
	req, err1 := http.NewRequest("GET", getApiUrl(devMode, "api/checks"), nil)

	if err1 != nil {
		return nil, fmt.Errorf("problem preparing request: %w", err1)
	}
	addHeaders(req, authToken)

	resp, err2 := client.Do(req)
	if err2 != nil {
		return nil, err2
	}
	defer resp.Body.Close()

	err3 := handleRestResponse("Checks", resp)
	if err3 != nil {
		return nil, err3
	}

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		panic(readErr)
	}

	var checks []Check
	jsonErr := json.Unmarshal(body, &checks)
	if jsonErr != nil {
		return nil, fmt.Errorf("unable to parse value: %q, error: %s", string(body), jsonErr.Error())
	}

	return checks, nil
}
