package network

import (
	"fmt"
	"net"
	"os"
	"slices"
)

type Client struct {
	serverAddress string
	port          int
}

func NewClient(address string, port int) *Client {
	return &Client{serverAddress: address, port: port}
}

func (c *Client) SendQuery(query []byte) []byte {

	// Resolve UDP address
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", c.serverAddress, c.port))
	if err != nil {
		fmt.Printf("failed to resolve address: %v", err)
		os.Exit(-1)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Printf("Dial err %v\n", err)
		os.Exit(-1)
	}
	defer conn.Close()

	if _, err = conn.Write(query); err != nil {
		fmt.Printf("Write err %v\n", err)
		os.Exit(-1)
	}

	response := make([]byte, 512) // Standard DNS message size
	lengthOfTheResponse, _, err := conn.ReadFromUDP(response)
	if err != nil {
		fmt.Printf("Read err %v\n", err)
		os.Exit(-1)
	}

	if !hasTheSameID(query, response) {
		fmt.Printf("Response doesn't have the same ID of the query q:%v, r:%v\n", query, response)
		os.Exit(-1)
	}

	return response[:lengthOfTheResponse]
}

func hasTheSameID(query, response []byte) bool {
	return slices.Equal(query[:2], response[:2])
}
