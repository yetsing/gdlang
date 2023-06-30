package parser

import (
	"fmt"
	"testing"
	"weilang/ast"
	"weilang/lexer"
)

func TestVarStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"var x = 5;", "x", 5},
		{"var y = true;", "y", true},
		{"var foobar = y;", "foobar", "y"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program, err := p.ParseProgram()
		if err != nil {
			t.Fatalf("%v", err)
		}

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d",
				len(program.Statements))
		}

		stmt := program.Statements[0]
		if !testVarStatement(t, stmt, tt.expectedIdentifier) {
			return
		}

		val := stmt.(*ast.VarStatement).Value
		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}
	}
}

func TestConStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"con x = 5;", "x", 5},
		{"con y = true;", "y", true},
		{"con foobar = y;", "foobar", "y"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program, err := p.ParseProgram()
		if err != nil {
			t.Fatalf("%v", err)
		}

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d",
				len(program.Statements))
		}

		stmt := program.Statements[0]
		if !testConStatement(t, stmt, tt.expectedIdentifier) {
			return
		}

		val := stmt.(*ast.ConStatement).Value
		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}
	}
}

func TestAssignStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"x = 5;", "x", 5},
		{"y = true;", "y", true},
		{"foobar = y;", "foobar", "y"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program, err := p.ParseProgram()
		if err != nil {
			t.Fatalf("%v", err)
		}

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d",
				len(program.Statements))
		}

		stmt := program.Statements[0]
		if !testAssignStatement(t, stmt, tt.expectedIdentifier) {
			return
		}

		val := stmt.(*ast.AssignStatement).Value
		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"return 5;", 5},
		{"return true;", true},
		{"return foobar;", "foobar"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program, err := p.ParseProgram()
		if err != nil {
			t.Fatalf("%v", err)
		}

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d",
				len(program.Statements))
		}

		stmt := program.Statements[0]
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Fatalf("stmt not *ast.ReturnStatement. got=%T", stmt)
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Fatalf("returnStmt.TokenLiteral not 'return', got %q",
				returnStmt.TokenLiteral())
		}
		if testLiteralExpression(t, returnStmt.ReturnValue, tt.expectedValue) {
			return
		}
	}
}

func TestParsingBinaryOpExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},

		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 % 5;", 5, "%", 5},

		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"5 <= 5;", 5, "<=", 5},
		{"5 >= 5;", 5, ">=", 5},

		{"15 >> 0x15", 15, ">>", 0x15},
		{"15 << 0x15", 15, "<<", 0x15},
		{"15 & 15", 15, "&", 15},
		{"15 ^ 15", 15, "^", 15},
		{"15 | 15", 15, "|", 15},

		{"20 and 20", 20, "and", 20},
		{"20 or 20", 20, "or", 20},

		{"foobar + barfoo;", "foobar", "+", "barfoo"},
		{"foobar - barfoo;", "foobar", "-", "barfoo"},
		{"foobar * barfoo;", "foobar", "*", "barfoo"},
		{"foobar / barfoo;", "foobar", "/", "barfoo"},
		{"foobar % barfoo;", "foobar", "%", "barfoo"},
		{"foobar > barfoo;", "foobar", ">", "barfoo"},
		{"foobar < barfoo;", "foobar", "<", "barfoo"},
		{"foobar == barfoo;", "foobar", "==", "barfoo"},
		{"foobar != barfoo;", "foobar", "!=", "barfoo"},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},

		{"foobar >> barfoo", "foobar", ">>", "barfoo"},
		{"foobar << barfoo", "foobar", "<<", "barfoo"},
		{"foobar & barfoo", "foobar", "&", "barfoo"},
		{"foobar ^ barfoo", "foobar", "^", "barfoo"},
		{"foobar | barfoo", "foobar", "|", "barfoo"},

		{"foobar and barfoo", "foobar", "and", "barfoo"},
		{"foobar or barfoo", "foobar", "or", "barfoo"},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program, err := p.ParseProgram()
		if err != nil {
			t.Fatalf("%s", err)
		}

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		if !testBinaryOpExpression(t, stmt.Expression, tt.leftValue,
			tt.operator, tt.rightValue) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"-a * -b",
			"((-a) * (-b))",
		},
		{
			"not-a",
			"(not(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"(((5 > 4) == 3) < 4)",
		},
		{
			"5 < 4 != 3 > 4",
			"(((5 < 4) != 3) > 4)",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"(5 + 5) * 2 * (5 + 5)",
			"(((5 + 5) * 2) * (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"not(true == true)",
			"(not(true == true))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
		{
			"a * [1, 2, 3, 4][b * c] * d",
			"((a * ([1, 2, 3, 4][(b * c)])) * d)",
		},
		{
			"add(a * b[2], b[1], 2 * [1, 2][1])",
			"add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program, err := p.ParseProgram()
		if err != nil {
			t.Fatalf("%s", err)
		}

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func TestParsingEmptyArrayLiterals(t *testing.T) {
	input := "[]"

	l := lexer.New(input)
	p := New(l)
	program, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("%v", err)
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	array, ok := stmt.Expression.(*ast.ListLiteral)
	if !ok {
		t.Fatalf("exp not ast.ListLiteral. got=%T", stmt.Expression)
	}

	if len(array.Elements) != 0 {
		t.Errorf("len(array.Elements) not 0. got=%d", len(array.Elements))
	}
}

func TestParsingArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

	l := lexer.New(input)
	p := New(l)
	program, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("%v", err)
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	array, ok := stmt.Expression.(*ast.ListLiteral)
	if !ok {
		t.Fatalf("exp not ast.ListLiteral. got=%T", stmt.Expression)
	}

	if len(array.Elements) != 3 {
		t.Fatalf("len(array.Elements) not 3. got=%d", len(array.Elements))
	}

	testIntegerLiteral(t, array.Elements[0], 1)
	testBinaryOpExpression(t, array.Elements[1], 2, "*", 2)
	testBinaryOpExpression(t, array.Elements[2], 3, "+", 3)
}

func TestParsingEmptyHashLiteral(t *testing.T) {
	input := "{}"

	l := lexer.New(input)
	p := New(l)
	program, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("%v", err)
	}

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.DictLiteral)
	if !ok {
		t.Fatalf("exp is not ast.DictLiteral. got=%T", stmt.Expression)
	}

	if len(hash.Pairs) != 0 {
		t.Errorf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
	}
}

