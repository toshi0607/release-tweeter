#!/bin/bash
set -eu

echo current version: $(gobump show -r)

GOOS=linux GOARCH=amd64 go build -o release-tweeter
zip release-tweeter.zip ./release-tweeter

ghr v$(gobump show -r) release-tweeter.zip
