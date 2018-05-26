package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	a "github.com/ChimeraCoder/anaconda"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	repo      = os.Getenv("REPO")
	repoURL   = "https://github.com/" + repo
	latestURL = repoURL + "/releases/latest"

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

type tweeter interface {
	tweet(message string) (string, error)
}

type api a.TwitterApi

func (api api) tweet(message string) (string, error) {
	tweet, err := (a.TwitterApi(api)).PostTweet(message, nil)
	return tweet.Text, err
}

func main() {
	lambda.Start(handler)
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if request.HTTPMethod != "POST" {
		return response(
			http.StatusBadRequest,
			"use POST request",
		), nil
	}

	if hasEmptyEnvVar() {
		return response(
			http.StatusInternalServerError,
			"env vars not set",
		), nil
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", latestURL, nil)
	if err != nil {
		return response(
			http.StatusInternalServerError,
			err.Error(),
		), nil
	}

	resp, err := client.Do(req)
	if err != nil {
		return response(
			http.StatusInternalServerError,
			err.Error(),
		), nil
	}
	defer resp.Body.Close()

	tag, err := getLatestTag(resp.Request.URL.Path)
	if err != nil {
		return response(
			http.StatusNotFound,
			err.Error(),
		), nil
	}
	log.Printf("tag: %s, ID: %s\n", tag, request.RequestContext.RequestID)

	message := fmt.Sprintf("%s %s released! check the new features on GitHub.\n%s", repo, tag, repoURL)
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

func hasEmptyEnvVar() bool {
	return twitterAccessToken == "" ||
		twitterAccessTokenSecret == "" ||
		twitterConsumerKey == "" ||
		twitterConsumerKeySecret == "" ||
		repo == ""
}

func getLatestTag(url string) (string, error) {
	s := strings.Split(url, "/")
	if !strings.Contains(url, "v") {
		return "", fmt.Errorf("%s has no tag", repo)
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
