package interpreter

import (
	"fmt"
	"weilang/evaluator"
	"weilang/lexer"
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

	env := evaluator.NewWeiEnvironment(filename)
	evaluated := evaluator.Eval(program, env)
	if evaluator.IsError(evaluated) {
		fmt.Println(evaluated.String())
	}
}
