#!/bin/bash

export GIT_HASH=$(git rev-parse HEAD)
docker compose -f ./compose-test.yml up --build -d
echo "Git hash: $GIT_HASH"
