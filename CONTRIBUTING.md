**POC in development**

# Introduction

As this is a POC, a lot of thing will change on the core and contributions areas.

So, if you want to share any feedback or even propose ideas/code, I would suggest first to discuss and see how and where you could participate in this idea.

In any case, you can of course provide a PR and suggest an idea that way.

# Why should I create a new FORJJ plugin?

Because you need it??? :)

FORJJ core is currently managing only git repositories (git push, mainly) but time to time, I believe this could be moved to the SCM part, definitely...
And FORJJ call 3/4 differents actions to create and manage the infra from code perspective.

For details about FORJJ actions, [read this page](https://github.hpe.com/christophe-larsonneur/forjj/blob/master/CONTRIBUTION.md)

For now, I'm focused on jenkins and github. So, I'm writing those drivers and the link between them (ie github PR flow)
If you want to participate in this, but with more tools or others CI/SCM Upstream tools, then you will need to write your FORJJ plugin.

# How to write a new FORJJ plugin?

The plugin notion in FORJJ is simple and basic.

FORJJ do his task by starting a FORJJ container, with a simple task name and flags.
- Flags are defined by the [yaml description file](https://github.hpe.com/christophe-larsonneur/forjj-contribs#description-of-yaml)
- Tasks are defined [here](https://github.hpe.com/christophe-larsonneur/forjj-contribs#description-of-yaml)

When FORJJ container is started, it will interpret task and flags and call the driver.

The drivers should implement following tasks:
  - `check`   : It provide a status information in json format.
  - `create`  : It create requested resource. Create is used once. Return a json data.
  - `update`  : It update the infra repository. Returns a json data
  - `maintain`: It update the infra with the data from the infra repository.

**Plugin json data format:** 
  - `repos/[]/name`     : repository configured.
  - `repos/[]/upstream` : upstream data updated/created
  - `repos/[]/config`   : driver configuration filename updated
  - `services/[]url`    : web service url. ex: https://github.hpe.com
  - `status`            : Driver task status message
  - `state_code`        : driver task status code. 200 OK
  - `error_message`     : driver error message

**In details, what your plugin needs to do?**

## Plugin check task

### Input
The plugin check task get input data from the command line flags.
The output is conform to the plugin json data format described below.

### Role
The main role of `check` is to ensure that the application can be used properly for the complete DevOps solution.
The upstream check task check if the service up and running and if repository given already exist on the upstream server.
The others plugin type check task check the service up and running. If it needs to add some configuration test, it should do it here as well.

### output
If found, the json must return at least `state_code = 200` and `repos/[]/name`
If not found; it should respond with `state_code = 200` and `error_message`

The upstream plugin, should provide at least `state_code = 200`, collection of repos with `repos/[]/name`, `repos/[]/upstream`, `repos/[]/config` and the collection of services url `services/[]url`

## Plugin create task

### Input
The plugin create task get input data from the command line flags.
The output is conform to the plugin json data format described below.

### Role
`create` is a little bit special case.
It is mainly related to the upstream (for now), in order to install the solution the first time.
So, it will install the upstream application if needed, then configure it, with organization, infra repos, optionaly credentials (users/groups).
At the end of that, the plugin should report the list of repositories, clone input and services.
If at least `infra_repos` gets created successfully, with the --repos added, `create` will also create additional repositories.

Next time, if there is need to create new repos, delete them, it has to follow the update/maintain process.

### internal data

The plugin can store any kind of data that is fully understandable by the plugin to maintain the upstream. FORJJ won't read it anymore, but the file will be considered as source and commited.
This one is going to be the source code of the upstream.
It must be stored under `repos/<infra_repos>/<any_kind_of_file>`. The file name can be anything the plugin want. But it has be stable and reported in `repos/[]/config` output.

### output
If at least the infra repos was already created, the plugin should return an error message.
A successful infra repos creation will return at least `state_code = 200`, collection of repos with `repos/[]/name`, `repos/[]/upstream`, `repos/[]/config` and the collection of services url `services/[]url`

## Plugin update task

### input
The plugin `update` task get input data from the command line flags.

### role

The plugin role in update task is mainly to update the infra repository with appropriate source to start and configure any application.
The plugin should NEVER update the application itself. It will only update the source.

### output

The output is conform to the json data format.

## Plugin maintain task

### input
The plugin maintain task get input data from the command line flags. Usually, the supported flags is quitely limited to where is the infra repository.

### role

The main role in maintain context is to instantiate and configure the application like jenkins for the jenkins-ci plugin. 
This activity is typically made by any kind of orchestration tools, like ansible/puppet or any other kind of tools you like to use.

Ansible or puppet are currently not installed in the FORJJ container. But this could be done, if FORJJ can support different containers.

### output

The docker container has been created to ensure your plugin environment is identified.


## docker image

Today, there is no cli option to choose the docker image. So it uses docker.hos.hpecorp.net/devops/forjj:latest built from [this Docker directory](https://github.hpe.com/christophe-larsonneur/forjj/tree/master/docker)
The forjj cli start only one container and FORJJ itself is not designed to start several containers, one per plugin. This is something that could be changed later so it will simplify the implementation of such feature.

# FORJJ container internal

when the container starts, it will do the following:

*Create*
A create is possible and succeed if the upstream driver as created a repository.

- Ensure the infra repository is in the workspace. The upstream driver is called to check the existence of the infra repository.
  If the infra repo exist and has already defined infra repos (yaml exists), create fails.
  Locally, FORJJ can :
  - TODO: Clone from a URL
  - Create a new one, empty with 2 directories : `repos/<organization>-infra`
  - TODO: Keep an existing directory to migrate to git. 
    TODO: We can imagine to introduce a migration step here to get source code migrated to git.
  - TODO: Keep an unknown existing cloned repository to migrate from an external GIT upstream environment to the one FORJJ will manage.
    TODO: We can imagine to introduce a migratioin step here to cleanup GIT commits.

  A second repository is going to be created with service information at least. It will be named suffixed by `-state`: `repos/<organization>-infra-state`
- Start the upstream driver with 'create' task to create the GIT upstream configuration
  - It depends on the upstream driver to properly make a valid clonable/fetchable remote upstream.
    the driver can install the upstream application (gerrit/gitlab/...) if it is required.
    The driver should return a json data, to the standard output as follow:
    It should write his configuration file in the infra repo directory.

    a --infra is passed to properly create those 2 infra repositories.

    ```json
{ "repos": { 
     "name": "<organization>-infra", 
     "upstream": "git@...", 
     "config": "<organization>-infra/github.yaml" 
     }, {
     "name": "<organization>-infra-state",
     "upstream": "git@...", 
     "config": "<organization>-infra-state/github.yaml" 
     },
  "services": {
      "upstream": "https://github.hpe.com"
     }
  "state_code": "200", 
  "status": "2 repositories, 1 organization created." }
    ```

    If an issue occurs, the standard out is used formatted in json:

    ```json
{ "state_code": "404", "error_message": "An error occured..." }
    ```


- update/create the origin remote then pull it.
Next will occur only if the repository is empty. (no commit with `repos/<organization>-infra/forj.yaml`)
- Read the json data driver output, save it in `repos/<organization>-infra/forj.yaml` if no `error_message` is reported
- The service is stored in the repository suffixed with `-state`. It will contains the current known service configured links. Jobs logs could be stored here.
- Add it, add the driver configuration file and commit it all.
- Then push it.

As soon as this commit has been created in the upstream, the create is successful and is definitely done. If any change needs to occur, it will be done by the couple 'update'/'maintain' tasks.

If any additionnal repositories are requested. (--repos), the 'create' task will pursue by creating those on the upstream. like it was with the infra.

*update*

when the container is started the *update* task, the driver is called to update the upstream configuration file ONLY.
It should return the json output.

*maintain* 

When the container is started the *maintain* task, the driver is called to update the upstream service configuration to reflect the configuration data.
The container will do:
- git clone
- git pull to get latest update
- call the driver to update the upstream service. json returned.

FORJ Team
