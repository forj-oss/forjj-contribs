#!/usr/bin/bash

if [ "$(type yum 2>/dev/null)" != "" ]
then
   yum install openssl -y && \
   yum clean all
   exit $?
fi

if [ "$(type apt-get 2>/dev/null)" != "" ]
then
   apt-get update && \
   apt-get install openssl -y && \
   rm -fr /var/lib/apt/lists/*
   exit $?
fi

if [ "$(type apk 2>/dev/null)" != "" ]
then
   apk --no-cache add openssl -y
   exit $?
fi
