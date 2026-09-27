package main

import (
	"fmt"

	"liamg2810.github.io/csjs/ast"
	"liamg2810.github.io/csjs/lexer"
)

func main() {
	tokens := lexer.Parse()

	built, err := ast.Build(tokens)

	if (err != nil) {
		panic(err)
	}

	for _, ast := range built {
		fmt.Printf("%T: %+v\n", ast, ast)
	}
}