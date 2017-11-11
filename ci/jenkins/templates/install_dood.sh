#!/usr/bin/env bash

function group_add() {
if [[ -f /etc/redhat-release ]]
then
   group_add_redhat $1
   return
fi
echo "Unable to add group in this image. Unsupported linux release. Think to contribute to fix it."
exit 1
}

function group_add_redhat() {
if [[ ! -x /usr/sbin/groupadd ]]
then
   yum install /usr/sbin/groupadd -y
   yum clean all
fi

groupadd docker -g $1
echo "group docker ($1) added."
}

if [[ "$1" = "" ]]
then
   echo "No dood installed: DOOD_DOCKER_GROUP is missing"
   exit
fi

GID=$1

set -e

GROUP=$(grep ":$1:" /etc/group | awk -F: '{printf $1}')

if [[ "$GROUP" = "" ]]
then
   group_add $GID
   GROUP=docker
fi

usermod jenkins -a -G $GROUP

echo "Added $GROUP($GID) to jenkins groups - used by /var/run/docker.sock"
