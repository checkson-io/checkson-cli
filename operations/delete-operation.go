package operations

import (
	"errors"
	"github.com/checkson-io/checkson-cli/operations/auth"
	"github.com/checkson-io/checkson-cli/services"
)

type DeleteCheckFlags struct {
	DevMode bool
}

type DeleteCheckOperation struct {
}

func (operation *DeleteCheckOperation) DeleteCheckOperation(checkName string, flags DeleteCheckFlags) error {
	authToken, err := auth.ReadAuthToken()
	if err != nil {
		return errors.New("you are not logged in. Login with: 'checkson login'")
	}

	return services.DeleteCheck(checkName, authToken, flags.DevMode)
}
