package lexer

type Lexer struct {
	program      string
	position     int
	readPosition int
	ch           byte
}

func NewLexer(program string) *Lexer {
	l := &Lexer{program: program}
	l.readChar()
	return l
}

func (l *Lexer) Lex() []Token {
	var tokens []Token
	for l.ch != 0 {
		tok := l.NextToken()
		tokens = append(tokens, tok)
	}

	tokens = append(tokens, Token{Type: EOF, Literal: ""})
	return tokens
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.program) {
		l.ch = 0
	} else {
		l.ch = l.program[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.program) {
		return 0
	} else {
		return l.program[l.readPosition]
	}
}

func (l *Lexer) NextToken() Token {
	var tok Token
	l.eatWhitespace()
	switch l.ch {
	case '+':
		tok = NewToken(ADD, l.ch)
	case '-':
		tok = l.readDouble(SUB, '>', ARROW)
	case '*':
		tok = NewToken(MUL, l.ch)
	case '/':
		tok = NewToken(DIV, l.ch)
	case '(':
		tok = NewToken(LPAREN, l.ch)
	case ')':
		tok = NewToken(RPAREN, l.ch)
	case '{':
		tok = NewToken(LBRACE, l.ch)
	case '}':
		tok = NewToken(RBRACE, l.ch)
	case '?':
		tok = NewToken(QUESTION, l.ch)
	case ';':
		tok = NewToken(SEMICOLON, l.ch)
	case ',':
		tok = NewToken(COMMA, l.ch)
	case ':':
		tok = l.readDouble(COLON, '=', ASSIGN)
	case '=':
		tok = l.readDouble(EQ, '=', EE)
	case '>':
		tok = l.readDouble(GT, '=', GTE)
	case '<':
		tok = l.readDouble(LT, '=', LTE)
	case '!':
		tok = l.readDouble(NOT, '=', NE)
	case '"':
		tok.Literal = l.readString()
		tok.Type = STRING
	case 0:
		tok = NewToken(EOF, l.ch)
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = lookupIdentifier(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = NewToken(ILLEGAL, l.ch)
		}
	}
	l.readChar()
	return tok
}

func (l *Lexer) eatWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.program[position:l.position]
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}

	return l.program[position:l.position]
}

func (l *Lexer) readString() string {
	l.readChar()
	position := l.position

	for l.ch != '"' {
		l.readChar()
	}

	return l.program[position:l.position]
}

func (l *Lexer) readDouble(firstType string, second byte, secondType string) Token {
	ch := l.ch
	if l.peekChar() == second {
		l.readChar()
		return Token{Type: secondType, Literal: string(ch) + string(l.ch)}
	} else {
		return Token{Type: firstType, Literal: string(ch)}
	}
}
