**This is a POC**

# Introduction

Forjj project is the [next generation of Forj](https://github.hpe.com/christophe-larsonneur/forjj). 

The core of Forjj is very limited to manage local repositories, mainly infra and have simple pre-defined tasks (create/update/maintain)

But it have to be concretely implemented against application the end user want to use (ie jenkins + github for example). To do so, Forjj rely on a collection of tools that implement those tasks on each application. Those tools are named `FORJJ plugins`

It defines as well the link between application that have to be implemented, to properly configure a DevOps flow wanted.
Currently the task is not well defined.

For now, this repository are focused on 2 kind of tools, known as `plugins type`:
- Continuous integration identified by `ci`
- GIT upstream, ie git backend SCM application (like github) identified by `upstream`

# FORJJ plugins

A Forjj plugins is stored in this repo under a type of plugins directory (ci/upstream).
The name of the plugin is stored as :
- a subdirectory of his type 
- as a yaml file name (`<PluginName>.yaml`)
- as a binary/script (`<plugin>`)

The plugin binary HAVE to be with executable bits in place. It can be any kind of binary or script.

The plugin will run in the context of a Forjj container [described here](https://github.hpe.com/christophe-larsonneur/forjj/docker/Dockerfile)
> Currently there is no way to execute a plugin in a different container. But for me it can make sense to do so.
> We can define a dedicated plugin container with your needed packages. This one could be derived from a more generic one that has the basic generic code to help Forjj to work properly (entrypoint.sh)
> So, if you want to get that now, consider a contribution to implement it.

This container is started by `forjj` automatically. For details about this part of forjj, [read it here](https://github.hpe.com/christophe-larsonneur/forjj)

## Plugin tasks

Each plugin must implement the 3 Forjj typical tasks (create/update/maintain). Forjj will call the plugin binary file as follow:

    <Forjj-contribs clone>/<plugin type/<plugin> create/update/maintain [flags]

## Description of <Plugin>.yaml

This file, in yaml format, describe the list of valid flags that the plugin requires to do the task properly.

It has the following format:

```yaml
---
flags:                           # common regroups options for all tasks
  common:                        # list of options for all tasks.
  - <optionName>:                # This named the option (without --).
    help : "option help"         # Describe the optionName help.
    required: <true/false>       # false by default.
    default : "my default value" # Define the default value.
  create: # create options list
  - ...:
    [...]
  update:
  - ...:
    [...]
  maintain:
  - ...:
    [...]
    
```

By convention, an option that is dedicated to the plugin is prefixed by the name of the plugin.

Ex: --github-server => github-server 

A global option can be already defined by forj. For example --debug. If the plugin requires it, you need to define it.

You can also define another plugin option, like jenkins-ci-server. It will work, but I would suspect this to be part of a link case instead which is currently not properly defined. (**TBD**)


Those options will be visible when the user will use the --ci or --git-us combined with --help to the `forjj` tool.

Ex: 
  List forjj common options + github driver common options
  `forjj --git-us github --help`

  List long help options for github driver
  `forjj --git-us github --help-long`

  create task list of options for github driver
  `forjj create --git-us github --help-long`

# Using you own version of forjj-contribs

By default forjj is defined to get the driver options definition from [github entreprise](https://github.hpe.com/forj/forjj-contribs)
So, getting the list of option can take some few seconds (time to read the yaml file from github)

You can change and use another url or even a local path, with `forjj --contrib-repo <base url/forjj-contribs path>`

If you use another url, it must expose the description file as `<contrib-repo base url>/<branchName>/<DriverType>/<DriverName>/<DriverName>.yaml`
If you use a path, the yaml file must be located in `<forjj-contribs path>/<DriverType>/<DriverName>/<DriverName>.yaml`

# More information

As this is a POC, not everything is fully published. So, feel create to ask me anything, or even contribute.

This repository has been defined to properly help you to contribute on any places where

FORJ Team
