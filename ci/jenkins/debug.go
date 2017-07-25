package main

import (
	"fmt"
	"github.com/kr/text"
)

func (r *CreateReq) String() string {
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

		txt += text.Indent(fmt.Sprintf("DeployTo               : %s\n", app.To), "    ")
		txt += text.Indent(fmt.Sprintf("DeployServiceAddr      : %s\n", app.ServiceAddr), "    ")
		txt += text.Indent(fmt.Sprintf("DeployServicePort      : %s\n", app.ServicePort), "    ")
		txt += text.Indent(fmt.Sprintf("From Image             : %s\n", app.FromImage), "    ")
		txt += text.Indent(fmt.Sprintf("From Image version     : %s\n", app.FromImageVersion), "    ")
		txt += text.Indent(fmt.Sprintf("Maintainer             : %s\n", app.Maintainer), "    ")
		txt += text.Indent(fmt.Sprintf("Image name             : %s\n", app.Name), "    ")
		txt += text.Indent(fmt.Sprintf("Registry name          : %s\n", app.RegistryAuth), "    ")
		txt += text.Indent(fmt.Sprintf("Registry server        : %s\n", app.RegistryServer), "    ")
		txt += text.Indent(fmt.Sprintf("Generated image version: %s\n", app.Version), "    ")
	}
	return txt
}
