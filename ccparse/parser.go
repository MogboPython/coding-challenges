package main

import (
	"errors"
	"fmt"
)

func check_validity(content []byte) (any, error) {
	contentStr := string(content)

	tokens := lexicalAnalysis(contentStr)
	if len(tokens) == 0 {
		return nil, errors.New("empty JSON string")
	}

	// Check if JSON starts with {
	if tokens[0].Type != LeftBrace {
		return nil, errors.New("JSON string does not start with {")
	}

	// Check if JSON ends with }
	if tokens[len(tokens)-1].Type != RightBrace {
		return nil, errors.New("JSON string does not end with }")
	}

	// Check for proper key-value structure and quote requirements
	for i := 1; i < len(tokens); i++ {
		if tokens[i].Type == Colon {
			if i-1 < 0 || tokens[i-1].Type != String {
				return nil, fmt.Errorf("invalid JSON: key must be a quoted string, found %q", tokens[i-1].Value)
			}
			if i+1 >= len(tokens) || tokens[i+1].Type == RightBrace || tokens[i+1].Type == Comma {
				return nil, errors.New("invalid JSON: colon must be followed by a value")
			}
			if i+1 >= len(tokens) || tokens[i+1].Type == Identifier {
				return nil, errors.New("invalid JSON: value unexpected")
			}
		}

		if tokens[i].Type == LeftBracket {
			if tokens[i+1].Type != String && tokens[i+1].Type == RightBrace {
				return nil, errors.New("invalid JSON: comma must be followed by a field or object end")
			}
		}

		// Check for comma placement - they should be followed by either a string (for a key)
		if tokens[i].Type == Comma && i+1 < len(tokens) {
			if tokens[i+1].Type != String && tokens[i+1].Type == RightBrace {
				return nil, errors.New("invalid JSON: comma must be followed by a field or object end")
			}
		}
	}
	// data, _ := serializeTokensToJSON(tokens)

	return tokens, nil
}

// TODO: implement parse array
// func parseArray(tokens []Token) bool {

// 	return false
// }

// func parseArray1(arr []string) bool {
// 	if arr[0] != "[" || arr[len(arr)-1] != "]" {
// 		return false
// 	}

// 	for _, item := range arr[1 : len(arr)-1] {
// 		if len(item) < 2 || item[0] != '"' || item[len(item)-1] != '"' {
// 			return false
// 		}
// 	}

// 	return true
// }
