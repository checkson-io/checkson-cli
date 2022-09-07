package operations

import (
	"errors"
	"github.com/stefan-hudelmaier/checkson-cli/operations/auth"
	"github.com/stefan-hudelmaier/checkson-cli/services"
)

// TODO: Update the flags, add notification channel

type CreateCheckFlags struct {
	DockerImage            string
	CheckIntervalInMinutes int16
	FailureThreshold       int16
	DevMode                bool
	Environment            map[string]string
	WebHookUrl             string
	EmailAddress           string
	NotificationChannels   []string
}

type CreateCheckOperation struct {
}

func (operation *CreateCheckOperation) CreateCheckOperation(checkName string, flags CreateCheckFlags) error {

	authToken, err := auth.ReadAuthToken()
	if err != nil {
		return errors.New("you are not logged in. Login with: 'checkson login'")
	}

	notificationChannels := flags.NotificationChannels
	if notificationChannels == nil {
		notificationChannels = make([]string, 0)
	}

	if flags.WebHookUrl != "" {
		channelName := checkName + "-webhook"
		channel := services.NotificationChannel{
			Name:       channelName,
			Type:       "webhook",
			WebhookUrl: flags.WebHookUrl,
		}

		err1 := services.CreateChannel(channel, authToken, flags.DevMode)
		if err1 != nil {
			return err1
		}

		notificationChannels = append(notificationChannels, channelName)
	}

	if flags.EmailAddress != "" {
		channelName := checkName + "-email"
		channel := services.NotificationChannel{
			Name:         channelName,
			Type:         "email",
			EmailAddress: flags.EmailAddress,
		}

		err1 := services.CreateChannel(channel, authToken, flags.DevMode)
		if err1 != nil {
			return err1
		}

		notificationChannels = append(notificationChannels, channelName)
	}

	// TODO: Validate that all channel names exist

	check := services.Check{
		Name:                   checkName,
		DockerImage:            flags.DockerImage,
		CheckIntervalInMinutes: flags.CheckIntervalInMinutes,
		FailureThreshold:       flags.FailureThreshold,
		Environment:            flags.Environment,
		NotificationChannels:   notificationChannels,
	}

	return services.CreateCheck(check, authToken, flags.DevMode)
}
