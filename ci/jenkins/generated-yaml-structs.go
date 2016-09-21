// This file is autogenerated by "go generate". Do not modify it.
// It has been generated from your 'jenkins.yaml' file.
// To update those structure, update the 'jenkins.yaml' and run 'go generate'
package main

import "github.hpe.com/christophe-larsonneur/goforjj"

// Common group of data between create/update actions
type DeployStruct struct {
    DeployTo string `json:"deploy-to"` // Where this jenkins source code will be deployed. Supports 'docker'. Future would be 'marathon', 'dcos' and 'host'
    ServiceAddr string `json:"service-addr"` // CNAME or IP address of the expected jenkins instance
    ServicePort string `json:"service-port"` // Expected jenkins instance port number.
}

type SourceStruct struct {
    DockerImage string `json:"docker-image"` // Base docker image name to use in Dockerfile
    DockerImageVersion string `json:"docker-image-version"` // Base docker image version to use in Dockerfile
    DockerRepoimage string `json:"docker-repoimage"` // Repository name containing your base docker image name to use in Dockerfile
    Features string `json:"features"` // List of features to add to jenkins features.lst.
    ForjjInstanceName string `json:"forjj-instance-name"` // Name of the jenkins instance given by forjj.
    ForjjOrganization string `json:"forjj-organization"` // Organization name used in the docker repo name if --docker-repoimage not set.
    Maintainer string `json:"maintainer"` // Jenkins image maintainer
}

type CreateReq struct {
    Args CreateArgReq `json:"args"`
    ReposData map[string]goforjj.PluginRepoData
}

type CreateArgReq struct {
    DeployStruct
    SourceStruct

    // common flags
    ForjjInfra string `json:"forjj-infra"` // Name of the Infra repository to use
    ForjjSourceMount string `json:"forjj-source-mount"` // Where the source dir is located for jenkins plugin.
    JenkinsDebug string `json:"jenkins-debug"` // To activate jenkins debug information
}

type UpdateReq struct {
    Args UpdateArgReq `json:"args"`
    ReposData map[string]goforjj.PluginRepoData
}

type UpdateArgReq struct {
    DeployStruct
    SourceStruct

    FeaturesAdd string `json:"features-add"` // List of features to add to jenkins.

    // common flags
    ForjjInfra string `json:"forjj-infra"` // Name of the Infra repository to use
    ForjjSourceMount string `json:"forjj-source-mount"` // Where the source dir is located for jenkins plugin.
    JenkinsDebug string `json:"jenkins-debug"` // To activate jenkins debug information
}

type MaintainReq struct {
    Args MaintainArgReq `json:"args"`
    ReposData map[string]goforjj.PluginRepoData
}

type MaintainArgReq struct {
    ForjjInstanceName string `json:"forjj-instance-name"` // Name of the jenkins instance given by forjj.

    // common flags
    ForjjInfra string `json:"forjj-infra"` // Name of the Infra repository to use
    ForjjSourceMount string `json:"forjj-source-mount"` // Where the source dir is located for jenkins plugin.
    JenkinsDebug string `json:"jenkins-debug"` // To activate jenkins debug information
}

// Function which adds maintain options as part of the plugin answer in create/update phase.
// forjj won't add any driver name because 'maintain' phase read the list of drivers to use from forjj-maintain.yml
// So --git-us is not available for forjj maintain.
func (r *CreateArgReq)SaveMaintainOptions(ret *goforjj.PluginData) {
    if ret.Options == nil {
        ret.Options = make(map[string]goforjj.PluginOption)
    }

}

func (r *UpdateArgReq)SaveMaintainOptions(ret *goforjj.PluginData) {
    if ret.Options == nil {
        ret.Options = make(map[string]goforjj.PluginOption)
    }

}

func addMaintainOptionValue(options map[string]goforjj.PluginOption, option, value, defaultv, help string) (goforjj.PluginOption){
    opt, ok := options[option]
    if ok && value != "" {
        opt.Value = value
        return opt
    }
    if ! ok {
        opt = goforjj.PluginOption { Help: help }
        if value == "" {
            opt.Value = defaultv
        } else {
            opt.Value = value
        }
    }
    return opt
}

