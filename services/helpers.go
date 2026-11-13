package services

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

func getApiUrl(devMode bool, path string) string {
	var baseUrl = "https://api.checkson.io"

	if devMode {
		baseUrl = "http://127.0.0.1:8080"
	}

	return fmt.Sprintf("%s/%s", baseUrl, path)
}

func handleRestResponse(entityName string, resp *http.Response) error {

	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusUnauthorized:
		return errors.New("you are not logged in, please login using 'checkson login'")
	case http.StatusBadRequest:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.New("validation error: Cannot read response")
		}
		return errors.New("validation error: " + string(body[:]))
	case http.StatusNotFound:
		return fmt.Errorf("%s not found", entityName)
	default:
		return fmt.Errorf("unexpected response: %s", resp.Status)
	}
}

func addHeaders(req *http.Request, authToken string) {
	req.Header.Add("Authorization", "Bearer "+authToken)
	req.Header.Add("Content-Type", "application/json")
}
