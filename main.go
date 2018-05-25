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
	"github.com/pkg/errors"
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
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
		}, errors.New("use POST request")
	}

	if hasEmptyEnvVar() {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
		}, errors.New("Env Vars not set")
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", latestURL, nil)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
		}, err
	}
	defer resp.Body.Close()

	tag, err := getLatestTag(resp.Request.URL.Path)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotFound,
		}, err
	}
	log.Printf("tag: %s, ID: %s\n", tag, request.RequestContext.RequestID)

	message := fmt.Sprintf("%s %s released! check the new features on GitHub.\n%s", repo, tag, repoURL)
	msg, err := twitterClient.tweet(message)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
		}, err
	}
	if !strings.Contains(msg, tag) {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
		}, errors.New("failed to tweet: " + msg)
	}
	log.Printf("message tweeted: %s, ID: %s\n", msg, request.RequestContext.RequestID)

	return events.APIGatewayProxyResponse{
		Body:       tag,
		StatusCode: http.StatusOK,
	}, nil
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
		return "", errors.Errorf("%s has no tag", repo)
	}
	return s[len(s)-1], nil
}
