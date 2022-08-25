package operations

import (
	"fmt"
	"github.com/stefan-hudelmaier/checkson-cli/operations/auth"
	"github.com/stefan-hudelmaier/checkson-cli/output"
	"github.com/stefan-hudelmaier/checkson-cli/services"
)

type ShowOperation struct {
}

type ShowOperationFlags struct {
	DevMode bool
}

func (operation *ShowOperation) ShowOperation(checkName string, flags ShowOperationFlags) error {
	authToken, _ := auth.ReadAuthToken()

	check, err := services.GetCheck(checkName, authToken, flags.DevMode)
	if err != nil {
		return err
	}

	// TODO: Property output check
	fmt.Println(check)
	output.PrintStrings(check.Name, check.DockerImage)

	return nil
}
