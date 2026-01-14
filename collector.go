package main

import (
	"os"
	"regexp"
	"strconv"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
)

type OpenVPNCollector struct {
	conf   *Config
	logger log.Logger
}

var (
	ovpnInfo          = prometheus.NewDesc("ovpn_info", "Software info", []string{"product", "version"}, nil)
	ovpnSessTotal     = prometheus.NewDesc("ovpn_sessions_total", "Total number of active sessions", nil, nil)
	ovpnBytesInTotal  = prometheus.NewDesc("ovpn_bytes_in_total", "Total number of bytes received", []string{"client"}, nil)
	ovpnBytesOutTotal = prometheus.NewDesc("ovpn_bytes_out_total", "Total number of bytes sent", []string{"client"}, nil)
	ovpnProbeSuccess  = prometheus.NewDesc("probe_success", "OpenVPN Status", []string{"version"}, nil)
)

// Implement prometheus.Collector Describe method
func (c *OpenVPNCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- ovpnProbeSuccess
	ch <- ovpnInfo
	ch <- ovpnSessTotal
	ch <- ovpnBytesInTotal
	ch <- ovpnBytesOutTotal
}

// getVersionFromStatusFile extracts version info from OpenVPN status file
func getVersionFromStatusFile(statusFile string, logger log.Logger) (product string, version string, err error) {
	reVersion := regexp.MustCompile(`(?m)^TITLE\s+(?P<product>\S+)\s+(?P<version>\S+)\s`)
	
	if statusFile == "" {
		return "", "", nil
	}
	
	// Check if file exists
	if _, err := os.Stat(statusFile); os.IsNotExist(err) {
		return "", "", err
	}

	ovpnstats, err := os.ReadFile(statusFile)
	if err != nil {
		_ = level.Error(logger).Log("task", "read status file for version", "file", statusFile, "err", err.Error())
		return "", "", err
	}

	ver := reVersion.FindAllSubmatch(ovpnstats, -1)
	if ver != nil && len(ver[0]) == 3 {
		product = string(ver[0][1])
		version = string(ver[0][2])
	}

	return product, version, nil
}

// Implement prometheus.Collector Collect method
func (c *OpenVPNCollector) Collect(ch chan<- prometheus.Metric) {
	_ = level.Debug(c.logger).Log("task", "Collecting OpenVPN metrics")

	var allSessions []SessionExport
	var products []string
	var versions []string
	var probeSuccess float64 = 0
	var versionFound string = ""

	// Try to get version from TCP status file first
	if c.conf.OvpnTCPStatus != "" {
		product, version, err := getVersionFromStatusFile(c.conf.OvpnTCPStatus, c.logger)
		if err == nil && version != "" {
			versionFound = version
			if product != "" {
				products = append(products, product)
				versions = append(versions, version)
			}
		}
	}

	// Try to get version from UDP status file if not found yet
	if versionFound == "" && c.conf.OvpnUDPStatus != "" {
		product, version, err := getVersionFromStatusFile(c.conf.OvpnUDPStatus, c.logger)
		if err == nil && version != "" {
			versionFound = version
			if product != "" {
				products = append(products, product)
				versions = append(versions, version)
			}
		}
	}

	// Collect TCP sessions
	if c.conf.OvpnTCPStatus != "" {
		_ = level.Debug(c.logger).Log("task", "Collecting TCP", "target", c.conf.OvpnTCPStatus)
		sess, product, version, err := getOpenVPNSessions("", c.conf.OvpnTCPStatus, "tcp", c.logger)
		if err != nil {
			_ = level.Error(c.logger).Log("task", "Collecting TCP", "status", "ERROR", "msg", err)
		} else {
			allSessions = append(allSessions, sess...)
			if product != "" && version != "" {
				// Add to products/versions if not already there
				alreadyExists := false
				for i := range products {
					if products[i] == product && versions[i] == version {
						alreadyExists = true
						break
					}
				}
				if !alreadyExists {
					products = append(products, product)
					versions = append(versions, version)
				}
			}
		}
	}

	// Collect UDP sessions
	if c.conf.OvpnUDPStatus != "" {
		_ = level.Debug(c.logger).Log("task", "Collecting UDP", "target", c.conf.OvpnUDPStatus)
		sess, product, version, err := getOpenVPNSessions("", c.conf.OvpnUDPStatus, "udp", c.logger)
		if err != nil {
			_ = level.Error(c.logger).Log("task", "Collecting UDP", "status", "ERROR", "msg", err)
		} else {
			allSessions = append(allSessions, sess...)
			if product != "" && version != "" {
				// Add to products/versions if not already there
				alreadyExists := false
				for i := range products {
					if products[i] == product && versions[i] == version {
						alreadyExists = true
						break
					}
				}
				if !alreadyExists {
					products = append(products, product)
					versions = append(versions, version)
				}
			}
		}
	}

	// Determine probe success - if we have version info, it means status file is readable
	if versionFound != "" {
		probeSuccess = 1
	}

	// Always export probe_success with version (even if empty when both files are unreadable)
	ch <- prometheus.MustNewConstMetric(ovpnProbeSuccess, prometheus.GaugeValue, probeSuccess, versionFound)

	// Always export total sessions (even if 0)
	ch <- prometheus.MustNewConstMetric(ovpnSessTotal, prometheus.GaugeValue, float64(len(allSessions)))

	// Software info for each unique product/version combination
	seen := make(map[string]bool)
	for i := range products {
		key := products[i] + ":" + versions[i]
		if !seen[key] {
			ch <- prometheus.MustNewConstMetric(ovpnInfo, prometheus.CounterValue, float64(1), products[i], versions[i])
			seen[key] = true
		}
	}

	// Per-client metrics (only if sessions exist)
	for _, v := range allSessions {
		clientLabel := v.RemoteID + "_" + v.Protocol
		if bi, err := strconv.ParseFloat(v.BytesIn, 64); err == nil {
			ch <- prometheus.MustNewConstMetric(ovpnBytesInTotal, prometheus.CounterValue, float64(bi), clientLabel)
		}
		if bo, err := strconv.ParseFloat(v.BytesOut, 64); err == nil {
			ch <- prometheus.MustNewConstMetric(ovpnBytesOutTotal, prometheus.CounterValue, float64(bo), clientLabel)
		}
	}

	_ = level.Debug(c.logger).Log("task", "Collection complete", "total_sessions", len(allSessions), "version", versionFound)
}
```

## Key Changes Explained:

1. **Added `getVersionFromStatusFile()` function**: This extracts version information from the status file without parsing sessions. This allows us to get the version even when there are no active connections.

2. **Always export baseline metrics**: The collector now always exports:
   - `probe_success{version="X.X.X"} 1` (or `0` if status files are completely unreadable)
   - `ovpn_sessions_total 0` (or actual count)
   - `ovpn_info{product="OpenVPN",version="X.X.X"} 1`

3. **Version detection priority**: 
   - First tries to get version from TCP status file
   - Falls back to UDP status file if TCP doesn't have version info
   - Uses empty string if neither file is readable

4. **Probe success logic**: 
   - `probe_success = 1` when version is found (meaning status file is readable)
   - `probe_success = 0` when no version is found (meaning status files are unreadable or missing)

## Expected Output With No Connections:
```
# HELP probe_success OpenVPN Status
# TYPE probe_success gauge
probe_success{version="2.6.15"} 1

# HELP ovpn_info Software info
# TYPE ovpn_info counter
ovpn_info{product="OpenVPN",version="2.6.15"} 1

# HELP ovpn_sessions_total Total number of active sessions
# TYPE ovpn_sessions_total gauge
ovpn_sessions_total 0