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

		txt += text.Indent(fmt.Sprint("* Add:\n"), "  ")
		txt += text.Indent(fmt.Sprintf("DeployTo               : %s\n", app.Add.DeployTo), "    ")
		txt += text.Indent(fmt.Sprintf("From Image             : %s\n", app.Add.FromImage), "    ")
		txt += text.Indent(fmt.Sprintf("From Image version     : %s\n", app.Add.FromImageVersion), "    ")
		txt += text.Indent(fmt.Sprintf("Maintainer             : %s\n", app.Add.Maintainer), "    ")
		txt += text.Indent(fmt.Sprintf("Image name             : %s\n", app.Add.Name), "    ")
		txt += text.Indent(fmt.Sprintf("Registry name          : %s\n", app.Add.RegistryAuth), "    ")
		txt += text.Indent(fmt.Sprintf("Registry server        : %s\n", app.Add.RegistryServer), "    ")
		txt += text.Indent(fmt.Sprintf("Generated image version: %s\n", app.Add.Version), "    ")

		txt += text.Indent(fmt.Sprint("* Change:\n"), "  ")
		txt += text.Indent(fmt.Sprintf("DeployTo               : %s\n", app.Change.DeployTo), "    ")
		txt += text.Indent(fmt.Sprintf("From Image             : %s\n", app.Change.FromImage), "    ")
		txt += text.Indent(fmt.Sprintf("From Image version     : %s\n", app.Change.FromImageVersion), "    ")
		txt += text.Indent(fmt.Sprintf("Maintainer             : %s\n", app.Change.Maintainer), "    ")
		txt += text.Indent(fmt.Sprintf("Image name             : %s\n", app.Change.Name), "    ")
		txt += text.Indent(fmt.Sprintf("Registry name          : %s\n", app.Change.RegistryAuth), "    ")
		txt += text.Indent(fmt.Sprintf("Registry server        : %s\n", app.Change.RegistryServer), "    ")
		txt += text.Indent(fmt.Sprintf("Generated image version: %s\n", app.Change.Version), "    ")
	}
	txt += "Deploy: \n"
	for name, deploy := range o.Deployment {
		txt += text.Indent(fmt.Sprintf("%s\n", name), "")
		txt += text.Indent(fmt.Sprint("* Add:\n"), "  ")
		txt += text.Indent(fmt.Sprintf("Name         : %s\n", deploy.Add.Name), "    ")
		txt += text.Indent(fmt.Sprintf("Deploy to    : %s\n", deploy.Add.DeployTo), "    ")
		txt += text.Indent(fmt.Sprintf("Service addr : %s\n", deploy.Add.ServiceAddr), "    ")
		txt += text.Indent(fmt.Sprintf("Service port : %s\n", deploy.Add.ServicePort), "    ")

		txt += text.Indent(fmt.Sprint("* Change:\n"), "  ")
		txt += text.Indent(fmt.Sprintf("Name         : %s\n", deploy.Change.Name), "    ")
		txt += text.Indent(fmt.Sprintf("Deploy to    : %s\n", deploy.Change.DeployTo), "    ")
		txt += text.Indent(fmt.Sprintf("Service addr : %s\n", deploy.Change.ServiceAddr), "    ")
		txt += text.Indent(fmt.Sprintf("Service port : %s\n", deploy.Change.ServicePort), "    ")
	}
	return txt
}
