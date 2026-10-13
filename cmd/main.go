package main

import (
	"bufio"
	"fmt"
	"os"

	"liamg2810.github.io/csjs/ast"
	"liamg2810.github.io/csjs/lexer"
)

func main() {
	contents, err := extractContents("./exampleScript.csjs")

	if err != nil {
		panic(err)
	}

	tokens, err := lexer.Parse(*contents)

	fmt.Printf("%+v\n", tokens)

	if err != nil {
		panic(err)
	}

	built, err := ast.Build(tokens)

	if err != nil {
		panic(err)
	}

	for _, ast := range built {
		fmt.Printf("%T: %+v\n", ast, ast)
	}

	fmt.Printf("%+v\n", built)
}

func extractContents(path string) (*string, error) {
	file, err := os.Open(path);
	
	if err != nil {
		return nil, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	output := ""

	for scanner.Scan() {
		line := scanner.Text()

		output = fmt.Sprintf("%s%s\n", output, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &output, nil
}