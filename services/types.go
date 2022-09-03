package services

type Check struct {
	Name                   string            `json:"name"`
	Enabled                bool              `json:"enabled"`
	DockerImage            string            `json:"dockerImage"`
	DockerCredentials      DockerCredentials `json:"dockerCredentials"`
	CheckIntervalInMinutes int16             `json:"checkIntervalInMinutes"`
	Environment            map[string]string `json:"environment"`
	FailureThreshold       int16             `json:"failureThreshold"`
	FailureCount           int16             `json:"failureCount"`
	LastRunOutcome         string            `json:"LastRunOutcome"`
	Status                 string            `json:"status"`
	LastStatusChange       string            `json:"lastStatusChange"`
	NextRun                string            `json:"nextRun"`
	NotificationChannels   []string          `json:"notificationChannels"`
}

type DockerCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
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
	Outcome           string `json:"outcome"`
	DurationInSeconds int    `json:"durationInSeconds"`
}
