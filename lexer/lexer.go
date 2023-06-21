package lexer

import (
	"unicode"
	"weilang/token"
)

type Lexer struct {
	input string
	// current index of ch in ucodes
	index int
	// current char
	ch rune
	// unicode 列表
	ucodes   []rune
	position token.Position
	// 标记索引和位置，方便计算 Token 的 start end
	markIndex    int
	markPosition token.Position
}

func New(input string) *Lexer {
	l := &Lexer{input: input, index: -1}
	l.init()
	l.readChar()
	return l
}

func (l *Lexer) NextToken() token.Token {
	var ttype token.TokenType

	l.skipWhitespace()
	l.mark()

	switch l.ch {
	case '=':
		if l.peekCharIs('=') {
			l.readChar()
			ttype = token.EQ
		} else {
			ttype = token.ASSIGN
		}
	case '+':
		ttype = token.PLUS
	case '-':
		ttype = token.MINUS
	case '!':
		if l.peekCharIs('=') {
			l.readChar()
			ttype = token.NOT_EQ
		} else {
			ttype = token.BANG
		}
	case '/':
		if l.peekCharIs('/') {
			return l.readComment()
		} else {
			ttype = token.SLASH
		}
	case '*':
		ttype = token.ASTERISK
	case '%':
		ttype = token.MODULO
	case '<':
		if l.peekCharIs('=') {
			l.readChar()
			ttype = token.LESS_EQUAL_THAN
		} else {
			ttype = token.LESS_THAN
		}
	case '>':
		if l.peekCharIs('=') {
			l.readChar()
			ttype = token.GREAT_EQUAL_THAN
		} else {
			ttype = token.GREAT_THAN
		}
	case ';':
		ttype = token.SEMICOLON
	case ':':
		ttype = token.COLON
	case ',':
		ttype = token.COMMA
	case '{':
		ttype = token.LBRACE
	case '}':
		ttype = token.RBRACE
	case '(':
		ttype = token.LPAREN
	case ')':
		ttype = token.RPAREN
	case '"':
		return l.readString(l.ch)
	case '\'':
		return l.readString(l.ch)
	case '[':
		ttype = token.LBRACKET
	case ']':
		ttype = token.RBRACKET
	case '.':
		ttype = token.DOT
	case 0:
		ttype = token.EOF
	default:
		if isIdentifierStart(l.ch) {
			return l.readIdentifier()
		} else if isDigit(l.ch) {
			return l.readNumber()
		} else {
			ttype = token.ILLEGAL
		}
	}

	l.readChar()
	return l.buildToken(ttype)
}

func (l *Lexer) init() {
	for _, u := range l.input {
		l.ucodes = append(l.ucodes, u)
	}
	l.position.Line = 0
	l.position.Column = -1
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readChar() {
	if l.index >= len(l.ucodes)-1 {
		l.ch = 0
	} else {
		if l.ch == '\n' {
			l.position.Line++
			l.position.Column = -1
		}
		l.index += 1
		l.ch = l.ucodes[l.index]
		l.position.Column++
	}
}

func (l *Lexer) peekCharIs(ch rune) bool {
	nextIndex := l.index + 1
	if nextIndex >= len(l.ucodes) {
		return 0 == ch
	} else {
		return l.ucodes[nextIndex] == ch
	}
}

// 标记一个位置
func (l *Lexer) mark() {
	l.markIndex = l.index
	l.markPosition.Line = l.position.Line
	l.markPosition.Column = l.position.Column
}

func (l *Lexer) buildToken(ttype token.TokenType) token.Token {
	start := l.markPosition
	end := l.position
	startIndex := l.markIndex
	endIndex := l.index
	switch ttype {
	case token.STRING:
		// 移除首尾的引号
		startIndex++
		endIndex--
	case token.EOF:
		start.Line = -1
		start.Column = -1
		end.Line = -1
		end.Column = -1
	}
	tok := token.Token{Type: ttype, Literal: string(l.ucodes[startIndex:endIndex]), Start: start, End: end}
	return tok
}

func (l *Lexer) readIdentifier() token.Token {
	for isIdentifierContinue(l.ch) {
		l.readChar()
	}
	tok := l.buildToken(token.IDENT)
	tok.Type = token.LookupIdent(tok.Literal)
	return tok
}

func (l *Lexer) readNumber() token.Token {
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.buildToken(token.INT)
}

func (l *Lexer) readString(end rune) token.Token {
	// todo 判断引号是否成对
	// 跳过开始的引号
	l.readChar()
	for {
		l.readChar()
		if l.ch == end || l.ch == 0 {
			break
		}
	}
	// 跳过末尾的引号
	l.readChar()
	return l.buildToken(token.STRING)
}

func (l *Lexer) readComment() token.Token {
	// 跳过开头的 // 两个字符
	l.readChar()
	l.readChar()
	l.mark()
	for {
		if l.ch == '\n' || l.ch == 0 {
			break
		}
		l.readChar()
	}
	return l.buildToken(token.COMMENT)
}

// 参考 Python 的规则 https://docs.python.org/3/reference/lexical_analysis.html#identifiers
var idStartCategorys = map[string]uint8{
	"Lu": 1,
	"Ll": 1,
	"Lm": 1,
	"Lt": 1,
	"Lo": 1,
	"Nl": 1,
}
var idContinueCategorys = map[string]uint8{
	"Lu": 1,
	"Ll": 1,
	"Lm": 1,
	"Lt": 1,
	"Lo": 1,
	"Nl": 1,
	"Mn": 1,
	"Mc": 1,
	"Nd": 1,
	"Pc": 1,
}

func isIdentifierStart(ch rune) bool {

	switch ch {
	case '_':
		return true
	default:
		cat := UnicodeCategory(ch)
		_, ok := idStartCategorys[cat]
		return ok
	}
}

func isIdentifierContinue(ch rune) bool {
	switch ch {
	case '_':
		return true
	default:
		cat := UnicodeCategory(ch)
		_, ok := idContinueCategorys[cat]
		return ok
	}
}

func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

// UnicodeCategory returns the Unicode Character Category of the given rune.
// code from https://stackoverflow.com/a/53507592
func UnicodeCategory(r rune) string {
	for name, table := range unicode.Categories {
		if len(name) == 2 && unicode.Is(table, r) {
			return name
		}
	}
	return "Cn"
}
