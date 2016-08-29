# This file is managed by Forjj jenkins plugin at Maintain phase. Do not update it manually!

FROM docker.hos.hpecorp.net/devops/jenkins-dood:0.1_2.7

COPY local_plugins/*  $JENKINS_DATA_REF/plugins/
COPY features.lst /tmp

USER root

ARG JENKINS_INSTALL_INITS_URL=https://github.hpe.com/RnDLabsIT/jenkins-install-inits/raw/master/

RUN /usr/local/bin/jenkins-install.sh /tmp/features.lst

RUN chown -R jenkins:jenkins $JENKINS_HOME /usr/share/jenkins/ref

USER jenkins
