package main

import (
	"net/url"
	"testing"

	a "github.com/ChimeraCoder/anaconda"
)

func TestHandler(t *testing.T) {
	postTweetFunc = func(api a.TwitterApi) postTweet {
		return postTweetHelper
	}
	Handler()
}

func postTweetHelper(s string, v url.Values) (a.Tweet, error) {
	return a.Tweet{Text: s}, nil
}