// YamlDesc has been created from your 'jenkins.yaml' file.
const YamlDesc="---\n" +
   "plugin: \"jenkins\"\n" +
   "version: \"0.1\"\n" +
   "description: \"CI jenkins plugin for FORJJ.\"\n" +
   "runtime:\n" +
   "  docker:\n" +
   "    image: \"hub.docker.hpecorp.net/devops/forjj-jenkins\"\n" +
   "    dood: true\n" +
   "  service_type: \"REST API\"\n" +
   "  service:\n" +
   "    #socket: \"jenkins.sock\"\n" +
   "    parameters: [ \"service\", \"start\", \"--templates\", \"/templates\"]\n" +
   "created_flag_file: \"{{ .InstanceName }}/forjj-{{ .Name }}.yaml\"\n" +
   "actions:\n" +
   "  common:\n" +
   "    flags:\n" +
   "      forjj-infra:\n" +
   "        help: \"Name of the Infra repository to use\"\n" +
   "      jenkins-debug:\n" +
   "        help: \"To activate jenkins debug information\"\n" +
   "      forjj-source-mount:\n" +
   "        help: \"Where the source dir is located for jenkins plugin.\"\n" +
   "  create:\n" +
   "    help: \"Create a jenkins instance source code.\"\n" +
   "    flags:\n" +
   "      # Options related to source code\n" +
   "      forjj-instance-name:\n" +
   "        help: \"Name of the jenkins instance given by forjj.\"\n" +
   "        group: \"source\"\n" +
   "      forjj-organization:\n" +
   "        help: \"Organization name used in the docker repo name if --docker-repoimage not set.\"\n" +
   "        group: \"source\"\n" +
   "      docker-image:\n" +
   "        help: \"Base docker image name to use in Dockerfile\"\n" +
   "        default: \"jenkins\"\n" +
   "        group: \"source\"\n" +
   "      docker-image-version:\n" +
   "        help: \"Base docker image version to use in Dockerfile\"\n" +
   "        group: \"source\"\n" +
   "      docker-repoimage:\n" +
   "        help: \"Repository name containing your base docker image name to use in Dockerfile\"\n" +
   "        group: \"source\"\n" +
   "      maintainer:\n" +
   "        help: \"Jenkins image maintainer\"\n" +
   "        group: \"source\"\n" +
   "      features:\n" +
   "        help: \"List of features to add to jenkins features.lst.\"\n" +
   "        group: \"source\"\n" +
   "      # Options related to deployment\n" +
   "      deploy-to:\n" +
   "        default: \"docker\"\n" +
   "        help: \"Where this jenkins source code will be deployed. Supports 'docker'. Future would be 'marathon', 'dcos' and 'host'\"\n" +
   "        group: \"deploy\"\n" +
   "      service-addr:\n" +
   "        required: true\n" +
   "        help: \"CNAME or IP address of the expected jenkins instance\"\n" +
   "        group: \"deploy\"\n" +
   "      service-port:\n" +
   "        default: \"8080\"\n" +
   "        help: \"Expected jenkins instance port number.\"\n" +
   "        group: \"deploy\"\n" +
   "  update:\n" +
   "    help: \"update a jenkins instance source code\"\n" +
   "    flags:\n" +
   "      forjj-instance-name:\n" +
   "        help: \"Name of the jenkins instance given by forjj.\"\n" +
   "        group: \"source\"\n" +
   "      forjj-organization:\n" +
   "        help: \"Organization name used in the docker repo name if --docker-repoimage not set.\"\n" +
   "        group: \"source\"\n" +
   "      docker-image-version:\n" +
   "        help: \"Base docker image version to use in Dockerfile\"\n" +
   "        group: \"source\"\n" +
   "      docker-image:\n" +
   "        help: \"Base docker image name to use in Dockerfile\"\n" +
   "        default: \"jenkins\"\n" +
   "        group: \"source\"\n" +
   "      docker-repoimage:\n" +
   "        help: \"Repository name containing your base docker image name to use in Dockerfile\"\n" +
   "        group: \"source\"\n" +
   "      maintainer:\n" +
   "        help: \"Jenkins image maintainer\"\n" +
   "        group: \"source\"\n" +
   "      features-add:\n" +
   "        help: \"List of features to add to jenkins.\"\n" +
   "  maintain:\n" +
   "    help: \"Instantiate jenkins thanks to source code.\"\n" +
   "    flags:\n" +
   "      forjj-instance-name:\n" +
   "        help: \"Name of the jenkins instance given by forjj.\"\n" +
   ""

