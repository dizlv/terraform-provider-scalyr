package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	scalyr "github.com/ansoni/terraform-provider-scalyr/scalyr-go"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"region": &schema.Schema{Type: schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SCALYR_REGION", "us"),
				Description: "Scalyr Region",
			},
			"endpoint": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SCALYR_SERVER", nil),
				Description: "Scalyr Server",
			},
			"read_log_token": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SCALYR_READLOG_TOKEN", nil),
				Description: "Scalyr ReadLog API Token",
			},
			"write_log_token": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SCALYR_WRITELOG_TOKEN", nil),
				Description: "Scalyr WriteLog API Token",
			},
			"read_config_token": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SCALYR_READCONFIG_TOKEN", nil),
				Description: "Scalyr ReadConfig API Token",
			},
			"write_config_token": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SCALYR_WRITECONFIG_TOKEN", nil),
				Description: "Scalyr WriteConfig API Token",
			},
			"team": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SCALYR_TEAM", nil),
				Description: "Scalyr Team Identifier",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"scalyr_event": resourceEvent(),
			"scalyr_file":  resourceFile(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"scalyr_file":   datasourceFile(),
			"scalyr_query":  datasourceQuery(),
			"scalyr_teams":  datasourceTeams(),
			"scalyr_tokens": datasourceTokens(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	region := d.Get("region").(string)
	endpoint := d.Get("endpoint").(string)
	readLogToken := d.Get("read_log_token").(string)
	writeLogToken := d.Get("write_log_token").(string)
	readConfigToken := d.Get("read_config_token").(string)
	writeConfigToken := d.Get("write_config_token").(string)
	tokens := scalyr.ScalyrTokens{ReadLog: readLogToken, WriteLog: writeLogToken, ReadConfig: readConfigToken, WriteConfig: writeConfigToken}
	team := d.Get("team").(string)
	return scalyr.NewClient(&scalyr.ScalyrConfig{Endpoint: endpoint, Region: region, Team: team, Tokens: tokens})
}
