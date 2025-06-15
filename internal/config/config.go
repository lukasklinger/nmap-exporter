package config

import (
	"fmt"
	"os"
	"strconv"
)

func ReadEnvironment() (metricsPath, scanNetwork string, metricsPort, scanInterval int, err error) {
	metrics_path := os.Getenv("METRICS_PATH")
	metrics_port := os.Getenv("METRICS_PORT")

	scan_network := os.Getenv("SCAN_NETWORK")
	scan_interval := os.Getenv("SCAN_INTERVAL")

	if metrics_path == "" || scan_network == "" {
		return "", "", 0, 0, fmt.Errorf("environment variables not set correctly")
	}

	port, err := strconv.ParseInt(metrics_port, 10, 64)
	if err != nil {
		return "", "", 0, 0, fmt.Errorf("could not parse metrics port variable")
	}

	interval, err := strconv.ParseInt(scan_interval, 10, 64)
	if err != nil {
		return "", "", 0, 0, fmt.Errorf("could not parse scan interval variable")
	}

	return metrics_path, scan_network, int(port), int(interval), nil
}
