package ast

import (
	"fmt"
	"strconv"

	"liamg2810.github.io/csjs/data"
)

func Build(context *data.FileContext) ([]data.INode, *data.AstError) {
	ret := []data.INode{}

	tokens := context.Tokens

	// Keep track of variables
	variables := map[string]data.DataType{}

	for i := 0; i < len(tokens); {
		token := tokens[i]

		switch token.Type {
		case data.DATA_TYPE:
			node, nextIx, err := parseAssign(tokens, i, variables)

			if err != nil {
				return nil, err
			}

			ret = append(ret, *node)

			i = nextIx + 1
		default:
			return nil, token.CreateAstError(fmt.Sprintf("Unknown token of type %q and value %s", token.Type, token.Value))
		}
	}

	return ret, nil
}

func parseAssign(tokens []data.Token, startIx int, variables map[string]data.DataType) (*data.Assignment, int, *data.AstError) {
	// Need format of DATA_TYPE, IDENT, ASSIGN, [... LIT | OP], EOL
	if startIx > len(tokens)-4 {
		return nil, startIx, tokens[startIx].CreateAstError("Invalid format for assignment")
	}

	name := tokens[startIx+1]

	if name.Type != data.IDENT {
		return nil, 0, tokens[startIx].CreateAstError("Invalid variable name")
	}

	assign := tokens[startIx+2]

	if assign.Type != data.ASSIGN {
		return nil, 0, assign.CreateAstError(fmt.Sprintf("Expected \"=\" got %q", assign.Value))
	}

	eolPos := -1

	for i := startIx + 3; i < len(tokens); i++ {
		if tokens[i].Type == data.EOL {
			eolPos = i
			break
		}
	}

	if eolPos == -1 {
		lastToken := tokens[len(tokens)-1]
		return nil, 0, lastToken.CreateAstError(fmt.Sprintf("Expected \";\" got %q", lastToken.Value))
	}

	slice := tokens[startIx+3 : eolPos]

	if len(slice) == 0 {
		return nil, 0, assign.CreateAstError("Expected expression after \"=\"")
	}

	binaryExpression, err := buildExpression(slice, variables)

	if err != nil {
		return nil, 0, err
	}

	variables[name.Value] = binaryExpression.GetDataType()

	return &data.Assignment{Name: name.Value, Value: binaryExpression, DataType: binaryExpression.GetDataType()}, eolPos, nil
}

func buildExpression(equation []data.Token, variables map[string]data.DataType) (data.IExpression, *data.AstError) {
	// first find all symbol types in the equation

	// Multiple + Divide
	md := []int{}
	// Add + Subtract
	as := []int{}

	// Parenthesis depth
	depth := 0

	for i, token := range equation {
		switch token.Type {
		case data.MULTIPLY, data.DIVIDE:
			if depth == 0 {
				md = append(md, i)
			}
		case data.ADD, data.SUBTRACT:
			if depth == 0 {
				as = append(as, i)
			}
		case data.LPAREN:
			depth += 1
		case data.RPAREN:
			depth -= 1
		default:
			continue
		}

		if depth < 0 {
			return nil, token.CreateAstError(fmt.Sprintf("Invalid parenthesis level (%d)", depth))
		}
	}

	if depth > 0 {
		return nil, equation[len(equation)-1].CreateAstError(fmt.Sprintf("Invalid parenthesis level (%d)", depth))
	}

	// then recursively solve equation (excluding parenthesis and exponents)
	i := 0

	if len(as) > 0 {
		i = as[len(as)-1]
	} else if len(md) > 0 {
		i = md[len(md)-1]
	} else {
		// Literal expression checks before erroring
		// Single number?
		if len(equation) == 1 {
			token := (equation)[0]

			if token.Type == data.IDENT {
				dataType, ok := variables[token.Value]

				if !ok {
					return nil, token.CreateAstError(fmt.Sprintf("Variable %q not defined.", token.Value))
				}

				return data.VariableExpression{Name: token.Value, DataType: dataType}, nil
			}

			lit, err := getLiteralExpression(&token)

			if err != nil {
				return nil, err
			}

			return lit, nil
		}

		if (equation)[0].Type == data.LPAREN && (equation)[len(equation)-1].Type == data.RPAREN {
			if len(equation) == 2 {
				return nil, equation[0].CreateAstError("Empty parentheses")
			}
			
			slice := equation[1 : len(equation)-1]
			exp, err := buildExpression(slice, variables)

			if err != nil {
				return nil, err
			}

			return exp, nil
		}

		return nil, equation[0].CreateAstError("Equation with no operations")
	}

	if i == 0 || i == len(equation)-1 {
		return nil, equation[i].CreateAstError("Operation out of bounds")
	}

	token := equation[i]

	operand, err := getOperandFromToken(&token)

	if err != nil {
		return nil, err
	}

	left, err := buildExpression(equation[:i], variables)

	if err != nil {
		return nil, err
	}

	right, err := buildExpression(equation[i+1:], variables)

	if err != nil {
		return nil, err
	}

	dataType := data.Int

	// Assumes that the expression only contains numbers
	if left.GetDataType() == data.Float || right.GetDataType() == data.Float {
		dataType = data.Float
	}

	expression := data.BinaryExpression{Operand: operand, Left: left, Right: right, DataType: dataType}

	return expression, nil
}

func getOperandFromToken(token *data.Token) (data.BinaryOperator, *data.AstError) {
	switch (*token).Type {
	case data.ADD:
		return data.Add, nil
	case data.SUBTRACT:
		return data.Subtract, nil
	case data.DIVIDE:
		return data.Divide, nil
	case data.MULTIPLY:
		return data.Multiply, nil
	default:
		return data.Add, token.CreateAstError(fmt.Sprintf("Invalid operand type %+v", token))
	}
}

func getLiteralExpression(token *data.Token) (data.IExpression, *data.AstError) {
	switch token.Type {
	case data.INT:
		i, err := strconv.ParseInt(token.Value, 10, 0)

		if err != nil {
			return nil, token.CreateAstError(fmt.Sprintf("Invalid integer type %q", token.Value))
		}

		expression := data.LiteralExpression[int64]{Data: i, DataType: data.Int}
		return expression, nil
	case data.FLOAT:
		i, err := strconv.ParseFloat(token.Value, 0)

		if err != nil {
			return nil, token.CreateAstError(fmt.Sprintf("Invalid float type %q", token.Value))
		}

		expression := data.LiteralExpression[float64]{Data: i, DataType: data.Float}
		return expression, nil
	case data.STRING:
		expression := data.LiteralExpression[string]{Data: token.Value, DataType: data.String}
		return expression, nil
	default:
		return nil, token.CreateAstError(fmt.Sprintf("Invalid data type: %q", token.Type))
	}
}
