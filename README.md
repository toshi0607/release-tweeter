release-tweeter
====

[![Build Status](https://travis-ci.org/toshi0607/release-tweeter.svg?branch=master)](https://travis-ci.org/toshi0607/release-tweeter)
[![Go Report Card](https://goreportcard.com/badge/github.com/toshi0607/release-tweeter)](https://goreportcard.com/report/github.com/toshi0607/release-tweeter)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](https://github.com/toshi0607/release-tweeter/blob/master/LICENSE)
[![Codecov](https://codecov.io/github/toshi0607/release-tweeter/coverage.svg?branch=master)](https://codecov.io/github/toshi0607/release-tweeter?branch=master)

## Description
API to tweet the latest release tag version using AWS Lambda(Golang) & API Gateway.

![](https://user-images.githubusercontent.com/7035446/40571106-f5be6f00-60cd-11e8-935d-0a6c9311d9d1.png)


## Local
You can test AWS Lambda & API Gateway locally with [AWS SAM (Serverless Application Model)](https://github.com/awslabs/serverless-application-model).

### prerequisites

* install [aws-sam-cli](https://github.com/awslabs/aws-sam-cli). Docker is also required. Follow the instruction [here](https://github.com/awslabs/aws-sam-cli#installation).
* set environment variables to [sample-env.json](sample-env.json).
  * You have to go to the [Twitter Apps page](https://apps.twitter.com/) and create new app by yourself. Then you can get *Access Token*, *Access Token Secret*, *Consumer Key (API Key)* and *Consumer Secret (API Secret)*.
  * *REPO* is your repository name on the GitHub. For example, if you want to notice the latest release on [this repository](https://github.com/toshi0607/gig) to your followers, you have to set **toshi0607/gig** to *REPO*.

### hosting

```
# start up API Gateway and Lambda with Docker
$ sam local start-api --env-vars sample-env.json

2018/05/26 11:46:50 Connected to Docker 1.37
2018/05/26 11:46:51 Fetching lambci/lambda:go1.x image for go1.x runtime...
go1.x: Pulling from lambci/lambda
Digest: sha256:79a06e00e85fd2951ad60ee2290590241e8173df5c307134fda061f4a3272bac
Status: Image is up to date for lambci/lambda:go1.x

Mounting main (go1.x) at http://127.0.0.1:3000/ [POST]

You can now browse to the above endpoints to invoke your functions.
You do not need to restart/reload SAM CLI while working on your functions,
changes will be reflected instantly/automatically. You only need to restart
SAM CLI if you update your AWS SAM template.


# request to the localhost from another tab
$ curl -XPOST http://127.0.0.1:3000/
v0.1.4 # and post twitter account's timeline
```


## Production

### prerequisites

You have to prepare credentials with proper policies.

* AWSLambdaFullAccess (should be limited)
* AmazonS3FullAccess (should be limited)
* CloudWatchLogsFullAccess (should be limited)
* AmazonAPIGatewayAdministrator (should be limited)
* AWSXrayFullAccess (should be limited)

In addition to the above, I added group policy for AWS CloudFormation. It's because `sam deploy` command is alias of `aws cloudformation deploy`.

Also, `sam package` command generates CloudFormation template from *template.yml*.

```
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "Stmt**********",
            "Effect": "Allow",
            "Action": [
                "cloudformation:CreateChangeSet",
                "cloudformation:ExecuteChangeSet",
                "cloudformation:DescribeStackEvents",
                "cloudformation:DeleteStack",
                "iam:CreateRole",
                "iam:AttachRolePolicy",
                "iam:DetachRolePolicy",
                "iam:DeleteRole"
            ],
            "Resource": [
                "*"
            ]
        }
    ]
}
```

if you're using credentials you can't change policies by yourself by default, you can add another user profiles.

```
$ aws configure --profile user2
AWS Access Key ID [None]: AKIAI44QH8DHBEXAMPLE
AWS Secret Access Key [None]: je7MtGbClwBF/2Zp9Utk/h3yCo8nvbEXAMPLEKEY
Default region name [None]: us-east-1
Default output format [None]: text
```

if you skip entering output format, default output (json) is used.

```
$ cat ~/.aws/credentials
[default]
aws_access_key_id = AKI********
aws_secret_access_key = ********
[user2]
aws_access_key_id = AKIAI44QH8DHBEXAMPLE
aws_secret_access_key = je7MtGbClwBF/2Zp9Utk/h3yCo8nvbEXAMPLEKEY
region=us-east-1
output=text
``` 

and you can specify profiles explicitly

```
$ aws [some sub command] --profile user2
```

check the [official site](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-started.html) in detail.


### deploy

1. make a s3 bucket

```
$ aws s3 mb s3://lambda-release-tweeter --profile private
```

* an S3 bucket is required. so make a bucket for the release-tweeter you're going to host for the first time
* AWS management console is also avilable to make a bucket instead of the AWS Cli
* a buekcet name (`lambda-release-tweeter` in this case) should be unique globally.

2. prepare artifacts for deploy

```
$ sam package --profile private \
  --template-file template.yml \
  --s3-bucket lambda-release-tweeter \
  --output-template-file packaged.yml
```

* if you use default credentials, you don't have to use *--profile private*


3. deploy artifacts

```
$ sam deploy --profile private \
  --template-file ./packaged.yml \
  --stack-name stack-release-tweeter \
  --capabilities CAPABILITY_IAM \
  --parameter-overrides REPO=$REPO \
    TWITTER_ACCESS_TOKEN=$TWITTER_ACCESS_TOKEN \
    TWITTER_ACCESS_TOKEN_SECRET=$TWITTER_ACCESS_TOKEN_SECRET \
    TWITTER_CONSUMER_KEY=$TWITTER_CONSUMER_KEY \
    TWITTER_CONSUMER_SECRET=$TWITTER_CONSUMER_SECRET
```

* *--parameter-overrides* option is requred because we shouldn't write qredentials to the *template.yml*.
* you have to set 5 environment variables in this case. you can pass value of *REPO*, *TWITTER_ACCESS_TOKEN*, *TWITTER_ACCESS_TOKEN_SECRET*, *TWITTER_CONSUMER_KEY* and *TWITTER_CONSUMER_SECRET* directly instead.
* a stack name (`stack-release-tweeter` in this case) should be unique globally.
