package lexer

import (
	"errors"
	"fmt"
	"regexp"
)

var floatRegex *regexp.Regexp;
var intRegex *regexp.Regexp;

func Parse(input string) ([]Token, error) {
	tokens := []Token {};

	parseStr := false
	parsingStr := ""
	lit := ""
	skip := false

	floatRegex = regexp.MustCompile(`^-?\d+\.\d+$`)
	intRegex = regexp.MustCompile(`^-?\d*$`)

	for ix, char := range input  {
		if (skip) {
			skip = false
			continue
		}

		if (parseStr) {
			if (char == '\\') {
				escaped, err := escapeCharacter(input, ix);

				if (err != nil) {
					return nil, err
				}

				skip = true
				parsingStr = parsingStr + escaped

				continue
			}

			if (char == '"') {
				tokens = append(tokens, Token{Type: STRING, Value: parsingStr})
				parseStr = false
				parsingStr = ""
				continue
			}

			parsingStr = fmt.Sprintf("%s%c", parsingStr, char)
			continue
		}

		
		switch char {
		case '\n':
			continue

		case '"':
			AddToken(&tokens, &lit, nil)
			parseStr = true
		
		case '=':
			AddToken(&tokens, &lit, &Token{Type: ASSIGN, Value: "="})
		
		case ';':
			AddToken(&tokens, &lit, &Token{Type: EOL, Value: ";"})

		case ' ':
			AddToken(&tokens, &lit, nil)

		case '+':
			AddToken(&tokens, &lit, &Token{Type: ADD, Value: "+"})

		case '-':
			if lit == "" && intRegex.MatchString(string(input[ix+1])) {
					lit = fmt.Sprintf("%s%c", lit, char)
			} else {
				AddToken(&tokens, &lit, &Token{Type: SUBTRACT, Value: "-"})
			}

		case '*':
			AddToken(&tokens, &lit, &Token{Type: MULTIPLY, Value: "*"})

		case '/':
			AddToken(&tokens, &lit, &Token{Type: DIVIDE, Value: "/"})

		default:
			lit = fmt.Sprintf("%s%c", lit, char)
		}
	}

	AddToken(&tokens, &lit, nil)

	return tokens, nil
}

func escapeCharacter(input string, index int) (string, error) {
	if (index + 1 >= len(input)) {
		return "", errors.New("Invalid escape charater at end of input")
	}

	char := input[index + 1]

	switch (char) {
		case '\\':
			return "\\", nil;
		case 'n':
			return "\n", nil;
		case '"':
			return "\"", nil;
		default:
			return "", fmt.Errorf("Unsupported escape code: %c", char)
	}
}


func AddToken(tokens *[]Token, lit *string, extra *Token) {

	if *lit != "" {
		if *lit == "str" {
			*tokens = append(*tokens, Token{Type: DATA_TYPE, Value: "str"})
		} else if *lit == "int" {
			*tokens = append(*tokens, Token{Type: DATA_TYPE, Value: "int"})
		} else if *lit == "float" {
			*tokens = append(*tokens, Token{Type: DATA_TYPE, Value: "float"})
		} else if intRegex.MatchString(*lit) {
			*tokens = append(*tokens, Token{Type: INT, Value: *lit})
		} else if floatRegex.MatchString(*lit) {
			*tokens = append(*tokens, Token{Type: FLOAT, Value: *lit})
		} else {
			*tokens = append(*tokens, Token{Type: IDENT, Value: *lit})
		}
	}

	if extra != nil {
		*tokens = append(*tokens, *extra)
	}

	*lit = ""
}