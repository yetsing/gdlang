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

// expression ::= atom (("+" | "-") atom)*
func (p *Parser) expression() (ast.Expression, error) {
	tok := p.currToken
	expr, err := p.atom()
	if err != nil {
		return nil, err
	}
	for p.currTokenIn(token.PLUS, token.MINUS) {
		op := p.currToken.Literal
		_ = p.eatIn(token.PLUS, token.MINUS)
		right, err := p.atom()
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

// atom ::= IDENT | INT_LIT | STRING_LIT
// 语法中的最小单元
func (p *Parser) atom() (ast.Expression, error) {
	var expr ast.Expression
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
		n, err := strconv.ParseInt(literal[start:], base, bitSize)
		if err != nil {
			return nil, err
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
	default:
		return nil, p.syntaxError("invalid syntax")
	}
	return expr, nil
}

func (p *Parser) eatIn(ts ...token.TokenType) error {
	if p.currTokenIn(ts...) {
		p.nextToken()
		return nil
	} else {
		return p.expectError(ts[0])
	}
}

func (p *Parser) eat(t token.TokenType) error {
	if p.currTokenIs(t) {
		p.nextToken()
		return nil
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
	return fmt.Errorf("expected %s, but got %s", expected, p.currToken.Type)
}

func (p *Parser) syntaxError(msg string) error {
	return fmt.Errorf(msg)
}
