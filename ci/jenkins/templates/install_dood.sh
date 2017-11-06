#!/usr/bin/env bash

if [[ "$1" = "" ]]
then
   echo "No dood installed: DOOD_DOCKER_GROUP is missing"
   exit
fi

GID=$1

usermod jenkins -a -G $GID

echo "Added $GID to jenkins groups."
