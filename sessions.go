package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

type SessionExport struct {
	Server      string `json:"server"`
	Protocol    string `json:"protocol"`
	P1Uniqueid  string `json:"p1uniqueid"`
	P2Uniqueid  string `json:"p2uniqueid"`
	State       string `json:"state"`
	RemoteHost  string `json:"remotehost"`
	RemotePort  string `json:"remoteport"`
	RemoteID    string `json:"remoteid"`
	RemoteTs    string `json:"remotets"`
	Established string `json:"established"`
	BytesIn     string `json:"bytesin"`
	BytesOut    string `json:"bytesout"`
	PacketsIn   string `json:"packetsin"`
	PacketsOut  string `json:"packetsout"`
}

func sessionsHandler(w http.ResponseWriter, r *http.Request, conf *Config, logger log.Logger) {
	serverName := conf.ServerName
	if serverName == "" {
		serverName = "openvpn-server"
	}

	connExport, err := getAllOpenVPNSessions(serverName, conf, logger)
	if err != nil {
		_ = level.Warn(logger).Log("task", "sessions handler", "msg", err.Error())
		connExport = []SessionExport{}
	}

	jsondata, err := json.Marshal(connExport)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprint(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsondata)
	if err != nil {
		_ = level.Error(logger).Log("task", "HTTP", "write", err.Error())
	}
}

func getAllOpenVPNSessions(server string, conf *Config, logger log.Logger) ([]SessionExport, error) {
	var allSessions []SessionExport

	// 2.7.x: single unified status file — protocol detected from address prefix
	if conf.OvpnStatus != "" {
		sessions, _, _, err := getOpenVPNSessions(server, conf.OvpnStatus, "", logger)
		if err != nil {
			_ = level.Warn(logger).Log("task", "read unified status", "err", err.Error())
		} else {
			allSessions = append(allSessions, sessions...)
			if conf.Debug {
				_ = level.Debug(logger).Log("task", "unified sessions", "count", len(sessions))
			}
		}
	}

	// 2.6.x: separate TCP status file, protocol forced to "tcp"
	if conf.OvpnTCPStatus != "" {
		sessions, _, _, err := getOpenVPNSessions(server, conf.OvpnTCPStatus, "tcp", logger)
		if err != nil {
			_ = level.Warn(logger).Log("task", "read TCP status", "err", err.Error())
		} else {
			allSessions = append(allSessions, sessions...)
			if conf.Debug {
				_ = level.Debug(logger).Log("task", "TCP sessions", "count", len(sessions))
			}
		}
	}

	// 2.6.x: separate UDP status file, protocol forced to "udp"
	if conf.OvpnUDPStatus != "" {
		sessions, _, _, err := getOpenVPNSessions(server, conf.OvpnUDPStatus, "udp", logger)
		if err != nil {
			_ = level.Warn(logger).Log("task", "read UDP status", "err", err.Error())
		} else {
			allSessions = append(allSessions, sessions...)
			if conf.Debug {
				_ = level.Debug(logger).Log("task", "UDP sessions", "count", len(sessions))
			}
		}
	}

	if len(allSessions) == 0 {
		return []SessionExport{}, fmt.Errorf("no OpenVPN sessions found in any status file")
	}

	return allSessions, nil
}

// parseRealAddress extracts the protocol (if detectable from prefix), host and port
// from the Real Address field in the OpenVPN status file.
//
// 2.6.x format:  "192.168.200.40:43154"
// 2.7.x formats: "tcp4-server:192.168.200.40:47054"
//                "udp4:192.168.200.10:45209"
//
// The detected protocol is returned as an empty string for 2.6.x entries since the
// caller already knows the protocol from which file it is reading.
func parseRealAddress(raw string) (detectedProto, host, port string) {
	// 2.7.x TCP: starts with "tcp"
	if strings.HasPrefix(raw, "tcp") {
		// raw = "tcp4-server:192.168.200.40:47054"
		// Strip everything up to and including the first colon to get "192.168.200.40:47054"
		rest := raw[strings.Index(raw, ":")+1:]
		host, port = splitHostPort(rest)
		return "tcp", host, port
	}

	// 2.7.x UDP: starts with "udp"
	if strings.HasPrefix(raw, "udp") {
		// raw = "udp4:192.168.200.10:45209"
		rest := raw[strings.Index(raw, ":")+1:]
		host, port = splitHostPort(rest)
		return "udp", host, port
	}

	// 2.6.x: plain "IP:PORT", no prefix
	host, port = splitHostPort(raw)
	return "", host, port
}

// splitHostPort splits "IP:PORT" into host and port. It uses the last colon as
// the separator so it is safe for both IPv4 and (bracketed) IPv6 addresses.
func splitHostPort(addr string) (host, port string) {
	idx := strings.LastIndex(addr, ":")
	if idx < 0 {
		return addr, ""
	}
	return addr[:idx], addr[idx+1:]
}

