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
	REPO       = os.Getenv("REPO")
	REPO_URL   = "https://github.com/" + REPO
	LATEST_URL = REPO_URL + "/releases/latest"

	TWITTER_ACCESS_TOKEN                = os.Getenv("TWITTER_ACCESS_TOKEN")
	TWITTER_ACCESS_TOKEN_SECRET         = os.Getenv("TWITTER_ACCESS_TOKEN_SECRET")
	TWITTER_CONSUMER_KEY                = os.Getenv("TWITTER_CONSUMER_KEY")
	TWITTER_CONSUMER_SECRET             = os.Getenv("TWITTER_CONSUMER_SECRET")
	twitterClient               tweeter = api(*a.NewTwitterApiWithCredentials(
		TWITTER_ACCESS_TOKEN,
		TWITTER_ACCESS_TOKEN_SECRET,
		TWITTER_CONSUMER_KEY,
		TWITTER_CONSUMER_SECRET,
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
	lambda.Start(Handler)
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if hasEmptyEnvVar() {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
		}, errors.New("Env Vars not set")
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", LATEST_URL, nil)
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

	tag := getLatestTag(resp.Request.URL.Path)
	log.Printf("tag: %s, ID: %s\n", tag, request.RequestContext.RequestID)

	message := fmt.Sprintf("%s %s released! check the new features on GitHub.\n%s", REPO, tag, REPO_URL)
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
	return TWITTER_ACCESS_TOKEN == "" ||
		TWITTER_ACCESS_TOKEN_SECRET == "" ||
		TWITTER_CONSUMER_KEY == "" ||
		TWITTER_CONSUMER_SECRET == "" ||
		REPO == ""
}

func getLatestTag(url string) string {
	s := strings.Split(url, "/")
	return s[len(s)-1]
}
