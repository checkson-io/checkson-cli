package operations

import (
	"fmt"
	"github.com/huditech/checkson/operations/auth"
	"github.com/huditech/checkson/output"
	"io/ioutil"
	"net/http"
)

type LogsOperation struct {
}

type LogsOperationFlags struct {
}

func (operation *LogsOperation) LogsOperation(checkName string, runId string, flags LogsOperationFlags) error {
	authToken, err := auth.ReadAuthToken()
	if err != nil {
		fmt.Println("You are not logged in. Login with: 'checkson login'")
		return nil
	}

	client := &http.Client{}
	url := fmt.Sprintf("http://127.0.0.1:8080/api/checks/%s/runs/%s/log", checkName, runId)
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Bearer "+authToken)

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
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
