package packet

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResponse(t *testing.T) {
	t.Run("Should create a record from a response", func(t *testing.T) {
		response, _ := hex.DecodeString("00168080000100020000000003646e7306676f6f676c6503636f6d0000010001c00c0001000100000214000408080808c00c0001000100000214000408080404")
		const RECORD_STARTING_POINT = 32
		message := response[RECORD_STARTING_POINT:]
		record, bytesRead := ParseMessageRecord(response, message)

		assert.NotEmpty(t, record)
		assert.Equal(t, []byte("dns.google.com"), record.Name)
		assert.Equal(t, TYPE_A, record.Type)
		assert.Equal(t, CLASS_IN, record.Class)
		assert.Greater(t, record.TTL, uint32(0))
		assert.Greater(t, record.RdLength, uint16(0))
		assert.Equal(t, "8.8.8.8", record.RData)

		offset := RECORD_STARTING_POINT + bytesRead
		message = response[offset:]
		record, _ = ParseMessageRecord(response, message)

		assert.NotEmpty(t, record)
		assert.Equal(t, TYPE_A, record.Type)
		assert.Equal(t, CLASS_IN, record.Class)
		assert.Greater(t, record.TTL, uint32(0))
		assert.Greater(t, record.RdLength, uint16(0))
		assert.Equal(t, "8.8.4.4", record.RData)

	})
}
