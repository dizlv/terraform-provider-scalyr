package scalyr

import (
	scalyr "ansoni/terraform-provider-scalyr/scalyr-go"
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/exp/slices"
	"log"
)

const (
	TypeArg            = "monitor_type"
	AwsRegionArg       = "aws_region"
	IamRoleToAssumeArg = "iam_role_to_assume"
	QueueUrlArg        = "queue_url"
	FileFormatArg      = "file_format"
	HostNameArg        = "host_name"
	ParserArg          = "parser"
	LabelArg           = "label"
)

var Schema = map[string]*schema.Schema{
	TypeArg: {
		Type:        schema.TypeString,
		Description: "Type of the monitor. Currently, supported only `s3Bucket`",
		ForceNew:    true,
		Required:    true,
	},

	AwsRegionArg: {
		Type:        schema.TypeString,
		Description: "The AWS region where sink SQS queue is located. e.g. `us-east-1`",
		ForceNew:    true,
		Required:    true,
	},

	IamRoleToAssumeArg: {
		Type:        schema.TypeString,
		Description: "ARN of the IAM role",
		ForceNew:    false,
		Required:    true,
	},

	QueueUrlArg: {
		Type:        schema.TypeString,
		Description: "Name of the SQS queue where bucket sends new-object notifications",
		ForceNew:    false,
		Required:    true,
	},

	FileFormatArg: {
		Type:        schema.TypeString,
		Description: "`text_gzip` or `text_zstd` or `text`",
		ForceNew:    false,
		Required:    true,
	},

	HostNameArg: {
		Type:        schema.TypeString,
		Description: "Name of a host under which bucket logs will appear in the UI",
		ForceNew:    false,
		Required:    true,
	},

	ParserArg: {
		Type:        schema.TypeString,
		Description: "Name of a parser that should be applied to logs",
		ForceNew:    false,
		Required:    true,
	},

	LabelArg: {
		Type:        schema.TypeString,
		Description: "Monitor Unique ID",
		ForceNew:    false,
		Required:    true,
	},
}

func resourceMonitor() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceMonitorRead,
		CreateContext: resourceMonitorCreate,
		DeleteContext: resourceMonitorDelete,
		UpdateContext: resourceMonitorUpdate,

		Schema: Schema,
	}
}

// resourceMonitorCreate executes API call on the File api to provision monitor with provided
// data. As a side effect this function triggers another request to retrieve information on just
// submitted file.
func resourceMonitorCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	monitorType := data.Get(TypeArg).(string)
	region := data.Get(AwsRegionArg).(string)
	roleToAssume := data.Get(IamRoleToAssumeArg).(string)
	queueUrl := data.Get(QueueUrlArg).(string)
	fileFormat := data.Get(FileFormatArg).(string)
	hostName := data.Get(HostNameArg).(string)
	parser := data.Get(ParserArg).(string)
	label := data.Get(LabelArg).(string)

	client := meta.(*scalyr.ScalyrConfig)

	var diagnostics diag.Diagnostics

	mutexKV.Lock("monitors")
	defer mutexKV.Unlock("monitors")

	// Configure new monitor with appropriate data.
	newMonitor := scalyr.NewMonitor(
		monitorType,
		region,
		roleToAssume,
		queueUrl,
		fileFormat,
		hostName,
		parser,
		label,
	)

	file, err := getMonitorConfigurationFile(client)
	if err != nil {
		return diag.FromErr(err)
	}

	index := findMonitorByLabel(label, file.Monitors)

	if index != -1 {
		return diag.Errorf("Monitor with configured label already exist")
	}

	// Insert new element
	file.Monitors = append(file.Monitors, newMonitor)

	err = updateMonitorsConfigurationFile(file, client)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(fmt.Sprintf("%s", label))

	return diagnostics
}

func resourceMonitorDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*scalyr.ScalyrConfig)

	// File api does not has specific ID provided, so we create one by combining file path + version.
	label := data.Id()

	mutexKV.Lock("monitors")
	defer mutexKV.Unlock("monitors")

	file, err := getMonitorConfigurationFile(client)
	if err != nil {
		return diag.FromErr(err)
	}

	index := findMonitorByLabel(label, file.Monitors)

	if index == -1 {
		return diag.Errorf("Monitor is not present in submitted configuration file")
	}

	file.Monitors = slices.Delete(file.Monitors, index, index+1)

	log.Printf("Deleting {%v}", file.Monitors)

	err = updateMonitorsConfigurationFile(file, client)
	if err != nil {
		return diag.FromErr(err)
	}

	// Set terraform id to "" which means that resource has been deleted successfully and can
	// be dropped from terraform state.
	data.SetId("")

	return nil
}

func resourceMonitorRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*scalyr.ScalyrConfig)
	label := data.Id()

	mutexKV.Lock("monitors")
	defer mutexKV.Unlock("monitors")

	var diagnostics diag.Diagnostics

	file, err := getMonitorConfigurationFile(client)
	if err != nil {
		return diag.FromErr(err)
	}

	index := findMonitorByLabel(label, file.Monitors)

	log.Printf("Looking for label %s in monitors %v", label, file.Monitors)

	if index == -1 && !data.IsNewResource() {
		data.SetId("")

		return nil
	}

	monitor := file.Monitors[index]

	// todo: handle this stuff in diags
	data.Set(TypeArg, monitor.Type)
	data.Set(AwsRegionArg, monitor.Region)
	data.Set(IamRoleToAssumeArg, monitor.RoleToAssume)
	data.Set(QueueUrlArg, monitor.QueueUrl)
	data.Set(FileFormatArg, monitor.FileFormat)
	data.Set(HostNameArg, monitor.HostName)
	data.Set(ParserArg, monitor.Parser)
	data.Set(LabelArg, monitor.Label)

	data.SetId(fmt.Sprintf("%v", monitor.Label))

	return diagnostics
}

func resourceMonitorUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*scalyr.ScalyrConfig)
	id := data.Id()

	mutexKV.Lock("monitors")
	defer mutexKV.Unlock("monitors")

	monitorType := data.Get(TypeArg).(string)
	region := data.Get(AwsRegionArg).(string)
	roleToAssume := data.Get(IamRoleToAssumeArg).(string)
	queueUrl := data.Get(QueueUrlArg).(string)
	fileFormat := data.Get(FileFormatArg).(string)
	hostName := data.Get(HostNameArg).(string)
	parser := data.Get(ParserArg).(string)
	label := data.Get(LabelArg).(string)

	file, err := getMonitorConfigurationFile(client)
	if err != nil {
		return diag.FromErr(err)
	}

	index := findMonitorByLabel(id, file.Monitors)

	updatedMonitor := scalyr.NewMonitor(
		monitorType,
		region,
		roleToAssume,
		queueUrl,
		fileFormat,
		hostName,
		parser,
		label,
	)

	file.Monitors[index] = updatedMonitor

	if err := updateMonitorsConfigurationFile(file, client); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(fmt.Sprintf("%s", label))

	return nil
}
