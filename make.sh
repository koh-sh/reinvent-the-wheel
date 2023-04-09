#!/usr/bin/env bash

set -eux

GOOS=linux GOARCH=arm64 go build ./chatgpt.go
docker build -t chatgpt .

if docker container ls | grep -q chatgpt ; then
  docker stop chatgpt
  docker container rm chatgpt
fi
docker run -d -p 8080:8080 -e OPENAI_API_KEY=$OPENAI_API_KEY --name chatgpt chatgpt
rm chatgpt
