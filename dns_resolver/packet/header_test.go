package packet

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeader(t *testing.T) {
	t.Run("Should encode an header into bytes", func(t *testing.T) {
		header := NewHeader(22, RECURSION_FLAG, 1, 0, 0, 0)

		encodedHeader := header.ToBytes()

		expected, err := hex.DecodeString("0016010000010000000000000") // Would return a odd-length error
		assert.NotNil(t, err)
		assert.Equal(t, expected, encodedHeader)
	})

	t.Run("Should decode bytes into an header", func(t *testing.T) {
		headerBytes, _ := hex.DecodeString("0016808000010002000000000")

		header, err := ParseMessageHeader(headerBytes)

		assert.Nil(t, err)
		assert.Equal(t, uint16(22), header.Id)
		assert.Equal(t, uint16(1), (header.Flags>>15)&0x1) // check QR bit is 1 (response)
		assert.Equal(t, uint16(0), header.Flags&0b1111)    // check RCODE is 0 (no error)
		assert.Equal(t, uint16(1), header.QdCount)
		assert.Equal(t, uint16(2), header.AnCount)
		assert.Equal(t, uint16(0), header.NsCount)
		assert.Equal(t, uint16(0), header.ArCount)
	})

	// Copied from https://github.com/dlion/unnije.git
	t.Run("Should return an error if the header flags contains a query error", func(t *testing.T) {
		headerBytes, _ := hex.DecodeString("001680810001000200000000")

		header, err := ParseMessageHeader(headerBytes)

		assert.Nil(t, header)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "error with the query")
	})

	t.Run("Should return an error if the header flags contains a server error", func(t *testing.T) {
		headerBytes, _ := hex.DecodeString("001680820001000200000000")

		header, err := ParseMessageHeader(headerBytes)

		assert.Nil(t, header)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "error with the server")
	})

	t.Run("Should return an error if the header flags contains a domain not exist error", func(t *testing.T) {
		headerBytes, _ := hex.DecodeString("001680830001000200000000")

		header, err := ParseMessageHeader(headerBytes)

		assert.Nil(t, header)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "the domain doesn't exist")
	})

}
