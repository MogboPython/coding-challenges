package main

import (
	"fmt"
	"os"

	"github.com/MogboPython/coding-challenges/dns_resolver/network"
	"github.com/MogboPython/coding-challenges/dns_resolver/packet"
)

type DNSPacket struct {
	header      *packet.Header
	questions   []*packet.Question
	answers     []*packet.Record
	additionals []*packet.Record
	authorities []*packet.Record
}

func main() {
	domains := os.Args[1:]
	if len(domains) < 1 {
		fmt.Println("Usage: ./dns <domain>")
		os.Exit(0)
	}

	fmt.Println("Hello World to my DNS Resolver!")
	for _, domain := range domains {
		fmt.Println(resolve(domain, packet.TYPE_A))
	}

	// dnsResponse := resolve("www.google.com", packet.TYPE_A)
	// fmt.Println(dnsResponse)
}

func resolve(domainName string, questionType uint16) string {
	nameServer := "198.41.0.4"
	for {
		// Send a DNS query to the server
		fmt.Printf("Querying %s for %s\n", nameServer, domainName)
		dnsResponse := sendQuery(nameServer, domainName, questionType)
		dnsPacket := getDnsPacketFromResponse(dnsResponse)

		if ip := getRecordData(dnsPacket.answers); ip != "" {
			return ip
		}

		if nsIp := getRecordData(dnsPacket.additionals); nsIp != "" {
			nameServer = nsIp
			continue
		}

		if nsDomain := getRecordData(dnsPacket.authorities); nsDomain != "" {
			nameServer = resolve(nsDomain, packet.TYPE_A)
		}
	}
}

func sendQuery(nameServer, domainName string, questionType uint16) []byte {
	query := packet.NewQuery(
		packet.NewHeader(22, 0, 1, 0, 0, 0),
		packet.NewQuestion(domainName, questionType, packet.CLASS_IN),
	)

	client := network.NewClient(nameServer, 53)
	return client.SendQuery(query)
}

// getDnsPacketFromResponse parses a DNS response message and returns a DNS packet
func getDnsPacketFromResponse(dnsResponse []byte) *DNSPacket {
	var (
		header      *packet.Header
		questions   []*packet.Question
		answers     []*packet.Record
		authorities []*packet.Record
		additionals []*packet.Record
	)

	headerBytes := dnsResponse[:packet.HEADER_SIZE_BYTES]
	header, err := packet.ParseMessageHeader(headerBytes)
	if err != nil {
		fmt.Printf("could not parse header: %v", err)
		os.Exit(-1)
	}

	offset := 12
	for range header.QdCount {
		question, bytesRead := packet.ParseMessageQuestion(dnsResponse, dnsResponse[offset:])
		questions = append(questions, question)
		offset += bytesRead
	}

	for range header.AnCount {
		answer, bytesRead := packet.ParseMessageRecord(dnsResponse, dnsResponse[offset:])
		answers = append(answers, answer)
		offset += bytesRead
	}

	for range header.NsCount {
		authority, bytesRead := packet.ParseMessageRecord(dnsResponse, dnsResponse[offset:])
		authorities = append(authorities, authority)
		offset += bytesRead
	}

	for range header.ArCount {
		additional, bytesRead := packet.ParseMessageRecord(dnsResponse, dnsResponse[offset:])
		additionals = append(additionals, additional)
		offset += bytesRead
	}

	return &DNSPacket{
		header:      header,
		questions:   questions,
		answers:     answers,
		authorities: authorities,
		additionals: additionals,
	}
}

func getRecordData(records []*packet.Record) string {
	for _, record := range records {
		if record.Type == packet.TYPE_A || record.Type == packet.TYPE_NS {
			return record.RData
		}
	}
	return ""
}
