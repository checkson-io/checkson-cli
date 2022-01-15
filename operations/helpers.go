package operations

import "fmt"

func getApiUrl(devMode bool, path string) string {
	var baseUrl = "https://api.checkson.io"

	if devMode {
		baseUrl = "http://127.0.0.1:8080"
	}

	return fmt.Sprintf("%s/%s", baseUrl, path)
}
