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
	Location *FileLocation
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
func (pe *UnaryExpression) GetFileLocation() *FileLocation {
	return pe.Location
}

// BinaryOpExpression 二元操作表达式，如 "1 + 2"
type BinaryOpExpression struct {
	Location *FileLocation
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
func (be *BinaryOpExpression) GetFileLocation() *FileLocation {
	return be.Location
}

type SubscriptionExpression struct {
	Location *FileLocation
	Token    token.Token
	Left     Expression
	Index    Expression
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
func (se *SubscriptionExpression) GetFileLocation() *FileLocation {
	return se.Location
}

type AttributeExpression struct {
	Location  *FileLocation
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
func (ae *AttributeExpression) GetFileLocation() *FileLocation {
	return ae.Location
}

type CallExpression struct {
	Location  *FileLocation
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
func (ce *CallExpression) GetFileLocation() *FileLocation {
	return ce.Location
}

type FunctionLiteral struct {
	Location   *FileLocation
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
func (fl *FunctionLiteral) GetFileLocation() *FileLocation {
	return fl.Location
}

type ListLiteral struct {
	Location *FileLocation
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
func (al *ListLiteral) GetFileLocation() *FileLocation {
	return al.Location
}

type DictLiteral struct {
	Location *FileLocation
	Token    token.Token // the '{' token
	Pairs    map[Expression]Expression
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
func (hl *DictLiteral) GetFileLocation() *FileLocation {
	return hl.Location
}

type Identifier struct {
	Location *FileLocation
	Token    token.Token // the token.IDENT token
	Value    string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }
func (i *Identifier) GetFileLocation() *FileLocation {
	return i.Location
}

type Boolean struct {
	Location *FileLocation
	Token    token.Token
	Value    bool
}

func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string       { return b.Token.Literal }
func (b *Boolean) GetFileLocation() *FileLocation {
	return b.Location
}

type NullLiteral struct {
	Location *FileLocation
	Token    token.Token
}

func (n *NullLiteral) expressionNode()      {}
func (n *NullLiteral) TokenLiteral() string { return n.Token.Literal }
func (n *NullLiteral) String() string       { return n.Token.Literal }
func (n *NullLiteral) GetFileLocation() *FileLocation {
	return n.Location
}

type IntegerLiteral struct {
	Location *FileLocation
	Token    token.Token
	Value    int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }
func (il *IntegerLiteral) GetFileLocation() *FileLocation {
	return il.Location
}

type StringLiteral struct {
	Location *FileLocation
	Token    token.Token
	Value    string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return sl.Token.Literal }
func (sl *StringLiteral) GetFileLocation() *FileLocation {
	return sl.Location
}

type WeiAttributeExpression struct {
	Location  *FileLocation
	Token     token.Token
	Attribute *Identifier
}

func (wa *WeiAttributeExpression) expressionNode()      {}
func (wa *WeiAttributeExpression) TokenLiteral() string { return wa.Token.Literal }
func (wa *WeiAttributeExpression) String() string {
	return fmt.Sprintf("(wei.%s)", wa.Attribute.String())
}
func (wa *WeiAttributeExpression) GetFileLocation() *FileLocation {
	return wa.Location
}

type WeiImportExpression struct {
	Location *FileLocation
	Token    token.Token
	Filename string
}

func (w *WeiImportExpression) expressionNode() {}
func (w *WeiImportExpression) TokenLiteral() string {
	return w.Token.Literal
}
func (w *WeiImportExpression) String() string {
	return fmt.Sprintf("(wei.import(%s))", w.Filename)
}
func (w *WeiImportExpression) GetFileLocation() *FileLocation {
	return w.Location
}
