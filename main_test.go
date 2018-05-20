package main

import (
	"testing"
)

func TestHandler(t *testing.T) {
	prepare()
	Handler()
}

func prepare() {
	twitterClient = mockClient("mock")
}

type mockClient string

func (m mockClient) tweet(s string) (string, error) {
	return s, nil
}
