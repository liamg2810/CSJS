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
	eolPos := 0

	// Ignoring BODMAS for now
	isEquation := false
	tokensToRead := []lexer.Token{}

	// Grabbing tokens to fucking figure out
	for i := startIx + 3; i < len(*tokens); i++ {
		token := (*tokens)[i]

		switch token.Type {
		case lexer.EOL:
			hitEOL = true
			eolPos = i
		case lexer.ADD, lexer.SUBTRACT:
			isEquation = true
		case lexer.DIVIDE, lexer.MULTIPLY:
			isEquation = true
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
			case lexer.INT, lexer.FLOAT, lexer.STRING:
				e, err := getLiteralExpression(&token)

				if err != nil {
					return nil, err;
				}

				expression = e
				dataType = e.GetDataType()
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

	slice := (*tokens)[startIx + 3:eolPos]
		for _, token := range slice {
			fmt.Printf("%+v", token)
		}
	binaryExpression, err := buildBinaryExpression(&slice)

	if (err != nil) {
		return nil, err
	}

	return &Assignment{ Name: name.Value, Value: *binaryExpression, DataType: binaryExpression.GetDataType() }, nil
}

// for now will not allow variables in the equations
func buildBinaryExpression(equation *[]lexer.Token) (*BinaryExpression, error) {
	// first find all symbol types in the equation

	// Multiple + Divide
	md := []int{}
	// Add + Subtract
	as := []int{}

	for i, token := range *equation {
		switch (token.Type) {
		case lexer.MULTIPLY, lexer.DIVIDE:
			md = append(md, i)
		case lexer.ADD, lexer.SUBTRACT:
			as = append(as, i)
		default:
			continue
		}
	}

	// then recursively solve equation (excluding parenthesis and exponents)
	i := 0

	if len(as) > 0 {
		i = as[len(as)-1]
	} else if len(md) > 0 {
		i = md[len(md)-1]
	} else {
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
		fmt.Printf("entering LHS recursion")
		slice := (*equation)[0:i]
		leftExp, err := buildBinaryExpression(&slice)

		if err != nil {
			return nil, err
		}

		left = *leftExp
	} else {
		token := (*equation)[i - 1]
		left, err = getLiteralExpression(&token)

		if err != nil {
			return nil, err
		}
	}

	if (i < len(*equation) - 2) {
		fmt.Printf("entering RHS recursion")
		slice := (*equation)[i+1:len(*equation)]
		rightExp, err := buildBinaryExpression(&slice)

		if err != nil {
			return nil, err
		}

		right = *rightExp
	} else {
		token := (*equation)[i + 1]
		right, err = getLiteralExpression(&token)

		if err != nil {
			return nil, err
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

	fmt.Printf("Left: %+v\n", left)
	fmt.Printf("Right: %+v\n", right)

	expression := BinaryExpression{Operand: operand, Left: left, Right: right, DataType: dataType}

	return &expression, nil
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