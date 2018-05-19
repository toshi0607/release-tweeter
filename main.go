package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"net/http"
	"github.com/prometheus/common/log"
	"fmt"
	"strings"
)

var (
	REPO = "toshi0607/gig"
	LATEST_URL = "https://github.com/" + REPO + "/releases/latest"
)

func main() {
	lambda.Start(Handler)
}

func Handler() {
	client := &http.Client{}
	req, err := http.NewRequest("GET", LATEST_URL, nil)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	fmt.Println(getLatestTag(resp.Request.URL.Path))
}

func getLatestTag(url string) string {
	s := strings.Split(url, "/")
	return s[len(s)-1]
}
