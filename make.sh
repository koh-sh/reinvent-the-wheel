#!/usr/bin/env bash

set -eux

docker build -t chatgpt .

if docker container ls | grep -q chatgpt ; then
  docker stop chatgpt
  docker container rm chatgpt
fi
docker run -d -p 8080:8080 -e OPENAI_API_KEY=$OPENAI_API_KEY --name chatgpt chatgpt
