package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stefan-hudelmaier/checkson-cli/config"
	"github.com/stefan-hudelmaier/checkson-cli/output"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type FirebaseResponse struct {
	IdToken string `json:"idToken"`
}

func ReadAuthToken() (string, error) {
	bytes, fileError := os.ReadFile(GetAuthFile())
	if fileError != nil {
		return "", fileError
	}
	return strings.TrimSpace(string(bytes)), nil
}

func GetAuthFile() string {
	return filepath.Join(config.GetConfigPath(), "auth")
}

func exchangeCustomAuthTokenForFirebaseToken(customAuthToken string) (string, error) {

	firebaseApiKey := "AIzaSyA7Uc_3kxi2bx9FyBFeho7fdQkyHQG18Gs"
	var jsonStr = []byte(fmt.Sprintf(`{"token": "%s", "returnSecureToken": "true"}`, customAuthToken))
	resp, err := http.Post("https://identitytoolkit.googleapis.com/v1/accounts:signInWithCustomToken?key="+firebaseApiKey, "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	output.Debugf("Response status: %s", resp.Status)

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		panic(readErr)
	}

	var firebaseResponse FirebaseResponse
	jsonErr := json.Unmarshal(body, &firebaseResponse)
	if jsonErr != nil {
		log.Fatalf("unable to parse value: %q, error: %s", string(body), jsonErr.Error())
	}

	output.Debugf("Received the Firebase auth token: %s", firebaseResponse.IdToken)
	return firebaseResponse.IdToken, nil
}

func persistAuthToken(firebaseAuthToken string) error {
	configPath := config.GetConfigPath()

	err := os.MkdirAll(configPath, os.ModePerm)
	if err != nil {
		panic(err)
	}

	writeErr := os.WriteFile(GetAuthFile(), []byte(firebaseAuthToken), os.ModePerm)
	if writeErr != nil {
		panic(writeErr)
	}

	return nil
}

func RemovePersistedAuthToken() error {
	configPath := config.GetConfigPath()

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil
	}

	err := os.Remove(GetAuthFile())
	if err != nil {
		return err
	}

	return nil
}

func getCloudFunctionUrl(devMode bool, function string) string {

	var baseUrl = "https://europe-west1-checkson-cadf1.cloudfunctions.net"

	if devMode {
		baseUrl = "http://localhost:5001/checkson-cadf1/europe-west1"
	}

	return fmt.Sprintf("%s/%s", baseUrl, function)
}

func getUiBaseUrl(devMode bool) string {

	var baseUrl = "https://app.checkson.io"

	if devMode {
		baseUrl = "http://localhost:3000"
	}

	return baseUrl
}
