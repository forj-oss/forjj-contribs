#!/bin/bash
#
#

REPO=$LOGNAME
IMAGE_NAME={{ .Docker.Name }}
IMAGE_VERSION=test

sudo docker rm -f {{ .Docker.Name }}-dood

if [ "$http_proxy" != "" ]
then
   PROXY=" --env http_proxy=$http_proxy --env https_proxy=$https_proxy --env no_proxy=$no_proxy"
   echo "Using your local proxy setting : $http_proxy"
   if [ "$no_proxy" != "" ]
   then
      PROXY="$PROXY -e no_proxy=$no_proxy"
      echo "no_proxy : $http_proxy"
   fi
fi

if [ -e jenkins_credentials.sh ]
then
   CREDS="-v $(pwd)/jenkins_credentials.sh:/tmp/jenkins_credentials.sh"
fi

# For production case, expect
# $LOGNAME set to {{ .Forjj.OrganizationName }}
if [ -f run_opts.sh ]
then
   echo "loading run_opts.sh..."
   source run_opts.sh
fi

TAG_NAME=docker.hos.hpecorp.net/$LOGNAME/$IMAGE_NAME:$IMAGE_VERSION

sudo docker run -p 8080:{{ .Deploy.ServicePort }} -it --name {{ .Docker.Name }}-dood $CREDS $PROXY $DOCKER_OPTS $TAG_NAME
