# This file is managed by Forjj jenkins plugin at Maintain phase. Do not update it manually!

FROM {{ .Dockerfile.FromImage }}{{ if .Dockerfile.FromImageVersion }}:{{ .Dockerfile.FromImageVersion }}{{ end }}

# COPY local_plugins/*  $JENKINS_DATA_REF/plugins/
COPY features.lst /tmp

ENV JENKINS_SLAVE_AGENT_PORT=50000{{ if .Projects }} \
    SEED_JOBS_REPO={{ .Projects.DslRepo }}/{{ .Projects.DslPath }}{{ end }}

USER root

ARG JENKINS_INSTALL_INITS_URL=https://github.com/forj-oss/jenkins-install-inits/raw/master/

RUN /usr/local/bin/jenkins-install.sh /tmp/features.lst

RUN chown -R jenkins:jenkins $JENKINS_HOME /usr/share/jenkins/ref

USER jenkins
