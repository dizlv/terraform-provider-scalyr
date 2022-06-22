package scalyr

import (
	scalyr "ansoni/terraform-provider-scalyr/scalyr-go"
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const monitorsLockName = "/scalyr/monitors"

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
	input := &scalyr.CreateMonitorInput{
		Type:         monitorType,
		Region:       region,
		RoleToAssume: roleToAssume,
		QueueUrl:     queueUrl,
		FileFormat:   fileFormat,
		HostName:     hostName,
		Parser:       parser,
		Label:        label,
	}

	unlock := synchronizer.Lock(monitorsLockName)
	defer unlock()

	if output, err := client.CreateMonitor(ctx, input); err != nil {
		return diag.FromErr(err)
	} else {
		data.SetId(output.Monitor.Label)
	}

	return nil
}

func resourceMonitorRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*scalyr.ScalyrConfig)

	input := &scalyr.ReadMonitorInput{
		Label: data.Id(),
	}

	unlock := synchronizer.Lock(monitorsLockName)
	defer unlock()

	if output, err := client.ReadMonitor(ctx, input); err != nil {
		if err == scalyr.MonitorNotFound && !data.IsNewResource() {
			data.SetId("")

			return nil
		}

		return diag.FromErr(err)
	} else {
		monitor := output.Monitor

		data.Set(TypeArg, monitor.Type)
		data.Set(AwsRegionArg, monitor.Region)
		data.Set(IamRoleToAssumeArg, monitor.RoleToAssume)
		data.Set(QueueUrlArg, monitor.QueueUrl)
		data.Set(FileFormatArg, monitor.FileFormat)
		data.Set(HostNameArg, monitor.HostName)
		data.Set(ParserArg, monitor.Parser)
		data.Set(LabelArg, monitor.Label)

		data.SetId(fmt.Sprintf("%v", monitor.Label))
	}

	return nil
}

func resourceMonitorUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	monitorType := data.Get(TypeArg).(string)
	region := data.Get(AwsRegionArg).(string)
	roleToAssume := data.Get(IamRoleToAssumeArg).(string)
	queueUrl := data.Get(QueueUrlArg).(string)
	fileFormat := data.Get(FileFormatArg).(string)
	hostName := data.Get(HostNameArg).(string)
	parser := data.Get(ParserArg).(string)
	label := data.Get(LabelArg).(string)

	client := meta.(*scalyr.ScalyrConfig)
	input := &scalyr.UpdateMonitorInput{
		Type:         monitorType,
		Region:       region,
		RoleToAssume: roleToAssume,
		QueueUrl:     queueUrl,
		FileFormat:   fileFormat,
		HostName:     hostName,
		Parser:       parser,
		Label:        label,
	}

	unlock := synchronizer.Lock(monitorsLockName)
	defer unlock()

	if _, err := client.UpdateMonitor(ctx, input); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceMonitorDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*scalyr.ScalyrConfig)

	input := &scalyr.DeleteMonitorInput{
		Label: data.Id(),
	}

	unlock := synchronizer.Lock(monitorsLockName)
	defer unlock()

	if _, err := client.DeleteMonitor(ctx, input); err != nil {
		return diag.FromErr(err)
	} else {
		data.SetId("")
	}

	return nil
}
