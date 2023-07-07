package object

import (
	"fmt"
	"weilang/ast"
)

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func NewFunction(fl *ast.FunctionLiteral, env *Environment) *Function {
	return &Function{
		Parameters: fl.Parameters,
		Body:       fl.Body,
		Env:        env,
	}
}

func (f *Function) Type() ObjectType {
	return FUNCTION_OBJ
}

func (f *Function) TypeIs(objectType ObjectType) bool {
	return f.Type() == objectType
}

func (f *Function) TypeNotIs(objectType ObjectType) bool {
	return f.Type() != objectType
}

func (f *Function) String() string {
	return fmt.Sprintf("<function at %p>", f)
}
