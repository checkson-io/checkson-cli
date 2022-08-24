package operations

import (
	"encoding/json"
	"fmt"
	"github.com/stefan-hudelmaier/checkson-cli/operations/auth"
	"github.com/stefan-hudelmaier/checkson-cli/output"
	"io/ioutil"
	"log"
	"net/http"
)

type ShowOperation struct {
}

type ShowOperationFlags struct {
	DevMode bool
}

func (operation *ShowOperation) ShowOperation(checkName string, flags ShowOperationFlags) error {
	authToken, _ := auth.ReadAuthToken()

	client := &http.Client{}
	req, err := http.NewRequest("GET", getApiUrl(flags.DevMode, "api/checks/")+checkName, nil)
	if err != nil {
		return fmt.Errorf("problem preparing request: %w", err)
	}
	req.Header.Add("Authorization", "Bearer "+authToken)

	resp, err1 := client.Do(req)
	if err1 != nil {
		return fmt.Errorf("problem performing request: %w", err1)
	}

	defer resp.Body.Close()
	output.PrintStrings("Response status:", resp.Status)

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		panic(readErr)
	}

	var check Check
	jsonErr := json.Unmarshal(body, &check)
	if jsonErr != nil {
		log.Fatalf("unable to parse value: %q, error: %s", string(body), jsonErr.Error())
	}

	fmt.Println(check)
	output.PrintStrings(check.Name, check.DockerImage)

	return nil
}
