package parser

import (
	"fmt"
	"strconv"
	"strings"
	"weilang/ast"
	"weilang/lexer"
	"weilang/token"
)

// 换行符（ NEWLINE token ）处理规则
//   () 括号内可以随意换行
//   [] {} 列表和字典字面量内可以随意换行
//   语句前后可以随意换行
//   语句块（{} 括号包裹的，包括 {}）前后可以随意换行

// 用来保存解析器 dump 后的信息
type dumpInfo struct {
	index int
	token token.Token
}

type Parser struct {
	l         *lexer.Lexer
	currToken token.Token
	filename  string
	lines     []string
	// parenCount 进入的括号数量，用于判断当前解析是否在括号内
	parenCount int
	// whileStack while 层级栈，用来检查 continue break 是否 while 块中
	// 之所以用栈，而不是用整数，是因为有下面这种情况， while 里面套函数定义
	// while (1) {
	//   var foo = fn() {
	//     continue
	//   }
	// }
	// 上面这种是非法的，但是下面这种是合法的
	// while (1) {
	//   var foo = fn() {
	//     while (2) {
	//       continue
	//     }
	//   }
	// }
	whileStack []int
	// 保存 token 用于回溯
	backupTokens       []token.Token
	backupTokenEnabled bool
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:          l,
		filename:   l.Filename,
		lines:      l.GetLines(),
		parenCount: 0,
		whileStack: []int{0},
	}
	p.nextToken()
	return p
}

/*
使用 lsbasi 的解析实现 (https://github.com/rspivak/lsbasi)
相比 《用Go语言自制解释器》 的实现更容易理解
*/

func (p *Parser) ParseProgram() (*ast.Program, error) {
	return p.program()
}

// program ::= (statement)* EOF
func (p *Parser) program() (*ast.Program, error) {
	location := p.currFileLocation()
	program := &ast.Program{Location: location}
	program.Statements = []ast.Statement{}

	for !p.currTokenIs(token.EOF) {
		stmt, err := p.statement()
		if err != nil {
			return nil, err
		}
		program.Statements = append(program.Statements, stmt)
	}

	err := p.eat(token.EOF)
	if err != nil {
		return nil, err
	}

	return program, nil
}

// block_statement ::= "{" (statement)* "}"
func (p *Parser) blockStatement() (*ast.BlockStatement, error) {
	p.skipNewline()

	tok := p.currToken
	location := p.currFileLocation()
	err := p.eat(token.LBRACE)
	if err != nil {
		return nil, err
	}
	block := &ast.BlockStatement{Location: location, Token: tok}
	// 处理空的语句块的情况，例如 "if (1) {}"
	p.skipNewline()

	for !p.currTokenIs(token.RBRACE) {
		stmt, err := p.statement()
		if err != nil {
			return nil, err
		}
		block.Statements = append(block.Statements, stmt)
	}
	err = p.eat(token.RBRACE)
	if err != nil {
		return nil, err
	}

	return block, nil
}

func (p *Parser) isStatementEnd() bool {
	switch p.currToken.Type {
	case token.SEMICOLON:
		p.nextToken()
		return true
	case token.NEWLINE, token.RBRACE, token.EOF:
		return true
	default:
		return false
	}
}

// statement ::= var_statement | con_statement
// | for_in_statement
// | if_statement
// | return_statement
// | expression_statement
// | assign_statement
// | while_statement
// | continue_statement
// | break_statement
// | wei_export_statement
func (p *Parser) statement() (ast.Statement, error) {
	p.skipNewline()
	defer func() { p.skipNewline() }()
	info := p.dump()
	switch p.currToken.Type {
	case token.VAR:
		return p.varStatement()
	case token.CON:
		return p.conStatement()
	case token.FOR:
		return p.forInStatement()
	case token.RETURN:
		return p.returnStatement()
	case token.IDENT:
		stmt, err := p.expressionStatement()
		if err == nil {
			return stmt, nil
		}
		p.restore(info)
		return p.assignStatement()
	case token.IF:
		return p.ifStatement()
	case token.WHILE:
		return p.whileStatement()
	case token.CONTINUE:
		return p.continueStatement()
	case token.BREAK:
		return p.breakStatement()
	case token.WEI:
		stmt, err := p.weiExportStatement()
		if err == nil {
			return stmt, nil
		}
		p.restore(info)
		return p.expressionStatement()
	default:
		return p.expressionStatement()
	}
}

