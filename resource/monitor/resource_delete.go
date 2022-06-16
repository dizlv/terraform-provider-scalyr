package monitor

import (
	"context"
	scalyr "github.com/ansoni/terraform-provider-scalyr/scalyr-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/exp/slices"
)

func resourceMonitorDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*scalyr.ScalyrConfig)

	// File api does not has specific ID provided, so we create one by combining file path + version.
	label := data.Get("label").(string)

	monitors, err := getMonitorConfigurationFile(client)
	if err != nil {
		return diag.FromErr(err)
	}

	index := findMonitorByLabel(label, monitors.Monitors)

	if index == -1 {
		return diag.Errorf("Monitor is not present in submitted configuration file")
	}

	slices.Delete(monitors.Monitors, index, index+1)

	err = updateMonitorsConfigurationFile(monitors, client)
	if err != nil {
		return diag.FromErr(err)
	}

	// Set terraform id to "" which means that resource has been deleted successfully and can
	// be dropped from terraform state.
	data.SetId("")

	return nil
}
