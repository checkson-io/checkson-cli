package operations

import (
	"github.com/pkg/errors"
	"github.com/stefan-hudelmaier/checkson-cli/operations/auth"
	"github.com/stefan-hudelmaier/checkson-cli/services"
)

type CreateChannelFlags struct {
	DevMode                 bool
	Type                    string
	WebHookUrl              string
	EmailAddress            string
	PagerDutyServiceKey     string
	SlackIncomingWebhookUrl string
}

type CreateChannelOperation struct {
}

func (operation *CreateChannelOperation) CreateChannelOperation(channelName string, flags CreateChannelFlags) error {

	if flags.Type == "email" && len(flags.EmailAddress) == 0 {
		return errors.New("When type 'email' is used, the parameter --email must be specified")
	}

	if flags.Type == "pager-duty" && len(flags.PagerDutyServiceKey) == 0 {
		return errors.New("When type 'pager-duty' is used, the option --pager-duty-incoming-webhook must be specified")
	}

	if flags.Type == "slack" && len(flags.SlackIncomingWebhookUrl) == 0 {
		return errors.New("When type 'slack' is used, the option --slack-incoming-webhook-url must be specified")
	}

	if flags.Type == "webhook" && len(flags.WebHookUrl) == 0 {
		return errors.New("When type 'webhook' is used, the option --webhook must be specified")
	}

	authToken, err := auth.ReadAuthToken()
	if err != nil {
		return errors.New("you are not logged in. Login with: 'checkson login'")
	}

	channel := services.NotificationChannel{
		Name:                    channelName,
		Type:                    flags.Type,
		WebhookUrl:              flags.WebHookUrl,
		SlackIncomingWebhookUrl: flags.SlackIncomingWebhookUrl,
		EmailAddress:            flags.EmailAddress,
		PagerDutyServiceKey:     flags.PagerDutyServiceKey,
	}

	return services.CreateChannel(channel, authToken, flags.DevMode)
}
