package parser

import (
	"testing"
	"weilang/ast"
	"weilang/lexer"
)

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

		{"null and true", "null", "and", true},
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

func TestFunctionLiteralParsing(t *testing.T) {
	input := `fn(x, y) { x + y; }`

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
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got=%T",
			stmt.Expression)
	}

	if len(function.Parameters) != 2 {
		t.Fatalf("function literal parameters wrong. want 2, got=%d\n",
			len(function.Parameters))
	}

	testLiteralExpression(t, function.Parameters[0], "x")
	testLiteralExpression(t, function.Parameters[1], "y")

	if len(function.Body.Statements) != 1 {
		t.Fatalf("function.Body.Statements has not 1 statements. got=%d\n",
			len(function.Body.Statements))
	}

	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("function body stmt is not ast.ExpressionStatement. got=%T",
			function.Body.Statements[0])
	}

	testBinaryOpExpression(t, bodyStmt.Expression, "x", "+", "y")

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

func TestParsingEmptyDictLiteral(t *testing.T) {
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

func TestParsingDictLiteralsStringKeys(t *testing.T) {
	testCases := []struct {
		input    string
		expected map[string]int64
	}{
		{
			`{}`,
			map[string]int64{},
		},
		{
			`{"one": 1}`,
			map[string]int64{
				"one": 1,
			},
		},
		{
			`{"one": 1,}`,
			map[string]int64{
				"one": 1,
			},
		},
		{
			`{"one": 1, "two":2}`,
			map[string]int64{
				"one": 1,
				"two": 2,
			},
		},
		{
			`{"one": 1, "two":2,}`,
			map[string]int64{
				"one": 1,
				"two": 2,
			},
		},
		{
			`{
					"one": 1, 
					"two":2, 
					"three": 3
					}`,
			map[string]int64{
				"one":   1,
				"two":   2,
				"three": 3,
			},
		},
		{
			`{
					"one": 1, 
					"two":2, 
					"three": 3,
					}`,
			map[string]int64{
				"one":   1,
				"two":   2,
				"three": 3,
			},
		},
	}
	for _, tc := range testCases {
		l := lexer.New(tc.input)
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

		expected := tc.expected

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
}

func TestParsingDictLiteralsBooleanKeys(t *testing.T) {
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

func TestParsingDictLiteralsIntegerKeys(t *testing.T) {
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

func TestParsingDictLiteralsWithExpressions(t *testing.T) {
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
		{"`hello\n\t\\\"world3`", "hello\n\t\\\"world3"},
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

func TestWeiExpression(t *testing.T) {
	tests := []struct {
		input     string
		attribute string
	}{
		{
			`con a = wei.import("abc")`,
			"import",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program, err := p.ParseProgram()
		if err != nil {
			t.Errorf("got error: %v", err)
			t.FailNow()
		}
		if len(program.Statements) != 1 {
			t.Errorf("wrong number of statements. got=%d, want=1", len(program.Statements))
			t.FailNow()
		}
		stmt, ok := program.Statements[0].(*ast.ConStatement)
		if !ok {
			t.Errorf("want ExpressionStatement, but got=%T", program.Statements[0])
			t.FailNow()
		}
		call, ok := stmt.Value.(*ast.CallExpression)
		if !ok {
			t.Errorf("want CallExpression, but got=%T", stmt.Value)
			t.FailNow()
		}
		weiAttr, ok := call.Function.(*ast.WeiAttributeExpression)
		if !ok {
			t.Errorf("want WeiAttributeExpression, but got=%T", call.Function)
			t.FailNow()
		}
		testIdentifier(t, weiAttr.Attribute, tt.attribute)
	}
}