// for_in_statement ::= "for" "(" ("var" | "con") IDENT ("," IDENT)* "in" expression ")" block_statement  (";" | NEWLINE)
func (p *Parser) forInStatement() (*ast.ForInStatement, error) {
	tk := p.currToken
	location := p.currFileLocation()
	err := p.eatContinuously(token.FOR, token.LPAREN)
	if err != nil {
		return nil, err
	}
	con := p.currTokenIs(token.CON)
	err = p.eatIn(token.VAR, token.CON)
	if err != nil {
		return nil, err
	}
	var targets []*ast.Identifier
	target, err := p.ident()
	targets = append(targets, target)
	if err != nil {
		return nil, err
	}
	for p.currTokenIs(token.COMMA) {
		p.nextToken()
		target, err = p.ident()
		targets = append(targets, target)
		if err != nil {
			return nil, err
		}
	}
	err = p.eat(token.IN)
	if err != nil {
		return nil, err
	}
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	err = p.eat(token.RPAREN)
	if err != nil {
		return nil, err
	}
	body, err := p.blockStatement()
	if err != nil {
		return nil, err
	}
	if !p.isStatementEnd() {
		return nil, p.expectError(token.SEMICOLON)
	}
	stmt := &ast.ForInStatement{
		Location: location,
		Token:    tk,
		Con:      con,
		Targets:  targets,
		Expr:     expr,
		Body:     body,
	}
	return stmt, nil
}

// wei_export_statement ::= "wei" "." "export" "(" [wei_export_args] ")"
// wei_export_args      ::= IDENT (,IDENT)*
func (p *Parser) weiExportStatement() (*ast.WeiExportStatement, error) {
	location := p.currFileLocation()
	tk := p.currToken
	err := p.eat(token.WEI)
	if err != nil {
		return nil, err
	}
	err = p.eat(token.DOT)
	if err != nil {
		return nil, err
	}
	literal := p.currToken.Literal
	if literal != "export" {
		return nil, p.syntaxError("expected 'export'")
	}
	err = p.eat(token.IDENT)
	if err != nil {
		return nil, err
	}
	p.parenCount++
	err = p.eat(token.LPAREN)
	if err != nil {
		return nil, err
	}

	// 解析参数
	var names []*ast.Identifier
	if p.currTokenIs(token.IDENT) {
		name, _ := p.ident()
		names = append(names, name)
	}
	for p.currTokenIs(token.COMMA) {
		p.nextToken()
		// 处理这种情况 wei.export(a,) ，最后一个参数后面有个逗号
		if p.currTokenIs(token.RPAREN) {
			break
		}
		name, err := p.ident()
		names = append(names, name)
		if err != nil {
			return nil, err
		}
	}

	p.parenCount--
	err = p.eat(token.RPAREN)
	if err != nil {
		return nil, err
	}

	stmt := &ast.WeiExportStatement{
		Location: location,
		Token:    tk,
		Names:    names,
	}
	return stmt, nil
}

// var_statement ::= "var" IDENT "=" expression (";" | NEWLINE)
func (p *Parser) varStatement() (*ast.VarStatement, error) {
	location := p.currFileLocation()
	tok := p.currToken
	err := p.eat(token.VAR)
	if err != nil {
		return nil, err
	}
	name, err := p.ident()
	if err != nil {
		return nil, err
	}
	err = p.eat(token.ASSIGN)
	if err != nil {
		return nil, err
	}
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	if !p.isStatementEnd() {
		return nil, p.expectError(token.SEMICOLON)
	}
	varStmt := &ast.VarStatement{
		Location: location,
		Token:    tok,
		Name:     name,
		Value:    expr,
	}
	return varStmt, nil
}

// con_statement ::= "con" IDENT "=" expression (";" | NEWLINE)
func (p *Parser) conStatement() (*ast.ConStatement, error) {
	location := p.currFileLocation()
	tok := p.currToken
	err := p.eat(token.CON)
	if err != nil {
		return nil, err
	}
	name, err := p.ident()
	if err != nil {
		return nil, err
	}
	err = p.eat(token.ASSIGN)
	if err != nil {
		return nil, err
	}
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	if !p.isStatementEnd() {
		return nil, p.expectError(token.SEMICOLON)
	}
	stmt := &ast.ConStatement{
		Location: location,
		Token:    tok,
		Name:     name,
		Value:    expr,
	}
	return stmt, nil
}

