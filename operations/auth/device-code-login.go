package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/stefan-hudelmaier/checkson-cli/output"
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

func DeviceCodeLogin(devMode bool) error {

	createDeviceCodeUrl := getCloudFunctionUrl(devMode, "createDeviceCode")

	output.Debugf("Creating device code: %s", createDeviceCodeUrl)

	resp, err := http.Get(createDeviceCodeUrl)
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

	fmt.Printf("Go to %s/signin-device-code?deviceCode=%s\n", getUiBaseUrl(devMode), deviceCodeCreationResult.DeviceCode)

	for i := 0; i < 20; i++ {
		time.Sleep(1 * time.Second)
		deviceCodeCreationResult := checkDeviceCodeStatus(devMode, deviceCodeCreationResult.DeviceCode)
		if deviceCodeCreationResult.Confirmed {
			fmt.Println("Device code has been confirmed")

			firebaseAuthToken, exchangeErr := exchangeCustomAuthTokenForFirebaseToken(deviceCodeCreationResult.AuthToken)
			if exchangeErr != nil {
				return errors.New("Login failure: Could not exchange custom auth token")
			}

			persistErr := persistAuthToken(firebaseAuthToken)
			if persistErr != nil {
				return persistErr
			}
			return nil
		} else {
			fmt.Println("Device code has not been confirmed yet, trying again")
		}
	}
	return errors.New("timeout waiting for login")
}

func checkDeviceCodeStatus(devMode bool, deviceCode string) DeviceCodeStatusResult {

	var jsonStr = []byte(fmt.Sprintf(`{"deviceCode":"%s"}`, deviceCode))
	resp, err := http.Post(getCloudFunctionUrl(devMode, "getDeviceCodeStatus"), "application/json", bytes.NewBuffer(jsonStr))

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
