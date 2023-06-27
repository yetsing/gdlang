package token

type TokenType string

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + literals
	IDENT  = "IDENT"  // add, foobar, x, y, ...
	INT    = "INT"    // 1343456
	STRING = "STRING" // "foobar"
	// COMMENT 注释
	COMMENT = "comment"

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"
	MODULO   = "%"
	DOT      = "."

	LESS_THAN        = "<"
	LESS_EQUAL_THAN  = "<="
	GREAT_THAN       = ">"
	GREAT_EQUAL_THAN = ">="

	EQ     = "=="
	NOT_EQ = "!="

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"

	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	// 位操作
	LEFT_SHIFT  = "<<"
	RIGHT_SHIFT = ">>"
	BITWISE_AND = "&"
	BITWISE_XOR = "^"
	BITWISE_OR  = "|"
	BITWISE_NOT = "~"

	// Keywords
	FUNCTION = "FUNCTION"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
	VAR      = "var"
	CON      = "con"
	NULL     = "null"
	WHILE    = "while"
	CONTINUE = "continue"
	BREAK    = "break"
)

type Position struct {
	Line   int
	Column int
}

func (p *Position) Equal(other *Position) bool {
	return p.Line == other.Line && p.Column == other.Column
}

func (p *Position) IsZero() bool {
	return p.Line == 0 && p.Column == 0
}

type Token struct {
	Type    TokenType
	Literal string
	Start   Position
	End     Position
}

func (t *Token) TypeIs(ttype TokenType) bool {
	return t.Type == ttype
}

func (t *Token) TypeIn(ttypes ...TokenType) bool {
	for _, ttype := range ttypes {
		if t.TypeIs(ttype) {
			return true
		}
	}
	return false
}

var keywords = map[string]TokenType{
	"fn":       FUNCTION,
	"var":      VAR,
	"con":      CON,
	"true":     TRUE,
	"false":    FALSE,
	"if":       IF,
	"else":     ELSE,
	"return":   RETURN,
	"null":     NULL,
	"while":    WHILE,
	"continue": CONTINUE,
	"break":    BREAK,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
