#!/bin/bash
#
#

TAG="-t docker.hos.hpecorp.net/forjj-us/github"

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

sudo docker build $PROXY $TAG .
