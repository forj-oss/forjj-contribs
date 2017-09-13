#!/bin/sh -x
#
#

REPO=$LOGNAME
IMAGE_NAME="{{ .JenkinsImage.Name }}"
IMAGE_VERSION=test


# For Docker Out Of Docker case, a docker run may provides the SRC to use in place of $(pwd)
# This is required in case we use the docker -v to mount a 'local' volume (from where the docker daemon run).
if [ "$SRC" != "" ]
then
    VOL_PWD="$SRC"
else
   VOL_PWD="$(pwd)"
fi

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

# A local volume is possible only when we are sure local mount is server dependent.
# Any docker cluster cannot be used with local mount because, the container can be started any where.
# This concerns swarm/ucp/mesos at least.

if [ "$DOCKER_CLUSTER_SRC_VOL" != "" ]
then
    CREDS="-v $DOCKER_CLUSTER_SRC_VOL/jenkins_credentials.sh:/tmp/jenkins_credentials.sh"
else
    # The following works on native docker, Dood and DinD.
    if [ -f jenkins_credentials.sh ] && [ -r jenkins_credentials.sh ]
    then
       CREDS="-v $VOL_PWD/jenkins_credentials.sh:/tmp/jenkins_credentials.sh"
    fi
fi

# For production case, expect
# $LOGNAME set to {{ .Forjj.OrganizationName }}
if [ -f run_opts.sh ]
then
   echo "loading run_opts.sh..."
   source run_opts.sh
fi

# Loading deployment environment ($1)
if [ -f source_$1.sh ]
then
   echo "Loading deployment environment '$1'"
   source source_$1.sh
fi

if [ "$SERVICE_ADDR" = "" ]
then
   echo "SERVICE_ADDR not defined by any deployment environment. Set to '{{.Deploy.Deployment.ServiceAddr}}'"
   SERVICE_ADDR="{{.Deploy.Deployment.ServiceAddr}}"
fi
if [ "$SERVICE_PORT" = "" ]
then
   SERVICE_PORT={{if and .Deploy.Ssl.Certificate (eq .Deploy.Deployment.ServicePort "8080")}}8443 # Default SSL port{{else}}{{.Deploy.Deployment.ServicePort}}{{end}}
   echo "SERVICE_PORT not defined by any deployment environment. Set to '$SERVICE_PORT'"
fi

TAG_NAME={{ .JenkinsImage.RegistryServer }}/$LOGNAME/$IMAGE_NAME:$IMAGE_VERSION

{{/* Docker uses go template for --format. So need to generate a template go string */}}\
CONTAINER_IMG="$(sudo docker ps -a -f name={{ .JenkinsImage.Name }}-dood --format "{{ "{{ .Image }}" }}")"

IMAGE_ID="$(sudo docker images --format "{{ "{{ .ID }}" }}" $IMAGE_NAME)"

if [ "$CONTAINER_IMG" != "" ]
then
    if [ "$CONTAINER_IMG" != "$TAG_NAME" ] && [ "$CONTAINER_IMG" != "$IMAGE_ID" ]
    then
        # TODO: Find a way to stop it safely
        sudo docker rm -f {{ .JenkinsImage.Name }}-dood
    else
        echo "Nothing to re/start. Jenkins is still accessible at http://$SERVICE_ADDR:$SERVICE_PORT"
        exit 0
    fi
fi

if [[ "$ADMIN_PWD" != "" ]]
then
   ADMIN="-e SIMPLE_ADMIN_PWD=$ADMIN_PWD"
fi

{{ if .Deploy.Ssl.Certificate }}\
if [[ "$CERTIFICATE_KEY" = "" ]]
then
   echo "Unable to set jenkins certificate without his key. Aborted."
   exit 1
fi
echo "$CERTIFICATE_KEY" > .certificate.key

JENKINS_OPTS='JENKINS_OPTS=--httpPort=-1 --httpsPort=8443 --httpsCertificate=/tmp/certificate.crt --httpsPrivateKey=/tmp/certificate.key'
JENKINS_MOUNT="-v ${SRC}certificate.crt:/tmp/certificate.crt -v ${SRC}.certificate.key:/tmp/certificate.key"

sudo docker run --restart always -d -p $SERVICE_PORT:8443 -e "$JENKINS_OPTS" $JENKINS_MOUNT --name {{ .JenkinsImage.Name }}-dood $ADMIN $CREDS $PROXY $DOCKER_OPTS $TAG_NAME

{{ else }}
sudo docker run --restart always -d -p $SERVICE_PORT:8080 --name {{ .JenkinsImage.Name }}-dood $CREDS $PROXY $DOCKER_OPTS $TAG_NAME
{{ end }}\

if [ $? -ne 0 ]
then
    echo "Issue about jenkins startup."
    sudo docker logs {{ .JenkinsImage.Name }}-dood
    return 1
fi
echo "Jenkins has been started and should be accessible at http{{ if .Deploy.Ssl.Certificate }}s{{ end }}://$SERVICE_ADDR:$SERVICE_PORT"
