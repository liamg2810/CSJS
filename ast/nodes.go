package ast

type BinaryOperator int
type DataType int

type INode interface {
	node()
}

type IExpression interface {
	INode
	expression()
}

type IAssignment interface {
	INode
	assignment()
}

type LiteralExpression[T any] struct {
	Data T
	DataType DataType
}

type BinaryExpression struct {
	Left IExpression
	Operand BinaryOperator
	Right IExpression
	DataType DataType
}

type VariableExpression struct {
	Name string
	DataType DataType
}

type Assignment struct {
	Name string
	Value IExpression
	DataType DataType
}

func (LiteralExpression[T]) node() {}
func (LiteralExpression[T]) expression() {}

func (BinaryExpression) node() {}
func (BinaryExpression) expression() {}

func (Assignment) node() {}
func (Assignment) assignment() {}

func (VariableExpression) node() {}
func (VariableExpression) expression() {}

const (
	Add BinaryOperator = iota
	Subtract
	Multiply
	Divide
)

const (
	Int DataType = iota
	String
	Float
	Boolean
)