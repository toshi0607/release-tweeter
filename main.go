package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	a "github.com/ChimeraCoder/anaconda"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pkg/errors"
)

var (
	twitterAccessToken               = os.Getenv("TWITTER_ACCESS_TOKEN")
	twitterAccessTokenSecret         = os.Getenv("TWITTER_ACCESS_TOKEN_SECRET")
	twitterConsumerKey               = os.Getenv("TWITTER_CONSUMER_KEY")
	twitterConsumerKeySecret         = os.Getenv("TWITTER_CONSUMER_SECRET")
	twitterClient            tweeter = api(*a.NewTwitterApiWithCredentials(
		twitterAccessToken,
		twitterAccessTokenSecret,
		twitterConsumerKey,
		twitterConsumerKeySecret,
	))
)

const gitHubBaseURL = "https://github.com/"

type tweeter interface {
	tweet(message string) (string, error)
}

type api a.TwitterApi

func (api api) tweet(message string) (string, error) {
	tweet, err := (a.TwitterApi(api)).PostTweet(message, nil)
	return tweet.Text, err
}

type params struct {
	Owner string `json:"owner"`
	Repo  string `json:"repo"`
}

type gitHubClient struct {
	latestURL  string
	repoURL    string
	httpClient *http.Client
	fullRepo   string
}

func newGitHubClient(p params) *gitHubClient {
	repo := p.Owner + "/" + p.Repo
	repoURL := gitHubBaseURL + repo

	return &gitHubClient{
		latestURL:  repoURL + "/releases/latest",
		httpClient: &http.Client{},
		fullRepo:   repo,
		repoURL:    repoURL,
	}
}

func (g gitHubClient) getLatestTag() (string, error) {
	req, err := http.NewRequest("GET", g.latestURL, nil)
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

func main() {
	lambda.Start(handler)
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if hasEmptyEnvVar() {
		return response(
			http.StatusInternalServerError,
			"env vars not set",
		), nil
	}

	p, err := parseRequest(request)
	if err != nil {
		return response(
			http.StatusBadRequest,
			err.Error(),
		), nil
	}

	c := newGitHubClient(*p)
	tag, err := c.getLatestTag()
	if err != nil {
		return response(
			http.StatusInternalServerError,
			err.Error(),
		), nil
	}
	log.Printf("tag: %s, ID: %s\n", tag, request.RequestContext.RequestID)

	message := fmt.Sprintf("%s %s released! check the new features on GitHub.\n%s", c.fullRepo, tag, c.repoURL)
	msg, err := twitterClient.tweet(message)
	if err != nil {
		return response(
			http.StatusInternalServerError,
			err.Error(),
		), nil
	}
	if !strings.Contains(msg, tag) {
		return response(
			http.StatusInternalServerError,
			fmt.Sprintf("failed to tweet: %s", msg),
		), nil
	}
	log.Printf("message tweeted: %s, ID: %s\n", msg, request.RequestContext.RequestID)

	return response(http.StatusOK, tag), nil
}

func parseRequest(r events.APIGatewayProxyRequest) (*params, error) {
	if r.HTTPMethod != "POST" {
		return nil, fmt.Errorf("use POST request")
	}

	var p params
	err := json.Unmarshal([]byte(r.Body), &p)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse params")
	}
	return &p, nil
}
func hasEmptyEnvVar() bool {
	return twitterAccessToken == "" ||
		twitterAccessTokenSecret == "" ||
		twitterConsumerKey == "" ||
		twitterConsumerKeySecret == ""
}

func getLatestTag(url string) (string, error) {
	s := strings.Split(url, "/")
	if !strings.Contains(url, "v") {
		return "", fmt.Errorf("has no tag")
	}
	return s[len(s)-1], nil
}

func response(code int, msg string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: code,
		Body:       fmt.Sprintf("{\"message\":\"%s\"}", msg),
		Headers:    map[string]string{"Content-Type": "application/json"},
	}
}
