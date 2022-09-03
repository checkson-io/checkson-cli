package operations

import (
	"encoding/json"
	"fmt"
	"github.com/stefan-hudelmaier/checkson-cli/operations/auth"
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

	bytes, err1 := json.MarshalIndent(check, "", "    ")
	if err1 != nil {
		return err1
	}

	fmt.Println(string(bytes))

	return nil
}
