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
)

type CustomAuthTokenWrapper struct {
	AuthToken string `json:"authToken"`
}

func PersonalAccessTokenLogin(personalAccessToken string) error {

	var jsonStr = []byte(fmt.Sprintf(`{"personalAccessToken":"%s"}`, personalAccessToken))
	resp, err := http.Post(getCloudFunctionUrl(false, "getcustomauthtokenforpersonalaccesstoken"), "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	output.Debugf("Response status: %s", resp.Status)

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		panic(readErr)
	}

	var customAuthTokenWrapper CustomAuthTokenWrapper
	jsonErr := json.Unmarshal(body, &customAuthTokenWrapper)
	if jsonErr != nil {
		log.Fatalf("unable to parse value: %q, error: %s", string(body), jsonErr.Error())
	}

	output.Debugf("Received the custom auth token: %s", customAuthTokenWrapper.AuthToken)

	firebaseAuthToken, exchangeErr := exchangeCustomAuthTokenForFirebaseToken(customAuthTokenWrapper.AuthToken)
	if exchangeErr != nil {
		return errors.New("Login failure: Could not exchange custom auth token")
	}

	persistErr := persistAuthToken(firebaseAuthToken)
	if persistErr != nil {
		return persistErr
	}

	output.PrintStrings("Login successful")

	return nil
}
