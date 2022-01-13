package operations

type Check struct {
	Name                      string `json:"name"`
	WebHookUrl                string `json:"webHookUrl"`
	DockerImage               string `json:"dockerImage"`
	CheckIntervalInMinutes    int16  `json:"checkIntervalInMinutes"`
	LastCheckOutcome          string `json:"lastCheckOutcome"`
	LastOutcomeChange         string `json:"lastOutcomeChange"`
	LastOutcomeChangeDuration string `json:"lastOutcomeChangeDuration"`
}
