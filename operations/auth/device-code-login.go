package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/browser"
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

	createDeviceCodeUrl := getCloudFunctionUrl(devMode, "createdevicecode")

	output.Debugf("Creating device code: %s", createDeviceCodeUrl)

	resp, err := http.Get(createDeviceCodeUrl)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	output.Debugf("Response status:", resp.Status)

	if resp.StatusCode != 200 {
		return errors.New("Could not create device code")
	}

	body, err1 := ioutil.ReadAll(resp.Body)
	if err1 != nil {
		panic(err1)
	}

	var deviceCodeCreationResult DeviceCodeCreationResult
	err2 := json.Unmarshal(body, &deviceCodeCreationResult)
	if err2 != nil {
		log.Fatalf("unable to parse value: %q, error: %s", string(body), err2.Error())
	}

	url := fmt.Sprintf("%s/signin-device-code?deviceCode=%s", getUiBaseUrl(devMode), deviceCodeCreationResult.DeviceCode)
	fmt.Printf("Go to %s\n", url)
	err3 := browser.OpenURL(url)
	if err3 != nil {
		output.PrintStrings("Could not open browser, copy and paste the above URL")
	}

	for i := 0; i < 20; i++ {
		time.Sleep(1 * time.Second)
		deviceCodeCreationResult := checkDeviceCodeStatus(devMode, deviceCodeCreationResult.DeviceCode)
		if deviceCodeCreationResult.Confirmed {
			output.Debugf("Device code has been confirmed, login succeeded")

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
			output.Debugf("Device code has not been confirmed yet, trying again")
		}
	}
	return errors.New("timeout waiting for login")
}

func checkDeviceCodeStatus(devMode bool, deviceCode string) DeviceCodeStatusResult {

	var jsonStr = []byte(fmt.Sprintf(`{"deviceCode":"%s"}`, deviceCode))

	deviceCodeStatusUrl := getCloudFunctionUrl(devMode, "getdevicecodestatus")
	output.Debugf("Checking device code: %s", deviceCodeStatusUrl)

	resp, err := http.Post(deviceCodeStatusUrl, "application/json", bytes.NewBuffer(jsonStr))

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	output.Debugf("Response status:", resp.Status)

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		panic(readErr)
	}

	output.Debugf("Response body: %s", string(body))

	var deviceCodeStatusResult DeviceCodeStatusResult
	jsonErr := json.Unmarshal(body, &deviceCodeStatusResult)
	if jsonErr != nil {
		log.Fatalf("unable to parse value: %q, error: %s", string(body), jsonErr.Error())
	}

	return deviceCodeStatusResult
}
