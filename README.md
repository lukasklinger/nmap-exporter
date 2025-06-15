# nmap-exporter

nmap-exporter is a Prometheus exporter that periodically scans a specified network using nmap and exposes metrics about discovered hosts and scan results. It is designed to help you monitor network hosts and their availability using Prometheus.

## Features
- Periodically scans a configurable network range using nmap
- Exposes metrics via an HTTP endpoint for Prometheus scraping
- Metrics include:
  - Number of hosts up/down
  - List of discovered hosts (IP, hostname, status)
  - Scan warning count
  - Timestamp of the last scan

## Configuration
nmap-exporter is configured via environment variables:

| Variable         | Description                        | Example           | Required |
|------------------|------------------------------------|-------------------|----------|
| `METRICS_PATH`   | Path for Prometheus metrics        | `/metrics`        | Yes      |
| `METRICS_PORT`   | Port for metrics HTTP server       | `2112`            | Yes      |
| `SCAN_NETWORK`   | Network range to scan (CIDR)       | `192.168.1.0/24`  | Yes      |
| `SCAN_INTERVAL`  | Scan interval in seconds           | `300`             | Yes      |

Example usage:
```sh
export METRICS_PATH=/metrics
export METRICS_PORT=2112
export SCAN_NETWORK=192.168.1.0/24
export SCAN_INTERVAL=300

go run ./cmd/main.go
```

## Usage
1. **Set the required environment variables** as shown above.
2. **Run the exporter**:
   ```sh
   go run ./cmd/main.go
   ```
3. **Prometheus scraping**:
   - By default, metrics are exposed at `http://localhost:2112/metrics`.
   - Add a scrape config in your Prometheus configuration:
     ```yaml
     scrape_configs:
       - job_name: 'nmap-exporter'
         static_configs:
           - targets: ['localhost:2112']
     ```

## Metrics
- `nmap_exporter_hosts{ip,hostname,status}`: Gauge for each discovered host
- `nmap_exporter_host_count{state}`: Number of hosts up/down
- `nmap_exporter_scan_warning_count`: Number of warnings from the last scan
- `nmap_exporter_last_scan_timestamp`: UNIX timestamp of the last scan

## Requirements
- Go 1.18+
- nmap installed on the system

## License
MIT
