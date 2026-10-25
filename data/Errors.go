package data

import (
	"fmt"
	"strings"
)

type AstError struct {
	Line    int
	Column  int
	Message string
}

func (a AstError) Output(context FileContext) {

	fmt.Printf("Error: %s\n", a.Message)
	fmt.Printf(" --> %s:%d:%d\n", context.FileName, a.Line+1, a.Column)

	prefixSpacing := strings.Repeat(" ", len(fmt.Sprint(a.Line)))
	fmt.Printf("%s |\n", prefixSpacing)
	fmt.Printf("%d | %s\n", a.Line+1, context.Lines[a.Line])
	fmt.Printf("%s | %s^\n", prefixSpacing, strings.Repeat(" ", a.Column-1))
}
