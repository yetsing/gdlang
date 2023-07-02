package repl

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"weilang/evaluator"
	"weilang/lexer"
	"weilang/parser"
)

const (
	START_PROMPT = ">> "
	//CONTINUE_PROMPT = "..."
)

var PROMPT = START_PROMPT

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

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

		evaluated := evaluator.Eval(program)
		if evaluated != nil {
			if _, err := io.WriteString(out, evaluated.String()); err != nil {
				fmt.Println(err)
			}
			if _, err := io.WriteString(out, "\n"); err != nil {
				fmt.Println(err)
			}
		}
	}
}
