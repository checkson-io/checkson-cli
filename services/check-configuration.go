package services

// CheckConfiguration contains only writable check properties.
type CheckConfiguration struct {
	DockerImage            string             `json:"dockerImage"`
	Enabled                bool               `json:"enabled"`
	DockerCredentials      *DockerCredentials `json:"dockerCredentials"`
	CheckIntervalInMinutes int16              `json:"checkIntervalInMinutes"`
	Environment            map[string]string  `json:"environment"`
	Description            string             `json:"description"`
	FailureThreshold       int16              `json:"failureThreshold"`
	NotificationChannels   []string           `json:"notificationChannels"`
}
