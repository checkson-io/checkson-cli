package operations

import (
	"github.com/checkson-io/checkson-cli/operations/auth"
	"github.com/checkson-io/checkson-cli/output"
)

type LogoutOperation struct {
}

type LogoutOperationFlags struct {
}

func (operation *LogoutOperation) LogoutOperation(flags LogoutOperationFlags) error {

	err := auth.RemovePersistedAuthToken()
	if err != nil {
		output.PrintStrings("You were not logged in")
		return nil
	}

	output.PrintStrings("Logout successful")
	return nil
}
