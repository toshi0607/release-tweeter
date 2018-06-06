package main

import (
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestHandler(t *testing.T) {
	prepare()

	res, err := handler(events.APIGatewayProxyRequest{
		HTTPMethod: "POST",
		Body:       "{\"owner\": \"toshi0607\", \"repo\": \"gig\"}",
	})

	if err != nil {
		t.Errorf("got an error %v, want nothing happened", err)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("ExitStatus=%d, want %d", res.StatusCode, http.StatusOK)
	}
}

func prepare() {
	twitterClient = mockClient("mock")
}

type mockClient string

func (m mockClient) Tweet(s string) (string, error) {
	return s, nil
}
