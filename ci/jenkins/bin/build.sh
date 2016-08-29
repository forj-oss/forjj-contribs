#!/bin/bash
#
#

TAG="-t $(awk '$0 ~ /^ *image: ".*"$/ { print $0 }' jenkins.yaml | sed 's/^ *image: "*\(.*\)".*$/\1/g')"

echo "Local go build, then create a docker image..."
CGO_ENABLED=0 go build
if [ $? -ne 0 ]
then
   exit 1
fi

if [ "$http_proxy" != "" ]
then
   PROXY=" --build-arg http_proxy=$http_proxy --build-arg https_proxy=$https_proxy --build-arg no_proxy=$no_proxy"
   echo "Using your local proxy setting : $http_proxy"
   if [ "$no_proxy" != "" ]
   then
      PROXY="$PROXY --build-arg no_proxy=$no_proxy"
      echo "no_proxy : $http_proxy"
   fi
fi

sudo docker build $PROXY $DOCKERFILE $TAG .
