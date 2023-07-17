package ast

import (
	"bytes"
	"fmt"
	"strings"

	"weilang/token"
)

// Expression 下面的是表达式节点

// UnaryExpression 一元操作表达式
type UnaryExpression struct {
	Token    token.Token // The prefix token, e.g. !
	Operator string
	Operand  Expression
}

func (pe *UnaryExpression) expressionNode()      {}
func (pe *UnaryExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *UnaryExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Operand.String())
	out.WriteString(")")

	return out.String()
}

// BinaryOpExpression 二元操作表达式，如 "1 + 2"
type BinaryOpExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (be *BinaryOpExpression) expressionNode()      {}
func (be *BinaryOpExpression) TokenLiteral() string { return be.Token.Literal }
func (be *BinaryOpExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(be.Left.String())
	out.WriteString(" " + be.Operator + " ")
	out.WriteString(be.Right.String())
	out.WriteString(")")

	return out.String()
}

type SubscriptionExpression struct {
	Token token.Token
	Left  Expression
	Index Expression
}

func (se *SubscriptionExpression) expressionNode()      {}
func (se *SubscriptionExpression) TokenLiteral() string { return se.Token.Literal }
func (se *SubscriptionExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(se.Left.String())
	out.WriteString("[")
	out.WriteString(se.Index.String())
	out.WriteString("])")

	return out.String()
}

type AttributeExpression struct {
	Token     token.Token
	Left      Expression
	Attribute *Identifier
}

func (ae *AttributeExpression) expressionNode()      {}
func (ae *AttributeExpression) TokenLiteral() string { return ae.Token.Literal }
func (ae *AttributeExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ae.Left.String())
	out.WriteString(".")
	out.WriteString(ae.Attribute.String())
	out.WriteString(")")

	return out.String()
}

type CallExpression struct {
	Token     token.Token
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var out bytes.Buffer

	var args []string
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

type FunctionLiteral struct {
	Token      token.Token // The 'fn' token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	var params []string
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fl.Body.String())

	return out.String()
}

type ListLiteral struct {
	Token    token.Token // the '[' token
	Elements []Expression
}

func (al *ListLiteral) expressionNode()      {}
func (al *ListLiteral) TokenLiteral() string { return al.Token.Literal }
func (al *ListLiteral) String() string {
	var out bytes.Buffer

	elements := []string{}
	for _, el := range al.Elements {
		elements = append(elements, el.String())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

type DictLiteral struct {
	Token token.Token // the '{' token
	Pairs map[Expression]Expression
}

func (hl *DictLiteral) expressionNode()      {}
func (hl *DictLiteral) TokenLiteral() string { return hl.Token.Literal }
func (hl *DictLiteral) String() string {
	var out bytes.Buffer

	pairs := []string{}
	for key, value := range hl.Pairs {
		pairs = append(pairs, key.String()+":"+value.String())
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string       { return b.Token.Literal }

type NullLiteral struct {
	Token token.Token
}

func (n *NullLiteral) expressionNode()      {}
func (n *NullLiteral) TokenLiteral() string { return n.Token.Literal }
func (n *NullLiteral) String() string       { return n.Token.Literal }

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return sl.Token.Literal }

type WeiAttributeExpression struct {
	Token     token.Token
	Attribute *Identifier
}

func (wa *WeiAttributeExpression) expressionNode()      {}
func (wa *WeiAttributeExpression) TokenLiteral() string { return wa.Token.Literal }
func (wa *WeiAttributeExpression) String() string {
	return fmt.Sprintf("(wei.%s)", wa.Attribute.String())
}

type WeiImportExpression struct {
	Token    token.Token
	Filename Expression
}

func (w *WeiImportExpression) expressionNode() {}
func (w *WeiImportExpression) TokenLiteral() string {
	return w.Token.Literal
}
func (w *WeiImportExpression) String() string {
	return fmt.Sprintf("(wei.import(%s))", w.Filename.String())
}
