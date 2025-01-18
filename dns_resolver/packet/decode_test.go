package packet

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecoding(t *testing.T) {
	t.Run("Should read domain names from a question", func(t *testing.T) {
		response, _ := hex.DecodeString("00168080000100020000000003646e7306676f6f676c6503636f6d00") //truncated for readability
		message := response[QUESTION_STARTING_POINT:]
		dnsName, _ := DecodeDomainName(response, message)

		assert.NotEmpty(t, dnsName)
		assert.Equal(t, "dns.google.com", dnsName)
	})
}
