package github

import (
	"net/http"
	"testing"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		owner, repo string
		wantError   bool
	}{
		{"test", "test", false},
		{"", "", true},
	}

	for _, te := range tests {
		_, err := NewClient(te.owner, te.repo)
		if !te.wantError && err != nil {
			t.Errorf("want no error happen, got an error: %s", err)
		}
		if te.wantError && err == nil {
			t.Error("want error happen, got nothing")
		}
	}
}

func TestClient_GetLatestTag(t *testing.T) {
	tests := []struct {
		label, latestURL string
		wantError        bool
	}{
		{"normal", baseURL + "toshi0607/gig" + latestPath, false},
		{"no URL", "", true},
		{"no REPO", baseURL + "toshi0607" + latestPath, true},
	}

	fullRepo := "toshi0607/gig"
	repoURL := baseURL + fullRepo

	for _, te := range tests {
		c := &Client{
			httpClient: &http.Client{},
			RepoURL:    repoURL,
			latestURL:  te.latestURL,
			FullRepo:   fullRepo,
		}
		_, err := c.GetLatestTag()
		if !te.wantError && err != nil {
			t.Errorf("want no error happen, got an error: %s", err)
		}
		if te.wantError && err == nil {
			t.Error("want error happen, got nothing")
		}
	}
}
