package operations

import "github.com/huditech/checkson/operations/auth"

type LoginOperation struct {
}

type LoginOperationFlags struct {
	PersonalAccessToken string
}

func (operation *LoginOperation) LoginOperation(flags LoginOperationFlags) error {
	if flags.PersonalAccessToken != "" {
		return auth.PersonalAccessTokenLogin(flags.PersonalAccessToken)
	} else {
		return auth.DeviceCodeLogin()
	}
}
