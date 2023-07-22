package interpreter

import (
	"context"
	"fmt"
	"path/filepath"
	"weilang/evaluator"
	"weilang/lexer"
	"weilang/object"
	"weilang/parser"
)

func RunFile(filename string) {
	filename, _ = filepath.Abs(filename)
	l := lexer.NewWithFilename(filename)
	p := parser.New(l)
	program, err := p.ParseProgram()
	if err != nil {
		fmt.Println(err)
		return
	}

	mod := object.NewModule(filename)
	evaluator.CacheModule(mod)
	state := evaluator.NewWeiState(mod)
	state.CreateFrame(filename, "<module>")
	evaluated := evaluator.Eval(context.Background(), state, program, mod.GetEnv())
	if evaluator.IsError(evaluated) {
		if state.HasExc() {
			state.PrintExc()
		} else {
			fmt.Println(evaluated.String())
		}
	}
}
