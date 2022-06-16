package monitor

import (
	"context"
	"fmt"
	scalyr "github.com/ansoni/terraform-provider-scalyr/scalyr-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceMonitorRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*scalyr.ScalyrConfig)
	label := data.Id()

	var diagnostics diag.Diagnostics

	file, err := getMonitorConfigurationFile(client)
	if err != nil {
		return diag.FromErr(err)
	}

	index := findMonitorByLabel(label, file.Monitors)

	if index == -1 {
		return diag.Errorf("Monitor with configured label already exist")
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
