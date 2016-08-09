# Introduction

This forjj plugin is dedicated to build and maintain a jenkins system.

It is based on `docker.hos.hpecorp.net/devops/jenkins-ci` inner source image.

It provides orchestration services for forjj:

- create: infrastructure initial code + README.md in the `infra` repository under `ci` directory.
- build/maintain: Code to build and maintain your CI infra structure as described in the `infra/ci`.

# How to create your jenkins CI infrastructure for your organization?

Using forjj, it is quite easy.

When you create your Organization, you just need to add `--ci jenkins-ci` :

Ex:
```bash
$ forjj create <myorg> --ci jenkins-ci <jenkins_flags> ...
```

If your organization already exist and want to add jenkins:

```bash
$ forjj update <myorg> --ci jenkins-ci <jenkins_flags> ...
```

In case of update, you probably need to follow the organization flow to approve your change and apply.

forjj create/update will call this plugin to create a `ci/<CI Name>` in your infra repository.

This directory will contains source code generated as follow:
- Dockerfile (derived from `docker.hos.hpecorp.net/devops/jenkins-ci`)
- features.lst with basic features:
  - basic jenkins security (admin RW and anonymous R only)
  - proxy settings as found on your current host.
  - seed-job to support Jobs-dsl
  - groovy plugin installed to maintain jenkins configuration.
- Create jenkins-params.sh with any required jenkins parameters and credentials.

Depending on upstream flow choosed, jenkins could generate more files to create jobs (jobs-dsl)


## github upstream with pull-request flow setting.
The github integration will update your `infra/ci/jenkins-ci` with the following code.

- ghprb feature
- 3 Jobs DSL for each project identified under `infra/jobs-dsl/<project>/`
  - `<project>_PR` : Pull request job
  - `<project>_MASTER` : Pull request merge job
  - `<project>_RELEASE` : Master branch tagging job and build code.

### Options

Currently there is no other options.

## Other SCM?

Currently this jenkins FORJJ plugin do not have any other upstream integration.
But this CI orchestrator has been designed to easily add a new one, like gitlab or other flows.
If you want to add you SCM/Jenkins integration, consider contribution to this repository.

For details on contribution, see CONTRIBUTION.md


