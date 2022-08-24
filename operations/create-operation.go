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

type CreateCheckFlags struct {
	DockerImage            string
	CheckIntervalInMinutes int16
	DevMode                bool
	Environment            map[string]string
	WebHookUrl             string
	EmailAddress           string
}

type CreateCheckOperation struct {
}

func (operation *CreateCheckOperation) CreateCheckOperation(checkName string, flags CreateCheckFlags) error {

	authToken, _ := auth.ReadAuthToken()

	check := Check{
		Name:                   checkName,
		WebHookUrl:             flags.WebHookUrl,
		DockerImage:            flags.DockerImage,
		CheckIntervalInMinutes: flags.CheckIntervalInMinutes,
		Environment:            flags.Environment,
	}

	client := &http.Client{}

	jsonBytes, jsonErr := json.Marshal(check)
	output.Debugf("Sending:", string(jsonBytes))
	if jsonErr != nil {
		return errors.New("Cannot serialize check")
	}

	// TODO: Create notification channel

	url := getApiUrl(flags.DevMode, "api/checks/")
	req, err := http.NewRequest("PUT", url+checkName, bytes.NewBuffer(jsonBytes))

	if err != nil {
		return fmt.Errorf("problem preparing request: %w", err)
	}
	req.Header.Add("Authorization", "Bearer "+authToken)
	req.Header.Add("Content-Type", "application/json")

	resp, err1 := client.Do(req)
	if err1 != nil {
		return fmt.Errorf("problem performing request: %w", err1)
	}
	defer resp.Body.Close()
	output.PrintStrings("Response status:", resp.Status)

	return nil
}
