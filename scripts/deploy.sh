#!/bin/bash
set -eu

sam package --profile private \
  --template-file template.yml \
  --s3-bucket lambda-release-tweeter \
  --output-template-file packaged.yaml

sam deploy --profile private \
  --template-file ./packaged.yaml \
  --stack-name stack-release-tweeter \
  --capabilities CAPABILITY_IAM \
  --parameter-overrides REPO=toshi0607/release-tweeter \
    TWITTER_ACCESS_TOKEN=${TWITTER_ACCESS_TOKEN} \
    TWITTER_ACCESS_TOKEN_SECRET=${TWITTER_ACCESS_TOKEN_SECRET} \
    TWITTER_CONSUMER_KEY=${TWITTER_CONSUMER_KEY} \
    TWITTER_CONSUMER_SECRET=${TWITTER_CONSUMER_SECRET}
