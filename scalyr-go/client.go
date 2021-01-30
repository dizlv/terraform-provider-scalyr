package sdk

import (
	"fmt"
	"strings"
)

type ScalyrTokens struct {
	WriteLog    string
	ReadLog     string
	WriteConfig string
	ReadConfig  string
}

type ScalyrConfig struct {
	Endpoint   string
	Region     string
	Tokens     ScalyrTokens
	TeamTokens ScalyrTokens
	Team       string
}

func (c *ScalyrConfig) hasTeam() bool {
	return c.Team != ""
}

var regionEndpoints map[string]string = map[string]string{"us": "https://app.scalyr.com/", "eu": "https://app.eu.scalyr.com/"}

func initialize(config *ScalyrConfig) {
	if config.Region == "" {
		config.Region = getEnvWithDefault("SCALYR_REGION", "us")
	}
	regionEndpoint := regionEndpoints[config.Region]

	if config.Tokens.WriteLog == "" {
		config.Tokens.WriteLog = getEnvWithDefault("SCALYR_WRITELOG_TOKEN", "")
	}
	if config.Tokens.ReadLog == "" {
		config.Tokens.ReadLog = getEnvWithDefault("SCALYR_READLOG_TOKEN", "")
	}
	if config.Tokens.ReadConfig == "" {
		config.Tokens.ReadConfig = getEnvWithDefault("SCALYR_READCONFIG_TOKEN", "")
	}
	if config.Tokens.WriteConfig == "" {
		config.Tokens.WriteConfig = getEnvWithDefault("SCALYR_WRITECONFIG_TOKEN", "")
	}
	if config.Endpoint == "" {
		config.Endpoint = getEnvWithDefault("SCALYR_SERVER", regionEndpoint)
	}
	// Add an SSL prefix if you don't tell us
	if ! strings.HasPrefix(config.Endpoint, "http") {
		config.Endpoint = fmt.Sprintf("https://%s", config.Endpoint)
	}
}

func NewClient(config *ScalyrConfig) (*ScalyrConfig, error) {
	initialize(config)

	return config, nil
}
