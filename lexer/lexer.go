package lexer

import (
	"fmt"
	"regexp"
)

var floatRegex *regexp.Regexp;
var intRegex *regexp.Regexp;

func Parse(input string) []Token {
	tokens := []Token {};

	parseStr := false
	parsingStr := ""
	lit := ""

	floatRegex = regexp.MustCompile(`^\d*\.\d*$`)
	intRegex = regexp.MustCompile(`^\d*$`)

	for _, i := range input  {
		char := string(i)

		if (parseStr) {
			if (char == "\"") {
				tokens = append(tokens, Token{Type: STRING, Value: parsingStr})
				parseStr = false
				parsingStr = ""
				continue
			}

			parsingStr = fmt.Sprintf("%s%s", parsingStr, char)
			continue
		}

		switch char {
		case "\n":
			continue

		case "\"":
			AddToken(&tokens, &lit)
			parseStr = true
		
		case "=":
			tokens = append(tokens, Token{Type: ASSIGN, Value: "="})
		
		case ";":
			AddToken(&tokens, &lit)
			tokens = append(tokens, Token{Type: EOL, Value: ";"})

		case " ":
			AddToken(&tokens, &lit)

		default:
			lit = fmt.Sprintf("%s%s", lit, char)
		}
	}

	AddToken(&tokens, &lit)

	return tokens
}

func AddToken(tokens *[]Token, lit *string) {
	if (*lit == "") {
		return
	}

	if (*lit == "str") {
		*tokens = append(*tokens, Token{Type: DATA_TYPE, Value: "str"})
	} else if (*lit == "int") {
		*tokens = append(*tokens, Token{Type: DATA_TYPE, Value: "int"})
	} else if (*lit == "float") {
		*tokens = append(*tokens, Token{Type: DATA_TYPE, Value: "float"})
	} else if (intRegex.MatchString(*lit)) {
		*tokens = append(*tokens, Token{Type: INT, Value: *lit})
	} else if (floatRegex.MatchString(*lit)) {
		*tokens = append(*tokens, Token{Type: FLOAT, Value: *lit})
	} else {
		*tokens = append(*tokens, Token{Type: IDENT, Value: *lit})
	}

	*lit = ""
}