func getOpenVPNSessions(server string, statusFile string, forcedProtocol string, logger log.Logger) (connExport []SessionExport, product string, version string, err error) {
	reVersion := regexp.MustCompile(`(?m)^TITLE\s+(?P<product>\S+)\s+(?P<version>\S+)\s`)

	// Two regexes to handle the variable number of fields in 2.7 UDP lines.
	//
	// 2.6.x and 2.7 TCP line (13 tab-separated fields after CLIENT_LIST):
	//   CN  RealIP  vIPv4  vIPv6  rB  sB  startTime  unixTime  Username  ClientID  PeerID  Cipher
	//
	// 2.7 UDP line (12 fields — PeerID is absent):
	//   CN  RealIP  vIPv4  vIPv6  rB  sB  startTime  unixTime  Username  ClientID  Cipher
	//
	// Strategy: try the full 12-field regex first (TCP / 2.6.x). If it does not
	// produce a PeerID capture, fall back to the 11-field UDP regex.
	reClientFull := regexp.MustCompile(
		`(?m)^CLIENT_LIST\t` +
			`(?P<CN>\S+)\t` +
			`(?P<RealIP>\S+)\t` +
			`(?P<vIPv4>\S+)\t` +
			`(?P<vIPv6>\S*)\t` +
			`(?P<rB>\S+)\t` +
			`(?P<sB>\S+)\t` +
			`(?P<startTime>[^\t]+)\t` +
			`(?P<unixTime>\S+)\t` +
			`(?P<Username>[^\t]+)\t` +
			`(?P<ClientID>\S+)\t` +
			`(?P<PeerID>\S+)\t` +
			`(?P<Cipher>\S+)`,
	)

	// UDP-only regex without PeerID
	reClientUDP := regexp.MustCompile(
		`(?m)^CLIENT_LIST\t` +
			`(?P<CN>\S+)\t` +
			`(?P<RealIP>\S+)\t` +
			`(?P<vIPv4>\S+)\t` +
			`(?P<vIPv6>\S*)\t` +
			`(?P<rB>\S+)\t` +
			`(?P<sB>\S+)\t` +
			`(?P<startTime>[^\t]+)\t` +
			`(?P<unixTime>\S+)\t` +
			`(?P<Username>[^\t]+)\t` +
			`(?P<ClientID>\S+)\t` +
			`(?P<Cipher>\S+)`,
	)

	if statusFile == "" {
		return nil, "", "", fmt.Errorf("status file path is empty")
	}

	if _, err := os.Stat(statusFile); os.IsNotExist(err) {
		return nil, "", "", fmt.Errorf("status file does not exist: %s", statusFile)
	}

	ovpnstats, err := os.ReadFile(statusFile)
	if err != nil {
		_ = level.Error(logger).Log("task", "open status file", "file", statusFile, "err", err.Error())
		return nil, "", "", err
	}

	ver := reVersion.FindAllSubmatch(ovpnstats, -1)
	if ver != nil && len(ver[0]) == 3 {
		product = string(ver[0][1])
		version = string(ver[0][2])
	}

	// Build a set of byte offsets that are already matched by the full regex so
	// we can skip them when running the UDP fallback regex.
	fullMatches := reClientFull.FindAllSubmatchIndex(ovpnstats, -1)
	fullMatchedOffsets := make(map[int]bool, len(fullMatches))

	for _, loc := range fullMatches {
		fullMatchedOffsets[loc[0]] = true

		clientList := namedCaptures(reClientFull, ovpnstats, loc)
		detectedProto, host, port := parseRealAddress(clientList["RealIP"])

		proto := forcedProtocol
		if proto == "" {
			proto = detectedProto
		}

		cexp := SessionExport{
			Server:      server,
			Protocol:    proto,
			P1Uniqueid:  clientList["PeerID"],
			P2Uniqueid:  clientList["ClientID"],
			RemoteID:    clientList["CN"],
			State:       "ESTABLISHED",
			RemoteHost:  host,
			RemotePort:  port,
			RemoteTs:    clientList["vIPv4"],
			BytesIn:     clientList["rB"],
			BytesOut:    clientList["sB"],
			Established: clientList["startTime"],
			PacketsIn:   "0",
			PacketsOut:  "0",
		}
		connExport = append(connExport, cexp)
	}

	// Run the UDP fallback regex and process only lines not already matched above.
	udpMatches := reClientUDP.FindAllSubmatchIndex(ovpnstats, -1)
	for _, loc := range udpMatches {
		if fullMatchedOffsets[loc[0]] {
			continue // already processed by the full regex
		}

		clientList := namedCaptures(reClientUDP, ovpnstats, loc)
		detectedProto, host, port := parseRealAddress(clientList["RealIP"])

		proto := forcedProtocol
		if proto == "" {
			proto = detectedProto
		}

		cexp := SessionExport{
			Server:      server,
			Protocol:    proto,
			P1Uniqueid:  "", // PeerID not present in 2.7 UDP lines
			P2Uniqueid:  clientList["ClientID"],
			RemoteID:    clientList["CN"],
			State:       "ESTABLISHED",
			RemoteHost:  host,
			RemotePort:  port,
			RemoteTs:    clientList["vIPv4"],
			BytesIn:     clientList["rB"],
			BytesOut:    clientList["sB"],
			Established: clientList["startTime"],
			PacketsIn:   "0",
			PacketsOut:  "0",
		}
		connExport = append(connExport, cexp)
	}

	return connExport, product, version, nil
}

// namedCaptures returns a map of named subgroup -> matched string for a single
// FindAllSubmatchIndex result (loc) applied to data using the given regexp.
func namedCaptures(re *regexp.Regexp, data []byte, loc []int) map[string]string {
	result := make(map[string]string)
	for i, name := range re.SubexpNames() {
		if i == 0 || name == "" {
			continue
		}
		start, end := loc[2*i], loc[2*i+1]
		if start >= 0 && end >= 0 {
			result[name] = string(data[start:end])
		}
	}
	return result
}