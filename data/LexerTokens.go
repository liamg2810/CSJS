package data

import (
	"fmt"
)

type TokenType int

const (
	IDENT TokenType = iota

	DATA_TYPE

	INT
	FLOAT
	STRING

	ADD
	SUBTRACT
	MULTIPLY
	DIVIDE

	ASSIGN

	LPAREN
	RPAREN

	EOL
	EOF

	UNARY
)

type Token struct {
	Type   TokenType
	Value  string
	Column int
	Line   int
}

func (t Token) CreateAstError(message string) *AstError {
	return &AstError{t.Line, t.Column, message}
}

func (t Token) Expect(tokenType TokenType) *AstError {
	if t.Type != tokenType {
		return &AstError{t.Line, t.Column, fmt.Sprintf("Expected (%q) but got (%q)", tokenType, t.Type)}
	}

	return nil
}

func (t TokenType) String() string {
	switch t {
	case IDENT:
		return "IDENT"
	case DATA_TYPE:
		return "DATA_TYPE"
	case INT:
		return "INT"
	case FLOAT:
		return "FLOAT"
	case STRING:
		return "STRING"
	case ADD:
		return "ADD"
	case SUBTRACT:
		return "SUBTRACT"
	case MULTIPLY:
		return "MULTIPLY"
	case DIVIDE:
		return "DIVIDE"
	case ASSIGN:
		return "ASSIGN"
	case LPAREN:
		return "LPAREN"
	case RPAREN:
		return "RPAREN"
	case EOL:
		return "EOL"
	case EOF:
		return "EOF"
	case UNARY:
		return "UNARY"
	default:
		return fmt.Sprintf("TokenType(%d)", t)
	}
}
