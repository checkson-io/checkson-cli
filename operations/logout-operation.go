package operations

import (
	"github.com/huditech/checkson/operations/auth"
	"github.com/huditech/checkson/output"
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
