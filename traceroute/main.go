package main

import (
	"fmt"
	"net"
	"os"

	"github.com/MogboPython/coding-challenges/traceroute/tracer"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ./traceroute-clone <destination>")
		os.Exit(0)
	}
	destination := os.Args[1]

	ipAddr, err := net.ResolveIPAddr("ip4", destination)
	if err != nil {
		fmt.Printf("Error: failed to resolve IP address: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("traceroute to %s (%s), %d hops max, %d byte packets\n", destination, ipAddr, tracer.MAX_HOPS, tracer.PACKET_SIZE)

	if tracer.Trace(*ipAddr) != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

}
