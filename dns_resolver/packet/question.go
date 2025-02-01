package packet

import (
	"bytes"
	"encoding/binary"
)

const TYPE_A uint16 = 1   // Question TYPE field for host address
const TYPE_NS uint16 = 2  // Question TYPE field for authoritative name server
const CLASS_IN uint16 = 1 // Question CLASS field for internet
const QUESTION_STARTING_POINT = 12

type Question struct {
	QName  []byte
	QType  uint16
	QClass uint16
}

func NewQuestion(qname string, qtype, qclass uint16) *Question {
	return &Question{
		QName:  encodeDnsName(qname),
		QType:  qtype,
		QClass: qclass,
	}
}

func encodeDnsName(host_name string) []byte {
	var encoded []byte

	// Remove trailing dot if present
	if host_name[len(host_name)-1] == '.' {
		host_name = host_name[:len(host_name)-1]
	}

	parts := bytes.Split([]byte(host_name), []byte{'.'})

	for _, part := range parts {
		encoded = append(encoded, byte(len(part)))
		encoded = append(encoded, part...)
	}
	return append(encoded, 0)
}

func (q *Question) ToBytes() []byte {
	buf := new(bytes.Buffer)
	buf.Write(q.QName)
	binary.Write(buf, binary.BigEndian, q.QType)
	binary.Write(buf, binary.BigEndian, q.QClass)

	return buf.Bytes()
}

func ParseMessageQuestion(fullMessage []byte, messageBytes []byte) (*Question, int) {
	var (
		q      Question
		offset int
	)

	domainName, bytesRead := DecodeDomainName(fullMessage, messageBytes)
	offset += bytesRead

	q.QName = []byte(domainName)
	q.QType = binary.BigEndian.Uint16(messageBytes[offset : offset+2])
	q.QClass = binary.BigEndian.Uint16(messageBytes[offset+2 : offset+4])

	bytesRead += 4

	return &q, bytesRead
}
