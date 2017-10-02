**This is a POC**

# Introduction

Forjj project is the [next generation of Forj](https://github.com/forj-oss/forjj).

The core of Forjj is very limited to manage local repositories, mainly infra and have simple pre-defined tasks (create/update/maintain)

But it have to be concretely implemented against application the end user want to use (ie jenkins + github for example). To do so, Forjj rely on a collection of tools that implement those tasks on each application. Those tools are named `FORJJ plugins`

It defines as well the link between application that have to be implemented, to properly configure a DevOps flow wanted.
Currently the task is closed to be fully defined. See [Contributing document](CONTRIBUTING.md)

For now, this repository are focused on 2 kind of tools, known as `plugins type`:
- Continuous integration identified by `ci`. Under that one, one `jenkins` plugin has been created to handle jenkins CI.
- GIT upstream, ie git backend SCM application (like github) identified by `upstream`. Under that one, one `github` plugin has been created to handle github service.

# FORJJ plugins

A Forjj plugins is stored in this repo under a type of plugins directory (ci/upstream).
The name of the plugin is stored as :
- a subdirectory of his type
- as a yaml file name (`<PluginName>.yaml`)
- as a binary/script (`<plugin>`)

The plugin binary HAVE to be with executable bits in place. It can be any kind of binary or script.

The plugin will run in the context of a Forjj container. But there is no image type to use with forjj. Forjj uses only a socket file and source mount point. Then it communicates with it through the socket file.

This container is started by `forjj` automatically. For details about this part of forjj, [read the source code here](https://github.com/forj-oss/forjj/raw/master/driver.go)

## Plugin tasks

Each plugin must implement the 3 Forjj typical tasks (create/update/maintain). Forjj will call the plugin service as follow:

### For shell plugin:

```bash
<plugin script> create/update/maintain [flags]
```
A special flag --data is used by forjj to send out the reposdata information to an upstream plugin. This data is in json format.

If flags are valid, it must returns data in JSON and return code null.

**NOTE:** To be honest, today, even if the code is maintained to support this case, I do not have any code running like that. All plugins that I created implements the REST API plugin model in GO.

So, there is no guarantee this case to work properly.

### For REST API service plugin:

```bash
POST <plugin url or socket>/<create/update/maintain>
```

data posted are sent in `application/json` and result is in `application/json`.

## Description of `<Plugin>.yaml`

This file, in yaml format, describe the list of valid flags that the plugin requires to do the task properly.

It has the following format:

```yaml
---
plugin: <PluginName>                 # Required by the GO `goforjj` package. It must be the name of your yaml file.
version: <PluginVersion>             # Version of your plugin.
description: "<Desc>"                # Small sentence about what your FORJJ plugin do
runtime:                             # Define Runtime pqrqmeters for the plugin service/process
  docker_image: "<Docker image>"     # Docker image tag name. If set, the plugin will run through docker.
  service_type: "<service_type>"     # Type of service exposed to forjj. Can be 'REST API' or 'shell'.
  service:                           # Options for REST API service
    socket: "<file>.sock"            # Socket File name used.
    parameters: [ "", ...]           # Collection of flags to provide to the service daemon startup.
  command:                           # Options for shell service
    parameters: [ "", ...]           # Collection of first flags to provide to the service process execution.
actions:                             # Collection of actions to describe. common regroups options for all actions.
  common:                            # list of options for all tasks.
    help: "<Plugin cmd description>" # Sentence describing the plugin command tasks.
    flags:                           # Describe the collection of common flags.
      <optionName>:                  # This named the option (without --).
        help : "option help"         # Describe the optionName help.
        required: <true/false>       # false by default.
        default : "my default value" # Define the default value. The default value can be interpreted by forjj with go template ie '{{ ... }}'
                                     # template interpret on few forjj internal data. See Forjj internal data section for details
        secure: <true/false>         # False by default. See details about how secure data are managed between forjj and the plugin.
  create:                            # create options list
    help: "<Plugin cmd description>" # Sentence describing the plugin command tasks.
    flags:
      ...:
      [...]
  update:                            # Same data as common/create/... commands
    [...]
  maintain:                          # Same data as common/create/... commands
    [...]
objects:
  <object Name>:                     # Object defined in Forjfile as root
    default-actions:                 # Can be `add`, `remove`, `change`, `list` or `rename`
                                     # Except defined on field, all fields can be manipulated by forjj with default actions.
    flags:                           # Collection of Object fields
      <optionName>:                  # This named the option (without --).
        help : "option help"         # Describe the optionName help.
        required: <true/false>       # false by default.
        default : "my default value" # Define the default value.
        secure: <true/false>         # False by default. See details about how secure data are managed between forjj and the plugin.
        cli-exported-to-actions:     # Can be `create`, `update` or `maintain`.
                                     # Declare this field as action flag. The field is prefixed by the instance name
        actions:                     # List of actions where this field can be added. If not defined here, the fields actions list comes from default-actions.
    groups:                          # Collection of group of object fields
      <groupName>:                   # The name of group of fields. All field will bre prefixed by `<GroupName>_`
        flags:                       # Collection of Object fields
          ...                        # Same list of options as object fields. (ie those in objects/<objectName>/flags/...)
```

By convention, an option that is dedicated to the plugin is prefixed by the name of the plugin.

Ex: --github-server => github-server

A global option can be already defined by forj. For example --debug. If the plugin requires it, you need to define it.

You can also define another plugin option, like jenkins-ci-server. It will work, but I would suspect this to be part of a link case instead which is currently not properly defined. (**TBD**)


Those options will be visible when the user will use the --apps combined with --help to the `forjj` tool.

Ex:
  List forjj common options + github driver common options
  `forjj create --apps upstream:github --help`

  List long help options for github driver
  `forjj create --apps upstream:github --help-long`

  create task list of options for github driver
  `forjj create --apps upstream:github --help-long`

## secure flag

When a plugin requires some critical data like password or token, the plugin should never store it in the SCM.
If you declare this flag as secure, forjj will give a facility to store this information in his workspace or a file selected by the end user.
It facilitates to
- manage this data outside forjj
- To deliver automatically them to each stage though the workspace/forjj-creds.yaml or through an external file.

Use case example:

```bash
# github has a secure flag called github-token.
# The following command will sent the token to the plugin and forjj will get it back requested by the plugin to store it
# in the workspace/forjj-creds.yaml and declare those data in the infra repo
forjj create ~/tmp/test --apps upstream:github --github-token $MYTOKEN
#
# So, now, you can update or maintain it without need to provide it again as it iw stored in the local workspace
forjj update ~/tmp/test add-myrepo-branch --repos myrepo
forjj maintain ~/tmp/test
# github plugin will use the saved token from the workspace to deliver the --github-token for you.
```

You can use your own secure file, delivered from vault for example:
```bash
# github has a secure flag called github-token.
# The following command will sent the token to the plugin and forjj will get it back requested by the plugin to store it
# in the workspace/forjj-creds.yaml and declare those data in the infra repo
forjj create ~/tmp/test -c ~/vault/mycreds.yaml --apps upstream:github --github-token $MYTOKEN
#
# So, now, you can update or maintain it without need to provide it again as it iw stored in the local workspace
forjj update ~/tmp/test add-myrepo-branch -c ~/vault/mycreds.yaml --repos myrepo
forjj maintain ~/tmp/test -c ~/vault/mycreds.yaml
# github plugin will use the saved token from the workspace to deliver the --github-token for you.
```

# Using you own version of forjj-contribs

By default forjj is defined to get the driver options definition from [github](https://github.com/forj-oss/forjj-contribs)
So, getting the list of option can take some few seconds (time to read the yaml file from github)

You can change and use another url or even a local path, with `forjj --contribs-repo <base url/forjj-contribs path>`
You can also use the CONTRIBS_REPO Environment variable.

If you use another url, it must expose the description file as `<contribs-repo base url>/<branchName>/<DriverType>/<DriverName>/<DriverName>.yaml`

If you use a path, the yaml file must be located in `<forjj-contribs path>/<DriverType>/<DriverName>/<DriverName>.yaml`

# Want to write your own FORJJ plugin?

See [goforjj project in github](https://github.com/forj-oss/goforjj)

# More information

As this is a POC, not everything is fully published. So, feel free to ask me anything, or even contribute.

This repository has been defined to properly help you to contribute everywhere.


FORJ Team
