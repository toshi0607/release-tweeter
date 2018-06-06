package twitter

import (
	"errors"

	a "github.com/ChimeraCoder/anaconda"
)

// Tweeter is interface to tweet a message
type Tweeter interface {
	Tweet(message string) (string, error)
}

// Client represents twitter client
type Client struct {
	*a.TwitterApi
}

// Tweet posts message to the twitter time line
func (api Client) Tweet(message string) (string, error) {
	tweet, err := api.PostTweet(message, nil)
	return tweet.Text, err
}

// NewClient creates a new twitter client from credentials
func NewClient(twitterAccessToken, twitterAccessTokenSecret, twitterConsumerKey, twitterConsumerKeySecret string) (*Client, error) {
	if twitterAccessToken == "" ||
		twitterAccessTokenSecret == "" ||
		twitterConsumerKey == "" ||
		twitterConsumerKeySecret == "" {
		return nil, errors.New("twitter: all of the credential is required")
	}

	return &Client{a.NewTwitterApiWithCredentials(
		twitterAccessToken,
		twitterAccessTokenSecret,
		twitterConsumerKey,
		twitterConsumerKeySecret,
	)}, nil
}
