package operations

import (
	"fmt"
	"github.com/huditech/checkson/operations/auth"
	"github.com/huditech/checkson/output"
	"net/http"
)

type DeleteCheckFlags struct {
}

type DeleteCheckOperation struct {
}

func (operation *DeleteCheckOperation) DeleteCheckOperation(checkName string, flags DeleteCheckFlags) error {

	authToken, _ := auth.ReadAuthToken()

	client := &http.Client{}

	req, err := http.NewRequest("DELETE", "http://127.0.0.1:8080/api/checks/"+checkName, nil)
	req.Header.Add("Authorization", "Bearer "+authToken)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("problem performing request: %w", err)
	}
	defer resp.Body.Close()
	output.PrintStrings("Response status:", resp.Status)

	return nil
}
