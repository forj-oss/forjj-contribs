# Introduction

This forjj orchestrator is dedicated to build and maintain a jenkins CI system and any CI integration to other tools like github.

It is based on `docker.hos.hpecorp.net/devops/jenkins-ci` inner source image.

It provides orchestration service for forjj-orchestrator:

- create: infrastructure initial code + README.md in the `infra` repository under `ci` directory.
- build/maintain: Code to build and maintain your CI infra structure as described in the `infra/ci`.

# How to create your jenkins CI infrastructure for your organization?

To initialize a simple local jenkins CI, do the following:

    $ docker run -it --rm docker.hos.hpecorp.net/devops/forjj run ci jenkins-ci

This will create the basic source code in your infra repository under `ci/jenkins-ci`
with basic features:
- basic jenkins security (admin RW and anonymous R only)
- proxy settings as found on your current host.
- seed-job to support Jobs-dsl
- groovy plugin installed to maintain jenkins configuration.

It will create following files:
- Create a Dockerfile based on docker.hos.hpecorp.net/devops/jenkins-ci
- Create features.lst with features.lst
- Create jenkins-params.sh with any required jenkins parameters and credentials.
 
Depending on which repository upstream solution you have choose, it may add the integration code to your `infra` repository.

See [SCM integration][]



# SCM integration

## github
The github integration may update (or create if missing) your `infra/ci/jenkins-ci` with the following code.

- ghprb feature
- 3 Jobs DSL for each project identified under `infra/jobs-dsl/<project>/`
  - `<project>_PR` : Pull request job
  - `<project>_MASTER` : Pull request merge job
  - `<project>_RELEASE` : Master branch tagging job and build code.

### Options

Currently there is no other options.

## Other SCM?

Currently this jenkins orchestrator do not have any other SCM integration. 
But this CI orchestrator has been designed to easily add a new one, like gitlab or 
If you want to add you SCM/Jenkins integration, consider contribution to this repository.

For details on contribution, see CONTRIBUTION.md


