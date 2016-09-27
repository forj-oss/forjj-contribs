#!/bin/sh -x
#
#

REPO=$LOGNAME
IMAGE_NAME={{ .JenkinsImage.FinalDockerImage }}
IMAGE_VERSION=test

if [ "$http_proxy" != "" ]
then
   PROXY=" --env http_proxy=$http_proxy --env https_proxy=$https_proxy --env no_proxy=$no_proxy"
   echo "Using your local proxy setting : $http_proxy"
   if [ "$no_proxy" != "" ]
   then
      PROXY="$PROXY -e no_proxy=$no_proxy"
      echo "no_proxy : $no_proxy"
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

TAG_NAME={{ .JenkinsImage.FinalDockerRegistryServer }}/$LOGNAME/$IMAGE_NAME:$IMAGE_VERSION

{{/* Docker uses go template for --format. So need to generate a template go string */}}\
CONTAINER_IMG="$(sudo docker ps -f name={{ .JenkinsImage.FinalDockerImage }}-dood --format "{{ "{{ .Image }}" }}")"

IMAGE_ID="$(sudo docker images --format "{{ "{{ .ID }}" }}" $IMAGE_NAME)"

if [ "$CONTAINER_IMG" != "" ]
then
    if [ "$CONTAINER_IMG" != "$TAG_NAME" ] && [ "$CONTAINER_IMG" != "$IMAGE_ID" ]
    then
        # TODO: Find a way to stop it safely
        sudo docker rm -f {{ .JenkinsImage.FinalDockerImage }}-dood
    else
        echo "Nothing to re/start. Jenkins is still accessible at http://{{ .Deploy.ServiceAddr }}:{{ .Deploy.ServicePort }}"
        exit 0
    fi
fi

sudo docker run -d -p 8080:{{ .Deploy.ServicePort }} --name {{ .JenkinsImage.FinalDockerImage }}-dood $CREDS $PROXY $DOCKER_OPTS $TAG_NAME

if [ $? -ne 0 ]
then
    echo "Issue about jenkins startup."
    sudo docker logs {{ .JenkinsImage.FinalDockerImage }}-dood
    return 1
fi
echo "Jenkins has been started and should be accessible at http://{{ .Deploy.ServiceAddr }}:{{ .Deploy.ServicePort }}"
