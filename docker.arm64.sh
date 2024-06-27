#!/bin/bash

VERSION=$(grep -o '\"version\": \"[^\"]*\"' package.json | sed 's/[^0-9a-z.-]//g' | sed 's/version//g')
LATEST="latest"

# if branch is unstable in git for circle ci
if [ -n "$CIRCLE_BRANCH" ]; then
  if [ "$CIRCLE_BRANCH" != "master" ]; then
    LATEST="$LATEST-$CIRCLE_BRANCH"
  fi
fi

echo "Pushing madejackson/cosmos-server:$VERSION and madejackson/cosmos-server:$LATEST"

sh build.arm64.sh

docker build \
  -t madejackson/cosmos-server:$VERSION-arm64 \
  -t madejackson/cosmos-server:$LATEST-arm64 \
  -f dockerfile.arm64 \
  --platform linux/arm64 \
  .

docker push madejackson/cosmos-server:$VERSION-arm64
docker push madejackson/cosmos-server:$LATEST-arm64