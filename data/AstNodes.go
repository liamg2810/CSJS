package data

import (
	"fmt"
)

type BinaryOperator int
type DataType int

type INode interface {
	node()
}

type IExpression interface {
	INode
	expression()
	GetDataType() DataType
}

type IAssignment interface {
	INode
	assignment()
}

type LiteralExpression[T any] struct {
	Data     T
	DataType DataType
}

type BinaryExpression struct {
	Left     IExpression
	Operand  BinaryOperator
	Right    IExpression
	DataType DataType
}

type VariableExpression struct {
	Name     string
	DataType DataType
}

type UnaryExpression struct {
	Expression IExpression
	DataType   DataType
}

type Assignment struct {
	Name     string
	Value    IExpression
	DataType DataType
}

func (LiteralExpression[T]) node()       {}
func (LiteralExpression[T]) expression() {}
func (l LiteralExpression[T]) GetDataType() DataType {
	return l.DataType
}

func (BinaryExpression) node()       {}
func (BinaryExpression) expression() {}
func (b BinaryExpression) GetDataType() DataType {
	return b.DataType
}

func (Assignment) node()       {}
func (Assignment) assignment() {}

func (VariableExpression) node()       {}
func (VariableExpression) expression() {}
func (v VariableExpression) GetDataType() DataType {
	return v.DataType
}

func (UnaryExpression) node()       {}
func (UnaryExpression) expression() {}
func (e UnaryExpression) GetDataType() DataType {
	return e.DataType
}

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

func (d DataType) String() string {
	switch d {
	case Int:
		return "Int"
	case String:
		return "String"
	case Float:
		return "Float"
	case Boolean:
		return "Boolean"
	default:
		return fmt.Sprintf("DataType(%d)", d)
	}
}

func (o BinaryOperator) String() string {
	switch o {
	case Add:
		return "Add"
	case Subtract:
		return "Subtract"
	case Multiply:
		return "Multiply"
	case Divide:
		return "Divide"
	default:
		return fmt.Sprintf("BinaryOperator(%d)", o)
	}
}
