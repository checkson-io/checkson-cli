package operations

import (
	"errors"
	"github.com/stefan-hudelmaier/checkson-cli/operations/auth"
	"github.com/stefan-hudelmaier/checkson-cli/services"
)

// TODO: Update the flags

type CreateCheckFlags struct {
	DockerImage            string
	CheckIntervalInMinutes int16
	DevMode                bool
	Environment            map[string]string
	WebHookUrl             string
	EmailAddress           string
}

type CreateCheckOperation struct {
}

func (operation *CreateCheckOperation) CreateCheckOperation(checkName string, flags CreateCheckFlags) error {

	authToken, err := auth.ReadAuthToken()
	if err != nil {
		return errors.New("you are not logged in. Login with: 'checkson login'")
	}

	check := services.Check{
		Name:                   checkName,
		DockerImage:            flags.DockerImage,
		CheckIntervalInMinutes: flags.CheckIntervalInMinutes,
		Environment:            flags.Environment,
	}

	return services.CreateCheck(check, authToken, flags.DevMode)
}
