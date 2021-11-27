package parser

type ProgramNode struct {
	Type        string
	Expressions []interface{}
}

type ReturnNode struct {
	Type       string
	Expression interface{}
}

type FunctionDefenitionNode struct {
	Type        string
	Identifier  string
	Parameters  []interface{}
	Consequence ProgramNode
}

type ForNode struct {
	Type        string
	Identifier  string
	MinValue    int
	MaxValue    int
	Consequence ProgramNode
}

type WhileNode struct {
	Type        string
	Condition   []ConditionNode
	Consequence ProgramNode
}

type IfNode struct {
	Type        string
	Condition   []ConditionNode
	Consequence ProgramNode
	Alternate   ProgramNode
}

type IfConditionNode struct {
	Type        string
	Condition   ConditionNode
	Consequence ProgramNode
}

type ConditionNode struct {
	Type      string
	Seperator string
	Condition interface{}
}

type FunctionCallNode struct {
	Type       string
	Identifier string
	Parameters []interface{}
}

type AssignmentNode struct {
	Type       string
	Identifier string
	Value      interface{}
}

type ParameterNode struct {
	Type       string
	Identifier string
}

type BinaryOperationNode struct {
	Type  string
	Left  interface{}
	Op    string
	Right interface{}
}

type UnaryOpNode struct {
	Type  string
	Op    string
	Right interface{}
}

type VarAccessNode struct {
	Type  string
	Value string
}

type IntNode struct {
	Type  string
	Value int
}

type StringNode struct {
	Type  string
	Value string
}

type ErrorNode struct {
	Type string
}
