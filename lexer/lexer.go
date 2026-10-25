package lexer

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"liamg2810.github.io/csjs/data"
)

var floatRegex *regexp.Regexp
var intRegex *regexp.Regexp

func Parse(input string) (*data.FileContext, error) {
	tokens := []data.Token{}
	lines := strings.Split(input, "\n")
	columnCounter := 0
	lineCounter := 0

	parseStr := false
	parsingStr := ""
	lit := ""
	skip := false

	floatRegex = regexp.MustCompile(`^\d+\.\d+$`)
	intRegex = regexp.MustCompile(`^\d*$`)

	for ix, char := range input {
		columnCounter++

		if skip {
			skip = false
			continue
		}

		if parseStr {
			if char == '\\' {
				escaped, err := escapeCharacter(input, ix)

				if err != nil {
					return nil, err
				}

				skip = true
				parsingStr = parsingStr + escaped

				continue
			}

			if char == '"' {
				tokens = append(tokens, data.Token{Type: data.STRING, Value: parsingStr})
				parseStr = false
				parsingStr = ""
				continue
			}

			parsingStr = fmt.Sprintf("%s%c", parsingStr, char)
			continue
		}

		switch char {
		case '\n':
			lineCounter += 1
			columnCounter = 0
			continue

		case '"':
			AddToken(&tokens, &lit, nil, columnCounter, lineCounter)
			parseStr = true

		case '=':
			AddToken(&tokens, &lit, &data.Token{Type: data.ASSIGN, Value: "=", Column: columnCounter, Line: lineCounter}, columnCounter, lineCounter)

		case ';':
			AddToken(&tokens, &lit, &data.Token{Type: data.EOL, Value: ";", Column: columnCounter, Line: lineCounter}, columnCounter, lineCounter)

		case ' ':
			AddToken(&tokens, &lit, nil, columnCounter, lineCounter)

		case '+':
			AddToken(&tokens, &lit, &data.Token{Type: data.ADD, Value: "+", Column: columnCounter, Line: lineCounter}, columnCounter, lineCounter)

		case '-':
			// lit is something so there is guaranteed to be something before this sign that would be a variable
			if lit != "" || ix == 0 {
				AddToken(&tokens, &lit, &data.Token{Type: data.SUBTRACT, Value: "-", Column: columnCounter, Line: lineCounter}, columnCounter, lineCounter)
				continue
			}

			prevTokenType := tokens[len(tokens)-1].Type

			// is the thing before something that would be subtractable
			if prevTokenType == data.INT || prevTokenType == data.FLOAT || prevTokenType == data.IDENT || prevTokenType == data.RPAREN {
				AddToken(&tokens, &lit, &data.Token{Type: data.SUBTRACT, Value: "-", Column: columnCounter, Line: lineCounter}, columnCounter, lineCounter)
				continue
			}

			// fall back to unary
			AddToken(&tokens, &lit, &data.Token{Type: data.UNARY, Value: "-", Column: columnCounter, Line: lineCounter}, columnCounter, lineCounter)

		case '*':
			AddToken(&tokens, &lit, &data.Token{Type: data.MULTIPLY, Value: "*", Column: columnCounter, Line: lineCounter}, columnCounter, lineCounter)

		case '/':
			AddToken(&tokens, &lit, &data.Token{Type: data.DIVIDE, Value: "/", Column: columnCounter, Line: lineCounter}, columnCounter, lineCounter)

		case '(':
			AddToken(&tokens, &lit, &data.Token{Type: data.LPAREN, Value: "(", Column: columnCounter, Line: lineCounter}, columnCounter, lineCounter)

		case ')':
			AddToken(&tokens, &lit, &data.Token{Type: data.RPAREN, Value: ")", Column: columnCounter, Line: lineCounter}, columnCounter, lineCounter)

		default:
			lit = fmt.Sprintf("%s%c", lit, char)
		}
	}

	AddToken(&tokens, &lit, nil, columnCounter, len(lines)-1)

	context := data.FileContext{Lines: lines, Tokens: tokens}

	return &context, nil
}

func escapeCharacter(input string, index int) (string, error) {
	if index+1 >= len(input) {
		return "", errors.New("Invalid escape charater at end of input")
	}

	char := input[index+1]

	switch char {
	case '\\':
		return "\\", nil
	case 'n':
		return "\n", nil
	case '"':
		return "\"", nil
	default:
		return "", fmt.Errorf("Unsupported escape code: %c", char)
	}
}

func AddToken(tokens *[]data.Token, lit *string, extra *data.Token, column int, line int) {

	if *lit != "" {
		if *lit == "str" {
			*tokens = append(*tokens, data.Token{Type: data.DATA_TYPE, Value: "str", Line: line, Column: column})
		} else if *lit == "int" {
			*tokens = append(*tokens, data.Token{Type: data.DATA_TYPE, Value: "int", Line: line, Column: column})
		} else if *lit == "float" {
			*tokens = append(*tokens, data.Token{Type: data.DATA_TYPE, Value: "float", Line: line, Column: column})
		} else if intRegex.MatchString(*lit) {
			*tokens = append(*tokens, data.Token{Type: data.INT, Value: *lit, Line: line, Column: column})
		} else if floatRegex.MatchString(*lit) {
			*tokens = append(*tokens, data.Token{Type: data.FLOAT, Value: *lit, Line: line, Column: column})
		} else {
			*tokens = append(*tokens, data.Token{Type: data.IDENT, Value: *lit, Line: line, Column: column})
		}
	}

	if extra != nil {
		*tokens = append(*tokens, *extra)
	}

	*lit = ""
}
