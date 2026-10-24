package ast

import (
	"errors"
	"fmt"
	"strconv"

	"liamg2810.github.io/csjs/lexer"
)

func Build(tokens []lexer.Token) ([]INode, error) {
	ret := []INode {}

	// Keep track of variables
	variables := map[string]DataType{}

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

func parseAssign(tokens *[]lexer.Token, startIx int, skip *int, variables *map[string]DataType) (*Assignment, error) {
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
	eolPos := 0

	for i := startIx + 3; i < len(*tokens); i++ {
		token := (*tokens)[i]

		switch token.Type {
		case lexer.EOL:
			hitEOL = true
			eolPos = i
		case lexer.LPAREN, lexer.RPAREN, lexer.INT, lexer.STRING, lexer.FLOAT, lexer.IDENT, lexer.MULTIPLY, lexer.ADD, lexer.DIVIDE, lexer.SUBTRACT:
			// Catch all valid token types, for validation
		default:
			return nil, fmt.Errorf("parseAssign: Error reading value of assignment. Bad token type: %+v", (*tokens)[i])
		}

		(*skip)++

		if (hitEOL) {
			break
		}
	}

	if (!hitEOL) {
		return nil, errors.New("parseAssign: Never hit EOL")
	}

	slice := (*tokens)[startIx + 3:eolPos]
	
	unaryFix(&slice)

	binaryExpression, err := buildExpression(&slice, variables)

	if (err != nil) {
		return nil, err
	}

	(*variables)[name.Value] = binaryExpression.GetDataType()

	return &Assignment{ Name: name.Value, Value: binaryExpression, DataType: binaryExpression.GetDataType() }, nil
}

// used to be called buildBinaryExpression
// it will now handle any type of expression
// TODO: include variable map in the input and remove all expression code from parseAssign
func buildExpression(equation *[]lexer.Token, variables *map[string]DataType) (IExpression, error) {
	// first find all symbol types in the equation

	// Multiple + Divide
	md := []int{}
	// Add + Subtract
	as := []int{}

	// Parenthesis depth
	depth := 0

	for i, token := range *equation {
		switch (token.Type) {
		case lexer.MULTIPLY, lexer.DIVIDE:
			if (depth == 0) {
				md = append(md, i)
			}
		case lexer.ADD, lexer.SUBTRACT:
			if (depth == 0) {
				as = append(as, i)
			}
		case lexer.LPAREN:
			depth += 1
		case lexer.RPAREN:
			depth -= 1
		default:
			continue
		}

		if (depth < 0) {
			return nil, fmt.Errorf("buildBinaryExpression: Invalid parenthesis")
		}
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
		if len(*equation) == 1 {
			token := (*equation)[0]
			lit, err := getLiteralExpression(&token)

			if err !=  nil {
				return nil, err
			}

			return lit, nil
		}

		if (*equation)[0].Type == lexer.LPAREN && (*equation)[len(*equation) - 1].Type == lexer.RPAREN {
			slice := (*equation)[1:len(*equation)-1]
			exp, err := buildExpression(&slice, variables)

			if err != nil {
				return nil, err
			}

			return exp, nil
		}

		return nil, fmt.Errorf("buildBinaryExpression: Malformed equation")
	}

	if (i == 0 || i == len(*equation) - 1) {
		return nil, fmt.Errorf("buildBinaryExpression: Malformed equation, as operation out of bounds")
	}

	token := (*equation)[i]

	operand, err := getOperandFromToken(&token)

	if err != nil {
		return nil, err
	}

	var left IExpression
	var right IExpression

	if (i > 1) {
		slice := (*equation)[0:i]

		parenCheck := (*equation)[i - 1]

		// Left side of operation will be RPAREN if it is a PAREN
		if (parenCheck.Type == lexer.LPAREN) {
			return nil, fmt.Errorf("buildBinaryExpression: Wrong parenthesis on LHS of operand")
		}

		if (parenCheck.Type == lexer.RPAREN) {
			start, err := getLParen(equation, i - 1)
			
			if err != nil {
				return nil, err
			}

			slice = (*equation)[start+1:i-1]
		}

		leftExp, err := buildExpression(&slice, variables)

		if err != nil {
			return nil, err
		}

		left = leftExp
	} else {
		token := (*equation)[i - 1]

		if (token.Type == lexer.IDENT) {
			dataType, ok := (*variables)[token.Value]

			if !ok {
				return nil, fmt.Errorf("buildExpression: variable %s not defined", token.Value)
			}

			left = VariableExpression{Name: token.Value, DataType: dataType}
		} else {
			left, err = getLiteralExpression(&token)

			if err != nil {
				return nil, err
			}
		}
	}

	if (i < len(*equation) - 2) {
		slice := (*equation)[i+1:len(*equation)]

		parenCheck := (*equation)[i + 1]

		// Right side of operation will be LPAREN if it is a PAREN
		if (parenCheck.Type == lexer.RPAREN) {
			return nil, fmt.Errorf("buildBinaryExpression: Wrong parenthesis on RHS of operand")
		}

		if (parenCheck.Type == lexer.LPAREN) {
			end, err := getRParen(equation, i + 1)
			
			if err != nil {
				return nil, err
			}

			slice = (*equation)[i+2:end]
		}

		rightExp, err := buildExpression(&slice, variables)

		if err != nil {
			return nil, err
		}

		right = rightExp
	} else {
		token := (*equation)[i + 1]

		if (token.Type == lexer.IDENT) {
			dataType, ok := (*variables)[token.Value]

			if !ok {
				return nil, fmt.Errorf("buildExpression: variable %s not defined", token.Value)
			}

			right = VariableExpression{Name: token.Value, DataType: dataType}
		} else {
			right, err = getLiteralExpression(&token)

			if err != nil {
				return nil, err
			}
		}
	}

	if left.GetDataType() != Float && left.GetDataType() != Int {
		return nil, fmt.Errorf("buildBinaryExpression: LHS as not int or float")
	}

	if right.GetDataType() != Float && right.GetDataType() != Int {
		return nil, fmt.Errorf("buildBinaryExpression: LHS as not int or float")
	}

	dataType := Int

	// Assumes that the expression only contains numbers
	if (left.GetDataType() == Float || right.GetDataType() == Float) {
		dataType = Float
	}

	expression := BinaryExpression{Operand: operand, Left: left, Right: right, DataType: dataType}

	return expression, nil
}

func getOperandFromToken(token *lexer.Token) (BinaryOperator, error) {
	switch ((*token).Type) {
	case lexer.ADD:
		return Add, nil
	case lexer.SUBTRACT:
		return Subtract, nil
	case lexer.DIVIDE:
		return Divide, nil
	case lexer.MULTIPLY:
		return Multiply, nil
	default:
		return Add, fmt.Errorf("Invalid operand type %+v", token)
	}
}

func getLiteralExpression(token *lexer.Token) (IExpression, error) {
	switch (token.Type) {
	case lexer.INT:
		i, _ := strconv.ParseInt(token.Value, 10, 0)
		expression := LiteralExpression[int64]{ Data: i, DataType: Int }
		return expression, nil
	case lexer.FLOAT:
		i, _ := strconv.ParseFloat(token.Value, 0)
		expression := LiteralExpression[float64]{ Data: i, DataType: Float }
		return expression, nil
	case lexer.STRING:
		expression := LiteralExpression[string]{ Data: token.Value, DataType: String }
		return expression, nil
	default:
		return nil, errors.New("parseAssign: Bad data type on expression")
	}
}

func getLParen(tokens *[]lexer.Token, LParen int) (int, error) {
	depth := 1

	for i := LParen - 1; i >= 0; i-- {
		token := (*tokens)[i]

		switch token.Type {
		case lexer.RPAREN:
			depth += 1
		case lexer.LPAREN:
			depth -= 1
		}

		if (depth == 0) {
			return i, nil;
		}
	}

	return -1, fmt.Errorf("getLParen: could not find left parenthesis")
}

func getRParen(tokens *[]lexer.Token, RParen int) (int, error) {
	depth := 1

	for i := RParen + 1; i < len(*tokens); i++ {
		token := (*tokens)[i]

		switch token.Type {
		case lexer.RPAREN:
			depth -= 1
		case lexer.LPAREN:
			depth += 1
		}

		if (depth == 0) {
			return i, nil;
		}
	}

	return -1, fmt.Errorf("getRParen: could not find right parenthesis")
}

func unaryFix(tokens *[]lexer.Token) {
	// Places to insert a zero (the index of the subtract symbol to insert before)
	zeroPlaces := []int{}

	for i, token := range (*tokens) {
		if token.Type != lexer.SUBTRACT {
			continue
		}
		
		// not going to throw syntax error because that will be handled elsewhere
		if i == len(*tokens) - 1 {
			break
		}

		// Subtract at start of equation
		if i == 0 && token.Type == lexer.SUBTRACT {
			zeroPlaces = append(zeroPlaces, 0)
			continue
		}

		switch ((*tokens)[i - 1].Type) {
		case lexer.INT, lexer.FLOAT, lexer.IDENT:
			continue
		default:
			zeroPlaces = append(zeroPlaces, i)
		}
	}

	for i := len(zeroPlaces) - 1; i >= 0; i-- {
		pos := zeroPlaces[i]

		newTokens := make([]lexer.Token, 0, len(*tokens)+1)
		newTokens = append(newTokens, (*tokens)[:pos]...)
		newTokens = append(newTokens, lexer.Token{
			Type:  lexer.INT,
			Value: "0",
		})
		newTokens = append(newTokens, (*tokens)[pos:]...)

		*tokens = newTokens
	}
}