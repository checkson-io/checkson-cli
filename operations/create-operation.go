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
	WebHookUrl             string
	DockerImage            string
	CheckIntervalInMinutes int16
	DevMode                bool
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
	}

	client := &http.Client{}

	jsonBytes, jsonErr := json.Marshal(check)
	fmt.Println("Sending:", string(jsonBytes))
	if jsonErr != nil {
		return errors.New("Cannot serialize check")
	}

	url := getApiUrl(flags.DevMode, "api/checks/")
	req, err := http.NewRequest("PUT", url+checkName, bytes.NewBuffer(jsonBytes))
	req.Header.Add("Authorization", "Bearer "+authToken)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("problem performing request: %w", err)
	}
	defer resp.Body.Close()
	output.PrintStrings("Response status:", resp.Status)

	return nil
}
