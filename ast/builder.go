package ast

import (
	"errors"
	"fmt"
	"slices"
	"strconv"

	"liamg2810.github.io/csjs/lexer"
)

func Build(tokens []lexer.Token) ([]INode, error) {
	ret := []INode {}

	// Keep track of variables
	variables := []VariableExpression {}

	for i := 0; i < len(tokens); {
		skip := 1

		token := tokens[i]

		switch token.Type {
		case lexer.DATA_TYPE:
			node, err := parseAssign(&tokens, i, &skip, &variables)

			if (err != nil) {
				return nil, err
			}

			ret = append(ret, *node)
		default:
			return nil, fmt.Errorf("Unknown token of type %d and value %s", token.Type, token.Value)
		}

		i += skip
	}

	return ret, nil
}

func parseAssign(tokens *[]lexer.Token, startIx int, skip *int, variables *[]VariableExpression) (*Assignment, error) {
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

	(*skip)++

	assign := (*tokens)[startIx + 2]

	if (assign.Type != lexer.ASSIGN) {
		return nil, errors.New("parseAssign: Bad assignment type")
	}

	(*skip)++

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
		case lexer.IDENT:
			if (slices.IndexFunc(*variables, func (v VariableExpression) bool { return v.Name == token.Value }) == -1) {
				return nil, fmt.Errorf("parseAssign: Error on right hand side of assignment. Unknown variable %s", token.Value)
			}

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

			dataType := Int

			switch (token.Type) {
			case lexer.INT:
				i, _ := strconv.ParseInt(token.Value, 10, 0)
				dataType = Int
				expression = LiteralExpression[int64]{ Data: i, DataType: dataType }
			case lexer.FLOAT:
				i, _ := strconv.ParseFloat(token.Value, 0)
				dataType = Float
				expression = LiteralExpression[float64]{ Data: i, DataType: dataType }
			case lexer.STRING:
				dataType = String
				expression = LiteralExpression[string]{ Data: token.Value, DataType: dataType }
			case lexer.IDENT:
				// Already confirmed the ix exists earlier in the code (this code so messy dude) make this a map[string]DataType{}
				ix := slices.IndexFunc(*variables, func (v VariableExpression) bool { return v.Name == token.Value })

				variable := (*variables)[ix]

				dataType = variable.DataType
				expression = VariableExpression{ Name: token.Value, DataType: dataType }
			default:
				return nil, errors.New("parseAssign: Bad data type on expression")
			}

			*variables = append(*variables, VariableExpression{ Name: name.Value, DataType: dataType })
			return &Assignment{ Name: name.Value, Value: expression, DataType: dataType }, nil
		}
	}

	return nil, errors.New("parseAssign: Equations not yet supported")
}