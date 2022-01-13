package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stefan-hudelmaier/checkson-cli/output"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type DeviceCodeCreationResult struct {
	DeviceCode string `json:"deviceCode"`
}

type DeviceCodeStatusResult struct {
	Confirmed bool   `json:"confirmed"`
	AuthToken string `json:"authToken"`
}

func DeviceCodeLogin() error {

	resp, err := http.Get("https://europe-west1-checkson-dc8a5.cloudfunctions.net/createDeviceCode")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	output.PrintStrings("Response status:", resp.Status)

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		panic(readErr)
	}

	var deviceCodeCreationResult DeviceCodeCreationResult
	jsonErr := json.Unmarshal(body, &deviceCodeCreationResult)
	if jsonErr != nil {
		log.Fatalf("unable to parse value: %q, error: %s", string(body), jsonErr.Error())
	}

	fmt.Println("Go to checkson.io and enter the device code", deviceCodeCreationResult.DeviceCode)

	for i := 0; i < 20; i++ {
		time.Sleep(1 * time.Second)
		deviceCodeCreationResult := checkDeviceCodeStatus(deviceCodeCreationResult.DeviceCode)
		if deviceCodeCreationResult.Confirmed == true {
			fmt.Println("Device code has been confirmed, custom auth token is", deviceCodeCreationResult.AuthToken)
			// TODO: Exchange custom auth token for Firebase token
			return nil
		} else {
			fmt.Println("Device code has not been confirmed yet, trying again")
		}
	}
	return errors.New("timeout waiting for login")
}

func checkDeviceCodeStatus(deviceCode string) DeviceCodeStatusResult {

	var jsonStr = []byte(fmt.Sprintf(`{"deviceCode":"%s"}`, deviceCode))
	resp, err := http.Post("https://europe-west1-checkson-dc8a5.cloudfunctions.net/getDeviceCodeStatus", "application/json", bytes.NewBuffer(jsonStr))

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	output.PrintStrings("Response status:", resp.Status)

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		panic(readErr)
	}

	var deviceCodeStatusResult DeviceCodeStatusResult
	jsonErr := json.Unmarshal(body, &deviceCodeStatusResult)
	if jsonErr != nil {
		log.Fatalf("unable to parse value: %q, error: %s", string(body), jsonErr.Error())
	}

	return deviceCodeStatusResult
}
