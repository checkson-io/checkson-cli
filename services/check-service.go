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
		return errors.New("cannot serialize check")
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
	req, err1 := http.NewRequest("GET", getApiUrl(devMode, "api/checks-full"), nil)

	if err1 != nil {
		return nil, fmt.Errorf("problem preparing request: %w", err1)
	}
	addHeaders(req, authToken)

	resp, err2 := client.Do(req)
	if err2 != nil {
		return nil, fmt.Errorf("problem performing request: %w", err2)
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

	output.Debugf("Received status: %s", body)

	var checks []Check
	jsonErr := json.Unmarshal(body, &checks)
	if jsonErr != nil {
		return nil, fmt.Errorf("unable to parse value: %q, error: %s", string(body), jsonErr.Error())
	}

	return checks, nil
}

func ListRuns(authToken string, checkName string, devMode bool) ([]Run, error) {

	client := &http.Client{}
	path := "api/finished-runs"

	if checkName != "" {
		path = "api/checks/" + checkName + "/finished-runs"
	}
	req, err1 := http.NewRequest("GET", getApiUrl(devMode, path), nil)
	if err1 != nil {
		return nil, fmt.Errorf("problem preparing request: %w", err1)
	}
	addHeaders(req, authToken)

	resp, err2 := client.Do(req)
	if err2 != nil {
		return nil, fmt.Errorf("problem performing request: %w", err2)
	}
	defer resp.Body.Close()

	err3 := handleRestResponse("Runs", resp)
	if err3 != nil {
		return nil, err3
	}

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return nil, errors.New("cannot parse result")
	}

	var runs []Run
	jsonErr := json.Unmarshal(body, &runs)
	if jsonErr != nil {
		return nil, fmt.Errorf("unable to parse value: %q, error: %s", string(body), jsonErr.Error())
	}

	return runs, nil
}

func GetLog(checkName string, runId string, authToken string, devMode bool) (string, error) {

	path := fmt.Sprintf("api/checks/%s/runs/%s/log", checkName, runId)
	client := &http.Client{}
	req, err1 := http.NewRequest("GET", getApiUrl(devMode, path), nil)
	if err1 != nil {
		return "", fmt.Errorf("problem preparing request: %w", err1)
	}
	addHeaders(req, authToken)

	resp, err2 := client.Do(req)
	if err2 != nil {
		return "", err2
	}
	defer resp.Body.Close()

	err3 := handleRestResponse("Run", resp)
	if err3 != nil {
		return "", err3
	}

	body, err4 := ioutil.ReadAll(resp.Body)
	if err4 != nil {
		return "", err4
	}

	return string(body[:]), nil
}

func GetCheck(checkName string, authToken string, devMode bool) (*Check, error) {

	client := &http.Client{}
	req, err := http.NewRequest("GET", getApiUrl(devMode, "api/checks/")+checkName, nil)
	if err != nil {
		return nil, fmt.Errorf("problem preparing request: %w", err)
	}
	addHeaders(req, authToken)

	resp, err1 := client.Do(req)
	if err1 != nil {
		return nil, fmt.Errorf("problem performing request: %w", err1)
	}

	defer resp.Body.Close()

	err2 := handleRestResponse("Check", resp)
	if err2 != nil {
		return nil, err2
	}

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return nil, readErr
	}

	var check Check
	jsonErr := json.Unmarshal(body, &check)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return &check, nil
}
