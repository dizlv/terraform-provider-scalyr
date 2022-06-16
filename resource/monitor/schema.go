package monitor

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

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

func Resource() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceMonitorRead,
		CreateContext: resourceMonitorCreate,
		DeleteContext: resourceMonitorDelete,
		UpdateContext: resourceMonitorUpdate,

		Schema: Schema,
	}
}
