#!/bin/bash

VERSION=$(grep -o '\"version\": \"[^\"]*\"' package.json | sed 's/[^0-9a-z.-]//g'| sed 's/version//g')
LATEST="latest"

# if branch is unstable in git for circle ci
if [ -n "$CIRCLE_BRANCH" ]; then
  if [ "$CIRCLE_BRANCH" != "master" ]; then
    LATEST="$LATEST-$CIRCLE_BRANCH"
  fi
fi

echo "Pushing madejackson/cosmos-server:$VERSION and madejackson/cosmos-server:$LATEST"

sh build.sh

# Multi-architecture build
docker buildx build \
  --platform linux/amd64,linux/arm64 \
  --tag madejackson/cosmos-server:$VERSION \
  --tag madejackson/cosmos-server:$LATEST \
  --push \
  .
