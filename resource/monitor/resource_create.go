package monitor

import (
	"context"
	"fmt"
	scalyr "github.com/ansoni/terraform-provider-scalyr/scalyr-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

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

	data.SetId(fmt.Sprintf("%v", label))

	return diagnostics
}
