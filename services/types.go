package services

type Check struct {
	Name                      string            `json:"name"`
	WebHookUrl                string            `json:"webHookUrl"`
	DockerImage               string            `json:"dockerImage"`
	CheckIntervalInMinutes    int16             `json:"checkIntervalInMinutes"`
	LastCheckOutcome          string            `json:"lastCheckOutcome"`
	LastOutcomeChange         string            `json:"lastOutcomeChange"`
	LastOutcomeChangeDuration string            `json:"lastOutcomeChangeDuration"`
	Environment               map[string]string `json:"environment"`
}

type NotificationChannel struct {
	Name                    string `json:"name"`
	Type                    string `json:"type"`
	WebhookUrl              string `json:"webhookUrl"`
	EmailAddress            string `json:"emailAddress"`
	SlackIncomingWebhookUrl string `json:"slackIncomingWebhookUrl"`
	PagerDutyServiceKey     string `json:"pagerDutyServiceKey"`
}

type Run struct {
	Id                string `json:"id"`
	CheckName         string `json:"checkName"`
	StartTime         string `json:"startTime"`
	EndTime           string `json:"endTime"`
	Success           bool   `json:"success"`
	DurationInSeconds int    `json:"durationInSeconds"`
}
