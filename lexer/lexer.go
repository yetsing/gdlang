package lexer

import (
	"errors"
	"strconv"
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
			ch := l.ch
			l.readChar()
			tok := l.buildToken(token.ILLEGAL)
			tok.Literal = "invalid char " + string(ch)
			return tok
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
		l.index++
		l.ch = l.ucodes[l.index]
		l.position.Column++
	}
}

func (l *Lexer) advance(n int) {
	for i := 0; i < n; i++ {
		l.readChar()
	}
}

func (l *Lexer) getString(n int) (string, error) {
	if l.index+n > len(l.ucodes) {
		return "", errors.New("not enough char")
	}
	s := string(l.ucodes[l.index : l.index+n])
	return s, nil
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

var escapeMap = map[rune]rune{
	'\\': '\\',
	'\'': '\'',
	'"':  '"',
	'a':  '\a',
	'b':  '\b',
	'f':  '\f',
	'n':  '\n',
	'r':  '\r',
	't':  '\t',
	'v':  '\v',
}

func parseRune(s string, base int, bitSize int) (rune, error) {
	n, err := strconv.ParseInt(s, base, bitSize)
	if err != nil {
		return 0, err
	}
	r := rune(n)
	return r, nil
}

func (l *Lexer) readString(end rune) token.Token {
	var buf []rune
	// 跳过开始的引号
	l.readChar()
	for {
		// 参考 Python 的转义字符 https://docs.python.org/3/reference/lexical_analysis.html#string-and-bytes-literals
		// 处理转义字符
		ch := l.ch
		if ch == '\\' {
			l.readChar()
			if actual, ok := escapeMap[l.ch]; ok {
				buf = append(buf, actual)
				l.readChar()
				continue
			}
			// 解析 Unicode 转义字符
			var ucode rune
			switch l.ch {
			case 'x':
			case 'u':
			case 'U':
			case '0', '1', '2', '3', '4', '5', '6', '7':
				// 解析八进制转义
				//"\ooo" o 代表八进制字符，最大为 "\377" (255)
				s, err := l.getString(3)
				if err != nil {
					tok := l.buildToken(token.ILLEGAL)
					tok.Literal = "unknown escape sequence"
					return tok
				}
				ucode, err = parseRune(s, 8, 8)
				if err != nil {
					tok := l.buildToken(token.ILLEGAL)
					tok.Literal = "escape sequence is invalid Unicode code point"
					return tok
				}
				l.advance(3)
			default:
				tok := l.buildToken(token.ILLEGAL)
				tok.Literal = "unknown escape sequence"
				return tok
			}
			// 不是转义字符
			buf = append(buf, ucode)
			continue
		}

		if l.ch == end {
			break
		}
		if l.ch == 0 || l.ch == '\n' {
			tok := l.buildToken(token.ILLEGAL)
			tok.Literal = "string literal not terminated"
			return tok
		}
		buf = append(buf, l.ch)
		l.readChar()
	}
	// 跳过末尾的引号
	l.readChar()
	tok := l.buildToken(token.STRING)
	tok.Literal = string(buf)
	return tok
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
