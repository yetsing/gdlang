package parser

import (
	"fmt"
	"strconv"
	"weilang/ast"
	"weilang/lexer"
	"weilang/token"
)

type Parser struct {
	l *lexer.Lexer

	currToken token.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l: l,
	}
	p.nextToken()
	return p
}

/*
使用 lsbasi 的解析实现 (https://github.com/rspivak/lsbasi)
相比 《用Go语言自制解释器》 的实现更容易理解
*/

func (p *Parser) ParseProgram() (*ast.Program, error) {
	node, err := p.program()
	if err != nil {
		return nil, err
	}
	if !p.currTokenIs(token.EOF) {
		return nil, p.expectError(token.EOF)
	}
	return node, nil
}

// program ::= (statement)*
func (p *Parser) program() (*ast.Program, error) {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.currTokenIs(token.EOF) {
		stmt, err := p.statement()
		if err != nil {
			return nil, err
		}
		program.Statements = append(program.Statements, stmt)
	}

	return program, nil
}

// statement ::= expr [";"]
func (p *Parser) statement() (ast.Statement, error) {
	stmt := &ast.ExpressionStatement{Token: p.currToken}
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	stmt.Expression = expr
	if p.currTokenIs(token.SEMICOLON) {
		_ = p.eat(token.SEMICOLON)
	}
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
// not_expression ::= ["not"] comparison_expression
func (p *Parser) notExpression() (ast.Expression, error) {
	if !p.currTokenIs(token.NOT) {
		return p.comparisonExpression()
	}
	tok := p.currToken
	op := p.currToken.Literal
	_ = p.eat(token.NOT)
	right, err := p.comparisonExpression()
	if err != nil {
		return nil, err
	}
	expr := &ast.UnaryExpression{
		Token:    tok,
		Operator: op,
		Right:    right,
	}
	return expr, nil
}

// comparisonExpression 解析关系表达式
//
// comparison_expression ::= shift_expression (("<" | "<=" | ">" | ">=" | "!=" | "==" ) shift_expression)*
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
// unary_expression ::= [("-" | "~")] atom
func (p *Parser) unaryExpression() (ast.Expression, error) {
	tok := p.currToken
	if !p.currTokenIn(token.MINUS, token.BITWISE_NOT) {
		return p.atom()
	}
	op := tok.Literal
	err := p.eatIn(token.MINUS, token.BITWISE_NOT)
	if err != nil {
		return nil, err
	}
	right, err := p.atom()
	if err != nil {
		return nil, err
	}
	expr := &ast.UnaryExpression{
		Token:    tok,
		Operator: op,
		Right:    right,
	}
	return expr, nil
}

// atom 解析表达式的基本单元
//
// atom ::= IDENT | INT_LIT | STRING_LIT | BOOL_INT | "(" expression ")"
func (p *Parser) atom() (ast.Expression, error) {
	var expr ast.Expression
	var err error
	switch p.currToken.Type {
	case token.IDENT:
		expr = &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
		_ = p.eat(token.IDENT)
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
	case token.LPAREN:
		_ = p.eat(token.LPAREN)
		expr, err = p.expression()
		if err != nil {
			return nil, err
		}
		err = p.eat(token.RPAREN)
		if err != nil {
			return nil, err
		}
	case token.ILLEGAL:
		return nil, p.syntaxError(p.currToken.Literal)
	default:
		return nil, p.syntaxError("invalid syntax")
	}
	return expr, nil
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

func (p *Parser) nextToken() {
	p.currToken = p.l.NextToken()
}

func (p *Parser) currTokenIn(ts ...token.TokenType) bool {
	return p.currToken.TypeIn(ts...)
}

func (p *Parser) currTokenIs(t token.TokenType) bool {
	return p.currToken.TypeIs(t)
}

func (p *Parser) expectError(expected token.TokenType) error {
	msg := fmt.Sprintf("expected %s, but got %s", expected, p.currToken.Type)
	return p.syntaxError(msg)
}

func (p *Parser) syntaxError(msg string) error {
	return fmt.Errorf("SyntaxError: %s", msg)
}
