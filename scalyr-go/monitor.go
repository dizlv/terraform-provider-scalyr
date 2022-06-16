package sdk

import (
	"encoding/json"
	"fmt"
)

const ConfigurationFilePath = "/scalyr/monitors"

type Monitor struct {
	Path string `json:"path"`

	// Type of the monitor. Currently, supported only `s3Bucket`.
	Type string `json:"type"`

	// The AWS region where sink SQS queue is located. e.g. `us-east-1`
	Region string `json:"region"`

	// ARN of the IAM role.
	RoleToAssume string `json:"roleToAssume"`

	// Name of the SQS queue where bucket sends new-object notifications.
	QueueUrl string `json:"queueUrl"`

	// `text_gzip` or `text_zstd` or `text`.
	FileFormat string `json:"fileFormat"`

	// Name of a host under which bucket logs will appear in the UI.
	HostName string `json:"hostname"`

	// Name of a parser that should be applied to logs.
	Parser string `json:"parser"`

	// Label of monitor, used as unique ID for this integration
	Label string `json:"label"`
}

type Monitors []*Monitor

type MonitorsConfigurationFile struct {
	ExecutionIntervalMinutes float64  `json:"executionIntervalMinutes"`
	TimeoutSeconds           float64  `json:"timeoutSeconds"`
	Monitors                 Monitors `json:"monitors"`
}

func NewMonitor(monitorType, region, roleToAssume, queueUrl, fileFormat, hostName, parser, label string) *Monitor {
	return &Monitor{
		Type:         monitorType,
		Region:       region,
		RoleToAssume: roleToAssume,
		QueueUrl:     queueUrl,
		FileFormat:   fileFormat,
		HostName:     hostName,
		Parser:       parser,
		Label:        label,
	}
}

func (scalyr *ScalyrConfig) UpdateMonitors(file *MonitorsConfigurationFile) error {
	data, err := json.Marshal(file)
	if err != nil {
		return err
	}

	_, err = scalyr.PutFile(ConfigurationFilePath, fmt.Sprintf("%s", data))
	if err != nil {
		return err
	}

	return nil
}

func (scalyr *ScalyrConfig) GetMonitor(path string) (*GetFileResponse, error) {
	return scalyr.GetFile(path)
}
