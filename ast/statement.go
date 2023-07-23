package ast

import (
	"bytes"
	"fmt"
	"strings"
	"weilang/token"
)

type VarStatement struct {
	Location *FileLocation
	Token    token.Token
	Name     *Identifier
	Value    Expression
}

func (vs *VarStatement) statementNode()       {}
func (vs *VarStatement) TokenLiteral() string { return vs.Token.Literal }
func (vs *VarStatement) String() string {
	var out bytes.Buffer

	out.WriteString(vs.TokenLiteral() + " ")
	out.WriteString(vs.Name.String())
	out.WriteString(" = ")

	out.WriteString(vs.Value.String())

	out.WriteString(";")

	return out.String()
}
func (vs *VarStatement) GetFileLocation() *FileLocation {
	return vs.Location
}

type ConStatement struct {
	Location *FileLocation
	Token    token.Token
	Name     *Identifier
	Value    Expression
}

func (cs *ConStatement) statementNode()       {}
func (cs *ConStatement) TokenLiteral() string { return cs.Token.Literal }
func (cs *ConStatement) String() string {
	var out bytes.Buffer

	out.WriteString(cs.TokenLiteral() + " ")
	out.WriteString(cs.Name.String())
	out.WriteString(" = ")

	out.WriteString(cs.Value.String())

	out.WriteString(";")

	return out.String()
}
func (cs *ConStatement) GetFileLocation() *FileLocation {
	return cs.Location
}

type AssignStatement struct {
	Location *FileLocation
	Token    token.Token
	Left     Expression
	Value    Expression
}

func (as *AssignStatement) statementNode()       {}
func (as *AssignStatement) TokenLiteral() string { return as.Token.Literal }
func (as *AssignStatement) String() string {
	var out bytes.Buffer

	out.WriteString(as.Left.String())
	out.WriteString(" = ")

	out.WriteString(as.Value.String())

	out.WriteString(";")

	return out.String()
}
func (as *AssignStatement) GetFileLocation() *FileLocation {
	return as.Location
}

type ReturnStatement struct {
	Location    *FileLocation
	Token       token.Token // the 'return' token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
}
func (rs *ReturnStatement) GetFileLocation() *FileLocation {
	return rs.Location
}

type ExpressionStatement struct {
	Location   *FileLocation
	Token      token.Token // the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}
func (es *ExpressionStatement) GetFileLocation() *FileLocation {
	return es.Location
}

type BlockStatement struct {
	Location   *FileLocation
	Token      token.Token // the { token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	out.WriteString("{\n")
	for _, s := range bs.Statements {
		out.WriteString(s.String())
		out.WriteString("\n")
	}
	out.WriteString("}")

	return out.String()
}
func (bs *BlockStatement) GetFileLocation() *FileLocation {
	return bs.Location
}

type IfBranch struct {
	Location *FileLocation
	// Token "if" token
	Condition Expression
	Body      *BlockStatement
}

func (i *IfBranch) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString("(")
	out.WriteString(i.Condition.String())
	out.WriteString(")")
	out.WriteString(i.Body.String())
	return out.String()
}

type IfStatement struct {
	Location *FileLocation
	// Token "if" token
	Token token.Token
	// Cases 对应多个 "if" "else if" 等多个条件分支
	IfBranches []*IfBranch
	ElseBody   *BlockStatement
}

func (is *IfStatement) statementNode()       {}
func (is *IfStatement) TokenLiteral() string { return is.Token.Literal }
func (is *IfStatement) String() string {
	var out bytes.Buffer

	var ms []string
	for _, ifCase := range is.IfBranches {
		ms = append(ms, ifCase.String())
	}

	out.WriteString(strings.Join(ms, " else "))
	if is.ElseBody != nil {
		out.WriteString(" else ")
		out.WriteString(is.ElseBody.String())
	}

	return out.String()
}
func (is *IfStatement) GetFileLocation() *FileLocation {
	return is.Location
}

type WhileStatement struct {
	Location *FileLocation
	// Token "while" token
	Token     token.Token
	Condition Expression
	Body      *BlockStatement
}

func (ws *WhileStatement) statementNode()       {}
func (ws *WhileStatement) TokenLiteral() string { return ws.Token.Literal }
func (ws *WhileStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ws.Token.Literal)
	out.WriteString("(")
	out.WriteString(ws.Condition.String())
	out.WriteString(")")
	out.WriteString(ws.Body.String())

	return out.String()
}
func (ws *WhileStatement) GetFileLocation() *FileLocation {
	return ws.Location
}

type ContinueStatement struct {
	Location *FileLocation
	Token    token.Token
}

func (c *ContinueStatement) statementNode()       {}
func (c *ContinueStatement) TokenLiteral() string { return c.Token.Literal }
func (c *ContinueStatement) String() string {
	return "continue"
}
func (c *ContinueStatement) GetFileLocation() *FileLocation {
	return c.Location
}

type BreakStatement struct {
	Location *FileLocation
	Token    token.Token
}

func (b *BreakStatement) statementNode()       {}
func (b *BreakStatement) TokenLiteral() string { return b.Token.Literal }
func (b *BreakStatement) String() string {
	return "break"
}
func (b *BreakStatement) GetFileLocation() *FileLocation {
	return b.Location
}

