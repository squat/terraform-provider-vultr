package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/squat/terraform-provider-vultr/vultr"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: vultr.Provider,
	})
}
