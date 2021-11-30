package evaluator

import (
	"fmt"
	"strconv"
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
		return parseProgramNode(n, e)
	case parser.BinaryOperationNode:
		return parseBinOpNode(n, e)
	case parser.UnaryOpNode:
		return parseUnaryOpNode(n, e)
	case parser.AssignmentNode:
		return parseAssignNode(n, e)
	case parser.IfNode:
		return parseIfNode(n, e)
	case parser.WhileNode:
		return parseWhileNode(n, e)
	case parser.ForNode:
		return parseForNode(n, e)
	case parser.VarAccessNode:
		return e.Variables[n.Identifier]
	case parser.FunctionCallNode:
		return parseFunctionCallNode(n, e)
	case parser.FunctionDefenitionNode:
		return parseFunctionDefenitionNode(n, e)
	case parser.IntNode:
		return n.Value
	case parser.StringNode:
		return n.Value
	}
	return -1
}

func parseProgramNode(n parser.ProgramNode, e *Environment) interface{} {
	for _, node := range n.Expressions {
		switch n := node.(type) {
		case parser.ReturnNode:
			return n
		}
		Eval(node, e)
	}
	return -1
}

func parseForNode(n parser.ForNode, e *Environment) interface{} {
	// TODO: FIX IT COZ U STUPID, CREATE IDENTIFIER FOR RETURN NODE
	for i := Eval(n.MinValue, e).(int); i < Eval(n.MaxValue, e).(int); i++ {
		var returned interface{}
		e.Variables[n.Identifier] = i
		for _, node := range n.Consequence.Expressions {
			returned = Eval(node, e)

			if isReturn(returned) {
				return returned
			}
			if isReturn(node) {
				return node
			}
		}
	}
	return -1
}

func parseWhileNode(n parser.WhileNode, e *Environment) interface{} {
	for parseConditions(n.Condition, e) {
		var returned interface{}
		for _, node := range n.Consequence.Expressions {
			returned = Eval(node, e)

			if isReturn(returned) {
				return returned
			}
			if isReturn(node) {
				return node
			}
		}
	}
	return -1
}

func parseIfNode(n parser.IfNode, e *Environment) interface{} {
	if parseConditions(n.Condition, e) {
		for _, node := range n.Consequence.Expressions {
			returned := Eval(node, e)

			if isReturn(returned) {
				return returned
			}
			if isReturn(node) {
				return node
			}
		}
	} else {
		for _, node := range n.Alternate.Expressions {
			returned := Eval(node, e)

			if isReturn(returned) {
				return returned
			}
			if isReturn(node) {
				return node
			}
		}
	}
	return -1
}

func parseConditions(conditions []parser.ConditionNode, e *Environment) bool {
	result := true
	for _, condition := range conditions {
		evaluated := toBool(Eval(condition.Condition, e).(int))
		switch condition.Seperator {
		case lexer.AND:
			result = result && evaluated
		case lexer.OR:
			result = result || evaluated
		}
	}

	return result
}

func parseAssignNode(n parser.AssignmentNode, e *Environment) interface{} {
	e.Variables[n.Identifier] = Eval(n.Value, e)
	return Eval(n.Value, e)
}

func parseUnaryOpNode(n parser.UnaryOpNode, e *Environment) int {
	switch n.Op {
	case lexer.SUB:
		return -Eval(n.Right, e).(int)
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
	switch n.Op {
	case lexer.ADD:
		return Eval(n.Left, e).(int) + Eval(n.Right, e).(int)
	case lexer.SUB:
		return Eval(n.Left, e).(int) - Eval(n.Right, e).(int)
	case lexer.MUL:
		return Eval(n.Left, e).(int) * Eval(n.Right, e).(int)
	case lexer.DIV:
		return Eval(n.Left, e).(int) / Eval(n.Right, e).(int)
	case lexer.MOD:
		return Eval(n.Left, e).(int) % Eval(n.Right, e).(int)

	case lexer.EE:
		return toBinary(Eval(n.Left, e).(int) == Eval(n.Right, e).(int))
	case lexer.NE:
		return toBinary(Eval(n.Left, e).(int) != Eval(n.Right, e).(int))
	case lexer.GT:
		return toBinary(Eval(n.Left, e).(int) > Eval(n.Right, e).(int))
	case lexer.LT:
		return toBinary(Eval(n.Left, e).(int) < Eval(n.Right, e).(int))
	case lexer.GTE:
		return toBinary(Eval(n.Left, e).(int) >= Eval(n.Right, e).(int))
	case lexer.LTE:
		return toBinary(Eval(n.Left, e).(int) <= Eval(n.Right, e).(int))
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

func isReturn(node interface{}) bool {
	switch node.(type) {
	case parser.ReturnNode:
		return true
	}
	return false
}

func contained(list []string, element string) bool {
	for _, e := range list {
		if e == element {
			return true
		}
	}
	return false
}

func parseFunctionCallNode(n parser.FunctionCallNode, e *Environment) interface{} {
	var FUNCTIONS = []string{
		"print",
	}
	if contained(FUNCTIONS, n.Identifier) {
		switch n.Identifier {
		case "print":
			handlePrint(n, e)
		}
	}
	return -1
}

func parseFunctionDefenitionNode(n parser.FunctionDefenitionNode, e *Environment) interface{} {
	return nil
}

func handlePrint(n parser.FunctionCallNode, e *Environment) string {
	str := ""
	for _, param := range n.Parameters {
		result := Eval(param, e)
		switch res := result.(type) {
		case int:
			str += strconv.Itoa(res)
		case string:
			str += res
		}
	}
	fmt.Println(str)
	return str
}