type ForInStatement struct {
	Location *FileLocation
	Token    token.Token
	Con      bool
	Targets  []*Identifier
	Expr     Expression
	Body     *BlockStatement
}

func (f *ForInStatement) statementNode()       {}
func (f *ForInStatement) TokenLiteral() string { return f.Token.Literal }
func (f *ForInStatement) String() string {
	var out bytes.Buffer

	out.WriteString("for (")
	if f.Con {
		out.WriteString("con ")
	} else {
		out.WriteString("var ")
	}
	var elements []string
	for _, target := range f.Targets {
		elements = append(elements, target.String())
	}
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("in ")
	out.WriteString(f.Expr.String())
	out.WriteString(") ")
	out.WriteString(f.Body.String())
	return out.String()
}
func (f *ForInStatement) GetFileLocation() *FileLocation {
	return f.Location
}

type FunctionDefineStatement struct {
	Location *FileLocation
	Token    token.Token
	Function *FunctionLiteral
}

func (fs *FunctionDefineStatement) statementNode() {}

func (fs *FunctionDefineStatement) TokenLiteral() string { return fs.Token.Literal }

func (fs *FunctionDefineStatement) String() string {
	return fs.Function.String()
}
func (fs *FunctionDefineStatement) GetFileLocation() *FileLocation {
	return fs.Location
}

type ClassVariableDeclarationStatement struct {
	Location *FileLocation
	Token    token.Token
	Con      bool
	Class    bool
	Name     *Identifier
	Expr     Expression
}

func (cv *ClassVariableDeclarationStatement) statementNode() {}

func (cv *ClassVariableDeclarationStatement) TokenLiteral() string { return cv.Token.Literal }

func (cv *ClassVariableDeclarationStatement) String() string {
	var out bytes.Buffer

	if cv.Con {
		out.WriteString("con ")
	} else {
		out.WriteString("var ")
	}
	if cv.Class {
		out.WriteString("class.")
	}
	out.WriteString(cv.Name.String())
	if cv.Expr != nil {
		out.WriteString(" = " + cv.Expr.String())
	}
	out.WriteString("\n")
	return out.String()
}
func (cv *ClassVariableDeclarationStatement) GetFileLocation() *FileLocation {
	return cv.Location
}

type ClassMethodDefineStatement struct {
	Location *FileLocation
	Token    token.Token
	Class    bool
	Function *FunctionLiteral
}

func (cm *ClassMethodDefineStatement) statementNode() {}

func (cm *ClassMethodDefineStatement) TokenLiteral() string { return cm.Token.Literal }

func (cm *ClassMethodDefineStatement) String() string {
	var out bytes.Buffer

	if cm.Class {
		out.WriteString("class.")
	}
	out.WriteString(cm.Function.Name)

	var params []string
	for _, p := range cm.Function.Parameters {
		params = append(params, p.String())
	}
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(cm.Function.Body.String())

	return out.String()
}
func (cm *ClassMethodDefineStatement) GetFileLocation() *FileLocation {
	return cm.Location
}

type ClassBlockStatement struct {
	Location   *FileLocation
	Token      token.Token
	Statements []Statement
}

func (cb *ClassBlockStatement) statementNode() {}

func (cb *ClassBlockStatement) TokenLiteral() string { return cb.Token.Literal }

func (cb *ClassBlockStatement) String() string {
	var out bytes.Buffer

	out.WriteString("{\n")
	var lines []string
	for _, statement := range cb.Statements {
		lines = append(lines, statement.String())
	}
	out.WriteString(strings.Join(lines, "\n"))
	out.WriteString("\n}")
	out.WriteString("\n")
	return out.String()
}
func (cb *ClassBlockStatement) GetFileLocation() *FileLocation {
	return cb.Location
}

type ClassDefineStatement struct {
	Location *FileLocation
	// "class" token
	Token token.Token
	Name  string
	Body  *ClassBlockStatement
}

func (cd *ClassDefineStatement) statementNode() {}

func (cd *ClassDefineStatement) TokenLiteral() string { return cd.Token.Literal }

func (cd *ClassDefineStatement) String() string {
	var out bytes.Buffer

	out.WriteString(
		fmt.Sprintf("class %s ", cd.Name))
	out.WriteString(cd.Body.String())
	return out.String()
}
func (cd *ClassDefineStatement) GetFileLocation() *FileLocation {
	return cd.Location
}

type WeiExportStatement struct {
	Location *FileLocation
	Token    token.Token
	Names    []*Identifier
}

func (w *WeiExportStatement) statementNode() {}

func (w *WeiExportStatement) TokenLiteral() string { return w.Token.Literal }

func (w *WeiExportStatement) String() string {
	var out bytes.Buffer

	out.WriteString("wei.export(")
	var names []string
	for _, n := range w.Names {
		names = append(names, n.String())
	}
	out.WriteString(strings.Join(names, ","))
	out.WriteString(")")
	out.WriteString("\n")
	return out.String()
}
func (w *WeiExportStatement) GetFileLocation() *FileLocation {
	return w.Location
}
