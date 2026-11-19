package services

import (
	"context"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"
)

type completionTransport func(*http.Request) (*http.Response, error)

func (f completionTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestListCompletionNames(t *testing.T) {
	original := http.DefaultClient
	t.Cleanup(func() { http.DefaultClient = original })
	for _, channels := range []bool{false, true} {
		http.DefaultClient = &http.Client{Transport: completionTransport(func(r *http.Request) (*http.Response, error) {
			expected := "/api/checks"
			if channels {
				expected = "/api/notification-channels"
			}
			if r.URL.Path != expected || r.Header.Get("Authorization") != "Bearer test-token" {
				t.Fatalf("unexpected request: %s", r.URL)
			}
			return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(`[{"name":"James Coq Admin"},{"name":"check-two"}]`))}, nil
		})}
		names, err := ListCompletionNames(context.Background(), "test-token", false, channels)
		if err != nil || !reflect.DeepEqual(names, []string{"James Coq Admin", "check-two"}) {
			t.Fatalf("names=%v err=%v", names, err)
		}
	}
	for _, body := range []string{"not json", `{"error":"unauthorized"}`} {
		http.DefaultClient = &http.Client{Transport: completionTransport(func(r *http.Request) (*http.Response, error) {
			return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(body))}, nil
		})}
		if _, err := ListCompletionNames(context.Background(), "test", false, false); err == nil {
			t.Fatal("accepted invalid response")
		}
	}
	http.DefaultClient = &http.Client{Transport: completionTransport(func(r *http.Request) (*http.Response, error) {
		<-r.Context().Done()
		return nil, r.Context().Err()
	})}
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()
	if _, err := ListCompletionNames(ctx, "test", false, false); err == nil {
		t.Fatal("expected deadline failure")
	}
}
