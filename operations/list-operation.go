package operations

import (
	"errors"
	"fmt"
	"github.com/checkson-io/checkson-cli/operations/auth"
	"github.com/checkson-io/checkson-cli/output"
	"github.com/checkson-io/checkson-cli/services"
	"os"
)
import "github.com/TwiN/go-color"

type ListOperation struct {
}

type ListOperationFlags struct {
	DevMode bool
}

func (operation *ListOperation) ListOperation(flags ListOperationFlags) error {
	authToken, err := auth.ReadAuthToken()
	if err != nil {
		return errors.New("you are not logged in. Login with: 'checkson login'")
	}

	checks, err1 := services.ListChecks(authToken, flags.DevMode)
	if err1 != nil {
		return err1
	}

	var data [][]string
	for _, check := range checks {
		statusChar := color.Ize(color.Green, "✔")
		if check.Status != "OK" {
			statusChar = color.Ize(color.Red, "✖")
		}
		recentFailures := fmt.Sprintf("%d / %d", check.FailureCount, check.FailureThreshold)
		data = append(data, []string{statusChar, check.Name, check.LastStatusChange, recentFailures})
	}

	return output.RenderTable(os.Stdout, []string{"STATUS", "NAME", "SINCE", "FAILURES / THRESH"}, data)
}