func TestParsingHashLiteralsStringKeys(t *testing.T) {
	input := `{"one": 1, "two": 2, "three": 3}`

	l := lexer.New(input)
	p := New(l)
	program, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("%v", err)
	}

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.DictLiteral)
	if !ok {
		t.Fatalf("exp is not ast.DictLiteral. got=%T", stmt.Expression)
	}

	expected := map[string]int64{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	if len(hash.Pairs) != len(expected) {
		t.Errorf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
	}

	for key, value := range hash.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not ast.StringLiteral. got=%T", key)
			continue
		}

		expectedValue := expected[literal.String()]
		testIntegerLiteral(t, value, expectedValue)
	}
}

func TestParsingHashLiteralsBooleanKeys(t *testing.T) {
	input := `{true: 1, false: 2}`

	l := lexer.New(input)
	p := New(l)
	program, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("%v", err)
	}

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.DictLiteral)
	if !ok {
		t.Fatalf("exp is not ast.DictLiteral. got=%T", stmt.Expression)
	}

	expected := map[string]int64{
		"true":  1,
		"false": 2,
	}

	if len(hash.Pairs) != len(expected) {
		t.Errorf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
	}

	for key, value := range hash.Pairs {
		boolean, ok := key.(*ast.Boolean)
		if !ok {
			t.Errorf("key is not ast.BooleanLiteral. got=%T", key)
			continue
		}

		expectedValue := expected[boolean.String()]
		testIntegerLiteral(t, value, expectedValue)
	}
}

func TestParsingHashLiteralsIntegerKeys(t *testing.T) {
	input := `{1: 1, 2: 2, 3: 3}`

	l := lexer.New(input)
	p := New(l)
	program, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("%v", err)
	}

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.DictLiteral)
	if !ok {
		t.Fatalf("exp is not ast.DictLiteral. got=%T", stmt.Expression)
	}

	expected := map[string]int64{
		"1": 1,
		"2": 2,
		"3": 3,
	}

	if len(hash.Pairs) != len(expected) {
		t.Errorf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
	}

	for key, value := range hash.Pairs {
		integer, ok := key.(*ast.IntegerLiteral)
		if !ok {
			t.Errorf("key is not ast.IntegerLiteral. got=%T", key)
			continue
		}

		expectedValue := expected[integer.String()]

		testIntegerLiteral(t, value, expectedValue)
	}
}

func TestParsingHashLiteralsWithExpressions(t *testing.T) {
	input := `{"one": 0 + 1, "two": 10 - 8, "three": 15 / 5}`

	l := lexer.New(input)
	p := New(l)
	program, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("%v", err)
	}

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.DictLiteral)
	if !ok {
		t.Fatalf("exp is not ast.DictLiteral. got=%T", stmt.Expression)
	}

	if len(hash.Pairs) != 3 {
		t.Errorf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
	}

	tests := map[string]func(ast.Expression){
		"one": func(e ast.Expression) {
			testBinaryOpExpression(t, e, 0, "+", 1)
		},
		"two": func(e ast.Expression) {
			testBinaryOpExpression(t, e, 10, "-", 8)
		},
		"three": func(e ast.Expression) {
			testBinaryOpExpression(t, e, 15, "/", 5)
		},
	}

	for key, value := range hash.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not ast.StringLiteral. got=%T", key)
			continue
		}

		testFunc, ok := tests[literal.String()]
		if !ok {
			t.Errorf("No test function for key %q found", literal.String())
			continue
		}

		testFunc(value)
	}
}

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

