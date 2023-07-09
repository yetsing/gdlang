package interpreter

import (
	"fmt"
	"weilang/evaluator"
	"weilang/lexer"
	"weilang/object"
	"weilang/parser"
)

func RunFile(filename string) {
	l := lexer.NewWithFilename(filename)
	p := parser.New(l)
	program, err := p.ParseProgram()
	if err != nil {
		fmt.Println(err)
		return
	}

	env := object.NewEnvironment()
	evaluated := evaluator.Eval(program, env)
	if evaluator.IsError(evaluated) {
		fmt.Println(evaluated.String())
	}
}
