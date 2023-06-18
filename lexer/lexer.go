package lexer

import "weilang/token"

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           rune // current char under examination
	// utf8 列表
	ucodes []rune
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.setup()
	l.readChar()
	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekCharIs('=') {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.EQ, Literal: literal}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '!':
		if l.peekCharIs('=') {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.NOT_EQ, Literal: literal}
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '/':
		if l.peekCharIs('/') {
			l.readChar()
			tok.Literal = l.readComment()
			tok.Type = token.COMMENT
			return tok
		} else {
			tok = newToken(token.SLASH, l.ch)
		}
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '%':
		tok = newToken(token.MODULO, l.ch)
	case '<':
		if l.peekCharIs('=') {
			l.readChar()
			tok = token.Token{Type: token.LESS_EQUAL_THAN, Literal: "<="}
		} else {
			tok = newToken(token.LESS_THAN, l.ch)
		}
	case '>':
		if l.peekCharIs('=') {
			l.readChar()
			tok = token.Token{Type: token.GREAT_EQUAL_THAN, Literal: ">="}
		} else {
			tok = newToken(token.GREAT_THAN, l.ch)
		}
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case ':':
		tok = newToken(token.COLON, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString(l.ch)
	case '\'':
		tok.Type = token.STRING
		tok.Literal = l.readString(l.ch)
	case '[':
		tok = newToken(token.LBRACKET, l.ch)
	case ']':
		tok = newToken(token.RBRACKET, l.ch)
	case '.':
		tok = newToken(token.DOT, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isIdentifier(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) setup() {
	for _, u := range l.input {
		l.ucodes = append(l.ucodes, u)
	}
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.ucodes) {
		l.ch = 0
	} else {
		l.ch = l.ucodes[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) peekCharIs(ch rune) bool {
	if l.readPosition >= len(l.ucodes) {
		return 0 == ch
	} else {
		return l.ucodes[l.readPosition] == ch
	}
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isIdentifier(l.ch) {
		l.readChar()
	}
	return string(l.ucodes[position:l.position])
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return string(l.ucodes[position:l.position])
}

func (l *Lexer) readString(end rune) string {
	// todo 判断引号是否成对
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == end || l.ch == 0 {
			break
		}
	}
	return string(l.ucodes[position:l.position])
}

func (l *Lexer) readComment() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '\n' || l.ch == 0 {
			break
		}
	}
	return string(l.ucodes[position:l.position])
}

func isIdentifier(ch rune) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

func newToken(tokenType token.TokenType, ch rune) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}
