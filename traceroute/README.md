# Traceroute Clone in Go

This is a simple traceroute clone implemented in Go. It sends UDP packets with increasing TTL (Time-To-Live) values and listens for ICMP "Time Exceeded" or "Destination Unreachable" messages to trace the route to a destination.

## Features

- Sends UDP packets with increasing TTL values.
- Listens for ICMP responses to determine the route.
- Supports DNS reverse lookup for intermediate hops.

## Prerequisites

- Go 1.16 or higher.
- Linux or macOS (requires raw socket access, which may require elevated privileges).

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/MogboPython/coding-challenges.git
   cd traceroute
   ```

2. Build the project:
   ```bash
   go build
   ```

## Usage

Run the traceroute tool with a destination IP address or hostname:
```bash
sudo ./traceroute-clone <destination>
```

Example:
```bash
sudo ./traceroute-clone dns.google.com
```

Output:
```
Tracing route to dns.google.com [8.8.8.8]
1: 192.168.1.1 (192.168.1.1) 1.234 ms
2: 10.0.0.1 (10.0.0.1) 2.345 ms
3: 203.0.113.1 (203.0.113.1) 3.456 ms
...
10: dns.google.com (8.8.8.8) 10.123 ms
Traceroute completed on hop 10
```
