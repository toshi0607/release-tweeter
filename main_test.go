package main

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestHandler(t *testing.T) {
	prepare()
	Handler(events.APIGatewayProxyRequest{})
}

func prepare() {
	twitterClient = mockClient("mock")
}

type mockClient string

func (m mockClient) tweet(s string) (string, error) {
	return s, nil
}
