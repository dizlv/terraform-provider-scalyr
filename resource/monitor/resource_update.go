package monitor

import (
	"context"
	"fmt"
	scalyr "github.com/ansoni/terraform-provider-scalyr/scalyr-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceMonitorUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*scalyr.ScalyrConfig)

	// File api does not has specific ID provided, so we create one by combining file path + version.
	label := data.Get("label").(string)

	monitors, err := getMonitorConfigurationFile(client)
	if err != nil {
		return diag.FromErr(err)
	}

	index := findMonitorByLabel(label, monitors.Monitors)
	fmt.Printf("%d", index)

	return diag.FromErr(updateMonitorsConfigurationFile(monitors, client))
}
