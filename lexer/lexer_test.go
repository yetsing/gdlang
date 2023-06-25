package lexer

import (
	"fmt"
	"strconv"
	"testing"

	"weilang/token"
)

func TestNextToken(t *testing.T) {
	input := `
var five = 5;
con ten = 10;

con add = fn(x, y) {
  x + y;
};

var result = add(five, ten);
!-/*5;
5 < 10 > 5;
5 <= 10 >= 5;

if (5 < 10) {
	return true;
} else {
	return false;
}

10 == 10;
10 != 9;
"foobar"
"foo bar"
[1, 2];
{"foo": "bar"}
a.foo
1 % 4
// 这是一段注释
'world'
'hello world'
 var 中文变量名 = "a开发b"
"中文字符串"
'转义\'\"\a\b\f\n\r\t\v\000\x00\xFF'
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
		expectedStart   token.Position
		expectedEnd     token.Position
	}{
		{
			token.VAR, "var",
			token.Position{Line: 1},
			token.Position{Line: 1, Column: 3},
		},
		{
			token.IDENT, "five",
			token.Position{Line: 1, Column: 4},
			token.Position{Line: 1, Column: 8},
		},
		{
			token.ASSIGN, "=",
			token.Position{Line: 1, Column: 9},
			token.Position{Line: 1, Column: 10},
		},
		{
			token.INT, "5",
			token.Position{Line: 1, Column: 11},
			token.Position{Line: 1, Column: 12},
		},
		{
			token.SEMICOLON, ";",
			token.Position{Line: 1, Column: 12},
			token.Position{Line: 1, Column: 13},
		},
		{
			token.CON, "con",
			token.Position{Line: 2},
			token.Position{Line: 2, Column: 3},
		},
		{
			token.IDENT, "ten",
			token.Position{Line: 2, Column: 4},
			token.Position{Line: 2, Column: 7},
		},
		{
			token.ASSIGN, "=",
			token.Position{Line: 2, Column: 8},
			token.Position{Line: 2, Column: 9},
		},
		{
			token.INT, "10",
			token.Position{Line: 2, Column: 10},
			token.Position{Line: 2, Column: 12},
		},
		{
			token.SEMICOLON, ";",
			token.Position{Line: 2, Column: 12},
			token.Position{Line: 2, Column: 13},
		},
		{
			token.CON, "con",
			token.Position{Line: 4},
			token.Position{Line: 4, Column: 3},
		},
		{
			token.IDENT, "add",
			token.Position{Line: 4, Column: 4},
			token.Position{Line: 4, Column: 7},
		},
		{
			token.ASSIGN, "=",
			token.Position{Line: 4, Column: 8},
			token.Position{Line: 4, Column: 9},
		},
		{
			token.FUNCTION, "fn",
			token.Position{Line: 4, Column: 10},
			token.Position{Line: 4, Column: 12},
		},
		{
			token.LPAREN, "(",
			token.Position{Line: 4, Column: 12},
			token.Position{Line: 4, Column: 13},
		},
		{
			token.IDENT, "x",
			token.Position{Line: 4, Column: 13},
			token.Position{Line: 4, Column: 14},
		},
		{
			token.COMMA, ",",
			token.Position{Line: 4, Column: 14},
			token.Position{Line: 4, Column: 15},
		},
		{
			token.IDENT, "y",
			token.Position{Line: 4, Column: 16},
			token.Position{Line: 4, Column: 17},
		},
		{
			token.RPAREN, ")",
			token.Position{Line: 4, Column: 17},
			token.Position{Line: 4, Column: 18},
		},
		{
			token.LBRACE, "{",
			token.Position{Line: 4, Column: 19},
			token.Position{Line: 4, Column: 20},
		},
		{
			token.IDENT, "x",
			token.Position{Line: 5, Column: 2},
			token.Position{Line: 5, Column: 3},
		},
		{
			token.PLUS, "+",
			token.Position{Line: 5, Column: 4},
			token.Position{Line: 5, Column: 5},
		},
		{
			token.IDENT, "y",
			token.Position{Line: 5, Column: 6},
			token.Position{Line: 5, Column: 7},
		},
		{
			token.SEMICOLON, ";",
			token.Position{Line: 5, Column: 7},
			token.Position{Line: 5, Column: 8},
		},
		{
			token.RBRACE, "}",
			token.Position{Line: 6},
			token.Position{Line: 6, Column: 1},
		},
		{
			token.SEMICOLON, ";",
			token.Position{Line: 6, Column: 1},
			token.Position{Line: 6, Column: 2},
		},
		{expectedType: token.VAR, expectedLiteral: "var"},
		{expectedType: token.IDENT, expectedLiteral: "result"},
		{expectedType: token.ASSIGN, expectedLiteral: "="},
		{expectedType: token.IDENT, expectedLiteral: "add"},
		{expectedType: token.LPAREN, expectedLiteral: "("},
		{expectedType: token.IDENT, expectedLiteral: "five"},
		{expectedType: token.COMMA, expectedLiteral: ","},
		{expectedType: token.IDENT, expectedLiteral: "ten"},
		{expectedType: token.RPAREN, expectedLiteral: ")"},
		{expectedType: token.SEMICOLON, expectedLiteral: ";"},
		{expectedType: token.BANG, expectedLiteral: "!"},
		{expectedType: token.MINUS, expectedLiteral: "-"},
		{expectedType: token.SLASH, expectedLiteral: "/"},
		{expectedType: token.ASTERISK, expectedLiteral: "*"},
		{expectedType: token.INT, expectedLiteral: "5"},
		{expectedType: token.SEMICOLON, expectedLiteral: ";"},
		{expectedType: token.INT, expectedLiteral: "5"},
		{expectedType: token.LESS_THAN, expectedLiteral: "<"},
		{expectedType: token.INT, expectedLiteral: "10"},
		{expectedType: token.GREAT_THAN, expectedLiteral: ">"},
		{expectedType: token.INT, expectedLiteral: "5"},
		{expectedType: token.SEMICOLON, expectedLiteral: ";"},
		{expectedType: token.INT, expectedLiteral: "5"},
		{expectedType: token.LESS_EQUAL_THAN, expectedLiteral: "<="},
		{expectedType: token.INT, expectedLiteral: "10"},
		{expectedType: token.GREAT_EQUAL_THAN, expectedLiteral: ">="},
		{expectedType: token.INT, expectedLiteral: "5"},
		{expectedType: token.SEMICOLON, expectedLiteral: ";"},
		{expectedType: token.IF, expectedLiteral: "if"},
		{expectedType: token.LPAREN, expectedLiteral: "("},
		{expectedType: token.INT, expectedLiteral: "5"},
		{expectedType: token.LESS_THAN, expectedLiteral: "<"},
		{expectedType: token.INT, expectedLiteral: "10"},
		{expectedType: token.RPAREN, expectedLiteral: ")"},
		{expectedType: token.LBRACE, expectedLiteral: "{"},
		{expectedType: token.RETURN, expectedLiteral: "return"},
		{expectedType: token.TRUE, expectedLiteral: "true"},
		{expectedType: token.SEMICOLON, expectedLiteral: ";"},
		{expectedType: token.RBRACE, expectedLiteral: "}"},
		{expectedType: token.ELSE, expectedLiteral: "else"},
		{expectedType: token.LBRACE, expectedLiteral: "{"},
		{expectedType: token.RETURN, expectedLiteral: "return"},
		{expectedType: token.FALSE, expectedLiteral: "false"},
		{expectedType: token.SEMICOLON, expectedLiteral: ";"},
		{expectedType: token.RBRACE, expectedLiteral: "}"},
		{expectedType: token.INT, expectedLiteral: "10"},
		{expectedType: token.EQ, expectedLiteral: "=="},
		{expectedType: token.INT, expectedLiteral: "10"},
		{expectedType: token.SEMICOLON, expectedLiteral: ";"},
		{expectedType: token.INT, expectedLiteral: "10"},
		{expectedType: token.NOT_EQ, expectedLiteral: "!="},
		{expectedType: token.INT, expectedLiteral: "9"},
		{expectedType: token.SEMICOLON, expectedLiteral: ";"},
		{expectedType: token.STRING, expectedLiteral: "foobar"},
		{expectedType: token.STRING, expectedLiteral: "foo bar"},
		{expectedType: token.LBRACKET, expectedLiteral: "["},
		{expectedType: token.INT, expectedLiteral: "1"},
		{expectedType: token.COMMA, expectedLiteral: ","},
		{expectedType: token.INT, expectedLiteral: "2"},
		{expectedType: token.RBRACKET, expectedLiteral: "]"},
		{expectedType: token.SEMICOLON, expectedLiteral: ";"},
		{expectedType: token.LBRACE, expectedLiteral: "{"},
		{expectedType: token.STRING, expectedLiteral: "foo"},
		{expectedType: token.COLON, expectedLiteral: ":"},
		{expectedType: token.STRING, expectedLiteral: "bar"},
		{expectedType: token.RBRACE, expectedLiteral: "}"},
		{expectedType: token.IDENT, expectedLiteral: "a"},
		{expectedType: token.DOT, expectedLiteral: "."},
		{expectedType: token.IDENT, expectedLiteral: "foo"},
		{expectedType: token.INT, expectedLiteral: "1"},
		{
			expectedType: token.MODULO, expectedLiteral: "%",
			expectedStart: token.Position{Line: 26, Column: 2},
			expectedEnd:   token.Position{Line: 26, Column: 3},
		},
		{
			expectedType: token.INT, expectedLiteral: "4",
			expectedStart: token.Position{Line: 26, Column: 4},
			expectedEnd:   token.Position{Line: 26, Column: 5},
		},
		{
			expectedType: token.COMMENT, expectedLiteral: " 这是一段注释",
			expectedStart: token.Position{Line: 27, Column: 2},
			expectedEnd:   token.Position{Line: 27, Column: 9},
		},
		{expectedType: token.STRING, expectedLiteral: "world"},
		{
			expectedType: token.STRING, expectedLiteral: "hello world",
			expectedStart: token.Position{Line: 29, Column: 0},
			expectedEnd:   token.Position{Line: 29, Column: 13},
		},
		{
			expectedType: token.VAR, expectedLiteral: "var",
			expectedStart: token.Position{Line: 30, Column: 1},
			expectedEnd:   token.Position{Line: 30, Column: 4},
		},
		{
			expectedType: token.IDENT, expectedLiteral: "中文变量名",
			expectedStart: token.Position{Line: 30, Column: 5},
			expectedEnd:   token.Position{Line: 30, Column: 10},
		},
		{
			expectedType: token.ASSIGN, expectedLiteral: "=",
			expectedStart: token.Position{Line: 30, Column: 11},
			expectedEnd:   token.Position{Line: 30, Column: 12},
		},
		{
			expectedType: token.STRING, expectedLiteral: "a开发b",
			expectedStart: token.Position{Line: 30, Column: 13},
			expectedEnd:   token.Position{Line: 30, Column: 19},
		},
		{
			expectedType: token.STRING, expectedLiteral: "中文字符串",
			expectedStart: token.Position{Line: 31, Column: 0},
			expectedEnd:   token.Position{Line: 31, Column: 7},
		},
		{
			expectedType: token.STRING,
			// go 大于 \x7f 的转义都会变成 65533
			expectedLiteral: "转义'\"\a\b\f\n\r\t\v\000\x00\u00ff",
		},
		{
			expectedType: token.EOF,
		},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if !tok.TypeIs(tt.expectedType) {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}

		// 方便测试，每个都判断行列太麻烦了
		if tt.expectedStart.IsZero() {
			continue
		}

		if !tok.Start.Equal(&tt.expectedStart) {
			t.Fatalf("tests[%d] - start wrong. expected=%+v, got=%+v",
				i, tt.expectedStart, tok.Start)
		}

		if !tok.End.Equal(&tt.expectedEnd) {
			t.Fatalf("tests[%d] - end wrong. expected=%+v, got=%+v",
				i, tt.expectedEnd, tok.End)
		}
	}
}

func TestIllegalToken(t *testing.T) {
	//	input := `
	//"abc
	//@
	//' 6月21日
	//'abc'
	//`
	tests := []struct {
		input           string
		expectedType    token.TokenType
		expectedLiteral string
		expectedStart   token.Position
		expectedEnd     token.Position
	}{
		{
			`"abc`,
			token.ILLEGAL, "string literal not terminated",
			token.Position{0, 0},
			token.Position{0, 3},
		},
		{
			`@`,
			token.ILLEGAL, "invalid char @",
			token.Position{0, 0},
			token.Position{0, 0},
		},
		{
			`' 6月21日`,
			token.ILLEGAL, "string literal not terminated",
			token.Position{0, 0},
			token.Position{0, 6},
		},
		{
			`'\d'`,
			token.ILLEGAL, "illegal escape sequence",
			token.Position{0, 0},
			token.Position{0, 2},
		},
		{
			`'\1'`,
			token.ILLEGAL, "illegal escape sequence",
			token.Position{0, 0},
			token.Position{0, 2},
		},
		{
			`'\777'`,
			token.ILLEGAL, "illegal escape sequence",
			token.Position{0, 0},
			token.Position{0, 2},
		},
		{
			`'\x'`,
			token.ILLEGAL, "illegal escape sequence",
			token.Position{0, 0},
			token.Position{0, 3},
		},
		{
			`'\x1'`,
			token.ILLEGAL, "illegal escape sequence",
			token.Position{0, 0},
			token.Position{0, 3},
		},
		{
			`'\xgg'`,
			token.ILLEGAL, "illegal escape sequence",
			token.Position{0, 0},
			token.Position{0, 3},
		},
	}

	n, err := strconv.ParseInt("8a", 16, 32)
	fmt.Printf("n: %d, err: %v\n", n, err)
	fmt.Printf("n: %s, err: %v\n", "\xff", nil)

	//tmp := `\t\123`
	//tmp2 := "\t\407"
	//fmt.Printf("tmp: %q, tmp2: %q\n", tmp, tmp2)

	for i, tt := range tests {
		l := New(tt.input)
		tok := l.NextToken()

		if !tok.TypeIs(tt.expectedType) {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q %q",
				i, tt.expectedType, tok.Type, tok.Literal)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}

		if !tok.Start.Equal(&tt.expectedStart) {
			t.Fatalf("tests[%d] - start wrong. expected=%+v, got=%+v",
				i, tt.expectedStart, tok.Start)
		}

		if !tok.End.Equal(&tt.expectedEnd) {
			t.Fatalf("tests[%d] - end wrong. expected=%+v, got=%+v",
				i, tt.expectedEnd, tok.End)
		}
	}
}

func TestUnicodeCategory(t *testing.T) {
	tests := []struct {
		ch               rune
		expectedCategory string
	}{
		{'A', "Lu"},
		{'a', "Ll"},
		{'中', "Lo"},
		{'1', "Nd"},
	}
	for i, tt := range tests {
		got := UnicodeCategory(tt.ch)
		if got != tt.expectedCategory {
			t.Fatalf("tests[%d] - category wrong. expected=%v, got=%v",
				i, tt.expectedCategory, got)
		}
	}
}