func TestParsingSubscriptionExpressions(t *testing.T) {
	input := "myArray[1 + 1]"

	l := lexer.New(input)
	p := New(l)
	program, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("%v", err)
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	subExpr, ok := stmt.Expression.(*ast.SubscriptionExpression)
	if !ok {
		t.Fatalf("exp not *ast.SubscriptionExpression. got=%T", stmt.Expression)
	}

	if !testIdentifier(t, subExpr.Left, "myArray") {
		return
	}

	if !testBinaryOpExpression(t, subExpr.Index, 1, "+", 1) {
		return
	}
}

func TestAttributeExpression(t *testing.T) {
	input := "a.b"

	l := lexer.New(input)
	p := New(l)
	program, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("%v", err)
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	dotExp, ok := stmt.Expression.(*ast.AttributeExpression)
	if !ok {
		t.Fatalf("exp not *ast.DotExpression. got=%T", stmt.Expression)
	}

	if !testIdentifier(t, dotExp.Left, "a") {
		return
	}

	if !testIdentifier(t, dotExp.Attribute, "b") {
		return
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5);"

	l := lexer.New(input)
	p := New(l)
	program, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("%v", err)
	}

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T",
			stmt.Expression)
	}

	if !testIdentifier(t, exp.Function, "add") {
		return
	}

	if len(exp.Arguments) != 3 {
		t.Fatalf("wrong length of arguments. got=%d", len(exp.Arguments))
	}

	testLiteralExpression(t, exp.Arguments[0], 1)
	testBinaryOpExpression(t, exp.Arguments[1], 2, "*", 3)
	testBinaryOpExpression(t, exp.Arguments[2], 4, "+", 5)
}

func TestCallExpressionParameterParsing(t *testing.T) {
	tests := []struct {
		input         string
		expectedIdent string
		expectedArgs  []string
	}{
		{
			input:         "add();",
			expectedIdent: "add",
			expectedArgs:  []string{},
		},
		{
			input:         "add(1);",
			expectedIdent: "add",
			expectedArgs:  []string{"1"},
		},
		{
			input:         "add(1,);",
			expectedIdent: "add",
			expectedArgs:  []string{"1"},
		},
		{
			input:         "add(1,a.b);",
			expectedIdent: "add",
			expectedArgs:  []string{"1", "(a.b)"},
		},
		{
			input:         "add(1,a.b,);",
			expectedIdent: "add",
			expectedArgs:  []string{"1", "(a.b)"},
		},
		{
			input:         "add(1, 2 * 3, 4 + 5);",
			expectedIdent: "add",
			expectedArgs:  []string{"1", "(2 * 3)", "(4 + 5)"},
		},
		{
			input:         "add(f(), f2(2), a.b.c[2](d).e, 1 + 2 * 3);",
			expectedIdent: "add",
			expectedArgs:  []string{"f()", "f2(2)", "((((a.b).c)[2])(d).e)", "(1 + (2 * 3))"},
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program, err := p.ParseProgram()
		if err != nil {
			t.Fatalf("%v", err)
		}

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		exp, ok := stmt.Expression.(*ast.CallExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T",
				stmt.Expression)
		}

		if !testIdentifier(t, exp.Function, tt.expectedIdent) {
			return
		}

		if len(exp.Arguments) != len(tt.expectedArgs) {
			t.Fatalf("wrong number of arguments. want=%d, got=%d",
				len(tt.expectedArgs), len(exp.Arguments))
		}

		for i, arg := range tt.expectedArgs {
			if exp.Arguments[i].String() != arg {
				t.Errorf("argument %d wrong. want=%q, got=%q", i,
					arg, exp.Arguments[i].String())
			}
		}
	}
}

func testVarStatement(t *testing.T, s ast.Statement, name string) bool {
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
		t.Errorf("varStmt.Name.Value not '%s'. got=%s", name, varStmt.Name.Value)
		return false
	}

	if varStmt.Name.TokenLiteral() != name {
		t.Errorf("varStmt.Name.TokenLiteral() not '%s'. got=%s",
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
		t.Errorf("stmt.Name.Value not '%s'. got=%s", name, stmt.Name.Value)
		return false
	}

	if stmt.Name.TokenLiteral() != name {
		t.Errorf("stmt.Name.TokenLiteral() not '%s'. got=%s",
			name, stmt.Name.TokenLiteral())
		return false
	}

	return true
}

func testAssignStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != name {
		t.Errorf("s.TokenLiteral not '%s'. got=%q", name, s.TokenLiteral())
		return false
	}

	stmt, ok := s.(*ast.AssignStatement)
	if !ok {
		t.Errorf("s not *ast.AssignStatement. got=%T", s)
		return false
	}

	if stmt.Name.Value != name {
		t.Errorf("stmt.Name.Value not '%s'. got=%s", name, stmt.Name.Value)
		return false
	}

	if stmt.Name.TokenLiteral() != name {
		t.Errorf("stmt.Name.TokenLiteral() not '%s'. got=%s",
			name, stmt.Name.TokenLiteral())
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

	//if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
	//	t.Errorf("integ.TokenLiteral not %d. got=%s", value,
	//		integ.TokenLiteral())
	//	return false
	//}

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
