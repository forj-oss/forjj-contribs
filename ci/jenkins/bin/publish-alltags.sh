// This file has been created by "go generate" as initial code. go generate will never update it, EXCEPT if you remove it.

// So, update it for your need.
#!/bin/bash
#
# This script is used to publish officially all released docker images (tagged)
#
# Release workflow is:
#
# - Someone fork and create a tag release then submit a PR.
# - GitHub jenkins can be started to start an 'ITG' image validation
# - The repo maintainer at some time will accept the new release.
# - Github should send a jenkins job to build officially this new release
#   I expect to get this info in $1 (Release number)

# Then this job should implement the following code in jenkins
# And jenkins-ci images for each flavors will be officially pushed to the internal registry.

TAG_BASE="$(awk '$1 ~ /image:/ { print $2 }' jenkins.yaml)"

if [ ! -f releases.lst ]
then
   echo "VERSION or releases.lst files not found. Please move to the repo root dir and call back this script."
   exit 1
fi

case "$1" in
  release-it )
    VERSION=$(eval "echo $(awk '$1 ~ /version:/ { print $2 }' jenkins.yaml)")
    if [ "$(git tag -l $VERSION)" = "" ]
    then
       echo "Unable to publish a release version. git tag missing"
       exit 1
    fi
    COMMIT="$(git log -1 --oneline| cut -d ' ' -f 1)"
    if [ "$(git tag -l --points-at $COMMIT | grep $VERSION)" = "" ]
    then
       echo "'$COMMIT' is not tagged with '$VERSION'. Only commit tagged can publish officially this tag as docker image."
       exit 1
    fi
    VERSION_TAG=${VERSION}_
    ;;
  latest )
    VERSION=latest
    VERSION_TAG=latest_
    ;;
  *)
    echo "Script used to publish release and latest code ONLY. If you want to test a fork, use build. It will create a local docker image jenkins-dood:test"
    exit 1
esac

if [ ! -f releases.lst ]
then
    echo "Missing releases.lst or not in the plugin source path."
    exit 1
fi

cat releases.lst | while read LINE
do
   [[ "$LINE" =~ ^# ]] && continue
   TAGS="$(echo "$LINE" | awk -F'|' '{ print $2 }' | sed 's/,/ /g')"
   echo "=============== Building forjj-us/github"
   $(dirname $0)/build.sh
   echo "=============== Publishing tags"
   for TAG in $TAGS
   do
      echo "=> $TAG_BASE:$TAG"
      sudo docker tag -f $TAG_BASE $TAG_BASE:$TAG
      sudo docker push $TAG_BASE:$TAG
   done
   echo "=============== DONE"
done
