package main

import (
	"os"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
)

type OpenVPNCollector struct {
	conf   *Config
	logger log.Logger

	// Background refresh cache
	mu          sync.RWMutex
	cachedState *collectedState
}

// collectedState holds the last successfully collected snapshot
type collectedState struct {
	sessions []SessionExport
	products []string
	versions []string
	version  string // first version found, used for probe_success label
	success  float64
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

// Implement prometheus.Collector Collect method — serves from cache
func (c *OpenVPNCollector) Collect(ch chan<- prometheus.Metric) {
	c.mu.RLock()
	state := c.cachedState
	c.mu.RUnlock()

	if state == nil {
		// No data yet — emit a failing probe and zeros
		ch <- prometheus.MustNewConstMetric(ovpnProbeSuccess, prometheus.GaugeValue, 0, "")
		ch <- prometheus.MustNewConstMetric(ovpnSessTotal, prometheus.GaugeValue, 0)
		return
	}

	ch <- prometheus.MustNewConstMetric(ovpnProbeSuccess, prometheus.GaugeValue, state.success, state.version)
	ch <- prometheus.MustNewConstMetric(ovpnSessTotal, prometheus.GaugeValue, float64(len(state.sessions)))

	seen := make(map[string]bool)
	for i := range state.products {
		key := state.products[i] + ":" + state.versions[i]
		if !seen[key] {
			ch <- prometheus.MustNewConstMetric(ovpnInfo, prometheus.CounterValue, 1, state.products[i], state.versions[i])
			seen[key] = true
		}
	}

	for _, v := range state.sessions {
		// Label uniquely identifies each session:
		// CN + virtual IP (unique per session on the server) + protocol
		clientLabel := v.RemoteID + "_" + v.RemoteTs + "_" + v.Protocol
		if bi, err := strconv.ParseFloat(v.BytesIn, 64); err == nil {
			ch <- prometheus.MustNewConstMetric(ovpnBytesInTotal, prometheus.CounterValue, bi, clientLabel)
		}
		if bo, err := strconv.ParseFloat(v.BytesOut, 64); err == nil {
			ch <- prometheus.MustNewConstMetric(ovpnBytesOutTotal, prometheus.CounterValue, bo, clientLabel)
		}
	}
}

// StartBackgroundRefresh launches the collection ticker at the configured interval.
// Interval is read from refresh_interval in exporter.yml (seconds); defaults to 15.
// Call this once after creating the collector (before registering with Prometheus).
func (c *OpenVPNCollector) StartBackgroundRefresh() {
	interval := c.conf.RefreshInterval
	if interval <= 0 {
		interval = 15
	}

	_ = level.Info(c.logger).Log("msg", "Starting background refresh", "interval_seconds", interval)

	// Collect immediately on start so the first scrape is never empty
	c.collect()

	go func() {
		ticker := time.NewTicker(time.Duration(interval) * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			c.collect()
		}
	}()
}

// collect does the actual work and stores the result in the cache
func (c *OpenVPNCollector) collect() {
	c.debugLog("task", "Collecting OpenVPN metrics")

	var allSessions []SessionExport
	var products []string
	var versions []string
	var versionFound string

	if c.conf.OvpnTCPStatus != "" {
		c.debugLog("task", "Collecting TCP", "target", c.conf.OvpnTCPStatus)
		sess, product, version, err := getOpenVPNSessions("", c.conf.OvpnTCPStatus, "tcp", c.logger)
		if err != nil {
			_ = level.Error(c.logger).Log("task", "Collecting TCP", "status", "ERROR", "msg", err)
		} else {
			allSessions = append(allSessions, sess...)
			if product != "" && version != "" {
				if versionFound == "" {
					versionFound = version
				}
				if !productVersionExists(products, versions, product, version) {
					products = append(products, product)
					versions = append(versions, version)
				}
			}
		}
	}

	if c.conf.OvpnUDPStatus != "" {
		c.debugLog("task", "Collecting UDP", "target", c.conf.OvpnUDPStatus)
		sess, product, version, err := getOpenVPNSessions("", c.conf.OvpnUDPStatus, "udp", c.logger)
		if err != nil {
			_ = level.Error(c.logger).Log("task", "Collecting UDP", "status", "ERROR", "msg", err)
		} else {
			allSessions = append(allSessions, sess...)
			if product != "" && version != "" {
				if versionFound == "" {
					versionFound = version
				}
				if !productVersionExists(products, versions, product, version) {
					products = append(products, product)
					versions = append(versions, version)
				}
			}
		}
	}

	var probeSuccess float64
	if versionFound != "" {
		probeSuccess = 1
	}

	c.debugLog("task", "Collection complete", "total_sessions", len(allSessions), "version", versionFound)

	c.mu.Lock()
	c.cachedState = &collectedState{
		sessions: allSessions,
		products: products,
		versions: versions,
		version:  versionFound,
		success:  probeSuccess,
	}
	c.mu.Unlock()
}

// debugLog emits a debug-level log only when debug is enabled in config
func (c *OpenVPNCollector) debugLog(keyvals ...interface{}) {
	if c.conf.Debug {
		_ = level.Debug(c.logger).Log(keyvals...)
	}
}

// productVersionExists checks for duplicate product/version pairs
func productVersionExists(products, versions []string, product, version string) bool {
	for i := range products {
		if products[i] == product && versions[i] == version {
			return true
		}
	}
	return false
}

func getVersionFromStatusFile(statusFile string, logger log.Logger) (product string, version string, err error) {
	reVersion := regexp.MustCompile(`(?m)^TITLE\s+(?P<product>\S+)\s+(?P<version>\S+)\s`)

	if statusFile == "" {
		return "", "", nil
	}

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