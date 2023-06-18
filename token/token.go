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

type Token struct {
	Type    TokenType
	Literal string
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
