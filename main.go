package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pkg/errors"
	"github.com/toshi0607/release-tweeter/github"
	"github.com/toshi0607/release-tweeter/twitter"
)

// Tweeter is interface to tweet a message
type tweeter interface {
	Tweet(message string) (string, error)
}

var twitterClient tweeter

func init() {
	twitterAccessToken := os.Getenv("TWITTER_ACCESS_TOKEN")
	twitterAccessTokenSecret := os.Getenv("TWITTER_ACCESS_TOKEN_SECRET")
	twitterConsumerKey := os.Getenv("TWITTER_CONSUMER_KEY")
	twitterConsumerKeySecret := os.Getenv("TWITTER_CONSUMER_SECRET")
	tc, err := twitter.NewClient(
		twitterAccessToken,
		twitterAccessTokenSecret,
		twitterConsumerKey,
		twitterConsumerKeySecret,
	)
	if err != nil {
		log.Fatal(err)
	}
	twitterClient = tc
}

type params struct {
	Owner string `json:"owner"`
	Repo  string `json:"repo"`
}

func main() {
	lambda.Start(handler)
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	p, err := parseRequest(request)
	if err != nil {
		return response(
			http.StatusBadRequest,
			err.Error(),
		), nil
	}

	c, err := github.NewClient(p.Owner, p.Repo)
	if err != nil {
		return response(
			http.StatusBadRequest,
			err.Error(),
		), nil
	}

	tag, err := c.GetLatestTag()
	if err != nil {
		return response(
			http.StatusInternalServerError,
			err.Error(),
		), nil
	}
	log.Printf("tag: %s, ID: %s\n", tag, request.RequestContext.RequestID)

	message := fmt.Sprintf("%s %s released! check the new features on GitHub.\n%s", c.FullRepo, tag, c.RepoURL)
	msg, err := twitterClient.Tweet(message)
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

func response(code int, msg string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: code,
		Body:       fmt.Sprintf("{\"message\":\"%s\"}", msg),
		Headers:    map[string]string{"Content-Type": "application/json"},
	}
}
