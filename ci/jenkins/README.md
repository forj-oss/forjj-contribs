# Introduction

Jenkins FORJJ plugin generate runnable source code to create and maintain a simple jenkins system.

By default, the source code generated from templates implements **Jenkins 2.x** with **pipeline**, several **jenkins plugins**, **basic security rights** and a collection of scripts to enhance jenkins management from code perspective.

Mainly it offers:

- Add/Maintain new project repository to generate jobs
  (Jenkinsfile + JobsDSL).
- Add/Maintain Jenkins plugins/Jenkins features
- Deploy Jenkins (master) on Docker/Swarm/UCP/Mesos/DCOS and scale it.

# How to create your jenkins CI infrastructure for your organization?

Using forjj, it is quite easy.

To create a new jenkins instance to your organization, you just need to add `--apps ci:jenkins` :

Ex:
```bash
$ forjj create --apps ci:jenkins <jenkins_flags> ...
```

If you want to update an existing jenkins instance, you can use:

```bash
$ forjj update --apps ci:jenkins <jenkins_flags> ...
```

If you have several jenkins instances, you can add it in the `--apps` flag:

Ex:
```bash
$ forjj ... --apps ci:jenkins:myinstance,ci:jenkins:anotherinstance <jenkins_flags> ...
```

Each instance requested will have his collection of jenkins flags prefixed by the instance name.

Ex:
```bash
$ forjj ... --apps ci:jenkins:myinstance,ci:jenkins:anotherinstance --myinstance-service-addr myinstance.com --anotherinstance-service-addr ...
```

In case of update, you probably need to follow the organization flow to approve your change and apply.

## Jenkins source Templates

All jenkins source files are generated from a collection of source templates (jenkins source model).
Currently, those templates are located under templates directory in forjj-jenkins container.

TODO: We can imagine having several templates directory as well as a different source of templates (git, tar, others...) to change jenkins sources model, but this has not been currently developped.

Feel free to contribute to add this feature!

The `templates/templates.yaml` defines how to generate the source model from a deploy perspective.

The template mechanisms implemented is based on [golang template](https://golang.org/pkg/text/template/).
The template data is given by the forjj-jenkins.yaml source file. The data structure is defined in [this go source file](jenkins_plugin.go#34)

You can update this file manually and ask forjj to update source files.

TODO: forjj-jenkins is a go binary exposing his service to forjj through a REST API. But we can image that this binary become available to simply regenerate source file from `forjj-jenkins.yaml`.
Today you must use forjj update --apps ci:jenkins to call the plugin and regenerate source files from `forjj-jenkins.yaml`.

Feel free to contribute to add this feature!

### forjj Jenkins source model

Currently, the embedded source model implements globally the following:

- A docker image built from `hub.docker.io/forj/jenkins-ci` [source](https://github.com/forj-oss/jenkins-ci)
- A collection of default features ([source](https://github.com/forj-oss/jenkins-install-inits))
  - Basic authentication (admin user with default password & anonymous has read access)
  - proxy setting (Set proxy from http_proxy env setting, found from the container)
  - seed-job (One job generated to populate the other collection of jobs/pipelines)
  - jenkins slave fixed port
- A collection of additional features and templates to add for a dedicated deployment
- A list of predefined deployment. ie:
  - docker - To deploy to your local docker environment. (Default deployment)
  - ucp - To deploy to a UCP system.
  - marathon - To deploy to dcos/mesos marathon.

This list of elements are not exhaustive and can be updated time to time. Please refer to the (templates.yaml)[templates/templates.yaml] for latest updates.

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


