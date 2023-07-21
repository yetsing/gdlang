package ast

import (
	"bytes"
)

type FileLocation struct {
	Filename string
	Lineno   int
}

func NewFileLocation(filename string, lineno int) *FileLocation {
	return &FileLocation{
		Filename: filename,
		Lineno:   lineno,
	}
}

// The base Node interface
type Node interface {
	GetFileLocation() *FileLocation
	TokenLiteral() string
	String() string
}

type Program struct {
	Location   *FileLocation
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}
func (p *Program) GetFileLocation() *FileLocation {
	return p.Location
}

// All statement nodes implement this
type Statement interface {
	Node
	statementNode()
}

// All expression nodes implement this
type Expression interface {
	Node
	expressionNode()
}
