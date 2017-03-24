package main

import (
	"github.com/kr/text"
	"fmt"
)

func (r *CreateReq)String() string {
	txt := "Forj\n"
	txt += text.Indent(fmt.Sprintf("debug : %s\n", r.Forj.Debug), "  ")
	txt += text.Indent(fmt.Sprintf("Infra : %s\n", r.Forj.ForjjInfra), "  ")
	txt += text.Indent(fmt.Sprintf("Instance : %s\n", r.Forj.ForjjInstanceName), "  ")
	txt += text.Indent(fmt.Sprintf("Organization : %s\n", r.Forj.ForjjOrganization), "  ")
	txt += text.Indent(fmt.Sprintf("Source mount : %s\n", r.Forj.ForjjSourceMount), "  ")
	txt += "\nObjects\n"
	txt += text.Indent(r.Objects.String(), "  ")
	return txt
}

func (o *CreateArgReq) String() string {
	txt := "APP:\n"
	for app_name, app := range o.App {
		txt += fmt.Sprintf("- %s\n", app_name)

		txt += text.Indent(fmt.Sprintf("DeployTo               : %s\n", app.DeployTo), "    ")
		txt += text.Indent(fmt.Sprintf("From Image             : %s\n", app.FromImage), "    ")
		txt += text.Indent(fmt.Sprintf("From Image version     : %s\n", app.FromImageVersion), "    ")
		txt += text.Indent(fmt.Sprintf("Maintainer             : %s\n", app.Maintainer), "    ")
		txt += text.Indent(fmt.Sprintf("Image name             : %s\n", app.Name), "    ")
		txt += text.Indent(fmt.Sprintf("Registry name          : %s\n", app.RegistryAuth), "    ")
		txt += text.Indent(fmt.Sprintf("Registry server        : %s\n", app.RegistryServer), "    ")
		txt += text.Indent(fmt.Sprintf("Generated image version: %s\n", app.Version), "    ")
	}
	txt += "Deploy: \n"
	for name, deploy := range o.Deployment {
		txt += text.Indent(fmt.Sprintf("%s\n", name), "")
		txt += text.Indent(fmt.Sprintf("Name         : %s\n", deploy.Name), "    ")
		txt += text.Indent(fmt.Sprintf("Deploy to    : %s\n", deploy.DeployTo), "    ")
		txt += text.Indent(fmt.Sprintf("Service addr : %s\n", deploy.ServiceAddr), "    ")
		txt += text.Indent(fmt.Sprintf("Service port : %s\n", deploy.ServicePort), "    ")
	}
	return txt
}
