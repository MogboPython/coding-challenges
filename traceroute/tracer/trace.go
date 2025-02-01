package tracer

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"os"
	"syscall"
	"time"
)

const (
	MAX_HOPS                       int           = 64
	PACKET_SIZE                    int           = 32
	UDP_PORT                       int           = 33434
	TIMEOUT                        time.Duration = 1000 * time.Millisecond
	ICMPTypeDestinationUnreachable int           = 3
	ICMPTypeTimeExceeded           int           = 11
)

// I don't need body field later so I didn't include it in the struct
type Message struct {
	Type     int
	Code     int
	Checksum int
}

func Trace(ipAddr net.IPAddr) error {

	// Raw socket for receiving ICMP messages
	icmpConn, err := net.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		return fmt.Errorf("failed to create ICMP socket: %v", err)
	}
	defer icmpConn.Close()

	// UDP socket for sending packets
	addr := &net.UDPAddr{
		IP:   ipAddr.IP,
		Port: int(UDP_PORT),
	}

	udpConn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return fmt.Errorf("failed to create UDP socket: %v", err)
	}
	defer udpConn.Close()

	// Get the underlying file descriptor
	f, err := udpConn.File()
	if err != nil {
		return fmt.Errorf("to get file descriptor: %v", err)
	}
	defer f.Close()

	for TTL := 1; TTL <= MAX_HOPS; TTL++ {
		// Set TTL option
		fd := int(f.Fd())
		err = syscall.SetsockoptInt(fd, syscall.IPPROTO_IP, syscall.IP_TTL, TTL)
		if err != nil {
			return fmt.Errorf("failed to set TTL: %v", err)
		}

		// Send message
		// I added retry sending message because I was getting an error "write udp 192.168.*.***:*****->8.8.8.8:*****: write: no route to host" while doing it only once
		startTime := time.Now()
		for attempt := 0; attempt < 3; attempt++ {
			_, err = udpConn.Write([]byte("TRACEROUTE"))
			if err != nil {
				time.Sleep(TIMEOUT) // Wait before retrying
				continue
			}
			break
		}

		// Wait for an ICMP response
		response := make([]byte, 1500)
		if err := icmpConn.SetReadDeadline(time.Now().Add(TIMEOUT)); err != nil {
			return fmt.Errorf("failed to set read deadline: %v", err)
		}

		n, from, err := icmpConn.ReadFrom(response)
		if err != nil {
			fmt.Printf("%d\t*\t*\t*\n", TTL)
			continue
		}

		// Parse the ICMP reply
		replyMsg, err := parseICMPResponse(response[:n])
		if err != nil {
			return fmt.Errorf("failed to parse ICMP reply: %v", err)
		}

		// Calculate round-trip time
		rtt := time.Since(startTime)

		if (replyMsg.Type == ICMPTypeTimeExceeded && replyMsg.Code == 0) ||
			(replyMsg.Type == ICMPTypeDestinationUnreachable && replyMsg.Code == 3) {
			fmt.Printf("%d: %s (%v) %v ms\n", TTL, reverseDnsLookup(from), from, rtt.Milliseconds())
		}

		if ipAddr.String() == from.String() {
			fmt.Printf("Traceroute completed on hop %d\n", TTL)
			os.Exit(0)
		}
	}

	return nil
}

func reverseDnsLookup(ip net.Addr) string {
	names, err := net.LookupAddr(ip.String())
	if err != nil {
		return ip.String()
	}
	return names[0]
}

// parseICMPResponse parses the received ICMP packet and extracts relevant information
func parseICMPResponse(packet []byte) (*Message, error) {
	if len(packet) < 4 {
		return nil, errors.New("invalid ICMP packet, too short")
	}

	return &Message{
		Type:     int(packet[0]),
		Code:     int(packet[1]),
		Checksum: int(binary.BigEndian.Uint16(packet[2:4])),
	}, nil
}
