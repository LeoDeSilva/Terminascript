package lexer

type Token struct {
	Type    string
	Literal string
}

func NewToken(tokenType string, ch byte) Token {
	return Token{Type: tokenType, Literal: string(ch)}
}

func lookupIdentifier(identifier string) string {
	if token, ok := keywords[identifier]; ok {
		return token
	}
	return IDENTIFIER
}

var keywords = map[string]string{
	"if":    IF,
	"while": WHILE,
	"for":   FOR,
	"func":  FUNC,
	"let":   LET,
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	IDENTIFIER = "IDENTIFIER"
	STRING     = "STRING"
	INT        = "INT"

	ADD = "ADD"
	SUB = "SUB"
	MUL = "MUL"
	DIV = "DIV"

	EE  = "EE"
	EQ  = "EQ"
	NOT = "NOT"
	NE  = "NE"

	LT  = "LT"
	GT  = "GT"
	LTE = "LTE"
	GTE = "GTE"

	QUESTION  = "QUESTION"
	COLON     = "COLON"
	COMMA     = "COMMA"
	SEMICOLON = "SEMICOLON"
	ASSIGN    = "ASSIGN"

	ARROW = "ARROW"

	LPAREN = "LPAREN"
	RPAREN = "RPAREN"
	LBRACE = "LBRACE"
	RBRACE = "RBRACE"

	WHILE = "WHILE"
	FOR   = "FOR"
	IF    = "IF"
	FUNC  = "FUNC"
	LET   = "LET"

	PROGRAM_NODE    = "PROGRAM_NODE"
	BIN_OP_NODE     = "BIN_OP_NODE"
	VAR_ACCESS_NODE = "VAR_ACCESS_NODE"
	INT_NODE        = "INT_NODE"
	STRING_NODE     = "STRING_NODE"
	UNARY_NODE      = "UNARY_NODE"
	ERROR_NODE      = "ERROR_NODE"
	FUNC_CALL_NODE  = "FUNC_CALL_NODE"
	PARAMETER_NODE  = "PARAMETER_NODE"
)
