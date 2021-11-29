package evaluator

import (
	"fmt"
	"terminascript/lexer"
	"terminascript/parser"
)

type Environment struct {
	Variables map[string]interface{}
	Functions map[string]parser.FunctionDefenitionNode
}

func NewEnvironment() *Environment {
	return &Environment{
		Variables: make(map[string]interface{}),
		Functions: make(map[string]parser.FunctionDefenitionNode),
	}
}

func Eval(node interface{}, e *Environment) interface{} {
	switch n := node.(type) {
	case parser.ProgramNode:
		 return parseProgramNode(n,e)
	case parser.BinaryOperationNode:
		return parseBinOpNode(n,e)
	case parser.UnaryOpNode:
		return parseUnaryOpNode(n,e)
	case parser.AssignmentNode:
		return parseAssignNode(n,e)
	case parser.IfNode:
		return parseIfNode(n,e)
	case parser.WhileNode:
		return parseWhileNode(n,e)
	case parser.ForNode:
		return parseForNode(n,e)
	case parser.VarAccessNode:
		return e.Variables[n.Identifier]
	case parser.IntNode:
		return n.Value
	case parser.StringNode:
		return n.Value
	default:
		fmt.Println(n)
	}
	return -1
}


func parseProgramNode(n parser.ProgramNode, e *Environment) int {
	for _,node := range n.Expressions {
		//TODO: If node.type === return : return Eval(returnNode.value)
		fmt.Println(Eval(node,e))
	}

	return 1
}

func parseForNode(n parser.ForNode, e *Environment) int {
	for i := Eval(n.MinValue,e).(int); i < Eval(n.MaxValue,e).(int); i++ {
		e.Variables[n.Identifier] = i
		Eval(n.Consequence,e)
	}
	return -1
}

func parseWhileNode(n parser.WhileNode, e *Environment) int {
	for parseConditions(n.Condition, e) {
		Eval(n.Consequence,e)
	}
	return -1
}

func parseIfNode(n parser.IfNode, e *Environment) interface{} {
	if parseConditions(n.Condition, e) {
		return Eval(n.Consequence,e)
	} else{
		return Eval(n.Alternate, e)
	}
}

func parseConditions(conditions []parser.ConditionNode, e *Environment) bool {
	result := true
	for _,condition := range conditions {
		evaluated := toBool(Eval(condition.Condition, e).(int))
		switch condition.Seperator{
		case lexer.AND:
			result = result && evaluated
		case lexer.OR:
			result = result || evaluated
		}
	}

	return result 
}

func parseAssignNode(n parser.AssignmentNode, e *Environment) interface{} {
	e.Variables[n.Identifier] = Eval(n.Value,e)
	return Eval(n.Value, e)
}

func parseUnaryOpNode(n parser.UnaryOpNode, e *Environment) int {
	switch n.Op {
	case lexer.SUB:
		return - Eval(n.Right, e).(int)
	case lexer.NOT:
		if Eval(n.Right, e).(int) == 1 {
			return 0
		} else {
			return 1
		}
	}

	return -1
}

func parseBinOpNode(n parser.BinaryOperationNode, e *Environment) int {
	switch n.Op{
	case lexer.ADD:
		return Eval(n.Left, e).(int) + Eval(n.Right,e).(int)
	case lexer.SUB:
		return Eval(n.Left,e).(int) - Eval(n.Right,e).(int)
	case lexer.MUL:
		return Eval(n.Left,e).(int) * Eval(n.Right,e).(int)
	case lexer.DIV:
		return Eval(n.Left,e).(int) / Eval(n.Right,e).(int)
	case lexer.MOD:
		return Eval(n.Left,e).(int) % Eval(n.Right,e).(int)

	case lexer.EE:
		return toBinary(Eval(n.Left,e).(int) == Eval(n.Right,e).(int))
	case lexer.NE:
		return toBinary(Eval(n.Left,e).(int) != Eval(n.Right,e).(int))
	case lexer.GT:
		return toBinary(Eval(n.Left,e).(int) > Eval(n.Right,e).(int))
	case lexer.LT:
		return toBinary(Eval(n.Left,e).(int) < Eval(n.Right,e).(int))
	case lexer.GTE:
		return toBinary(Eval(n.Left,e).(int) >= Eval(n.Right,e).(int))
	case lexer.LTE:
		return toBinary(Eval(n.Left,e).(int) <= Eval(n.Right,e).(int))
	}

	return -1
}

func toBool(value int) bool {
	if value == 1 {
		return true
	} else {
		return false
	}
}

func toBinary(value bool) int {
	if value {
		return 1
	} else {
		return 0
	}
}