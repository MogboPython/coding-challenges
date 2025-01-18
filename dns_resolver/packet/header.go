package packet

import (
	"encoding/binary"
	"errors"
	"fmt"
)

const RECURSION_FLAG uint16 = 1 << 8
const HEADER_SIZE_BYTES = 12

type Header struct {
	Id      uint16
	Flags   uint16
	QdCount uint16
	AnCount uint16
	NsCount uint16
	ArCount uint16
}

func NewHeader(id, flags, qdcount, ancount, nscount, arcount uint16) *Header {
	return &Header{
		Id:      id,
		Flags:   flags,
		QdCount: qdcount,
		AnCount: ancount,
		NsCount: nscount,
		ArCount: arcount,
	}
}

func (h *Header) ToBytes() []byte {
	buf := make([]byte, 12)
	binary.BigEndian.PutUint16(buf[0:2], h.Id)
	binary.BigEndian.PutUint16(buf[2:4], h.Flags)
	binary.BigEndian.PutUint16(buf[4:6], h.QdCount)
	binary.BigEndian.PutUint16(buf[6:8], h.AnCount)
	binary.BigEndian.PutUint16(buf[8:10], h.NsCount)
	binary.BigEndian.PutUint16(buf[10:12], h.ArCount)
	return buf
}

func ParseMessageHeader(headerBytes []byte) (*Header, error) {
	var header Header

	if len(headerBytes) != 12 {
		return nil, fmt.Errorf("expected header to be 12 bytes")
	}

	header.Id = binary.BigEndian.Uint16(headerBytes[0:2])
	header.Flags = binary.BigEndian.Uint16(headerBytes[2:4])
	// perform bitwise AND operation to get the last 4 bits of the flags field that represent the response code (RCODE)
	switch header.Flags & 0b1111 {
	case 1:
		return nil, errors.New("error with the query")
	case 2:
		return nil, errors.New("error with the server")
	case 3:
		return nil, errors.New("the domain doesn't exist")
	}
	header.QdCount = binary.BigEndian.Uint16(headerBytes[4:6])
	header.AnCount = binary.BigEndian.Uint16(headerBytes[6:8])
	header.NsCount = binary.BigEndian.Uint16(headerBytes[8:10])
	header.ArCount = binary.BigEndian.Uint16(headerBytes[10:12])

	return &header, nil
}
