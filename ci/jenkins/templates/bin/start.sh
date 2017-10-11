#!/bin/sh
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

if [[ "$ADMIN_PWD" != "" ]]
then
   ADMIN="-e SIMPLE_ADMIN_PWD=$ADMIN_PWD"
   unset ADMIN_PWD
   echo "Admin password set."
fi

if [[ "$GITHUB_USER_PASS" != "" ]]
then
   GITHUB_USER="-e GITHUB_PASS=$GITHUB_USER_PASS"
   unset GITHUB_USER_PASS
   echo "Github user password set."
fi

{{ if .Deploy.Ssl.Certificate }}\
if [[ "$CERTIFICATE_KEY" = "" ]]
then
   echo "Unable to set jenkins certificate without his key. Aborted."
   exit 1
fi
echo "$CERTIFICATE_KEY" > .certificate.key
unset CERTIFICATE_KEY
echo "Certificate set."

JENKINS_OPTS='JENKINS_OPTS=--httpPort=-1 --httpsPort=8443 --httpsCertificate=/tmp/certificate.crt --httpsPrivateKey=/tmp/certificate.key'
JENKINS_MOUNT="-v ${SRC}certificate.crt:/tmp/certificate.crt -v ${SRC}.certificate.key:/tmp/certificate.key"

{{ end }}\

if [ "$CONTAINER_IMG" != "" ]
then
    if [ "$CONTAINER_IMG" != "$TAG_NAME" ] && [ "$CONTAINER_IMG" != "$IMAGE_ID" ]
    then
        # TODO: Find a way to stop it safely - Using safe shutdown?
{{/* # Following code will be executed by default if there is no other event driven system (bot/stackstorm/...) */}}\
        echo "#!/bin/sh
sleep 30
docker rm -f {{ .JenkinsImage.Name }}-dood
sleep 2
{{ if .Deploy.Ssl.Certificate }}\
docker run --restart always $DOCKER_DOOD -d -p $SERVICE_PORT:8443 -e \"$JENKINS_OPTS\" $JENKINS_MOUNT --name {{ .JenkinsImage.Name }}-dood $GITHUB_USER $ADMIN $CREDS $PROXY $DOCKER_OPTS $TAG_NAME
{{ else }}
docker run --restart always $DOCKER_DOOD -d -p $SERVICE_PORT:8080 --name {{ .JenkinsImage.Name }}-dood $GITHUB_USER $ADMIN $CREDS $PROXY $DOCKER_OPTS $TAG_NAME
{{ end }}\
echo 'Service is restarted'
rm -f \$0" > do_restart.sh
        chmod +x do_restart.sh

        echo "The image has been updated. It will be restarted in about 30 seconds"
{{/* # End of this code to be executed by default if there is no other event driven system (bot/stackstorm/...) */}}\
        sudo docker run --rm -v $VOL_PWD/do_restart.sh:/tmp/do_restart.sh $DOCKER_DOOD alpine /tmp/do_restart.sh
    else
        echo "Nothing to re/start. Jenkins is still accessible at http://$SERVICE_ADDR:$SERVICE_PORT"
    fi
    exit 0
fi

# No container found. Start it.
{{ if .Deploy.Ssl.Certificate }}\
sudo docker run --restart always $DOCKER_DOOD -d -p $SERVICE_PORT:8443 -e "$JENKINS_OPTS" $JENKINS_MOUNT --name {{ .JenkinsImage.Name }}-dood $GITHUB_USER $ADMIN $CREDS $PROXY $DOCKER_OPTS $TAG_NAME
{{ else }}
sudo docker run --restart always $DOCKER_DOOD -d -p $SERVICE_PORT:8080 --name {{ .JenkinsImage.Name }}-dood $GITHUB_USER $ADMIN $CREDS $PROXY $DOCKER_OPTS $TAG_NAME
{{ end }}\

if [ $? -ne 0 ]
then
    echo "Issue about jenkins startup."
    sudo docker logs {{ .JenkinsImage.Name }}-dood
    return 1
fi
echo "Jenkins has been started and should be accessible at http{{ if .Deploy.Ssl.Certificate }}s{{ end }}://$SERVICE_ADDR:$SERVICE_PORT"
