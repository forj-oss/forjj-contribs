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

The plugin notion is FORJJ is simple and basic.

FORJJ do his task by starting a container, with the task name and flags.
- Flags are defined by the [yaml description file](https://github.hpe.com/christophe-larsonneur/forjj-contribs#description-of-yaml)
- The Tasks name are 'create', 'update', ... as [defined in this page](https://github.hpe.com/christophe-larsonneur/forjj/blob/master/CONTRIBUTION.md#forjj-cli)
- The docker container has been created to ensure your plugin environment is identified.

Today, there is no cli option to choose the docker image. So it uses docker.hos.hpecorp.net/devops/forjj:latest built from [this Docker directory](https://github.hpe.com/christophe-larsonneur/forjj/tree/master/docker)

If you need/want to build your image for your plugin, feel free to do it. It won't be a big deal to update forjj to take care of this.

FORJ Team
