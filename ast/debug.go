package ast

import (
	"fmt"
	"strings"

	"liamg2810.github.io/csjs/data"
)

func Debug(nodes []data.INode) string {
	var b strings.Builder

	for i, node := range nodes {
		if i > 0 {
			b.WriteByte('\n')
		}

		writeNode(&b, node, "", i == len(nodes)-1)
	}

	return b.String()
}

func writeNode(b *strings.Builder, node data.INode, prefix string, last bool) {
	if node == nil {
		fmt.Fprintf(b, "%snil\n", prefix)
		return
	}

	branch := "├── "
	nextPrefix := prefix + "│   "

	if last {
		branch = "└── "
		nextPrefix = prefix + "    "
	}

	switch n := node.(type) {
	case data.Assignment:
		fmt.Fprintf(b, "%sAssignment: %s\n", prefix+branch, n.Name)
		writeNode(b, n.Value, nextPrefix, true)

	case data.BinaryExpression:
		fmt.Fprintf(b, "%sBinaryExpression (DataType: %s)\n",
			prefix+branch, n.DataType)
		writeNode(b, n.Left, nextPrefix, false)
		fmt.Fprintf(b, "%s├── Operator: %s\n", nextPrefix, n.Operand)
		writeNode(b, n.Right, nextPrefix, true)

	case data.UnaryExpression:
		fmt.Fprintf(b, "%sUnaryExpression (DataType: %s)\n",
			prefix+branch, n.DataType)
		writeNode(b, n.Expression, nextPrefix, true)

	case data.VariableExpression:
		fmt.Fprintf(b, "%sVariableExpression: %s (DataType: %s)\n",
			prefix+branch, n.Name, n.DataType)

	case data.LiteralExpression[int64]:
		fmt.Fprintf(b, "%sLiteralExpression[int64]: %d (DataType: %s)\n",
			prefix+branch, n.Data, n.DataType)

	case data.LiteralExpression[int32]:
		fmt.Fprintf(b, "%sLiteralExpression[int32]: %d (DataType: %s)\n",
			prefix+branch, n.Data, n.DataType)

	case data.LiteralExpression[float32]:
		fmt.Fprintf(b, "%sLiteralExpression[float32]: %g (DataType: %s)\n",
			prefix+branch, n.Data, n.DataType)

	case data.LiteralExpression[string]:
		fmt.Fprintf(b, "%sLiteralExpression[string]: %q (DataType: %s)\n",
			prefix+branch, n.Data, n.DataType)

	case data.LiteralExpression[float64]:
		fmt.Fprintf(b, "%sLiteralExpression[float64]: %g (DataType: %s)\n",
			prefix+branch, n.Data, n.DataType)

	case data.LiteralExpression[bool]:
		fmt.Fprintf(b, "%sLiteralExpression[bool]: %t (DataType: %s)\n",
			prefix+branch, n.Data, n.DataType)

	default:
		fmt.Fprintf(b, "%s%T\n", prefix+branch, node)
	}
}
