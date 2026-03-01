package main

// Config holds the exporter configuration loaded from exporter.yml
type Config struct {
	OvpnTCPStatus   string `yaml:"ovpn_tcp_status"`   // 2.6.x separate TCP status file
	OvpnUDPStatus   string `yaml:"ovpn_udp_status"`   // 2.6.x separate UDP status file
	OvpnStatus      string `yaml:"ovpn_status"`       // 2.7.x unified status file (TCP+UDP combined)
	ServerName      string `yaml:"server_name"`
	Debug           bool   `yaml:"debug"`             // true/false, default false
	RefreshInterval int    `yaml:"refresh_interval"`  // in seconds, default 15
}