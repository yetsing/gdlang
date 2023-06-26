package parser

import (
	"testing"
	"weilang/ast"
	"weilang/lexer"
)

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.New(input)
	p := New(l)
	program, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp not *ast.Identifier. got=%T", stmt.Expression)
	}
	if ident.Value != "foobar" {
		t.Errorf("ident.Value not %s. got=%s", "foobar", ident.Value)
	}
	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral not %s. got=%s", "foobar",
			ident.TokenLiteral())
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue int64
	}{
		{"5;", 5},
		{"0b1111_1111", 255},
		{"0o01", 1},
		{"0o3_77", 255},
		{"0x0100", 256},
	}

	for i, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program, err := p.ParseProgram()
		if err != nil {
			t.Fatalf("[test %d]syntax error: %s", i, err)
		}

		if len(program.Statements) != 1 {
			t.Fatalf("[test %d]program has not enough statements. got=%d",
				i, len(program.Statements))
		}
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("[test %d]program.Statements[0] is not ast.ExpressionStatement. got=%T",
				i, program.Statements[0])
		}

		literal, ok := stmt.Expression.(*ast.IntegerLiteral)
		if !ok {
			t.Fatalf("[test %d]exp not *ast.IntegerLiteral. got=%T", i, stmt.Expression)
		}
		if literal.Value != tt.expectedValue {
			t.Errorf("[test %d]literal.Value not %d. got=%d", i, tt.expectedValue, literal.Value)
		}
	}
}

func TestStringLiteralExpression(t *testing.T) {
	tests := []struct {
		input           string
		expectedLiteral string
	}{
		{`"hello world"`, "hello world"},
		{`'hello world2';`, "hello world2"},
		{"`hello\nworld3`", "hello\nworld3"},
	}
	for i, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program, err := p.ParseProgram()
		if err != nil {
			t.Fatalf("[test %d]syntax error: %s", i, err)
		}

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		literal, ok := stmt.Expression.(*ast.StringLiteral)
		if !ok {
			t.Fatalf("[test %d]exp not *ast.StringLiteral. got=%T", i, stmt.Expression)
		}

		if literal.Value != tt.expectedLiteral {
			t.Errorf("[test %d]literal.Value not %q. got=%q", i, tt.expectedLiteral, literal.Value)
		}
	}

}
