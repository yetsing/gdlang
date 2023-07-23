package parser

import (
	"fmt"
	"testing"
	"weilang/ast"
	"weilang/lexer"
)

func TestSyntaxError(t *testing.T) {
	tests := []struct {
		input string
	}{
		{"var \na= 1"},
		{"var a\n= 1"},
		{"var a=\n 1"},
		{"con \nb= a"},
		{"con b\n= a"},
		{"con b=\n a"},
		{"a \n= 1"},
		{"a =\n 1"},
		{"a =\n 1"},
		{"a and\n 1"},
		{"if \n(a) {a}"},
		{"if (a) {a} \n else \n if (b) {b}"},
		{"while \n(a) {a}"},
		{"fn \n(a) {a}"},
		{"fn(a) {a} var b = 2"},
		{"if(a) {a} var b = 2"},
		{"if(a) {a} else {} var b = 2"},
		{"while(a) {a} var b = 2"},
		{"continue"},
		{"break"},
		{"var a = 1; break; a = 2"},
		{"var a = 1; continue; a = 2"},
		{`
while (1) {
  var foo = fn() {
    continue
  }
}
`},
		{
			`wei`,
		},
		{
			`wei.`,
		},
		{
			`wei.ddd d`,
		},
		{
			`wei."export"`,
		},
		{
			`class Foo`,
		},
		{
			`class Foo {} 
var a =`,
		},
		{
			`class Foo {
var a var b
} 
`,
		},
		{
			`class Foo {
var a
// some
fn init() {} con d
} 
`,
		},
		{
			`class Foo {
var class.d
} 
`,
		},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		_, err := p.ParseProgram()
		if err == nil {
			t.Errorf("expected error\n%s", tt.input)
		}
		//t.Logf("%v", err)
	}
}

func testVarStatement(t *testing.T, s ast.Statement, name string) bool {
	t.Helper()
	if s.TokenLiteral() != "var" {
		t.Errorf("s.TokenLiteral not 'var'. got=%q", s.TokenLiteral())
		return false
	}

	varStmt, ok := s.(*ast.VarStatement)
	if !ok {
		t.Errorf("s not *ast.VarStatement. got=%T", s)
		return false
	}

	if varStmt.Name.Value != name {
		t.Errorf("varStmt.Filename.Value not '%s'. got=%s", name, varStmt.Name.Value)
		return false
	}

	if varStmt.Name.TokenLiteral() != name {
		t.Errorf("varStmt.Filename.TokenLiteral() not '%s'. got=%s",
			name, varStmt.Name.TokenLiteral())
		return false
	}

	return true
}

func testConStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "con" {
		t.Errorf("s.TokenLiteral not 'var'. got=%q", s.TokenLiteral())
		return false
	}

	stmt, ok := s.(*ast.ConStatement)
	if !ok {
		t.Errorf("s not *ast.ConStatement. got=%T", s)
		return false
	}

	if stmt.Name.Value != name {
		t.Errorf("stmt.Filename.Value not '%s'. got=%s", name, stmt.Name.Value)
		return false
	}

	if stmt.Name.TokenLiteral() != name {
		t.Errorf("stmt.Filename.TokenLiteral() not '%s'. got=%s",
			name, stmt.Name.TokenLiteral())
		return false
	}

	return true
}

func testAssignStatement(t *testing.T, s ast.Statement, name string) bool {
	t.Helper()
	stmt, ok := s.(*ast.AssignStatement)
	if !ok {
		t.Errorf("s not *ast.AssignStatement. got=%T", s)
		return false
	}

	if stmt.Left.String() != name {
		t.Errorf("stmt.Left not '%s'. got=%s", name, stmt.Left.String())
		return false
	}

	return true
}

func testBinaryOpExpression(
	t *testing.T, exp ast.Expression, left interface{},
	operator string, right interface{},
) bool {
	opExp, ok := exp.(*ast.BinaryOpExpression)
	if !ok {
		t.Errorf("exp is not ast.InfixExpression. got=%T(%s)", exp, exp)
		return false
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, opExp.Operator)
		return false
	}

	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}

	return true
}

func testLiteralExpression(
	t *testing.T,
	exp ast.Expression,
	expected interface{},
) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		if v == "null" {
			return testNullLiteral(t, exp, v)
		}
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}

	if integ.Value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
		return false
	}

	return true
}

//goland:noinspection GoUnusedFunction
func testStringLiteral(t *testing.T, exp ast.Expression, value string) bool {
	sl, ok := exp.(*ast.StringLiteral)
	if !ok {
		t.Errorf("exp not *ast.StringLiteral. got=%T", exp)
		return false
	}

	if sl.Value != value {
		t.Errorf("stringLiteral.Value not %s. got=%s", value, sl.Value)
		return false
	}
	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s. got=%s", value,
			ident.TokenLiteral())
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	bo, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("exp not *ast.Boolean. got=%T", exp)
		return false
	}

	if bo.Value != value {
		t.Errorf("bo.Value not %t. got=%t", value, bo.Value)
		return false
	}

	if bo.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("bo.TokenLiteral not %t. got=%s",
			value, bo.TokenLiteral())
		return false
	}

	return true
}

func testNullLiteral(t *testing.T, exp ast.Expression, value string) bool {
	bo, ok := exp.(*ast.NullLiteral)
	if !ok {
		t.Errorf("exp not *ast.Boolean. got=%T", exp)
		return false
	}

	if bo.TokenLiteral() != value {
		t.Errorf("bo.TokenLiteral not %s. got=%s",
			value, bo.TokenLiteral())
		return false
	}

	return true
}
