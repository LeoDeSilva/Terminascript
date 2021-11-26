package parser

import (
	"strconv"
	"terminascript/lexer"
)

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
	var node interface{}

	switch p.token.Type {
	default:
		node = p.ParseArith()
	}

	return node
}

func (p *Parser) ParseArith() interface{} {
	leftNode := p.ParseTerm()
	if p.token.Type != lexer.SEMICOLON && p.token.Type != lexer.EOF {

		switch p.token.Type {
		case lexer.ADD:
			p.advance()
			rightNode := p.ParseArith()
			return BinaryOperationNode{Type: lexer.BIN_OP_NODE, Left: leftNode, Op: lexer.ADD, Right: rightNode}
		case lexer.SUB:
			p.advance()
			rightNode := p.ParseArith()
			return BinaryOperationNode{Type: lexer.BIN_OP_NODE, Left: leftNode, Op: lexer.SUB, Right: rightNode}
		}
	}
	return leftNode
}

func (p *Parser) ParseTerm() interface{} {
	leftNode := p.ParseFactor()
	if p.token.Type != lexer.SEMICOLON && p.token.Type != lexer.EOF {
		p.advance()

		switch p.token.Type {
		case lexer.MUL:
			p.advance()
			rightNode := p.ParseTerm()
			return BinaryOperationNode{Type: lexer.BIN_OP_NODE, Left: leftNode, Op: lexer.MUL, Right: rightNode}
		case lexer.DIV:
			p.advance()
			rightNode := p.ParseTerm()
			return BinaryOperationNode{Type: lexer.BIN_OP_NODE, Left: leftNode, Op: lexer.DIV, Right: rightNode}
		}
	}
	return leftNode
}

func (p *Parser) ParseFactor() interface{} {
	for p.token.Type != lexer.EOF && p.token.Type != lexer.SEMICOLON {
		if p.token.Type == lexer.IDENTIFIER {
			ID := p.token.Literal

			if p.peekToken().Type == lexer.LPAREN {
				p.advance()
				parameters := p.ParseParameters()
				return FunctionCallNode{lexer.FUNC_CALL_NODE, ID, parameters}

			} else {
				return VarAccessNode{lexer.VAR_ACCESS_NODE, ID}
			}

		} else if p.token.Type == lexer.INT {
			intValue, _ := strconv.Atoi(p.token.Literal)
			return IntNode{lexer.INT_NODE, intValue}

		} else if p.token.Type == lexer.STRING {
			return StringNode{lexer.STRING_NODE, p.token.Literal}

		} else if p.token.Type == lexer.LPAREN {
			p.advance()
			expr := p.ParseArith()
			if p.token.Type == lexer.RPAREN {
				return expr
			}

		} else if p.token.Type == lexer.SUB {
			p.advance()
			return UnaryOpNode{lexer.UNARY_NODE, lexer.SUB, p.ParseFactor()}
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
