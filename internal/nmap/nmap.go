package nmap

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Ullaakut/nmap/v3"
	"github.com/go-co-op/gocron/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type NMAP struct {
	Scanner   *nmap.Scanner
	Scheduler gocron.Scheduler
}

var (
	hosts = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "nmap_exporter_hosts",
		Help: "Currently available hosts, reset after each scan.",
	}, []string{"ip", "hostname", "status"})

	hostCount = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "nmap_exporter_host_count",
		Help: "Count of hosts, split by state.",
	}, []string{"state"})

	warningCount = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "nmap_exporter_scan_warning_count",
		Help: "Count of warnings for last scan, changes with each scan.",
	})

	timestamp = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "nmap_exporter_last_scan_timestamp",
		Help: "UNIX timestamp of the last scan.",
	})
)

// init nmap scanning support
func New(ctx context.Context, network string, intervalSeconds int) (NMAP, error) {
	n := NMAP{}

	// init nmap
	scanner, err := nmap.NewScanner(ctx, nmap.WithTargets(network), nmap.WithPingScan())
	if err != nil {
		return NMAP{}, fmt.Errorf("unable to create nmap scanner: %v", err)
	}

	n.Scanner = scanner

	// init scheduler
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		return n, err
	}

	// create job
	_, err = scheduler.NewJob(gocron.DurationJob(time.Duration(intervalSeconds)*time.Second), gocron.NewTask(
		func() {
			n.ScanAndExport()
		},
	))

	// start scheduler
	scheduler.Start()

	n.Scheduler = scheduler
	return n, err
}

// trigger a scan on the pre-configured scanner
func (n NMAP) scan() (result *nmap.Run, warnings *[]string, err error) {
	log.Default().Println("running scan")
	return n.Scanner.Run()
}

func (n NMAP) ScanAndExport() error {
	// run scan to get current values
	result, warnings, err := n.scan()
	if err != nil {
		return err
	}

	log.Default().Println("exporting metrics")

	// set metrics
	timestamp.SetToCurrentTime()
	warningCount.Set(float64(len(*warnings)))
	hostCount.WithLabelValues("up").Set(float64(result.Stats.Hosts.Up))
	hostCount.WithLabelValues("down").Set(float64(result.Stats.Hosts.Down))

	// reset host list
	hosts.Reset()

	// set host list
	for _, host := range result.Hosts {
		hostname := ""
		ip := ""

		if len(host.Hostnames) > 0 {
			hostname = host.Hostnames[0].Name
		}

		if len(host.Addresses) > 0 {
			ip = host.Addresses[0].Addr
		}

		hosts.WithLabelValues(ip, hostname, host.Status.State).Set(1)
	}

	return nil
}
