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

	mod := object.NewModule("")
	ctx := evaluator.NewModuleContext(mod)
	evaluated := evaluator.Eval(ctx, program, mod.GetEnv())
	if evaluator.IsError(evaluated) {
		fmt.Println(evaluated.String())
	}
}
