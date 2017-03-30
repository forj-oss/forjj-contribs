#!/bin/sh

if [ "$1" = "" ] || [ "$2" = "" ]
then
   echo "Missing UID and/or GID parameter
   Syntax is $0 UID GID"
   exit 1
fi

sed -i 's/\(devops:x:\)1000:1000/\1'"$1"':'"$2"'/g' /etc/passwd
sed -i 's/\(devops:x:\)1000/\1'"$2"'/g' /etc/group

echo "devops uid: $1, gid: $2"
