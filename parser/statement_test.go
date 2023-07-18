package parser

import (
	"testing"
	"weilang/ast"
	"weilang/lexer"
)

func TestProgram(t *testing.T) {
	input := `
// a

// b
// d
var foo = fn(a) { // abc
	if(1){
		var a = 2 // abc
	}
	return a; // abc
}
var b = foo(10) // abc
b // abc
// if
if(2) 
{} 
else {
}
b = b + 2

if(2) 
{} 
else if (4) {
}
b = b + 2

while(2){}
b = b + 2

while(2)
{
var d = 1}

b = b + 2
var a = 'abc'
// var m = [1, 2, 3, 4]
// print(m.append(5, 6, 7, 8, 9))
// m.append()
print(a)

wei.export(a, b)
`

	l := lexer.New(input)
	p := New(l)
	program, err := p.ParseProgram()
	if err != nil {
		t.Errorf("%v", err)
		t.FailNow()
	}

	if len(program.Statements) < 3 {
		t.Errorf("program.Statements does not contain 3 statements. got=%d",
			len(program.Statements))
		t.FailNow()
	}

	stmt := program.Statements[0].(*ast.VarStatement)
	if !testVarStatement(t, stmt, "foo") {
		return
	}
	function, ok := stmt.Value.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Value is not ast.FunctionLiteral. got=%T",
			stmt.Value)
	}
	if len(function.Parameters) != 1 {
		t.Fatalf("function literal parameters wrong. want 1, got=%d\n",
			len(function.Parameters))
	}
	testLiteralExpression(t, function.Parameters[0], "a")
	if len(function.Body.Statements) != 2 {
		t.Fatalf("function.Body.Statements has not 2 statements. got=%d\n",
			len(function.Body.Statements))
	}
	ifStmt, ok := function.Body.Statements[0].(*ast.IfStatement)
	if !ok {
		t.Fatalf("function body stmt is not ast.IfStatement. got=%T",
			function.Body.Statements[0])
	}
	if len(ifStmt.IfBranches) != 1 {
		t.Fatalf("if stmt has not 1 branches. got=%d", len(ifStmt.IfBranches))
	}
	if !testIntegerLiteral(t, ifStmt.IfBranches[0].Condition, 1) {
		return
	}
	if len(ifStmt.IfBranches[0].Body.Statements) != 1 {
		t.Fatalf("if branch has not 1 statement. got=%d", len(ifStmt.IfBranches[0].Body.Statements))
	}
	ifBody, ok := ifStmt.IfBranches[0].Body.Statements[0].(*ast.VarStatement)
	if !ok {
		t.Fatalf("if branch body is not ast.VarStatement. got=%T", ifStmt.IfBranches[0].Body.Statements[0])
	}
	if !testVarStatement(t, ifBody, "a") {
		return
	}
	if !testIntegerLiteral(t, ifBody.Value, 2) {
		return
	}
	if ifStmt.ElseBody != nil {
		t.Fatalf("if stmt has not else. got=%T", ifStmt.ElseBody)
	}

	returnStmt, ok := function.Body.Statements[1].(*ast.ReturnStatement)
	if !ok {
		t.Fatalf("function body stmt is not ast.ReturnStatement. got=%T",
			function.Body.Statements[0])
	}
	if !testIdentifier(t, returnStmt.ReturnValue, "a") {
		return
	}

	varStmt, ok := program.Statements[1].(*ast.VarStatement)
	if !ok {
		t.Fatalf("program.Statements[1] is not ast.VarStatement. got=%T", program.Statements[1])
	}
	if !testVarStatement(t, varStmt, "b") {
		return
	}
	if varStmt.Value.String() != "foo(10)" {
		t.Fatalf("value got=%s", varStmt.Value.String())
	}

	exprStmt, ok := program.Statements[2].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[2] is not ast.ExpressionStatement. got=%T", program.Statements[2])
	}
	if !testIdentifier(t, exprStmt.Expression, "b") {
		return
	}
}

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

