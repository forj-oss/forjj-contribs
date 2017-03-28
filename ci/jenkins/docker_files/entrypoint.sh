#!/bin/sh -x

/bin/update_user.sh $UID $GID
exec /bin/su devops -c "/bin/jenkins $*"
