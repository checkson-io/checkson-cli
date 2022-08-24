package operations

import (
	"errors"
	"github.com/stefan-hudelmaier/checkson-cli/operations/auth"
	"github.com/stefan-hudelmaier/checkson-cli/services"
)

type DeleteChannelFlags struct {
	DevMode bool
}

type DeleteChannelOperation struct {
}

func (operation *DeleteChannelOperation) DeleteChannelOperation(channelName string, flags DeleteChannelFlags) error {
	authToken, err := auth.ReadAuthToken()
	if err != nil {
		return errors.New("you are not logged in. Login with: 'checkson login'")
	}
	return services.DeleteChannel(channelName, authToken, flags.DevMode)
}
