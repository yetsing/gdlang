package repl

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"weilang/ast"
	"weilang/evaluator"
	"weilang/lexer"
	"weilang/object"
	"weilang/parser"
)

const (
	START_PROMPT = ">> "
	//CONTINUE_PROMPT = "..."
)

var PROMPT = START_PROMPT

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	mod := object.NewModule("")
	state := evaluator.NewWeiState(mod)
	ctx := context.Background()

	var buffer bytes.Buffer
	for {
		_, _ = fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		buffer.WriteString(line)
		input := buffer.String()
		buffer.Reset()

		l := lexer.New(input)
		p := parser.New(l)
		program, err := p.ParseProgram()
		if err != nil {
			fmt.Println(err)
			continue
		}

		evaluated := evaluator.Eval(ctx, state, program, mod.GetEnv())
		if evaluated != nil {
			if !evaluator.IsError(evaluated) {
				n := len(program.Statements)
				// 如果最后一条语句不是表达式语句，不要输出任何值
				if _, ok := program.Statements[n-1].(*ast.ExpressionStatement); !ok {
					continue
				}

				// 值为 null 不输出
				if evaluated == object.NULL {
					continue
				}

			}

			if _, err := io.WriteString(out, evaluated.String()); err != nil {
				fmt.Println(err)
			}
			if _, err := io.WriteString(out, "\n"); err != nil {
				fmt.Println(err)
			}
		}
	}
}
