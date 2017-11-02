#!/usr/bin/env bash

DOCKER_VERSION=$1

# setup Docker repository
apt-get -qq update
apt-get -qq install -y \
     apt-transport-https \
     ca-certificates \
     curl \
     gnupg2 \
     software-properties-common
curl -fsSL https://download.docker.com/linux/$(. /etc/os-release; echo "$ID")/gpg | apt-key add -
add-apt-repository \
   "deb [arch=amd64] https://download.docker.com/linux/$(. /etc/os-release; echo "$ID") \
   $(lsb_release -cs) \
   stable"
# install Docker CE
apt-get -q update -y
apt-get -q install -y docker-ce=$DOCKER_VERSION
