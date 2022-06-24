package sdk

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/exp/slices"
)

const ConfigurationFilePath = "/scalyr/monitors"

var MonitorNotFound = errors.New("")

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

	// Name of a parser that should be applied to log.
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

type SearchMonitorResult struct {
	Index    int
	Err      error
	Monitor  *Monitor
	Monitors Monitors
}

func GetMonitorByLabel(ctx context.Context, label string, client *ScalyrConfig) SearchMonitorResult {
	monitors, err := getMonitorsFromFile(ctx, client)
	if err != nil {
		return SearchMonitorResult{
			Index: -1,
			Err:   err,
		}
	}

	index := findMonitorByLabel(label, monitors)
	if index != -1 {
		return SearchMonitorResult{
			Index: -1,
			Err:   err,
		}
	}

	return SearchMonitorResult{
		Index:    index,
		Err:      nil,
		Monitor:  monitors[index],
		Monitors: monitors,
	}
}

func getMonitorsFromFile(ctx context.Context, client *ScalyrConfig) (Monitors, error) {
	response, err := client.GetFile(ctx, ConfigurationFilePath)
	if err != nil {
		return nil, err
	}

	var monitors Monitors

	err = json.Unmarshal([]byte(response.Content), &monitors)
	if err != nil {
		return nil, err
	}

	return monitors, nil
}

func updateMonitorsInFile(ctx context.Context, client *ScalyrConfig, monitors Monitors) error {
	// get files put them in string and pass to the file

	content, err := json.Marshal(monitors)
	if err != nil {
		return err
	}

	_, err = client.PutFile(ctx, ConfigurationFilePath, string(content))
	if err != nil {
		return err
	}

	return nil
}

func findMonitorByLabel(label string, monitors Monitors) int {
	index := slices.IndexFunc(monitors, func(m *Monitor) bool {
		return m.Label == label
	})

	return index
}

type CreateMonitorInput = Monitor

type CreateMonitorOutput struct {
	Monitor *Monitor
}

func (scalyr *ScalyrConfig) CreateMonitor(ctx context.Context, input *CreateMonitorInput) (*CreateMonitorOutput, error) {
	result := GetMonitorByLabel(ctx, input.Label, scalyr)
	if result.Err != nil {
		return nil, result.Err
	}

	if result.Monitors != nil {
		return nil, errors.New(fmt.Sprintf("monitor with provided label=%s already exist", input.Label))
	}

	newMonitor := NewMonitor(
		input.Type,
		input.Region,
		input.RoleToAssume,
		input.QueueUrl,
		input.FileFormat,
		input.HostName,
		input.Parser,
		input.Label,
	)

	result.Monitors = append(result.Monitors, newMonitor)

	err := updateMonitorsInFile(ctx, scalyr, result.Monitors)
	if err != nil {
		return nil, err
	}

	return &CreateMonitorOutput{Monitor: result.Monitor}, nil
}

type ReadMonitorInput struct {
	Label string
}

type ReadMonitorOutput struct {
	Monitor *Monitor
}

func (scalyr *ScalyrConfig) ReadMonitor(ctx context.Context, input *ReadMonitorInput) (*ReadMonitorOutput, error) {
	result := GetMonitorByLabel(ctx, input.Label, scalyr)
	if result.Err != nil {
		return nil, result.Err
	}

	return &ReadMonitorOutput{Monitor: result.Monitor}, nil
}

type UpdateMonitorInput = Monitor
type UpdateMonitorOutput struct{}

func (scalyr *ScalyrConfig) UpdateMonitor(ctx context.Context, input *UpdateMonitorInput) (*UpdateMonitorOutput, error) {
	result := GetMonitorByLabel(ctx, input.Label, scalyr)
	if result.Err != nil {
		return nil, result.Err
	}

	updatedMonitor := NewMonitor(
		input.Type,
		input.Region,
		input.RoleToAssume,
		input.QueueUrl,
		input.FileFormat,
		input.HostName,
		input.Parser,
		input.Label,
	)

	result.Monitors[result.Index] = updatedMonitor

	err := updateMonitorsInFile(ctx, scalyr, result.Monitors)
	if err != nil {
		return nil, err
	}

	return &UpdateMonitorOutput{}, nil
}

type DeleteMonitorInput struct {
	Label string
}

type DeleteMonitorOutput struct{}

func (scalyr *ScalyrConfig) DeleteMonitor(ctx context.Context, input *DeleteMonitorInput) (*DeleteMonitorOutput, error) {
	result := GetMonitorByLabel(ctx, input.Label, scalyr)
	if result.Err != nil {
		return nil, result.Err
	}

	result.Monitors = slices.Delete(result.Monitors, result.Index, result.Index+1)

	return &DeleteMonitorOutput{}, nil
}
