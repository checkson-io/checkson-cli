package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// ListCompletionNames fetches only names for shell completion within the caller's deadline.
func ListCompletionNames(ctx context.Context, token string, devMode, channels bool) ([]string, error) {
	path := "api/checks"
	if channels {
		path = "api/notification-channels"
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getApiUrl(devMode, path), nil)
	if err != nil {
		return nil, err
	}
	addHeaders(req, token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("completion lookup failed: %s", resp.Status)
	}
	var resources []struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&resources); err != nil {
		return nil, err
	}
	names := make([]string, 0, len(resources))
	for _, resource := range resources {
		names = append(names, resource.Name)
	}
	return names, nil
}
