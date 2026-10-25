package main

import (
	"bufio"
	"fmt"
	"os"

	"liamg2810.github.io/csjs/ast"
	"liamg2810.github.io/csjs/lexer"
)

const fileName string = "test.csjs"

func main() {
	contents, err := extractContents("./" + fileName)

	if err != nil {
		panic(err)
	}

	context, err := lexer.Parse(*contents)

	context.FileName = fileName

	if err != nil {
		panic(err)
	}

	built, astError := ast.Build(context)

	if astError != nil {
		astError.Output(*context)
		return
	}

	data := ast.Debug(built)

	if err := writeContents("./output.txt", []byte(data)); err != nil {
		panic(err)
	}

	fmt.Printf("output to output.txt\n")
}

func extractContents(path string) (*string, error) {
	file, err := os.Open(path)

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
