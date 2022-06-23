package scalyr

import (
	"context"
	scalyr "github.com/ansoni/terraform-provider-scalyr/scalyr-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"sync"
)

const (
	RegionEnvArg           = "SCALYR_REGION"
	ServerEnvArg           = "SCALYR_SERVER"
	ReadLogTokenEvnArg     = "SCALYR_READLOG_TOKEN"
	WriteLogTokenEnvArg    = "SCALYR_WRITELOG_TOKEN"
	ReadConfigTokenEnvArg  = "SCALYR_READCONFIG_TOKEN"
	WriteConfigTokenEnvArg = "SCALYR_WRITECONFIG_TOKEN"
	TeamEnvArg             = "SCALYR_TEAM"
)

type Synchronizer struct {
	mutexes sync.Map
}

func (synchronizer *Synchronizer) Lock(key string) func() {
	value, _ := synchronizer.mutexes.LoadOrStore(key, &sync.Mutex{})
	mutex := value.(*sync.Mutex)

	mutex.Lock()

	return func() {
		mutex.Unlock()
	}
}

var synchronizer = Synchronizer{}

func CreateProvider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(RegionEnvArg, "us"),
				Description: "Scalyr Region",
			},
			"endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(ServerEnvArg, nil),
				Description: "Scalyr Server",
			},
			"read_log_token": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(ReadLogTokenEvnArg, nil),
				Description: "Scalyr ReadLog API Token",
			},
			"write_log_token": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(WriteLogTokenEnvArg, nil),
				Description: "Scalyr WriteLog API Token",
			},
			"read_config_token": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(ReadConfigTokenEnvArg, nil),
				Description: "Scalyr ReadConfig API Token",
			},
			"write_config_token": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(WriteConfigTokenEnvArg, nil),
				Description: "Scalyr WriteConfig API Token",
			},
			"team": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(TeamEnvArg, nil),
				Description: "Scalyr Team Identifier",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"scalyr_event":                 resourceEvent(),
			"scalyr_file":                  resourceFile(),
			"scalyr_monitor":               resourceMonitor(),
			"scalyr_parser":                resourceParser(),
			"scalyr_parser_format_rewrite": resourceParserFormatRewrite(),
			"scalyr_parser_line_grouper":   resourceParserLineGrouper(),
			"scalyr_parser_pattern":        resourceParserPattern(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"scalyr_file":          datasourceFile(),
			"scalyr_query":         datasourceQuery(),
			"scalyr_teams":         datasourceTeams(),
			"scalyr_tokens":        datasourceTokens(),
			"scalyr_parser_format": dataSourceParserFormat(),
		},

		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	team := d.Get("team").(string)
	region := d.Get("region").(string)
	endpoint := d.Get("endpoint").(string)
	readLogToken := d.Get("read_log_token").(string)
	writeLogToken := d.Get("write_log_token").(string)
	readConfigToken := d.Get("read_config_token").(string)
	writeConfigToken := d.Get("write_config_token").(string)

	tokens := scalyr.ScalyrTokens{
		ReadLog:     readLogToken,
		WriteLog:    writeLogToken,
		ReadConfig:  readConfigToken,
		WriteConfig: writeConfigToken,
	}

	client, err := scalyr.NewClient(&scalyr.ScalyrConfig{
		Endpoint: endpoint,
		Region:   region,
		Team:     team,
		Tokens:   tokens,
	})

	if err != nil {
		return nil, diag.FromErr(err)
	}

	return client, nil
}
