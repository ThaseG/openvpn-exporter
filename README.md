# OpenVPN Exporter

Prometheus exporter for OpenVPN server metrics with JSON session API endpoint.

## Example usage

To see an example usage and e2e tests, please visit this [openvpn repository.](https://github.com/ThaseG/openvpn)

## Images

#### Main page
`http://<OpenVPN server IP>:9234`

![Exporter page](images/main.png)

#### Local clients
`http://<OpenVPN server IP>:9234/static`

![Local clients](images/local_sessions.png)

#### Exporter Metrics
```
# http://<OpenVPN server IP>:9234/metrics

# HELP ovpn_bytes_in_total Total number of bytes received
# TYPE ovpn_bytes_in_total counter
ovpn_bytes_in_total{client="bookworm.openvpn.com_tcp"} 3317
# HELP ovpn_bytes_out_total Total number of bytes sent
# TYPE ovpn_bytes_out_total counter
ovpn_bytes_out_total{client="bookworm.openvpn.com_tcp"} 3616
# HELP ovpn_info Software info
# TYPE ovpn_info counter
ovpn_info{product="OpenVPN",version="2.6.15"} 1
# HELP ovpn_sessions_total Total number of active sessions
# TYPE ovpn_sessions_total gauge
ovpn_sessions_total 1
# HELP probe_success OpenVPN Status
# TYPE probe_success gauge
probe_success{version="2.6.15"} 1
```
#### Internal Exporter Metrics
```
# http://<OpenVPN server IP>:9234/prom

# HELP go_gc_duration_seconds A summary of the wall-time pause (stop-the-world) duration in garbage collection cycles.
# TYPE go_gc_duration_seconds summary
go_gc_duration_seconds{quantile="0"} 0.000274894
go_gc_duration_seconds{quantile="0.25"} 0.000274894
go_gc_duration_seconds{quantile="0.5"} 0.000274894
go_gc_duration_seconds{quantile="0.75"} 0.000274894
go_gc_duration_seconds{quantile="1"} 0.000274894
go_gc_duration_seconds_sum 0.000274894
go_gc_duration_seconds_count 1
# HELP go_gc_gogc_percent Heap size target percentage configured by the user, otherwise 100. This value is set by the GOGC environment variable, and the runtime/debug.SetGCPercent function. Sourced from /gc/gogc:percent
# TYPE go_gc_gogc_percent gauge
go_gc_gogc_percent 100
# HELP go_gc_gomemlimit_bytes Go runtime memory limit configured by the user, otherwise math.MaxInt64. This value is set by the GOMEMLIMIT environment variable, and the runtime/debug.SetMemoryLimit function. Sourced from /gc/gomemlimit:bytes
# TYPE go_gc_gomemlimit_bytes gauge
go_gc_gomemlimit_bytes 9.223372036854776e+18
# HELP go_goroutines Number of goroutines that currently exist.
# TYPE go_goroutines gauge
go_goroutines 7
# HELP go_info Information about the Go environment.
# TYPE go_info gauge
go_info{version="go1.25.3"} 1
# HELP go_memstats_alloc_bytes Number of bytes allocated in heap and currently in use. Equals to /memory/classes/heap/objects:bytes.
# TYPE go_memstats_alloc_bytes gauge
go_memstats_alloc_bytes 2.269552e+06
# HELP go_memstats_alloc_bytes_total Total number of bytes allocated in heap until now, even if released already. Equals to /gc/heap/allocs:bytes.
# TYPE go_memstats_alloc_bytes_total counter
go_memstats_alloc_bytes_total 3.201328e+06
# HELP go_memstats_buck_hash_sys_bytes Number of bytes used by the profiling bucket hash table. Equals to /memory/classes/profiling/buckets:bytes.
# TYPE go_memstats_buck_hash_sys_bytes gauge
go_memstats_buck_hash_sys_bytes 1.444008e+06
# HELP go_memstats_frees_total Total number of heap objects frees. Equals to /gc/heap/frees:objects + /gc/heap/tiny/allocs:objects.
# TYPE go_memstats_frees_total counter
go_memstats_frees_total 6874
# HELP go_memstats_gc_sys_bytes Number of bytes used for garbage collection system metadata. Equals to /memory/classes/metadata/other:bytes.
# TYPE go_memstats_gc_sys_bytes gauge
go_memstats_gc_sys_bytes 2.107152e+06
# HELP go_memstats_heap_alloc_bytes Number of heap bytes allocated and currently in use, same as go_memstats_alloc_bytes. Equals to /memory/classes/heap/objects:bytes.
# TYPE go_memstats_heap_alloc_bytes gauge
go_memstats_heap_alloc_bytes 2.269552e+06
# HELP go_memstats_heap_idle_bytes Number of heap bytes waiting to be used. Equals to /memory/classes/heap/released:bytes + /memory/classes/heap/free:bytes.
# TYPE go_memstats_heap_idle_bytes gauge
go_memstats_heap_idle_bytes 4.456448e+06
# HELP go_memstats_heap_inuse_bytes Number of heap bytes that are in use. Equals to /memory/classes/heap/objects:bytes + /memory/classes/heap/unused:bytes
# TYPE go_memstats_heap_inuse_bytes gauge
go_memstats_heap_inuse_bytes 3.538944e+06
# HELP go_memstats_heap_objects Number of currently allocated objects. Equals to /gc/heap/objects:objects.
# TYPE go_memstats_heap_objects gauge
go_memstats_heap_objects 2724
# HELP go_memstats_heap_released_bytes Number of heap bytes released to OS. Equals to /memory/classes/heap/released:bytes.
# TYPE go_memstats_heap_released_bytes gauge
go_memstats_heap_released_bytes 3.555328e+06
# HELP go_memstats_heap_sys_bytes Number of heap bytes obtained from system. Equals to /memory/classes/heap/objects:bytes + /memory/classes/heap/unused:bytes + /memory/classes/heap/released:bytes + /memory/classes/heap/free:bytes.
# TYPE go_memstats_heap_sys_bytes gauge
go_memstats_heap_sys_bytes 7.995392e+06
# HELP go_memstats_last_gc_time_seconds Number of seconds since 1970 of last garbage collection.
# TYPE go_memstats_last_gc_time_seconds gauge
go_memstats_last_gc_time_seconds 1.762202458882669e+09
# HELP go_memstats_mallocs_total Total number of heap objects allocated, both live and gc-ed. Semantically a counter version for go_memstats_heap_objects gauge. Equals to /gc/heap/allocs:objects + /gc/heap/tiny/allocs:objects.
# TYPE go_memstats_mallocs_total counter
go_memstats_mallocs_total 9598
# HELP go_memstats_mcache_inuse_bytes Number of bytes in use by mcache structures. Equals to /memory/classes/metadata/mcache/inuse:bytes.
# TYPE go_memstats_mcache_inuse_bytes gauge
go_memstats_mcache_inuse_bytes 2416
# HELP go_memstats_mcache_sys_bytes Number of bytes used for mcache structures obtained from system. Equals to /memory/classes/metadata/mcache/inuse:bytes + /memory/classes/metadata/mcache/free:bytes.
# TYPE go_memstats_mcache_sys_bytes gauge
go_memstats_mcache_sys_bytes 15704
# HELP go_memstats_mspan_inuse_bytes Number of bytes in use by mspan structures. Equals to /memory/classes/metadata/mspan/inuse:bytes.
# TYPE go_memstats_mspan_inuse_bytes gauge
go_memstats_mspan_inuse_bytes 53920
# HELP go_memstats_mspan_sys_bytes Number of bytes used for mspan structures obtained from system. Equals to /memory/classes/metadata/mspan/inuse:bytes + /memory/classes/metadata/mspan/free:bytes.
# TYPE go_memstats_mspan_sys_bytes gauge
go_memstats_mspan_sys_bytes 65280
# HELP go_memstats_next_gc_bytes Number of heap bytes when next garbage collection will take place. Equals to /gc/heap/goal:bytes.
# TYPE go_memstats_next_gc_bytes gauge
go_memstats_next_gc_bytes 4.396978e+06
# HELP go_memstats_other_sys_bytes Number of bytes used for other system allocations. Equals to /memory/classes/other:bytes.
# TYPE go_memstats_other_sys_bytes gauge
go_memstats_other_sys_bytes 584952
# HELP go_memstats_stack_inuse_bytes Number of bytes obtained from system for stack allocator in non-CGO environments. Equals to /memory/classes/heap/stacks:bytes.
# TYPE go_memstats_stack_inuse_bytes gauge
go_memstats_stack_inuse_bytes 393216
# HELP go_memstats_stack_sys_bytes Number of bytes obtained from system for stack allocator. Equals to /memory/classes/heap/stacks:bytes + /memory/classes/os-stacks:bytes.
# TYPE go_memstats_stack_sys_bytes gauge
go_memstats_stack_sys_bytes 393216
# HELP go_memstats_sys_bytes Number of bytes obtained from system. Equals to /memory/classes/total:byte.
# TYPE go_memstats_sys_bytes gauge
go_memstats_sys_bytes 1.2605704e+07
# HELP go_sched_gomaxprocs_threads The current runtime.GOMAXPROCS setting, or the number of operating system threads that can execute user-level Go code simultaneously. Sourced from /sched/gomaxprocs:threads
# TYPE go_sched_gomaxprocs_threads gauge
go_sched_gomaxprocs_threads 2
# HELP go_threads Number of OS threads created.
# TYPE go_threads gauge
go_threads 4
# HELP ovpn_bytes_in_total Total number of bytes received
# TYPE ovpn_bytes_in_total counter
ovpn_bytes_in_total{client="bookworm.openvpn.com_tcp"} 2837
# HELP ovpn_bytes_out_total Total number of bytes sent
# TYPE ovpn_bytes_out_total counter
ovpn_bytes_out_total{client="bookworm.openvpn.com_tcp"} 3112
# HELP ovpn_info Software info
# TYPE ovpn_info counter
ovpn_info{product="OpenVPN",version="2.6.15"} 1
# HELP ovpn_sessions_total Total number of active sessions
# TYPE ovpn_sessions_total gauge
ovpn_sessions_total 1
# HELP probe_success OpenVPN Status
# TYPE probe_success gauge
probe_success{version="2.6.15"} 1
# HELP process_cpu_seconds_total Total user and system CPU time spent in seconds.
# TYPE process_cpu_seconds_total counter
process_cpu_seconds_total 0.02
# HELP process_max_fds Maximum number of open file descriptors.
# TYPE process_max_fds gauge
process_max_fds 1.048576e+06
# HELP process_network_receive_bytes_total Number of bytes received by the process over the network.
# TYPE process_network_receive_bytes_total counter
process_network_receive_bytes_total 18478
# HELP process_network_transmit_bytes_total Number of bytes sent by the process over the network.
# TYPE process_network_transmit_bytes_total counter
process_network_transmit_bytes_total 28386
# HELP process_open_fds Number of open file descriptors.
# TYPE process_open_fds gauge
process_open_fds 10
# HELP process_resident_memory_bytes Resident memory size in bytes.
# TYPE process_resident_memory_bytes gauge
process_resident_memory_bytes 1.4336e+07
# HELP process_start_time_seconds Start time of the process since unix epoch in seconds.
# TYPE process_start_time_seconds gauge
process_start_time_seconds 1.76220227316e+09
# HELP process_virtual_memory_bytes Virtual memory size in bytes.
# TYPE process_virtual_memory_bytes gauge
process_virtual_memory_bytes 1.264467968e+09
# HELP process_virtual_memory_max_bytes Maximum amount of virtual memory available in bytes.
# TYPE process_virtual_memory_max_bytes gauge
process_virtual_memory_max_bytes 1.8446744073709552e+19
# HELP promhttp_metric_handler_requests_in_flight Current number of scrapes being served.
# TYPE promhttp_metric_handler_requests_in_flight gauge
promhttp_metric_handler_requests_in_flight 1
# HELP promhttp_metric_handler_requests_total Total number of scrapes by HTTP status code.
# TYPE promhttp_metric_handler_requests_total counter
promhttp_metric_handler_requests_total{code="200"} 1
promhttp_metric_handler_requests_total{code="500"} 0
promhttp_metric_handler_requests_total{code="503"} 0
```
#### Local Sessions Local
```
# http://<OpenVPN server IP>:9234/sessions_local

[{"server":"openvpn-server","protocol":"tcp","p1uniqueid":"0","p2uniqueid":"0","state":"ESTABLISHED","remotehost":"78.99.236.15","remoteport":"53993","remoteid":"bookworm.openvpn.com","remotets":"10.0.0.2","established":"2025-11-03 20:40:03","bytesin":"3557","bytesout":"3868","packetsin":"0","packetsout":"0"}]
```

## Features

- **Prometheus Metrics**: Exposes OpenVPN session metrics in Prometheus format
- **JSON API**: Provides detailed session information via `/sessions` endpoint
- **Multi-Protocol Support**: Collects data from both TCP and UDP OpenVPN instances depends on `exporter.yaml` config
- **Flexible Configuration**: Use one or both status files

## Metrics Exposed

### Prometheus Metrics (`/metrics`)

- `probe_success` - OpenVPN status (1 = up, 0 = down) with version label
- `ovpn_info` - Software information (product, version)
- `ovpn_sessions_total` - Total number of active VPN sessions
- `ovpn_bytes_in_total` - Bytes received per client (labeled by client CN)
- `ovpn_bytes_out_total` - Bytes sent per client (labeled by client CN)

### JSON API (`/sessions`)

Returns detailed session information including:
- Server hostname
- Client remote IP and port
- Client certificate CN (Common Name)
- Virtual IP address assigned
- Connection state
- Bytes/packets in/out
- Session established time

## Prerequisites

### OpenVPN Server Configuration

Your OpenVPN server must be configured to generate a status file. Add this to your OpenVPN server config:

```
# For UDP status file
status /home/openvpn/logs/openvpn-udp-status
# For TCP status file
status /home/openvpn/logs/openvpn-tcp-status
status-version 3
```

The status file will be automatically updated by OpenVPN with current client connections.

## Installation

### Build from Source

```bash
# Clone the repository
git clone <your-repo-url>
cd openvpn-exporter

# Download dependencies
go mod download

# Build
go build -o openvpn-exporter

# Optional: Build with version info
go build -ldflags "-X main.version=1.0.0 -X main.buildDate=$(date -u +%Y-%m-%d)" -o openvpn-exporter
```

## Configuration

Create a `../Server/exporter.yaml` file:

```yaml
# Path to OpenVPN status file/s
ovpntcpstatus: "/home/openvpn/logs/openvpn-tcp-status"
ovpnudpstatus: "/home/openvpn/logs/openvpn-udp-status"
```

## Usage

### Basic Usage

```bash
./openvpn-exporter
```

### With Custom Configuration

```bash
./openvpn-exporter --config.file=/etc/openvpn-exporter/exporter.yaml
```

### Command-line Options

```
  --web.listen-address=":9234"
      Address to listen on for web interface and telemetry
  
  --web.telemetry-path="/metrics"
      Path under which to expose metrics
  
  --config.file="expoter.yaml"
      Path to configuration file
  
  --help, -h
      Show help
  
  --version
      Show version information
```

## Endpoints

- `/` - Landing page with links
- `/metrics` - Prometheus metrics endpoint
- `/sessions` - JSON API with detailed session information

## Prometheus Configuration

Add this to your `prometheus.yml`:

```yaml
scrape_configs:
  - job_name: 'openvpn'
    static_configs:
      - targets: ['localhost:9234']
```

## Example Outputs

### Prometheus Metrics

```
# HELP probe_success OpenVPN Status
# TYPE probe_success gauge
probe_success{version="2.6.15"} 1

# HELP ovpn_info Software info
# TYPE ovpn_info counter
ovpn_info{product="OpenVPN",version="2.6.15"} 1

# HELP ovpn_sessions_total Total number of active sessions
# TYPE ovpn_sessions_total gauge
ovpn_sessions_total 3

# HELP ovpn_bytes_in_total Total number of bytes received
# TYPE ovpn_bytes_in_total counter
ovpn_bytes_in_total{client="user1_tcp"} 1234567
ovpn_bytes_in_total{client="user2_udp"} 9876543

# HELP ovpn_bytes_out_total Total number of bytes sent
# TYPE ovpn_bytes_out_total counter
ovpn_bytes_out_total{client="user1_tcp"} 7654321
ovpn_bytes_out_total{client="user2_udp"} 3456789
```

### JSON Sessions API

```json
[
  {
    "server": "vpn-server-01",
    "protocol": "tcp",
    "p1uniqueid": "1234567890",
    "p2uniqueid": "123",
    "state": "ESTABLISHED",
    "remotehost": "203.0.113.45",
    "remoteport": "54321",
    "remoteid": "user1",
    "remotets": "10.8.0.6",
    "established": "2024-01-15 10:30:45",
    "bytesin": "1234567",
    "bytesout": "7654321",
    "packetsin": "0",
    "packetsout": "0"
  },
  {
    "server": "vpn-server-01",
    "protocol": "udp",
    "p1uniqueid": "9876543210",
    "p2uniqueid": "456",
    "state": "ESTABLISHED",
    "remotehost": "203.0.113.89",
    "remoteport": "12345",
    "remoteid": "user2",
    "remotets": "10.8.0.10",
    "established": "2024-01-15 11:15:22",
    "bytesin": "9876543",
    "bytesout": "3456789",
    "packetsin": "0",
    "packetsout": "0"
  }
]
```

## Security Considerations

1. **File Permissions**: Ensure the OpenVPN status file is readable by the exporter user
2. **Firewall**: Restrict access to port 9234 to authorized hosts only
3. **Sensitive Data**: The `/sessions` endpoint exposes client IPs and CNs - consider authentication

## Troubleshooting

### No metrics shown

- Check that at least one OpenVPN status file exists and is readable
- Verify the paths in `../Server/exporter.yaml` are correct
- Ensure OpenVPN is configured with `status` directive
- Check logs: the exporter will warn if it can't read a status file but continue if at least one is available

### Only TCP or UDP sessions showing

- Verify both status files exist if you configured both
- Check file permissions on both status files
- Look at exporter logs for warnings about which file couldn't be read

### Permission denied

```bash
# Give read permission to the status files
sudo chmod 644 /home/openvpn/logs/openvpn-tcp-status
sudo chmod 644 /home/openvpn/logs/openvpn-udp-status

# Or add exporter user to openvpn group
sudo usermod -aG openvpn prometheus
```

### Status file format not recognized

This exporter expects OpenVPN status version 3 format. Ensure your OpenVPN config has:

```
status-version 3
```

## License

This project maintains compatibility with the original exporter format.

## Contributing

Contributions are welcome! Please submit pull requests or open issues for bugs and feature requests.