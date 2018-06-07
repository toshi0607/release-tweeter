package github

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

const (
	baseURL    = "https://github.com/"
	latestPath = "/releases/latest"
)

// Client represents GitHub client
type Client struct {
	httpClient *http.Client
	RepoURL    string
	latestURL  string
	FullRepo   string
}

// NewClient creates a new GitHub client from a user name and a repo name
func NewClient(owner, repo string) (*Client, error) {
	if owner == "" || repo == "" {
		return nil, errors.New("github: owner and repo are required")
	}

	repoURL := baseURL + owner + "/" + repo

	return &Client{
		httpClient: &http.Client{},
		RepoURL:    repoURL,
		latestURL:  repoURL + latestPath,
		FullRepo:   repo,
	}, nil
}

// GetLatestTag returns the latest tag from a GitHub release
func (g Client) GetLatestTag() (string, error) {
	req, err := http.NewRequest("GET", g.latestURL, nil)
	if err != nil {
		return "", err
	}

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	tag, err := getLatestTag(resp.Request.URL.Path)
	if err != nil {
		return "", err
	}
	return tag, nil
}

func getLatestTag(url string) (string, error) {
	s := strings.Split(url, "/")
	if !strings.Contains(url, "v") {
		return "", fmt.Errorf("has no tag")
	}
	return s[len(s)-1], nil
}
