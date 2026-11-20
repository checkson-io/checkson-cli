package operations

import (
	"encoding/json"
	"fmt"
	"github.com/checkson-io/checkson-cli/operations/auth"
	"github.com/checkson-io/checkson-cli/services"
	"github.com/pkg/errors"
)

type ShowChannelFlags struct {
	DevMode bool
}

type ShowChannelOperation struct {
}

func (operation *ShowChannelOperation) ShowChannelOperation(channelName string, flags ShowChannelFlags) error {

	authToken, err := auth.ReadAuthToken()
	if err != nil {
		return errors.New("you are not logged in. Login with: 'checkson login'")
	}

	channel, err1 := services.GetChannel(channelName, authToken, flags.DevMode)
	if err1 != nil {
		return err1
	}

	bytes, err2 := json.MarshalIndent(channel, "", "   ")
	if err2 != nil {
		return err2
	}

	fmt.Println(string(bytes))

	return nil
}
