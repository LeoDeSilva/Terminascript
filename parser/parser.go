package parser

import (
	"strconv"
	"terminascript/lexer"
)

// TODO : if ... elif ... else statements

type Parser struct {
	tokens       []lexer.Token
	position     int
	readPosition int
	token        lexer.Token
}

func NewParser(tokens []lexer.Token) *Parser {
	p := &Parser{tokens: tokens}
	p.advance()
	return p
}

func (p *Parser) advance() {
	if p.readPosition >= len(p.tokens) {
		p.token = lexer.Token{}
	} else {
		p.token = p.tokens[p.readPosition]
	}

	p.position = p.readPosition
	p.readPosition++
}

func Includes(array []string, element string) bool {
	for _,e := range array {
		if e == element{
			return true
		}
	}
	return false
}

func (p *Parser) peekToken() lexer.Token {
	if p.readPosition >= len(p.tokens) {
		return lexer.Token{}
	} else {
		return p.tokens[p.readPosition]
	}
}

func (p *Parser) Parse() ProgramNode {
	var ast = ProgramNode{lexer.PROGRAM_NODE, make([]interface{}, 0)}

	for (p.token.Type != lexer.EOF && p.token != lexer.Token{}) {
		node := p.ParseExpr()
		ast.Expressions = append(ast.Expressions, node)
		p.advance()
	}

	return ast
}

func (p *Parser) ParseExpr() interface{} {
	switch p.token.Type {
	case lexer.LET:
		p.advance()
		return p.ParseAssignment()
	case lexer.RETURN:
		return p.ParseReturn()
	case lexer.IF:
		return p.ParseIf()
	case lexer.WHILE:
		return p.ParseWhile()
	case lexer.FOR:
		return p.ParseFor()
	case lexer.FUNC:
		return p.ParseFunction()
	default:
		if p.token.Type == lexer.IDENTIFIER {
			if p.peekToken().Type == lexer.EQ || p.peekToken().Type == lexer.ASSIGN {
				return p.ParseAssignment()
			}
		}
		return p.ParseComparison()
	}
}

func (p *Parser) ParseReturn() interface{} {
	p.advance()
	return ReturnNode{lexer.RETURN, p.ParseComparison()}
}

func (p *Parser) ParseFunction() interface{} {
	p.advance()
	
	if p.token.Type != lexer.IDENTIFIER { return nil }
	identifier := p.token.Literal
	p.advance()

	if p.token.Type != lexer.LPAREN { return nil }
	parameters := p.ParseParameters()
	p.advance()

	if p.token.Type != lexer.LBRACE { return nil }
	p.advance()
	return FunctionDefenitionNode{lexer.FUNCTION_DEFENITION_NODE,identifier,parameters,ProgramNode{lexer.PROGRAM_NODE,p.ParseMultiline()}}
}

func (p *Parser) ParseFor() interface{} {
	p.advance()

	if p.token.Type != lexer.LPAREN { return nil }
	p.advance()

	if p.token.Type != lexer.IDENTIFIER { return nil }
	identifier := p.token.Literal
	p.advance()

	if p.token.Type != lexer.ASSIGN && p.token.Type != lexer.EQ { return nil }
	p.advance()

	if p.token.Type != lexer.INT { return nil }
	min,_ := strconv.Atoi(p.token.Literal)
	p.advance()

	if p.token.Type != lexer.ARROW { return nil }
	p.advance()

	if p.token.Type != lexer.INT { return nil }
	max,_ := strconv.Atoi(p.token.Literal)
	p.advance()

	if p.token.Type != lexer.RPAREN { return nil }
	p.advance()

	if p.token.Type != lexer.LBRACE { return nil }
	p.advance()

	return ForNode{lexer.FOR_NODE, identifier, min, max,ProgramNode{lexer.PROGRAM_NODE,p.ParseMultiline()}}
}

func (p *Parser) ParseWhile() interface{} {
	p.advance()
	conditions := p.ParseConditions()

	p.advance()
	if p.token.Type != lexer.LBRACE { return nil }

	p.advance()
	consequence := ProgramNode{lexer.PROGRAM_NODE, p.ParseMultiline()}	

	return WhileNode{lexer.WHILE_NODE, conditions, consequence}
}

func (p *Parser) ParseIf() interface{} {
	p.advance()
	conditions := p.ParseConditions()

	p.advance()
	if p.token.Type != lexer.LBRACE { return nil }

	p.advance()
	prog := ProgramNode{lexer.PROGRAM_NODE, p.ParseMultiline()}	

	return IfNode{lexer.IF_NODE,conditions,prog,ProgramNode{}}
}

