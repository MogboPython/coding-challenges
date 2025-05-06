package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type Type string

type Token struct {
	Type  Type
	Value any
	// Literal string
	// Line    int
	// Start   int
	// End     int
}

var jsonNull = []rune("null")

const (
	LeftBrace    Type = "{"
	RightBrace   Type = "}"
	LeftBracket  Type = "["
	RightBracket Type = "]"
	Colon        Type = ":"
	Comma        Type = ","
	String       Type = "STRING"
	Number       Type = "NUMBER"
	Float        Type = "FLOAT"
	EOF          Type = "EOF"
	Null         Type = "NULL"
	Bool         Type = "BOOL"
	Identifier   Type = "IDENTIFIER" // For unquoted keys/values
)

func lexicalAnalysis(input string) []Token {
	tokens := []Token{}

	for i := 0; i < len(input); i++ {
		switch input[i] {
		case '{':
			tokens = append(tokens, Token{Type: LeftBrace, Value: "{"})
		case '}':
			tokens = append(tokens, Token{Type: RightBrace, Value: "}"})
		case '[':
			tokens = append(tokens, Token{Type: LeftBracket, Value: "["})
		case ']':
			tokens = append(tokens, Token{Type: RightBracket, Value: "]"})
		case ':':
			tokens = append(tokens, Token{Type: Colon, Value: ":"})
		case ',':
			tokens = append(tokens, Token{Type: Comma, Value: ","})
		case '"':
			value := ""
			i++
			for i < len(input) && input[i] != '"' {
				value += string(input[i])
				i++
			}
			tokens = append(tokens, Token{Type: String, Value: value})
		case '\'':
			value := ""
			i++
			for i < len(input) && input[i] != '\'' {
				value += string(input[i])
				i++
			}
			tokens = append(tokens, Token{Type: String, Value: value})
		default:
			if unicode.IsSpace(rune(input[i])) {
				continue
			}
			// Handle identifiers (unquoted keys)
			if unicode.IsLetter(rune(input[i])) || unicode.IsDigit(rune(input[i])) {
				identifier := ""
				for i < len(input) && (unicode.IsLetter(rune(input[i])) ||
					unicode.IsDigit(rune(input[i])) || input[i] == '_') {
					identifier += string(input[i])
					i++
				}
				value, tokenType := parseIdentifier(identifier)
				tokens = append(tokens, Token{Type: tokenType, Value: value})
				i--
			} else {
				fmt.Printf("Unexpected character: %c\n", input[i])
			}
		}
	}
	return tokens
}

func parseIdentifier(identifier string) (any, Type) {
	if identifier == "null" {
		return string(jsonNull), Null
	}

	if identifier == "true" {
		return true, Bool
	} else if identifier == "false" {
		return false, Bool
	}

	if strings.ContainsRune(identifier, '.') {
		if value, err := strconv.ParseFloat(identifier, 64); err == nil {
			return value, Float
		}
	} else {
		if value, err := strconv.ParseInt(identifier, 10, 64); err == nil {
			return int(value), Number
		}
	}

	// if value, err := strconv.ParseBool(identifier); err == nil {
	// 	return value, Bool
	// }
	// if value, err := strconv.ParseFloat(identifier, 64); err == nil {
	// 	fmt.Println(value)
	// 	return value, Float
	// }
	// if value, err := strconv.ParseInt(identifier, 10, 64); err == nil {
	// 	return int(value), Number
	// }

	return identifier, Identifier
}

// TODO: Implement this function
// map[string]any
// convert tokens to formatted JSON with indentation
// func SerializeTokensToJSON(tokens []Token) (string, error) {
// 	if len(tokens) == 0 {
// 		return "", nil
// 	}

// 	var jsonData = map[string]any{}
// 	i := 0

// 	// Skip the opening brace
// 	if tokens[0].Type == LeftBrace {
// 		i++
// 	}

// 	for i < len(tokens) {

// 		// Skip closing brace and commas
// 		if tokens[i].Type == RightBrace || tokens[i].Type == Comma {
// 			i++
// 			continue
// 		}

// 		// Check if we have a key (either String or Identifier)
// 		if (tokens[i].Type == String || tokens[i].Type == Identifier) && i+1 < len(tokens) && tokens[i+1].Type == Colon {
// 			key := tokens[i].Value
// 			i += 2
// 			if i >= len(tokens) {
// 				return "", fmt.Errorf("unexpected end of tokens after colon")
// 			}

// 			value := tokens[i].Value
// 			jsonData[key] = value
// 		}
// 	}
// 	j, err := json.Marshal(jsonData)
// 	if err != nil {
// 		return "", fmt.Errorf("error marshaling JSON: %v", err)
// 	}

// 	return string(j), nil
// }

// Helper function to write indentation spaces
// func writeIndentation(builder *strings.Builder, level int) {
// 	for i := 0; i < level*2; i++ {
// 		builder.WriteString(" ")
// 	}
// }
