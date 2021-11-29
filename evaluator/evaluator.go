package evaluator

import (
	"fmt"
	"terminascript/lexer"
	"terminascript/parser"
)

func Eval(node interface{}) interface{} {
	switch n := node.(type) {
	case parser.ProgramNode:
		 return parseProgramNode(n)
	case parser.BinaryOperationNode:
		return parseBinOpNode(n)
	case parser.UnaryOpNode:
		return parseUnaryOpNode(n)
	case parser.IntNode:
		return n.Value
	case parser.StringNode:
		return n.Value
	default:
		fmt.Println(n)
	}
	return -1
}

func parseProgramNode(n parser.ProgramNode) int {
	for _,node := range n.Expressions {
		fmt.Println(Eval(node))
	}

	return 1
}

func parseUnaryOpNode(n parser.UnaryOpNode) int {
	switch n.Op {
	case lexer.SUB:
		return - Eval(n.Right).(int)
	case lexer.NOT:
		if Eval(n.Right).(int) == 1 {
			return 0
		} else {
			return 1
		}
	}

	return -1
}

func parseBinOpNode(n parser.BinaryOperationNode) int {
	switch n.Op{
	case lexer.ADD:
		return Eval(n.Left).(int) + Eval(n.Right).(int)
	case lexer.SUB:
		return Eval(n.Left).(int) - Eval(n.Right).(int)
	case lexer.MUL:
		return Eval(n.Left).(int) * Eval(n.Right).(int)
	case lexer.DIV:
		return Eval(n.Left).(int) / Eval(n.Right).(int)
	case lexer.MOD:
		return Eval(n.Left).(int) % Eval(n.Right).(int)

	case lexer.EE:
		return toBinary(Eval(n.Left).(int) == Eval(n.Right).(int))
	case lexer.NE:
		return toBinary(Eval(n.Left).(int) != Eval(n.Right).(int))
	case lexer.GT:
		return toBinary(Eval(n.Left).(int) > Eval(n.Right).(int))
	case lexer.LT:
		return toBinary(Eval(n.Left).(int) < Eval(n.Right).(int))
	case lexer.GTE:
		return toBinary(Eval(n.Left).(int) >= Eval(n.Right).(int))
	case lexer.LTE:
		return toBinary(Eval(n.Left).(int) <= Eval(n.Right).(int))
	}

	return -1
}

func toBinary(value bool) int {
	if value {
		return 1
	} else {
		return 0
	}
}