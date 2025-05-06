package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_lexicalAnalysis(t *testing.T) {
	t.Run("Empty braces", func(t *testing.T) {
		input := "{}"
		expectedTokens := []Token{
			{Type: LeftBrace, Value: "{"},
			{Type: RightBrace, Value: "}"},
		}

		tokens := lexicalAnalysis(input)

		// Check that we got the expected number of tokens
		assert.Equal(t, len(expectedTokens), len(tokens), "Should have exactly 2 tokens")

		// Check each token matches expected type and value
		for i, token := range tokens {
			assert.Equal(t, expectedTokens[i].Type, token.Type, "Token %d should have correct type", i)
			assert.Equal(t, expectedTokens[i].Value, token.Value, "Token %d should have correct value", i)
		}
	})

	// t.Run("With whitespace", func(t *testing.T) {
	// 	input := " { } "
	// 	expectedTokens := []Token{
	// 		{Type: LeftBrace, Value: "{"},
	// 		{Type: RightBrace, Value: "}"},
	// 	}

	// 	tokens := lexicalAnalysis(input)

	// 	assert.Equal(t, expectedTokens, tokens, "Whitespace should be ignored")
	// })

	t.Run("JSON object containing string keys and string values", func(t *testing.T) {
		input := "{\"key\": \"value\"}"
		expectedTokens := []Token{
			{Type: LeftBrace, Value: "{"},
			{Type: String, Value: "key"},
			{Type: Colon, Value: ":"},
			{Type: String, Value: "value"},
			{Type: RightBrace, Value: "}"},
		}

		tokens := lexicalAnalysis(input)

		assert.Equal(t, expectedTokens, tokens)
	})

	t.Run("JSON object containing null, integer, and bool values", func(t *testing.T) {
		input := "{\"key1\": true, \"key2\": false, \"key3\": null, \"key4\": \"value\", \"key5\": 101}"
		expectedTokens := []Token{
			{Type: LeftBrace, Value: "{"},
			{Type: String, Value: "key1"},
			{Type: Colon, Value: ":"},
			{Type: Bool, Value: true},
			{Type: Comma, Value: ","},
			{Type: String, Value: "key2"},
			{Type: Colon, Value: ":"},
			{Type: Bool, Value: false},
			{Type: Comma, Value: ","},
			{Type: String, Value: "key3"},
			{Type: Colon, Value: ":"},
			{Type: Null, Value: string(jsonNull)},
			{Type: Comma, Value: ","},
			{Type: String, Value: "key4"},
			{Type: Colon, Value: ":"},
			{Type: String, Value: "value"},
			{Type: Comma, Value: ","},
			{Type: String, Value: "key5"},
			{Type: Colon, Value: ":"},
			{Type: Number, Value: 101},
			{Type: RightBrace, Value: "}"},
		}

		tokens := lexicalAnalysis(input)

		assert.Equal(t, expectedTokens, tokens)
	})

	t.Run("Empty string", func(t *testing.T) {
		input := ""
		tokens := lexicalAnalysis(input)

		assert.Empty(t, tokens, "Empty string should produce no tokens")
	})
}
