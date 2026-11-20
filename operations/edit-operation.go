package operations

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/stefan-hudelmaier/checkson-cli/operations/auth"
	"github.com/stefan-hudelmaier/checkson-cli/services"
)

func requireCheckExistence(name, token string, devMode, wantExists bool) error {
	checks, err := services.ListChecks(token, devMode)
	if err != nil {
		return err
	}
	for _, check := range checks {
		if check.Name == name {
			if !wantExists {
				return fmt.Errorf("check %q already exists; use edit to replace its configuration", name)
			}
			return nil
		}
	}
	if wantExists {
		return fmt.Errorf("check %q does not exist; use create first", name)
	}
	return nil
}

func ParseCheckConfiguration(data []byte) (services.CheckConfiguration, error) {
	var config services.CheckConfiguration
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return config, err
	}
	for _, name := range []string{"dockerImage", "enabled", "dockerCredentials", "checkIntervalInMinutes", "environment", "description", "failureThreshold", "notificationChannels"} {
		value, ok := fields[name]
		if !ok || (name != "dockerCredentials" && bytes.Equal(bytes.TrimSpace(value), []byte("null"))) {
			return config, fmt.Errorf("incomplete configuration: required property %q is missing or null", name)
		}
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&config); err != nil {
		return config, err
	}
	if err := decoder.Decode(new(interface{})); err != io.EOF {
		return config, fmt.Errorf("expected a single JSON object")
	}
	if config.DockerImage == "" {
		return config, fmt.Errorf("dockerImage must not be empty")
	}
	if config.CheckIntervalInMinutes < 10 || config.CheckIntervalInMinutes > 360 {
		return config, fmt.Errorf("checkIntervalInMinutes must be between 10 and 360")
	}
	if config.FailureThreshold < 1 || config.FailureThreshold > 20 {
		return config, fmt.Errorf("failureThreshold must be between 1 and 20")
	}
	if config.DockerCredentials != nil && (config.DockerCredentials.Username == "" || config.DockerCredentials.Password == "") {
		return config, fmt.Errorf("dockerCredentials requires username and password; use null for no credentials")
	}
	return config, nil
}

func EditCheck(name, filename string, devMode bool) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	config, err := ParseCheckConfiguration(data)
	if err != nil {
		return fmt.Errorf("invalid check configuration: %w", err)
	}
	token, err := auth.ReadAuthToken()
	if err != nil {
		return fmt.Errorf("you are not logged in. Login with: 'checkson login'")
	}
	if err := requireCheckExistence(name, token, devMode, true); err != nil {
		return err
	}
	return services.ReplaceCheck(name, config, token, devMode)
}