func TestConFunctionStatements(t *testing.T) {
	input := "con say = fn(name) { return }"
	l := lexer.New(input)
	p := New(l)
	program, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("%v", err)
	}

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statements. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ConStatement)
	if !ok {
		t.Fatalf("not ast.ConStatement. got=%T", program.Statements[0])
	}

	if !testConStatement(t, stmt, "say") {
		return
	}

	function, ok := stmt.Value.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got=%T",
			stmt.Value)
	}

	if len(function.Parameters) != 1 {
		t.Fatalf("function literal parameters wrong. want 1, got=%d\n",
			len(function.Parameters))
	}

	testLiteralExpression(t, function.Parameters[0], "name")

	if len(function.Body.Statements) != 1 {
		t.Fatalf("function.Body.Statements has not 1 statements. got=%d\n",
			len(function.Body.Statements))
	}

	returnStmt, ok := function.Body.Statements[0].(*ast.ReturnStatement)
	if !ok {
		t.Fatalf("function body stmt is not ast.ReturnStatement. got=%T",
			function.Body.Statements[0])
	}

	if returnStmt.ReturnValue != nil {
		t.Fatalf("return unexpected value %s", returnStmt.ReturnValue)
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
		{"foobar.a = y;", "(foobar.a)", "y"},
		{"foobar.b.a = y;", "((foobar.b).a)", "y"},
		{"foobar.c.b.a = y;", "(((foobar.c).b).a)", "y"},
		{"foobar[1] = y;", "(foobar[1])", "y"},
		{"foobar[1][2] = y;", "((foobar[1])[2])", "y"},
		{"foobar[1][2]['a'] = y;", `(((foobar[1])[2])[a])`, "y"},
		{"foobar[1].b.d['a'] = y;", `((((foobar[1]).b).d)[a])`, "y"},
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

func TestIfStatements(t *testing.T) {
	input := `if (x < y) { x }`

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

	stmt, ok := program.Statements[0].(*ast.IfStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	if len(stmt.IfBranches) != 1 {
		t.Fatalf("too many if branch: %d", len(stmt.IfBranches))
	}

	if stmt.ElseBody != nil {
		t.Fatalf("else not empty: \"%s\"", stmt.ElseBody.String())
	}

	ifCase := stmt.IfBranches[0]
	if !testBinaryOpExpression(t, ifCase.Condition, "x", "<", "y") {
		return
	}

	if len(ifCase.Body.Statements) != 1 {
		t.Errorf("body is not 1 statements. got=%d\n",
			len(ifCase.Body.Statements))
	}

	consequence, ok := ifCase.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			ifCase.Body.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if stmt.ElseBody != nil {
		t.Errorf("ElseBody was not nil. got=%+v", stmt.ElseBody)
	}
}

func TestIfStatementsWithElse(t *testing.T) {
	input := `
if (x < y) 
{ x }
else { y }
`

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

	stmt, ok := program.Statements[0].(*ast.IfStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	if len(stmt.IfBranches) != 1 {
		t.Fatalf("too many if branch: %d", len(stmt.IfBranches))
	}

	ifBranch := stmt.IfBranches[0]
	if !testBinaryOpExpression(t, ifBranch.Condition, "x", "<", "y") {
		return
	}

	if len(ifBranch.Body.Statements) != 1 {
		t.Errorf("body is not 1 statements. got=%d\n",
			len(ifBranch.Body.Statements))
	}

	consequence, ok := ifBranch.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			ifBranch.Body.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	elseExpr, ok := stmt.ElseBody.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			ifBranch.Body.Statements[0])
	}

	if !testIdentifier(t, elseExpr.Expression, "y") {
		return
	}
}

func TestIfStatementsWithMultiBranch(t *testing.T) {
	input := `
if (x < y){ x }
else if (x == y) { y }
else if (x >= y) { ddd; }
else if (x != y) { xy }
else { x + y }
`

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

	stmt, ok := program.Statements[0].(*ast.IfStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	if stmt.TokenLiteral() != "if" {
		t.Fatalf("expected if, but got=%s", stmt.TokenLiteral())
	}

	expecteds := []struct {
		ops  []string
		expr string
	}{
		{
			[]string{"x", "<", "y"},
			"x",
		},
		{
			[]string{"x", "==", "y"},
			"y",
		},
		{
			[]string{"x", ">=", "y"},
			"ddd",
		},
		{
			[]string{"x", "!=", "y"},
			"xy",
		},
	}

	if len(stmt.IfBranches) != len(expecteds) {
		t.Fatalf("wrong number if branch: %d", len(stmt.IfBranches))
	}

	for i, ifBranch := range stmt.IfBranches {
		expected := expecteds[i]
		if !testBinaryOpExpression(t, ifBranch.Condition, expected.ops[0], expected.ops[1], expected.ops[2]) {
			return
		}

		if len(ifBranch.Body.Statements) != 1 {
			t.Errorf("body is not 1 statements. got=%d\n",
				len(ifBranch.Body.Statements))
		}

		consequence, ok := ifBranch.Body.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
				ifBranch.Body.Statements[0])
		}

		if !testIdentifier(t, consequence.Expression, expected.expr) {
			return
		}
	}

	elseStmt, ok := stmt.ElseBody.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			stmt.ElseBody.Statements[0])
	}

	if !testBinaryOpExpression(t, elseStmt.Expression, "x", "+", "y") {
		return
	}
}