// return_statement ::= "return" ( expression ) (";" | NEWLINE)
func (p *Parser) returnStatement() (*ast.ReturnStatement, error) {
	location := p.currFileLocation()
	tok := p.currToken
	err := p.eat(token.RETURN)
	if err != nil {
		return nil, err
	}
	var expr ast.Expression
	if !p.isStatementEnd() {
		expr, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	if !p.isStatementEnd() {
		return nil, p.expectError(token.SEMICOLON)
	}
	stmt := &ast.ReturnStatement{
		Location:    location,
		Token:       tok,
		ReturnValue: expr,
	}
	return stmt, nil
}

// assignStatement 解析赋值语句
//
// assign_statement ::= primary "=" expression (";" | NEWLINE)
func (p *Parser) assignStatement() (*ast.AssignStatement, error) {
	location := p.currFileLocation()
	tok := p.currToken
	left, err := p.primary()
	if err != nil {
		return nil, err
	}
	err = p.eat(token.ASSIGN)
	if err != nil {
		return nil, err
	}
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	if !p.isStatementEnd() {
		return nil, p.expectError(token.SEMICOLON)
	}
	stmt := &ast.AssignStatement{
		Location: location,
		Token:    tok,
		Left:     left,
		Value:    expr,
	}
	return stmt, nil
}

// primary 解析标志符、属性访问、下标访问
//
// primary          ::= IDENT ( subscription | attribute)*
// subscription     ::= "[" expression "]"
// attribute        ::= "." IDENT
func (p *Parser) primary() (ast.Expression, error) {
	tok := p.currToken
	var expr ast.Expression
	expr, err := p.ident()
	if err != nil {
		return nil, err
	}
	for p.currTokenIn(token.LBRACKET, token.DOT) {
		tok = p.currToken
		switch p.currToken.Type {
		case token.LBRACKET:
			p.nextToken()
			index, err := p.expression()
			if err != nil {
				return nil, err
			}
			expr = &ast.SubscriptionExpression{
				Token: tok,
				Left:  expr,
				Index: index,
			}
			err = p.eat(token.RBRACKET)
			if err != nil {
				return nil, err
			}
		case token.DOT:
			p.nextToken()
			ident, err := p.ident()
			if err != nil {
				return nil, err
			}
			expr = &ast.AttributeExpression{
				Token:     tok,
				Left:      expr,
				Attribute: ident,
			}
		}
	}
	return expr, nil
}

// expression_statement ::= expression (";" | NEWLINE)
func (p *Parser) expressionStatement() (*ast.ExpressionStatement, error) {
	stmt := &ast.ExpressionStatement{Token: p.currToken}
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	stmt.Expression = expr
	if !p.isStatementEnd() {
		return nil, p.expectError(token.SEMICOLON)
	}
	return stmt, nil
}

// if_statement ::= if_branch (else_if_branch)* [else_branch]  (";" | NEWLINE)
// if_branch      ::= "if" "(" expression ")" block_statement
// else_if_branch ::= "else" if_branch
// else_branch    ::= "else" block_statement
func (p *Parser) ifStatement() (*ast.IfStatement, error) {
	tok := p.currToken
	var branches []*ast.IfBranch
	var elseBody *ast.BlockStatement
	branch, err := p.ifBranch()
	if err != nil {
		return nil, err
	}
	branches = append(branches, branch)
	endToken := p.currToken
	p.skipNewline()

	for p.currTokenIs(token.ELSE) {
		p.nextToken()
		if p.currTokenIs(token.IF) {
			branch, err = p.ifBranch()
			if err != nil {
				return nil, err
			}
			branches = append(branches, branch)
			endToken = p.currToken
			p.skipNewline()
		} else {
			elseBody, err = p.blockStatement()
			if err != nil {
				return nil, err
			}
			break
		}
	}
	// endToken 用来避免语句直接跟在 } 后面，而不是另起一行
	// 麻烦点在于 if 语句会有多个语句块
	// 需要处理下面三种非法情况
	// if(1) {} var a = 1
	// if(1) {} else if(2) {} var a = 1
	// if(1) {} else if(2) {} else {} var a = 1
	if !endToken.TypeIs(token.NEWLINE) && !p.isStatementEnd() {
		return nil, p.expectError(token.SEMICOLON)
	}
	stmt := &ast.IfStatement{
		Token:      tok,
		IfBranches: branches,
		ElseBody:   elseBody,
	}
	return stmt, nil
}

// if_branch      ::= "if" "(" expression ")" block_statement
func (p *Parser) ifBranch() (*ast.IfBranch, error) {
	err := p.eat(token.IF)
	if err != nil {
		return nil, err
	}
	p.parenCount++
	err = p.eat(token.LPAREN)
	if err != nil {
		return nil, err
	}
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	p.parenCount--
	err = p.eat(token.RPAREN)
	if err != nil {
		return nil, err
	}
	body, err := p.blockStatement()
	if err != nil {
		return nil, err
	}
	branch := &ast.IfBranch{
		Condition: condition,
		Body:      body,
	}
	return branch, nil
}

// while_statement ::= "while" "(" expression ")" block_statement  (";" | NEWLINE)
func (p *Parser) whileStatement() (*ast.WhileStatement, error) {
	tok := p.currToken
	err := p.eat(token.WHILE)
	if err != nil {
		return nil, err
	}
	p.parenCount++
	err = p.eat(token.LPAREN)
	if err != nil {
		return nil, err
	}
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	p.parenCount--
	err = p.eat(token.RPAREN)
	if err != nil {
		return nil, err
	}
	p.whileStack[len(p.whileStack)-1]++
	body, err := p.blockStatement()
	if err != nil {
		return nil, err
	}
	p.whileStack[len(p.whileStack)-1]--
	if !p.isStatementEnd() {
		return nil, p.expectError(token.SEMICOLON)
	}
	stmt := &ast.WhileStatement{
		Token:     tok,
		Condition: condition,
		Body:      body,
	}
	return stmt, nil
}

// continue_statement ::= "continue" (";" | NEWLINE)
func (p *Parser) continueStatement() (*ast.ContinueStatement, error) {
	if p.whileStack[len(p.whileStack)-1] == 0 {
		return nil, p.syntaxError("continue is not in a loop")
	}
	tok := p.currToken
	err := p.eat(token.CONTINUE)
	if err != nil {
		return nil, err
	}
	if !p.isStatementEnd() {
		return nil, p.expectError(token.SEMICOLON)
	}
	stmt := &ast.ContinueStatement{Token: tok}
	return stmt, nil
}

// break_statement ::= "break" (";" | NEWLINE)
func (p *Parser) breakStatement() (*ast.BreakStatement, error) {
	if p.whileStack[len(p.whileStack)-1] == 0 {
		return nil, p.syntaxError("break is not in a loop")
	}
	tok := p.currToken
	err := p.eat(token.BREAK)
	if err != nil {
		return nil, err
	}
	if !p.isStatementEnd() {
		return nil, p.expectError(token.SEMICOLON)
	}
	stmt := &ast.BreakStatement{Token: tok}
	return stmt, nil
}

/*
表达式语法规则
	每个优先级都有一个语法表示，里面包含本级的所有运算符和更高优先级的表示
	最后会有一个最底层的表示，代表表达式里面的基本单元，例如下面的 atom

运算符优先级参照 Python
https://docs.python.org/3/reference/expressions.html#operator-precedence
*/

// expression 解析表达式
//
// expression ::= orExpression
func (p *Parser) expression() (ast.Expression, error) {
	return p.orExpression()
}

// orExpression 解析 or 逻辑表达式
//
// or_expression ::= and_expression ("or" and_expression)*
func (p *Parser) orExpression() (ast.Expression, error) {
	tok := p.currToken
	expr, err := p.andExpression()
	if err != nil {
		return nil, err
	}
	for p.currTokenIs(token.OR) {
		op := p.currToken.Literal
		_ = p.eat(token.OR)
		right, err := p.andExpression()
		if err != nil {
			return nil, err
		}
		expr = &ast.BinaryOpExpression{
			Token:    tok,
			Left:     expr,
			Operator: op,
			Right:    right,
		}
	}
	return expr, nil
}

// andExpression 解析 and 逻辑表达式
//
// and_expression ::= not_expression ("and" not_expression)*
func (p *Parser) andExpression() (ast.Expression, error) {
	tok := p.currToken
	expr, err := p.notExpression()
	if err != nil {
		return nil, err
	}
	for p.currTokenIs(token.AND) {
		op := p.currToken.Literal
		_ = p.eat(token.AND)
		right, err := p.notExpression()
		if err != nil {
			return nil, err
		}
		expr = &ast.BinaryOpExpression{
			Token:    tok,
			Left:     expr,
			Operator: op,
			Right:    right,
		}
	}
	return expr, nil
}

// notExpression 解析 not 逻辑表达式
//
// not_expression ::= comparison_expression | "not" not_expression
func (p *Parser) notExpression() (ast.Expression, error) {
	if !p.currTokenIs(token.NOT) {
		return p.comparisonExpression()
	}
	tok := p.currToken
	op := p.currToken.Literal
	_ = p.eat(token.NOT)
	right, err := p.notExpression()
	if err != nil {
		return nil, err
	}
	expr := &ast.UnaryExpression{
		Token:    tok,
		Operator: op,
		Operand:  right,
	}
	return expr, nil
}

// comparisonExpression 解析关系表达式
//
// comparison_expression ::= bitwise_xor_expression (("<" | "<=" | ">" | ">=" | "!=" | "==" ) bitwise_xor_expression)*
func (p *Parser) comparisonExpression() (ast.Expression, error) {
	tok := p.currToken
	expr, err := p.bitwiseOrExpression()
	if err != nil {
		return nil, err
	}
	optypes := []token.TokenType{
		token.LESS_THAN, token.LESS_EQUAL_THAN, token.GREAT_THAN, token.GREAT_EQUAL_THAN,
		token.NOT_EQ, token.EQ,
	}
	for p.currTokenIn(optypes...) {
		op := p.currToken.Literal
		_ = p.eatIn(optypes...)
		right, err := p.bitwiseOrExpression()
		if err != nil {
			return nil, err
		}
		expr = &ast.BinaryOpExpression{
			Token:    tok,
			Left:     expr,
			Operator: op,
			Right:    right,
		}
	}
	return expr, nil
}

// bitwiseOrExpression 解析与操作表达式
//
// bitwise_or_expression ::= bitwise_xor_expression ( "|" bitwise_xor_expression)*
func (p *Parser) bitwiseOrExpression() (ast.Expression, error) {
	tok := p.currToken
	expr, err := p.bitwiseXorExpression()
	if err != nil {
		return nil, err
	}
	for p.currTokenIs(token.BITWISE_OR) {
		op := p.currToken.Literal
		_ = p.eat(token.BITWISE_OR)
		right, err := p.bitwiseXorExpression()
		if err != nil {
			return nil, err
		}
		expr = &ast.BinaryOpExpression{
			Token:    tok,
			Left:     expr,
			Operator: op,
			Right:    right,
		}
	}
	return expr, nil
}

// bitwiseXorExpression 解析与操作表达式
//
// bitwise_xor_expression ::= bitwise_and_expression ( "^" bitwise_and_expression)*
func (p *Parser) bitwiseXorExpression() (ast.Expression, error) {
	tok := p.currToken
	expr, err := p.bitwiseAndExpression()
	if err != nil {
		return nil, err
	}
	for p.currTokenIs(token.BITWISE_XOR) {
		op := p.currToken.Literal
		_ = p.eat(token.BITWISE_XOR)
		right, err := p.bitwiseAndExpression()
		if err != nil {
			return nil, err
		}
		expr = &ast.BinaryOpExpression{
			Token:    tok,
			Left:     expr,
			Operator: op,
			Right:    right,
		}
	}
	return expr, nil
}

// bitwiseAndExpression 解析与操作表达式
//
// bitwise_and_expression ::= shift_expression ( "&" shift_expression)*
func (p *Parser) bitwiseAndExpression() (ast.Expression, error) {
	tok := p.currToken
	expr, err := p.shiftExpression()
	if err != nil {
		return nil, err
	}
	for p.currTokenIs(token.BITWISE_AND) {
		op := p.currToken.Literal
		_ = p.eat(token.BITWISE_AND)
		right, err := p.shiftExpression()
		if err != nil {
			return nil, err
		}
		expr = &ast.BinaryOpExpression{
			Token:    tok,
			Left:     expr,
			Operator: op,
			Right:    right,
		}
	}
	return expr, nil
}

// shiftExpression 解析移位表达式
//
// shift_expression ::= plus_expression (( "<<" | ">>" ) plus_expression)*
func (p *Parser) shiftExpression() (ast.Expression, error) {
	tok := p.currToken
	expr, err := p.plusExpression()
	if err != nil {
		return nil, err
	}
	for p.currTokenIn(token.LEFT_SHIFT, token.RIGHT_SHIFT) {
		op := p.currToken.Literal
		_ = p.eatIn(token.LEFT_SHIFT, token.RIGHT_SHIFT)
		right, err := p.plusExpression()
		if err != nil {
			return nil, err
		}
		expr = &ast.BinaryOpExpression{
			Token:    tok,
			Left:     expr,
			Operator: op,
			Right:    right,
		}
	}
	return expr, nil
}

// plusExpression 解析加法类表达式
//
// plus_expression ::= multiplication_expression (("+" | "-") multiplication_expression)*
func (p *Parser) plusExpression() (ast.Expression, error) {
	tok := p.currToken
	expr, err := p.multiplyExpression()
	if err != nil {
		return nil, err
	}
	for p.currTokenIn(token.PLUS, token.MINUS) {
		op := p.currToken.Literal
		_ = p.eatIn(token.PLUS, token.MINUS)
		right, err := p.multiplyExpression()
		if err != nil {
			return nil, err
		}
		expr = &ast.BinaryOpExpression{
			Token:    tok,
			Left:     expr,
			Operator: op,
			Right:    right,
		}
	}
	return expr, nil
}

// multiplyExpression 解析乘法类表达式
//
// multiply_expression ::= unary_expression (("*" | "/" | "%") unary_expression)*
func (p *Parser) multiplyExpression() (ast.Expression, error) {
	tok := p.currToken
	expr, err := p.unaryExpression()
	if err != nil {
		return nil, err
	}
	for p.currTokenIn(token.ASTERISK, token.SLASH, token.MODULO) {
		op := p.currToken.Literal
		_ = p.eatIn(token.ASTERISK, token.SLASH, token.MODULO)
		right, err := p.unaryExpression()
		if err != nil {
			return nil, err
		}
		expr = &ast.BinaryOpExpression{
			Token:    tok,
			Left:     expr,
			Operator: op,
			Right:    right,
		}
	}
	return expr, nil
}

// unaryExpression 解析一元表达式
//
// unary_expression ::= primary_expression | ["-" | "+" | "~"] unary_expression
func (p *Parser) unaryExpression() (ast.Expression, error) {
	tok := p.currToken
	if !p.currTokenIn(token.MINUS, token.PLUS, token.BITWISE_NOT) {
		return p.primaryExpression()
	}
	op := tok.Literal
	err := p.eatIn(token.MINUS, token.PLUS, token.BITWISE_NOT)
	if err != nil {
		return nil, err
	}
	right, err := p.unaryExpression()
	if err != nil {
		return nil, err
	}
	expr := &ast.UnaryExpression{
		Token:    tok,
		Operator: op,
		Operand:  right,
	}
	return expr, nil
}

// primaryExpression 解析索引访问、属性访问、函数调用表达式
//
// primary_expression ::= atom ( subscription | attribute | call)*
// subscription       ::= "[" expression "]"
// attribute          ::= "." IDENT
// call               ::= "(" [argument_list] ")"
// argument_list      ::= expression ("," expression)* [","]
func (p *Parser) primaryExpression() (ast.Expression, error) {
	expr, err := p.atom()
	if err != nil {
		return nil, err
	}
	for p.currTokenIn(token.LBRACKET, token.DOT, token.LPAREN) {
		tok := p.currToken
		switch p.currToken.Type {
		case token.LBRACKET:
			_ = p.eat(token.LBRACKET)
			index, err := p.expression()
			if err != nil {
				return nil, err
			}
			expr = &ast.SubscriptionExpression{
				Token: tok,
				Left:  expr,
				Index: index,
			}
			err = p.eat(token.RBRACKET)
			if err != nil {
				return nil, err
			}
		case token.DOT:
			p.nextToken()
			ident, err := p.ident()
			if err != nil {
				return nil, err
			}
			expr = &ast.AttributeExpression{
				Token:     tok,
				Left:      expr,
				Attribute: ident,
			}
		case token.LPAREN:
			p.parenCount++
			p.nextToken()
			arguments, err := p.expressionList(token.RPAREN)
			if err != nil {
				return nil, err
			}
			p.parenCount--
			err = p.eat(token.RPAREN)
			if err != nil {
				return nil, err
			}
			expr = &ast.CallExpression{
				Token:     tok,
				Function:  expr,
				Arguments: arguments,
			}
		}
	}
	return expr, nil
}

// atom 解析表达式的基本单元
//
// atom ::= IDENT | INT_LIT | STRING_LIT | BOOL_LIT | NULL_LIT
// | list_literal | dict_literal | function_literal | "(" expression ")"
// | wei_expression
func (p *Parser) atom() (ast.Expression, error) {
	var expr ast.Expression
	var err error
	switch p.currToken.Type {
	case token.IDENT:
		return p.ident()
	case token.INT:
		var n int64
		literal := p.currToken.Literal
		prefix := ""
		if len(literal) > 2 {
			prefix = literal[:2]
		}
		start := 2
		base := 10
		bitSize := 64
		// 处理不同进制的数字
		switch prefix {
		case "0b":
			base = 2
		case "0o":
			base = 8
		case "0x":
			base = 16
		default:
			start = 0
		}
		n, err = strconv.ParseInt(literal[start:], base, bitSize)
		if err != nil {
			return nil, p.syntaxError(err.Error())
		}
		expr = &ast.IntegerLiteral{
			Token: p.currToken,
			Value: n,
		}
		_ = p.eat(token.INT)
	case token.STRING:
		expr = &ast.StringLiteral{
			Token: p.currToken,
			Value: p.currToken.Literal,
		}
		_ = p.eat(token.STRING)
	case token.TRUE, token.FALSE:
		expr = &ast.Boolean{
			Token: p.currToken,
			Value: p.currToken.Literal == "true",
		}
		_ = p.eatIn(token.TRUE, token.FALSE)
	case token.NULL:
		expr = &ast.NullLiteral{Token: p.currToken}
		p.nextToken()
	case token.LPAREN:
		p.parenCount++
		p.nextToken()
		expr, err = p.expression()
		if err != nil {
			return nil, err
		}
		p.parenCount--
		err = p.eat(token.RPAREN)
		if err != nil {
			return nil, err
		}
	case token.LBRACKET:
		return p.listLiteral()
	case token.LBRACE:
		return p.dictLiteral()
	case token.FUNCTION:
		return p.functionLiteral()
	case token.WEI:
		return p.weiExpression()
	case token.ILLEGAL:
		return nil, p.syntaxError(p.currToken.Literal)
	default:
		return nil, p.invalidError()
	}
	return expr, nil
}

// listLiteral 解析列表字面量
//
// list_literal ::= "[" [expression] ("," expression)* [","] "]"
func (p *Parser) listLiteral() (*ast.ListLiteral, error) {
	tok := p.currToken
	p.parenCount++
	err := p.eat(token.LBRACKET)
	if err != nil {
		return nil, err
	}
	elements, err := p.expressionList(token.RBRACKET)
	if err != nil {
		return nil, err
	}
	p.parenCount--
	err = p.eat(token.RBRACKET)
	if err != nil {
		return nil, err
	}
	expr := &ast.ListLiteral{
		Token:    tok,
		Elements: elements,
	}
	return expr, nil
}

// expression_list ::= [expression] ("," expression)* [","]
func (p *Parser) expressionList(end token.TokenType) ([]ast.Expression, error) {
	var elements []ast.Expression
	var ele ast.Expression
	var err error
	// 有下面几种情况
	//  <空>
	//  expr
	//  expr,
	//  expr1,expr2
	//  expr1,expr2,
	if p.currTokenNotIs(end) {
		ele, err = p.expression()
		if err != nil {
			return nil, err
		}
		elements = append(elements, ele)
	}
	for p.currTokenIs(token.COMMA) {
		p.nextToken()
		if p.currTokenIs(end) {
			break
		}
		ele, err = p.expression()
		if err != nil {
			return nil, err
		}
		elements = append(elements, ele)
	}
	if p.currTokenIs(token.COMMA) {
		p.nextToken()
	}
	return elements, nil
}

// dictLiteral 解析字典字面量
//
// dict_literal ::= "{" [ pairs ] "}"
// pairs        ::= [pair ("," pair)* [","]
// pair         ::= expression ":" expression
func (p *Parser) dictLiteral() (*ast.DictLiteral, error) {
	tok := p.currToken
	p.parenCount++
	err := p.eat(token.LBRACE)
	if err != nil {
		return nil, err
	}
	var key ast.Expression
	var val ast.Expression
	pairs := make(map[ast.Expression]ast.Expression)
	if !p.currTokenIs(token.RBRACE) {
		key, err = p.expression()
		if err != nil {
			return nil, err
		}
		err = p.eat(token.COLON)
		if err != nil {
			return nil, err
		}
		val, err = p.expression()
		if err != nil {
			return nil, err
		}
		pairs[key] = val
	}
	for p.currTokenIs(token.COMMA) {
		_ = p.eat(token.COMMA)
		if p.currTokenIs(token.RBRACE) {
			break
		}
		key, err = p.expression()
		if err != nil {
			return nil, err
		}
		err = p.eat(token.COLON)
		if err != nil {
			return nil, err
		}
		val, err = p.expression()
		if err != nil {
			return nil, err
		}
		pairs[key] = val
	}
	if p.currTokenIs(token.COMMA) {
		_ = p.eat(token.COMMA)
	}
	p.parenCount--
	err = p.eat(token.RBRACE)
	if err != nil {
		return nil, err
	}
	expr := &ast.DictLiteral{
		Token: tok,
		Pairs: pairs,
	}
	return expr, nil
}

// functionLiteral 解析函数定义
//
// function_literal ::= "fn" "(" parameter_list ")" block_statement
func (p *Parser) functionLiteral() (*ast.FunctionLiteral, error) {
	tok := p.currToken
	err := p.eat(token.FUNCTION)
	if err != nil {
		return nil, err
	}
	p.parenCount++
	err = p.eat(token.LPAREN)
	if err != nil {
		return nil, err
	}
	paramters, err := p.parameterList()
	if err != nil {
		return nil, err
	}
	p.parenCount--
	err = p.eat(token.RPAREN)
	if err != nil {
		return nil, err
	}
	p.whileStack = append(p.whileStack, 0)
	block, err := p.blockStatement()
	if err != nil {
		return nil, err
	}
	p.whileStack = p.whileStack[:len(p.whileStack)-1]
	fl := &ast.FunctionLiteral{
		Token:      tok,
		Parameters: paramters,
		Body:       block,
	}
	return fl, nil
}

// parameter_list ::= [IDENT] ("," IDENT)* [","]
func (p *Parser) parameterList() ([]*ast.Identifier, error) {
	// 主要有下面这些情况
	// ()
	// (a)
	// (a,)
	// (a,b)
	// (a,b,)
	var parameters []*ast.Identifier
	var param *ast.Identifier
	var err error
	if !p.currTokenIs(token.RPAREN) {
		param, err = p.ident()
		if err != nil {
			return nil, err
		}
		parameters = append(parameters, param)
	}
	for p.currTokenIs(token.COMMA) {
		p.nextToken()
		// 只有一个参数，后面有个逗号
		if !p.currTokenIs(token.IDENT) {
			break
		}
		param, err = p.ident()
		if err != nil {
			return nil, err
		}
		parameters = append(parameters, param)
	}
	if p.currTokenIs(token.COMMA) {
		p.nextToken()
	}
	return parameters, nil
}

// wei_expression ::= ( "wei" "." IDENT ) | ( "wei" "." "import" "(" STRING_LIT ")" )
func (p *Parser) weiExpression() (ast.Expression, error) {
	tk := p.currToken
	err := p.eat(token.WEI)
	if err != nil {
		return nil, err
	}
	err = p.eat(token.DOT)
	if err != nil {
		return nil, err
	}
	if p.currTokenLiteralIs("import") {
		p.nextToken()
		err = p.eat(token.LPAREN)
		if err != nil {
			return nil, err
		}
		filename := p.currToken.Literal
		err = p.eat(token.STRING)
		if err != nil {
			return nil, err
		}
		err = p.eat(token.RPAREN)
		if err != nil {
			return nil, err
		}
		expr := &ast.WeiImportExpression{
			Token:    tk,
			Filename: filename,
		}
		return expr, nil
	}
	attribute, err := p.ident()
	if err != nil {
		return nil, err
	}
	expr := &ast.WeiAttributeExpression{
		Token:     tk,
		Attribute: attribute,
	}
	return expr, nil
}

func (p *Parser) ident() (*ast.Identifier, error) {
	identifier := &ast.Identifier{
		Location: p.currFileLocation(),
		Token:    p.currToken,
		Value:    p.currToken.Literal,
	}
	err := p.eat(token.IDENT)
	if err != nil {
		return nil, err
	}
	return identifier, nil
}

func (p *Parser) eatIn(ts ...token.TokenType) error {
	if p.currTokenIn(ts...) {
		p.nextToken()
		return nil
	} else if p.currTokenIs(token.ILLEGAL) {
		return p.syntaxError(p.currToken.Literal)
	} else {
		return p.expectError(ts[0])
	}
}

func (p *Parser) eatContinuously(sequence ...token.TokenType) error {
	var err error
	for _, t := range sequence {
		err = p.eat(t)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Parser) eat(t token.TokenType) error {
	if p.currTokenIs(t) {
		p.nextToken()
		return nil
	} else if p.currTokenIs(token.ILLEGAL) {
		return p.syntaxError(p.currToken.Literal)
	} else {
		return p.expectError(t)
	}
}

func (p *Parser) dump() *dumpInfo {
	return &dumpInfo{
		index: p.l.Dump(),
		token: p.currToken,
	}
}

func (p *Parser) restore(info *dumpInfo) {
	p.l.Restore(info.index)
	p.currToken = info.token
}

func (p *Parser) nextToken() {
	p.doNextToken()
	if p.parenCount > 0 {
		p.skipNewline()
	}
}

func (p *Parser) doNextToken() {
	p.currToken = p.l.NextToken()
	//fmt.Printf("%d,%d-%d,%d:\t%s\t%q\n",
	//	p.currToken.Start.Line, p.currToken.Start.Column,
	//	p.currToken.End.Line, p.currToken.End.Column,
	//	p.currToken.Type, p.currToken.Literal,
	//)
	p.skipComment()
}

func (p *Parser) skipNewline() {
	for p.currTokenIs(token.NEWLINE) {
		p.doNextToken()
	}
}

func (p *Parser) skipComment() {
	for p.currTokenIs(token.COMMENT) {
		p.currToken = p.l.NextToken()
	}
}

func (p *Parser) currTokenLiteralIs(s string) bool {
	return p.currToken.LiteralIs(s)
}

func (p *Parser) currTokenIn(ts ...token.TokenType) bool {
	return p.currToken.TypeIn(ts...)
}

func (p *Parser) currTokenIs(t token.TokenType) bool {
	return p.currToken.TypeIs(t)
}

func (p *Parser) currTokenNotIs(t token.TokenType) bool {
	return p.currToken.TypeNotIs(t)
}

func (p *Parser) currFileLocation() *ast.FileLocation {
	return ast.NewFileLocation(p.filename, p.currToken.Start.Line)
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	peek := p.l.PeekToken()
	return peek.TypeIs(t)
}

func (p *Parser) expectError(expected token.TokenType) error {
	msg := fmt.Sprintf("expected %q, but got %q", expected, p.currToken.Type)
	return p.syntaxError(msg)
}

func (p *Parser) invalidError() error {
	return p.syntaxError("invalid syntax")
}

func (p *Parser) syntaxError(msg string) error {
	// 标注错误的位置
	template := `
File "%s", line %d
  %s
  %s
SyntaxError: %s`
	line := p.currToken.Start.Line
	column := p.currToken.Start.Column
	return fmt.Errorf(template, p.filename, line+1, p.lines[line], strRjust("^", column+1), msg)
}

func strRjust(s string, n int) string {
	if len(s) >= n {
		return s
	}
	return strings.Repeat(" ", n-len(s)) + s
}
