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
  - `create`  : It create requested resource. Create is used once. Return a json data.
  - `update`  : It update the infra repository. Returns a json data
  - `maintain`: It update the infra with the data from the infra repository.

A plugin implementing the REST API gets his input in json, and generate output in json.
A plugin implementing the shell gets this input in 2 forms:
`args` is passed as arguments of the plugin binary, prefixed by --. Ex: github-token => --github-token
`reposdata` is encoded in json and transmitted as a string to --data option.
The shell plguin generate the output in json, like the REST API plugin.

**Plugin Input format**
```yaml
args:        : Collection of actions (create/update/maintain) arguments as defined in your `<plugin>.yaml`
  <argumentName>:<ArgumentValue>
# Following Structure data is sent to an upstream driver.
reposdata:   : Repository Parameters for an upstream driver.
  templates: : List of templates to apply to the repository
  - "<templateName>"
  title      : Repo title
  Users:     : List of users and roles attached.
    <user>:"<roles>"
  Groups:    : List of users and roles attached.
    <group>:"<roles>"
  flow       : Flow name to apply
  instance   : Instance owning the upstream repo.
  options:   : Collection of options to add to the repo. This list is defined by the plugin.
    <Name>:<Value>
```

**Plugin json output data format:**
```yaml
data:                : Plugin data structure
Following `repos:` and `reposdata:` structures must be returned by an upstream driver.
  repos:             : List of Repos Name, containing Repo data.
    <Name>:          : repository name.
      <Name>         : repository name.
      remotes        : List of Remote name and remote url attached.
       <RemoteName>:<RemoteUrl>
      branchconnect:`: List of local branch attached to the upstream branch.
       <LocalBranchName>:<UpstreamBranchName>
     exist           : True if the resource exist. Otherwise false.
  reposdata:         : Repository Parameters for an upstream driver.
    templates:       : List of templates to apply to the repository
    - <templateName>
    title            : Repo title
    Users:           : List of users and roles attached.
      <user>:<roles>
      <group>:<roles>
    flow             : Flow name to apply
    instance         : Instance owning the upstream repo.
    options:         : Collection of options to add to the repo. This list is defined by the plugin.
      <Name>:<Value>
Following structures are returned by all plugins
  services:          : web service url. ex: https://github.hpe.com
    urls:
     <url>
files:               : List of driver files managed
- <file>
status               : Driver task status message
state_code           : driver task status code.
state_code: 200      : Tasks executed without issue. `status` should be non empty.
state_code: 419      : Task aborted. `error_message` must be non empty.
state_code: 422      : Task failure. `error_message` must be non empty.
error_message        : driver error message
```

**In details, what your plugin needs to do?**

## Plugin create task

### Input
The plugin create task get input data from the forjj command line flags and posted as data in json (REST API) or shell flags.
If the plugin is an upstream type, forjj add `reposdata` structure to the request.
The list of flags is defined in your plugin yaml file, under section `flags/create:`

### Role

#### Generic driver
`create` will generate some application source files used to install the application if needed, then configure it.

At the end of that, the plugin should report the list of
- generated source files and
- services.

#### Upstream driver

In case of an upstream driver, Forjj will send out a collection of repository that the usptream driver will need to create.
At the end of that, the plugin should report the list of
- generated source files,
- repositories created,
- upstream information and
- services.

**NOTE:**
When forjj creates a new DevOps Environment, it will create an `infra-repo` repository. This repo is the first repository that Forjj will need to create and will be used to stored all plugins source files.
Your upstream do not need to do any extra task to handle this. Forjj will call your driver create and maintain to do this.

Next time, if there is need to create new repos, delete them, it has to follow the update/maintain process.

### internal data

The plugin generate/store files in the `forjj-srcmount` place. Those files should be text files, human readable to permit any manual update from Dev/Ops teams.

When source code generated is completed, they must be listed in the output json data (`files`)

### output

The output is generated in json.

By default, the plugin must return `status`, `state_code = 200` and `services/urls/...`

An `upstream` plugin must return also `repos` and `reposdata`.

In case of abort situation, `state_code` must be `419`. `message` must not be empty.
In that case, forjj will not interrupt the create task and will move to the next driver call.

In case of errors situation, `state_code` must be `4xx` except `419`. `message` must not be empty.
In that case, forjj will interrupt the create task and exit with the message you returned.

## Plugin update task

### input
The plugin `update` task get input data from the command line flags and workspace yaml file. Then it posted it as data in json (REST API) or shell flags.

The list of flags is defined in your plugin yaml file, under section `flags/update:`

Usually, the `update` input flags is really close to what create has.

### role

The plugin role in update task is mainly to update the infra repository by updating the generated source files.

The plugin should NEVER update the application itself, as this is the role of maintain to update the real application.

### output

The output is conform to the json output data format.

## Plugin maintain task

### input
The plugin maintain task get input data from the command line flags, workspace yaml file, forjj option files and credential file. Then it posted it as data in json (REST API) or shell flags.

The maintain flags defined in your `<plugin>.yaml` file is not exposed automatically to forjj as a plugin flag (like create/update do)
When forjj starts maintain task, a short list of flags are authorized.
If your plugin requires some `critical` data to connect to the service (like credentials, token, etc...) you need to define then in the `flags/maintain` section of your `<plugin>.yaml` file and define the same option in `flags/create` or `flags/update` or both.

Forjj will detect this and stored the maintain flags to the `infra-repo` and create a local credential yaml file which will contains your critical flags. Usually, DevOps team will this this file in a secure place and passed it to the `forjj maintain` `--file` flag.

This is totally transparent for your plugin. The critical data will be stored in `args/` as usual.

### role

The main role in maintain task is to instantiate and configure the application like jenkins for the jenkins-ci plugin.
This activity is typically made by any kind of orchestration tools, like ansible/puppet or any other kind of tools you like to use.

If your plugin uses several tools to install and configure the application, ensure they are available in your docker image.

### output


## docker image

Forjj will start your plugin from a docker container. The image name must be defined in your `<plugin>.yaml` file under `runtime/docker/image`

There is no pre-define image to contains your plugin program. So, you are free to create anything you need. Forjj has no interaction with his content.

# FORJJ container internal

Forjj can start your container as a REST API or shell service.

A REST APi plugin is a foreground tool which usually create a unix socket in /lib/ and listen it.
So forjj starts it in daemon mode. (-d)
A shell plugin, is simply a command with flags call. Forjj starts it, and expect it to terminate shortly with the answer in json.

Forjj will mount several directories:

- /lib : Used to store the socket file.
- /src : Your plguin source directory. Forjj manages the repository where your source are located. You won't have access to any other source files.


FORJ Team
