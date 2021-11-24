package parser

type ProgramNode struct {
	Type        string
	Expressions []interface{}
}

type FunctionCallNode struct {
	Type       string
	Identifier string
	Parameters []interface{}
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
