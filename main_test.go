package main

import (
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestHandler(t *testing.T) {
	prepare()

	tests := []struct {
		label, owner, repo, method string
		status                     int
	}{
		{"normal", "toshi0607", "gig", "POST", http.StatusOK},
		{"invalid req", "", "gig", "POST", http.StatusBadRequest},
		{"invalid repo", "toshi0607", "nothing", "POST", http.StatusInternalServerError},
		{"invalid method", "toshi0607", "gig", "GET", http.StatusBadRequest},
	}

	for _, te := range tests {

		res, _ := handler(events.APIGatewayProxyRequest{
			HTTPMethod: te.method,
			Body:       "{\"owner\": \"" + te.owner + "\", \"repo\": \"" + te.repo + "\"}",
		})

		if res.StatusCode != te.status {
			t.Errorf("ExitStatus=%d, want %d", res.StatusCode, te.status)
		}
	}
}

func prepare() {
	twitterClient = mockClient("mock")
}

type mockClient string

func (m mockClient) Tweet(s string) (string, error) {
	return s, nil
}
