package lexer

type IDENT struct {
	name string
}

type LIT struct {
	data string
	dataType string
}

type OP struct {
	operand string
}

type DATATYPE struct {
	name string
}

type QUOTE struct {}

type EOL struct {}