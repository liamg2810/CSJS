package lexer

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
)

type Token struct {
	Type TokenType
	Value string
}