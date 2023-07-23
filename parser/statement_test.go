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
for (var i, a in b) {}
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

	if stmt.Con {
		t.Fatalf("expected var, but got con")
	}

	if len(stmt.Targets) != 2 {
		t.Fatalf("expected 2 targets, but got %d", len(stmt.Targets))
	}

	if !testIdentifier(t, stmt.Targets[0], "i") {
		return
	}

	if !testIdentifier(t, stmt.Targets[1], "a") {
		return
	}

	if !testIdentifier(t, stmt.Expr, "b") {
		return
	}

	if len(stmt.Body.Statements) > 0 {
		t.Fatalf("expected zero statement, but got=%d", len(stmt.Body.Statements))
	}
}

func TestFunctionDefineStatement(t *testing.T) {
	input := `fn ddd(x, y) { x + y; }`

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

	stmt, ok := program.Statements[0].(*ast.FunctionDefineStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.FunctionDefineStatement. got=%T",
			program.Statements[0])
	}

	function := stmt.Function
	expectedName := "ddd"
	if function.Name != expectedName {
		t.Fatalf("expected function name %q, but got=%q", expectedName, function.Name)
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

func TestClassDefineStatementEmptyBody(t *testing.T) {
	input := `class Foo{}`

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

	stmt, ok := program.Statements[0].(*ast.ClassDefineStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ClassDefineStatement. got=%T",
			program.Statements[0])
	}

	expectedName := "Foo"
	if stmt.Name != expectedName {
		t.Fatalf("expected name %q, got=%q", expectedName, stmt.Name)
	}

	if len(stmt.Body.Statements) != 0 {
		t.Fatalf("expected empty body, got=%d", len(stmt.Body.Statements))
	}
}

func TestClassDefineStatement(t *testing.T) {
	input := `class Foo{
var a
var b = 2
con c
con d = 3
var class.e = 4

fn __init__() {}
fn class.init() {}
}`

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

	stmt, ok := program.Statements[0].(*ast.ClassDefineStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ClassDefineStatement. got=%T",
			program.Statements[0])
	}

	expectedName := "Foo"
	if stmt.Name != expectedName {
		t.Fatalf("expected name %q, got=%q", expectedName, stmt.Name)
	}

	if len(stmt.Body.Statements) != 7 {
		t.Fatalf("expected body 7, got=%d", len(stmt.Body.Statements))
	}

	statements := stmt.Body.Statements

	{
		stmt0, ok := statements[0].(*ast.ClassVariableDeclarationStatement)
		if !ok {
			t.Fatalf("Body.Statements[0] is not ast.ClassVariableDeclarationStatement. got=%T", statements[0])
		}
		if !testIdentifier(t, stmt0.Name, "a") {
			t.FailNow()
		}
		if stmt0.Con || stmt0.Class {
			t.Fatalf("stmt0 expected normal var, got con=%v class=%v", stmt0.Con, stmt0.Class)
		}
		if stmt0.Expr != nil {
			t.Fatalf("stmt0 expected empty expr, got=%T", stmt0.Expr)
		}
	}
	{
		stmt1, ok := statements[1].(*ast.ClassVariableDeclarationStatement)
		if !ok {
			t.Fatalf("Body.Statements[1] is not ast.ClassVariableDeclarationStatement. got=%T", statements[0])
		}
		if !testIdentifier(t, stmt1.Name, "b") {
			t.FailNow()
		}
		if stmt1.Con || stmt1.Class {
			t.Fatalf("stmt1 expected normal var, got con=%v class=%v", stmt1.Con, stmt1.Class)
		}
		if stmt1.Expr == nil {
			t.Fatalf("stmt1 expected expr, got=%T", stmt1.Expr)
		}
		testIntegerLiteral(t, stmt1.Expr, 2)
	}
	{
		stmt2, ok := statements[2].(*ast.ClassVariableDeclarationStatement)
		if !ok {
			t.Fatalf("Body.Statements[2] is not ast.ClassVariableDeclarationStatement. got=%T", statements[0])
		}
		if !testIdentifier(t, stmt2.Name, "c") {
			t.FailNow()
		}
		if !stmt2.Con || stmt2.Class {
			t.Fatalf("stmt2 expected normal con, got con=%v class=%v", stmt2.Con, stmt2.Class)
		}
		if stmt2.Expr != nil {
			t.Fatalf("stmt2 expected empty expr, got=%T", stmt2.Expr)
		}
	}

	{
		stmt3, ok := statements[3].(*ast.ClassVariableDeclarationStatement)
		if !ok {
			t.Fatalf("Body.Statements[3] is not ast.ClassVariableDeclarationStatement. got=%T", statements[0])
		}
		if !testIdentifier(t, stmt3.Name, "d") {
			t.FailNow()
		}
		if !stmt3.Con || stmt3.Class {
			t.Fatalf("stmt3 expected normal con, got con=%v class=%v", stmt3.Con, stmt3.Class)
		}
		if stmt3.Expr == nil {
			t.Fatalf("stmt3 expected expr, got=%T", stmt3.Expr)
		}
		testIntegerLiteral(t, stmt3.Expr, 3)
	}

	{
		stmt4, ok := statements[4].(*ast.ClassVariableDeclarationStatement)
		if !ok {
			t.Fatalf("Body.Statements[4] is not ast.ClassVariableDeclarationStatement. got=%T", statements[0])
		}
		if !testIdentifier(t, stmt4.Name, "e") {
			t.FailNow()
		}
		if stmt4.Con || !stmt4.Class {
			t.Fatalf("stmt4 expected class con, got con=%v class=%v", stmt4.Con, stmt4.Class)
		}
		if stmt4.Expr == nil {
			t.Fatalf("stmt4 expected expr, got=%T", stmt4.Expr)
		}
		testIntegerLiteral(t, stmt4.Expr, 4)
	}

	{
		stmt5, ok := statements[5].(*ast.ClassMethodDefineStatement)
		if !ok {
			t.Fatalf("Body.Statements[5] is not ast.ClassMethodDefineStatement. got=%T", statements[0])
		}
		if stmt5.Class {
			t.Fatalf("stmt5 expected normal method, but class=%v", stmt5.Class)
		}
		function := stmt5.Function
		if function.Name != "__init__" {
			t.Fatalf("stmt5 expected __init__ name, but got=%q", function.Name)
		}
		if len(function.Parameters) != 0 {
			t.Fatalf("stmt5 expected empty parameter, but got=%v", len(function.Parameters))
		}
		if len(function.Body.Statements) != 0 {
			t.Fatalf("stmt5 expected empty body, but got=%v", len(function.Body.Statements))
		}
	}
	{
		stmt6, ok := statements[6].(*ast.ClassMethodDefineStatement)
		if !ok {
			t.Fatalf("Body.Statements[6] is not ast.ClassMethodDefineStatement. got=%T", statements[0])
		}
		if !stmt6.Class {
			t.Fatalf("stmt6 expected class method, but class=%v", stmt6.Class)
		}
		function := stmt6.Function
		if function.Name != "init" {
			t.Fatalf("stmt6 expected init name, but got=%q", function.Name)
		}
		if len(function.Parameters) != 0 {
			t.Fatalf("stmt6 expected empty parameter, but got=%v", len(function.Parameters))
		}
		if len(function.Body.Statements) != 0 {
			t.Fatalf("stmt6 expected empty body, but got=%v", len(function.Body.Statements))
		}
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
