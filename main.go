package main

import (
	"ansoni/terraform-provider-scalyr/scalyr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: scalyr.CreateProvider,
	})
}