func TestWhileStatements(t *testing.T) {
	input := `while (x < y) { x; continue; break; }`

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

	stmt, ok := program.Statements[0].(*ast.WhileStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.WhileStatement. got=%T",
			program.Statements[0])
	}

	if stmt.TokenLiteral() != "while" {
		t.Fatalf("expected while, but got=%s", stmt.TokenLiteral())
	}

	if !testBinaryOpExpression(t, stmt.Condition, "x", "<", "y") {
		return
	}

	exprStmt, ok := stmt.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("not ast.ExpressionStatement. got=%T", stmt.Body.Statements[0])
	}
	if !testIdentifier(t, exprStmt.Expression, "x") {
		return
	}

	continueStmt, ok := stmt.Body.Statements[1].(*ast.ContinueStatement)
	if !ok {
		t.Fatalf("not ast.ContinueStatement. got=%T", stmt.Body.Statements[1])
	}
	if continueStmt.TokenLiteral() != "continue" {
		t.Fatalf("expected continue, but got=%q", continueStmt.TokenLiteral())
	}

	breakStmt, ok := stmt.Body.Statements[2].(*ast.BreakStatement)
	if !ok {
		t.Fatalf("not ast.BreakStatement. got=%T", stmt.Body.Statements[2])
	}
	if breakStmt.TokenLiteral() != "break" {
		t.Fatalf("expected break, but got=%q", breakStmt.TokenLiteral())
	}

	input = `
while (1) {
  con foo = fn(){
    while (2) {break}
  }
  break
}
`
	l = lexer.New(input)
	p = New(l)
	program, err = p.ParseProgram()
	if err != nil {
		t.Fatalf("%v", err)
	}

	input = `
while (1) {
  continue
  con foo = fn(){
    while (2) {continue}
  }
}
`
	l = lexer.New(input)
	p = New(l)
	program, err = p.ParseProgram()
	if err != nil {
		t.Fatalf("%v", err)
	}
}

func TestForInStatement(t *testing.T) {
	input := `
for (i, a in b) {}
`

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

	stmt, ok := program.Statements[0].(*ast.ForInStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.WhileStatement. got=%T",
			program.Statements[0])
	}

	if stmt.TokenLiteral() != "for" {
		t.Fatalf("expected for, but got=%s", stmt.TokenLiteral())
	}

	if !testIdentifier(t, stmt.First, "i") {
		return
	}

	if !testIdentifier(t, stmt.Second, "a") {
		return
	}

	if !testIdentifier(t, stmt.Expr, "b") {
		return
	}

	if len(stmt.Body.Statements) > 0 {
		t.Fatalf("expected zero statement, but got=%d", len(stmt.Body.Statements))
	}
}

func TestWeiExportStatements(t *testing.T) {
	tests := []struct {
		input         string
		expectedNames []string
	}{
		{
			`wei.export()`,
			[]string{},
		},
		{
			`wei.export(a)`,
			[]string{"a"},
		},
		{
			`wei.export(a,)`,
			[]string{"a"},
		},
		{
			`wei.export(a,foo)`,
			[]string{"a", "foo"},
		},
		{
			`wei.export(a,foo,)`,
			[]string{"a", "foo"},
		},
		{
			`wei.export(a,foo,bar)`,
			[]string{"a", "foo", "bar"},
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
		stmt, ok := program.Statements[0].(*ast.WeiExportStatement)
		if !ok {
			t.Errorf("want WeiExportStatement, but got=%T", program.Statements[0])
			t.FailNow()
		}
		if len(stmt.Names) != len(tt.expectedNames) {
			t.Errorf("wrong number of name. got=%d, want=%d", len(stmt.Names), len(tt.expectedNames))
			t.FailNow()
		}
		for i, name := range stmt.Names {
			testIdentifier(t, name, tt.expectedNames[i])
		}
	}
}
