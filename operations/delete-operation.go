package operations

import (
	"fmt"
	"github.com/stefan-hudelmaier/checkson-cli/operations/auth"
	"github.com/stefan-hudelmaier/checkson-cli/output"
	"net/http"
)

type DeleteCheckFlags struct {
	DevMode bool
}

type DeleteCheckOperation struct {
}

func (operation *DeleteCheckOperation) DeleteCheckOperation(checkName string, flags DeleteCheckFlags) error {

	authToken, _ := auth.ReadAuthToken()

	client := &http.Client{}

	req, err := http.NewRequest("DELETE", getApiUrl(flags.DevMode, "api/checks/")+checkName, nil)
	req.Header.Add("Authorization", "Bearer "+authToken)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("problem performing request: %w", err)
	}
	defer resp.Body.Close()
	output.PrintStrings("Response status:", resp.Status)

	return nil
}
