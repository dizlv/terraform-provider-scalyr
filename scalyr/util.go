package scalyr

import (
	"encoding/json"
	scalyr "github.com/ansoni/terraform-provider-scalyr/scalyr-go"
	"golang.org/x/exp/slices"
)

// getMonitorConfigurationFile retrieves monitor configuration file. It is stored under
// `ConfigurationFilePath` constant.
func getMonitorConfigurationFile(client *scalyr.ScalyrConfig) (*scalyr.MonitorsConfigurationFile, error) {
	response, err := client.GetMonitor(scalyr.ConfigurationFilePath)
	if err != nil {
		return nil, err
	}

	file := &scalyr.MonitorsConfigurationFile{}

	err = json.Unmarshal([]byte(response.Content), file)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func findMonitorByLabel(label string, monitors scalyr.Monitors) int {
	index := slices.IndexFunc(monitors, func(m *scalyr.Monitor) bool {
		return m.Label == label
	})

	return index
}

func updateMonitorsConfigurationFile(file *scalyr.MonitorsConfigurationFile, client *scalyr.ScalyrConfig) error {
	err := client.UpdateMonitors(file)

	if err != nil {
		return err
	}

	return nil
}
