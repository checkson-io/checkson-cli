package operations

import (
	"fmt"
	"github.com/stefan-hudelmaier/checkson-cli/operations/auth"
	"github.com/stefan-hudelmaier/checkson-cli/output"
	"io/ioutil"
	"net/http"
)

type LogsOperation struct {
}

type LogsOperationFlags struct {
	DevMode bool
}

func (operation *LogsOperation) LogsOperation(checkName string, runId string, flags LogsOperationFlags) error {
	authToken, err := auth.ReadAuthToken()
	if err != nil {
		fmt.Println("You are not logged in. Login with: 'checkson login'")
		return nil
	}

	path := fmt.Sprintf("api/checks/%s/runs/%s/log", checkName, runId)
	client := &http.Client{}
	req, err1 := http.NewRequest("GET", getApiUrl(flags.DevMode, path), nil)
	if err1 != nil {
		return fmt.Errorf("problem preparing request: %w", err1)
	}
	req.Header.Add("Authorization", "Bearer "+authToken)

	resp, err2 := client.Do(req)
	if err2 != nil {
		panic(err2)
	}
	defer resp.Body.Close()
	output.Debugf("Response status: %s", resp.Status)

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		panic(readErr)
	}

	output.PrintStrings(string(body[:]))

	return nil
}
