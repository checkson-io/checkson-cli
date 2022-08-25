package operations

import (
	"fmt"
	"github.com/stefan-hudelmaier/checkson-cli/operations/auth"
	"github.com/stefan-hudelmaier/checkson-cli/output"
	"github.com/stefan-hudelmaier/checkson-cli/services"
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

	log, err1 := services.GetLog(checkName, runId, authToken, flags.DevMode)
	if err != nil {
		return err1
	}

	output.PrintStrings(log)
	return nil
}
