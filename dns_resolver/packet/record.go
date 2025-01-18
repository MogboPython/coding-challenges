package packet

import (
	"encoding/binary"
	"fmt"
)

type Record struct {
	Name     []byte
	Type     uint16
	Class    uint16
	TTL      uint32
	RdLength uint16
	RData    string
}

func ParseMessageRecord(fullMessage []byte, messageBytes []byte) (*Record, int) {
	var (
		record Record
		offset int
	)

	domainName, bytesRead := DecodeDomainName(fullMessage, messageBytes)
	record.Name = []byte(domainName)

	offset += bytesRead

	record.Type = binary.BigEndian.Uint16(messageBytes[offset : offset+2])
	record.Class = binary.BigEndian.Uint16(messageBytes[offset+2 : offset+4])
	record.TTL = binary.BigEndian.Uint32(messageBytes[offset+4 : offset+8])
	record.RdLength = binary.BigEndian.Uint16(messageBytes[offset+8 : offset+10])

	offset += 10
	lengthOfData := int(record.RdLength)

	rData := messageBytes[offset : offset+lengthOfData]

	// 	//TODO: Try to figure out for TYPE_IN and others
	switch record.Type {
	case TYPE_A:
		record.RData = fmt.Sprintf("%d.%d.%d.%d", rData[0], rData[1], rData[2], rData[3])
	case TYPE_NS:
		record.RData, _ = DecodeDomainName(fullMessage, rData)
	default:
		record.RData = string(rData)
	}

	offset += lengthOfData

	return &record, offset
}
