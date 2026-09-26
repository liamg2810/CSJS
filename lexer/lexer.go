package lexer

import (
	"fmt"
	"regexp"
)

const input string = "str x = \"Hello World\"; int y = 17; float z = 19.2";

var floatRegex *regexp.Regexp;
var intRegex *regexp.Regexp;

func Parse() {
	tokens := []any {};

	parseStr := false
	parsingStr := ""
	lit := ""

	floatRegex = regexp.MustCompile(`^\d*\.\d*$`)
	intRegex = regexp.MustCompile(`^\d*$`)

	for _, i := range input  {
		char := string(i)

		if (parseStr) {
			if (char == "\"") {
				tokens = append(tokens, LIT{data: parsingStr, dataType: "str"})
				parseStr = false
				parsingStr = ""
				continue
			}

			parsingStr = fmt.Sprintf("%s%s", parsingStr, char)
			continue
		}

		switch char {
		case "\"":
			AddToken(&tokens, &lit)
			parseStr = true
		
		case "=":
			tokens = append(tokens, OP{operand: "assign"})
		
		case ";":
			AddToken(&tokens, &lit)
			tokens = append(tokens, EOL{})

		case " ":
			AddToken(&tokens, &lit)

		default:
			lit = fmt.Sprintf("%s%s", lit, char)
		}
	}

	AddToken(&tokens, &lit)

	for _, item := range tokens {
		fmt.Printf("%T %+v\n", item, item)
	}
}

func AddToken(tokens *[]any, lit *string) {
	if (*lit == "") {
		return
	}

	if (*lit == "str") {
		*tokens = append(*tokens, DATATYPE{name: "str"})
	} else if (*lit == "int") {
		*tokens = append(*tokens, DATATYPE{name: "int"})
	} else if (*lit == "float") {
		*tokens = append(*tokens, DATATYPE{name: "float"})
	} else if (intRegex.MatchString(*lit)) {
		*tokens = append(*tokens, LIT{data: *lit, dataType: "int"})
	} else if (floatRegex.MatchString(*lit)) {
		*tokens = append(*tokens, LIT{data: *lit, dataType: "float"})
	}

	*lit = ""
}