#!/bin/sh
#
# Script to move the docker registry content to ~root/.docker/config.json
#

CUR_USER=`id -un`

if [ "$CUR_USER" != root ]
then
   echo "$1 requires to be root. Aborted."
   exit 1
fi

if [ ! -d /root/.docker ]
then
   mkdir /root/.docker
fi

DOCKER_FILE=/root/.docker/config.json
SRC_DOCKER_FILE=/tmp/docker_config.json

if [ -f $SRC_DOCKER_FILE ]
then
   cat $SRC_DOCKER_FILE > $DOCKER_FILE
   echo "$DOCKER_FILE updated."
else
   echo "No $SRC_DOCKER_FILE found. Aborted."
fi
