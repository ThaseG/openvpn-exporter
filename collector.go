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

func (c *OpenVPNCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- ovpnProbeSuccess
	ch <- ovpnInfo
	ch <- ovpnSessTotal
	ch <- ovpnBytesInTotal
	ch <- ovpnBytesOutTotal
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

func (c *OpenVPNCollector) Collect(ch chan<- prometheus.Metric) {
	_ = level.Debug(c.logger).Log("task", "Collecting OpenVPN metrics")

	var allSessions []SessionExport
	var products []string
	var versions []string
	var probeSuccess float64 = 0
	var versionFound string = ""

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

	if c.conf.OvpnTCPStatus != "" {
		_ = level.Debug(c.logger).Log("task", "Collecting TCP", "target", c.conf.OvpnTCPStatus)
		sess, product, version, err := getOpenVPNSessions("", c.conf.OvpnTCPStatus, "tcp", c.logger)
		if err != nil {
			_ = level.Error(c.logger).Log("task", "Collecting TCP", "status", "ERROR", "msg", err)
		} else {
			allSessions = append(allSessions, sess...)
			if product != "" && version != "" {
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

	if c.conf.OvpnUDPStatus != "" {
		_ = level.Debug(c.logger).Log("task", "Collecting UDP", "target", c.conf.OvpnUDPStatus)
		sess, product, version, err := getOpenVPNSessions("", c.conf.OvpnUDPStatus, "udp", c.logger)
		if err != nil {
			_ = level.Error(c.logger).Log("task", "Collecting UDP", "status", "ERROR", "msg", err)
		} else {
			allSessions = append(allSessions, sess...)
			if product != "" && version != "" {
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

	if versionFound != "" {
		probeSuccess = 1
	}

	ch <- prometheus.MustNewConstMetric(ovpnProbeSuccess, prometheus.GaugeValue, probeSuccess, versionFound)
	ch <- prometheus.MustNewConstMetric(ovpnSessTotal, prometheus.GaugeValue, float64(len(allSessions)))

	seen := make(map[string]bool)
	for i := range products {
		key := products[i] + ":" + versions[i]
		if !seen[key] {
			ch <- prometheus.MustNewConstMetric(ovpnInfo, prometheus.CounterValue, float64(1), products[i], versions[i])
			seen[key] = true
		}
	}

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