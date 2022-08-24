package operations

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/stefan-hudelmaier/checkson-cli/operations/auth"
	"github.com/stefan-hudelmaier/checkson-cli/output"
	"net/http"
)

type DeleteChannelFlags struct {
	DevMode      bool
	WebHookUrl   string
	EmailAddress string
}

type DeleteChannelOperation struct {
}

func (operation *DeleteChannelOperation) DeleteChannelOperation(checkName string, flags DeleteChannelFlags) error {

	authToken, _ := auth.ReadAuthToken()

	// TODO: Implement

	check := Check{
		Name:       checkName,
		WebHookUrl: flags.WebHookUrl,
	}

	client := &http.Client{}

	jsonBytes, jsonErr := json.Marshal(check)
	output.Debugf("Sending:", string(jsonBytes))
	if jsonErr != nil {
		return errors.New("Cannot serialize check")
	}

	url := getApiUrl(flags.DevMode, "api/checks/")
	req, err := http.NewRequest("PUT", url+checkName, bytes.NewBuffer(jsonBytes))
	req.Header.Add("Authorization", "Bearer "+authToken)
	req.Header.Add("Content-Type", "application/json")

	resp, err1 := client.Do(req)
	if err1 != nil {
		return fmt.Errorf("problem performing request: %w", err)
	}
	defer resp.Body.Close()
	output.PrintStrings("Response status:", resp.Status)

	return nil
}
