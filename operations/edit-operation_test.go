package operations

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

const completeConfig = `{"dockerImage":"example/image:latest","enabled":false,"dockerCredentials":null,"checkIntervalInMinutes":10,"environment":{},"description":"","failureThreshold":1,"notificationChannels":[]}`

func TestParseCheckConfiguration(t *testing.T) {
	config, err := ParseCheckConfiguration([]byte(completeConfig))
	if err != nil {
		t.Fatal(err)
	}
	if config.Enabled || config.DockerCredentials != nil || config.Environment == nil || config.NotificationChannels == nil {
		t.Fatalf("explicit empty values not preserved: %#v", config)
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal([]byte(completeConfig), &fields); err != nil {
		t.Fatal(err)
	}
	for key := range fields {
		t.Run(key, func(t *testing.T) {
			var copy map[string]json.RawMessage
			_ = json.Unmarshal([]byte(completeConfig), &copy)
			delete(copy, key)
			data, _ := json.Marshal(copy)
			if _, err := ParseCheckConfiguration(data); err == nil {
				t.Fatalf("accepted missing %s", key)
			}
			copy[key] = json.RawMessage("null")
			data, _ = json.Marshal(copy)
			if _, err := ParseCheckConfiguration(data); key != "dockerCredentials" && err == nil {
				t.Fatalf("accepted null %s", key)
			}
		})
	}
	for _, data := range []string{
		strings.Replace(completeConfig, `"failureThreshold":1`, `"failureThreshold":0`, 1),
		strings.Replace(completeConfig, `"enabled":false`, `"enabled":"false"`, 1),
		strings.Replace(completeConfig, `"description":""`, `"unexpected":true`, 1),
		completeConfig + `{}`,
	} {
		if _, err := ParseCheckConfiguration([]byte(data)); err == nil {
			t.Fatalf("accepted invalid input: %s", data)
		}
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestRequireCheckExistence(t *testing.T) {
	original := http.DefaultTransport
	defer func() { http.DefaultTransport = original }()
	http.DefaultTransport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.Method != "GET" {
			t.Fatalf("unexpected write: %s", r.Method)
		}
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(`[{"name":"existing"}]`)), Header: make(http.Header)}, nil
	})
	for _, test := range []struct {
		name                  string
		wantExists, wantError bool
	}{
		{"existing", false, true}, {"new", false, false}, {"existing", true, false}, {"new", true, true},
	} {
		err := requireCheckExistence(test.name, "token", false, test.wantExists)
		if (err != nil) != test.wantError {
			t.Fatalf("%+v: %v", test, err)
		}
	}
}
