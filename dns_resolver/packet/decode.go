package packet

import (
	"encoding/binary"
	"strings"
)

func DecodeDomainName(fullMessage []byte, data []byte) (string, int) {
	var nameParts []string
	var position int = 0

	for position < len(data) {
		length := int(data[position])

		if length == 0 {
			position++
			break
		}

		// checks if it is a pointer related to compression
		if length&0xC0 == 0xC0 {
			pointer := int(binary.BigEndian.Uint16(data[position:position+2]) & 0x3FFF)
			position += 2
			part, _ := DecodeDomainName(fullMessage, fullMessage[pointer:])
			nameParts = append(nameParts, part)
			break
		}
		position++
		nameParts = append(nameParts, string(data[position:position+length]))
		position += length
	}
	return strings.Join(nameParts, "."), position
}
