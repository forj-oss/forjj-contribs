#!/usr/bin/env bash

set -e

if [ "$BUILD_ENV_LOADED" != "true" ]
then
   echo "Please go to your project and load your build environment. 'source build-env.sh'"
   exit 1
fi

TAG="-t $(awk '$0 ~ /^ *image: ".*"$/ { print $0 }' github.yaml | sed 's/^ *image: "*\(.*\)".*$/\1/g')"

cd $BUILD_ENV_PROJECT

create-build-env.sh

if [ "$GOPATH" = "" ]
then
    echo "Unable to build without GOPATH. Please set it (build)env or your local personal '.be-gopath')"
    exit 1
fi

glide i

# Requires forjj to be static.
export CGO_ENABLED=0
go build

$BUILD_ENV_DOCKER build $PROXY $DOCKERFILE $TAG .
