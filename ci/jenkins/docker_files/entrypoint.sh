#!/bin/sh -x

if [ "$UID" = "" ] || [ "$GID" = "" ]
then
   echo "Missing UID/GID env variable"
   exit 1
fi

/bin/update_user.sh $UID $GID
exec /bin/su devops -c "/bin/jenkins $*"
