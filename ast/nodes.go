package ast

type BinaryOperator int

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
}

type BinaryExpression struct {
	Left IExpression
	Operand BinaryOperator
	Right IExpression
}

type Assignment struct {
	Name string
	Value IExpression
}

func (LiteralExpression[T]) node() {}
func (LiteralExpression[T]) expression() {}

func (BinaryExpression) node() {}
func (BinaryExpression) expression() {}

func (Assignment) node() {}
func (Assignment) assignment() {}

const (
	Add BinaryOperator = iota
	Subtract
	Multiply
	Divide
)