package main

// Config holds the exporter configuration loaded from exporter.yml
type Config struct {
	OvpnTCPStatus   string `yaml:"ovpn_tcp_status"`
	OvpnUDPStatus   string `yaml:"ovpn_udp_status"`
	ServerName      string `yaml:"server_name"`
	Debug           bool   `yaml:"debug"`.           // true/false default false
	RefreshInterval int    `yaml:"refresh_interval"` // in seconds, default 15
}