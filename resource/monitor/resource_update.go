package monitor

import (
	"context"
	scalyr "github.com/ansoni/terraform-provider-scalyr/scalyr-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceMonitorUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*scalyr.ScalyrConfig)
	id := data.Id()

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

	return diag.FromErr(updateMonitorsConfigurationFile(file, client))
}
