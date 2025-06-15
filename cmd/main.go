package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"cyaniccerulean.com/nmap-exporter/v2/internal/config"
	"cyaniccerulean.com/nmap-exporter/v2/internal/nmap"
	"cyaniccerulean.com/nmap-exporter/v2/internal/prometheus"
)

func main() {
	// signal for stopping
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGTERM)
	defer cancel()

	// read config
	metricsPath, scanNetwork, metricsPort, scanInterval, err := config.ReadEnvironment()
	if err != nil {
		log.Default().Fatalf("error reading config: %v\n", err)
	}

	// init Prometheus
	prometheus.New(metricsPath, metricsPort)

	// init NMAP
	scanner, err := nmap.New(ctx, scanNetwork, scanInterval)
	if err != nil {
		log.Default().Fatalf("error initializing nmap scanner: %v\n", err)
	}

	scanner.ScanAndExport()

	// wait for sigint
	<-ctx.Done()
	cancel()
}
