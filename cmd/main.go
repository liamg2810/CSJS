package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"liamg2810.github.io/csjs/ast"
	"liamg2810.github.io/csjs/lexer"
)

func main() {
	contents, err := extractContents("./variables.csjs")

	if err != nil {
		panic(err)
	}

	tokens, err := lexer.Parse(*contents)

	if err != nil {
		panic(err)
	}

	built, err := ast.Build(tokens)

	if err != nil {
		panic(err)
	}

	data, err := json.MarshalIndent(built, "", " ")

	if err != nil {
		panic(err)
	}

	if err := writeContents("./output.json", data); err != nil {
		panic(err)
	}

	fmt.Printf("output to json\n")
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

func writeContents(path string, contents []byte) error {
	return os.WriteFile(path, contents, 0666)
}