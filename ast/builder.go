package ast

import (
	"errors"
	"strconv"

	"liamg2810.github.io/csjs/lexer"
)

func Build(tokens []lexer.Token) ([]INode, error) {
	ret := []INode {}

	for i := 0; i < len(tokens); {
		skip := 1

		token := tokens[i]

		switch (token.Type) {
		case (lexer.DATA_TYPE):
			node, err := parseAssign(&tokens, i, &skip)

			if (err != nil) {
				return nil, err
			}

			ret = append(ret, *node)
		}

		i += skip
	}

	return ret, nil
}

func parseAssign(tokens *[]lexer.Token, startIx int, skip *int) (*Assignment, error) {
	// Read from startIx until EOL and track skip
	*skip = 1

	// Need format of DATA_TYPE, IDENT, ASSIGN, [... LIT | OP], EOL
	if (startIx > len(*tokens) - 4) {
		return nil, errors.New("parseAssign: Invalid format for assignment")
	}

	name := (*tokens)[startIx + 1]

	if (name.Type != lexer.IDENT) {
		return nil, errors.New("parseAssign: Bad token type on variable name")
	}

	assign := (*tokens)[startIx + 2]

	if (assign.Type != lexer.ASSIGN) {
		return nil, errors.New("parseAssign: Bad assignment type")
	}

	hitEOL := false

	// Ignoring BODMAS for now
	isEquation := false
	tokensToRead := []lexer.Token {}

	// Grabbing tokens to fucking figure out
	for i := startIx + 3; i < len(*tokens); i++ {
		token := (*tokens)[i]

		switch token.Type {
		case lexer.EOL:
			hitEOL = true
		case lexer.ADD, lexer.SUBTRACT, lexer.DIVIDE, lexer.MULTIPLY:
			isEquation = true
			tokensToRead = append(tokensToRead, token)
		case lexer.INT, lexer.STRING, lexer.FLOAT:
			tokensToRead = append(tokensToRead, token)
		default:
			return nil, errors.New("parseAssign: Error reading value of assignment. Bad token type")
		}

		// We are building the amount to skip
		(*skip)++

		if (hitEOL) {
			break
		}
	}

	if (!hitEOL) {
		return nil, errors.New("parseAssign: Never hit EOL")
	}

	if (!isEquation) {
		if (len(tokensToRead) != 1) {
			return nil, errors.New("parseAssign: Too many tokens on right hand side of assignment without equation")
		} else {
			var expression IExpression;

			token := tokensToRead[0]

			switch (token.Type) {
			case lexer.INT:
				i, _ := strconv.ParseInt(token.Value, 10, 0)
				expression = LiteralExpression[int64]{ Data: i }
			case lexer.FLOAT:
				i, _ := strconv.ParseFloat(token.Value, 0)
				expression = LiteralExpression[float64]{ Data: i }
			case lexer.STRING:
				expression = LiteralExpression[string]{ Data: token.Value }
			default:
				return nil, errors.New("parseAssign: Bad data type on expression")
			}

			return &Assignment{ Name: name.Value, Value: expression }, nil
		}
	}

	return nil, errors.New("parseAssign: Expressions not yet supported")
}