func (p *Parser) ParseConditions() []ConditionNode {
	var conditions []ConditionNode
	var seperators = []string{lexer.AND,lexer.OR}

	if p.token.Type != lexer.LPAREN { return nil }
	p.advance()

	if p.token.Type == lexer.RPAREN { return conditions }
	conditions = append(conditions, ConditionNode{lexer.CONDITION_NODE,"AND",p.ParseComparison()})
	
	currentSeperator := "AND"
	if (p.token != lexer.Token{}) {
		for (p.token != lexer.Token{} && p.token.Type != lexer.RPAREN && p.token.Type != lexer.SEMICOLON) {
			isSeperator := Includes(seperators,p.token.Type)
			if !isSeperator {
				conditions = append(conditions, ConditionNode{lexer.CONDITION_NODE,currentSeperator,p.ParseComparison()})
			} else {
				currentSeperator = p.token.Type
				p.advance()
			}
		}
	}

	return conditions
}

func (p *Parser) ParseMultiline() []interface{} {
	var nodes []interface{}
	for (p.token.Type != lexer.RBRACE) {
		if p.token.Type != lexer.SEMICOLON {
			expr := p.ParseExpr()
			p.advance()
			nodes = append(nodes,expr)
		} else {
			p.advance()
		}
	}
	return nodes
}

func (p *Parser) ParseAssignment() interface{} {
	if p.token.Type != lexer.IDENTIFIER { return nil }
	identifier := p.token.Literal

	p.advance()
	if p.token.Type != lexer.EQ && p.token.Type != lexer.ASSIGN { return nil }

	p.advance()
	return AssignmentNode{lexer.ASSIGN_NODE, identifier, p.ParseComparison()}
}

func (p *Parser) ParseComparison() interface{} {
	leftNode := p.ParseArith()
	if p.token.Type != lexer.SEMICOLON && p.token.Type != lexer.EOF {
		var operations = []string{lexer.EE,lexer.NE,lexer.GT,lexer.GTE,lexer.LT,lexer.LTE}

		if Includes(operations,p.token.Type){
			op := p.token.Type
			p.advance()
			return BinaryOperationNode{Type: lexer.BIN_OP_NODE, Left: leftNode, Op: op, Right: p.ParseComparison()}
		}
	}
	return leftNode
}

func (p *Parser) ParseArith() interface{} {
	leftNode := p.ParseTerm()
	if p.token.Type != lexer.SEMICOLON && p.token.Type != lexer.EOF {
		var operations = []string{lexer.ADD,lexer.SUB,lexer.MOD}

		if Includes(operations,p.token.Type){
			op := p.token.Type
			p.advance()
			return BinaryOperationNode{Type: lexer.BIN_OP_NODE, Left: leftNode, Op: op, Right: p.ParseArith()}
		}

	}
	return leftNode
}

func (p *Parser) ParseTerm() interface{} {
	leftNode := p.ParseFactor()
	if p.token.Type != lexer.SEMICOLON && p.token.Type != lexer.EOF {
		p.advance()
		var operations = []string{lexer.MUL,lexer.DIV}

		if Includes(operations,p.token.Type){
			op := p.token.Type
			p.advance()
			return BinaryOperationNode{Type: lexer.BIN_OP_NODE, Left: leftNode, Op: op, Right: p.ParseTerm()}
		}

	}
	return leftNode
}

func (p *Parser) ParseFactor() interface{} {
	for p.token.Type != lexer.EOF && p.token.Type != lexer.SEMICOLON {
		switch p.token.Type{
		case lexer.IDENTIFIER:
			ID := p.token.Literal

			if p.peekToken().Type == lexer.LPAREN {
				p.advance()
				parameters := p.ParseParameters()
				return FunctionCallNode{lexer.FUNC_CALL_NODE, ID, parameters}

			} else {
				return VarAccessNode{lexer.VAR_ACCESS_NODE, ID}
			}

		case lexer.INT:
			intValue, _ := strconv.Atoi(p.token.Literal)
			return IntNode{lexer.INT_NODE, intValue}

		case lexer.STRING:
			return StringNode{lexer.STRING_NODE, p.token.Literal}

		case lexer.LPAREN:
			p.advance()
			expr := p.ParseComparison()
			if p.token.Type == lexer.RPAREN {
				return expr
			}

		case lexer.SUB:
			p.advance()
			return UnaryOpNode{lexer.UNARY_NODE, lexer.SUB, p.ParseFactor()}
		
		case lexer.NOT:
			p.advance()
			return UnaryOpNode{lexer.UNARY_NODE, lexer.NOT, p.ParseFactor()}
		}
	}
	return ErrorNode{}
}

func (p *Parser) ParseParameters() []interface{} {
	parameters := make([]interface{}, 0)
	p.advance()

	if p.token.Type == lexer.RPAREN {
		return parameters
	}

	for (p.token != lexer.Token{} && p.token.Type != lexer.RPAREN && p.token.Type != lexer.SEMICOLON) {
		if p.token.Type != lexer.COMMA {
			parameters = append(parameters, p.ParseExpr())
		} else {
			p.advance()
		}
	}
	return parameters
}