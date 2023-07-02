package repl

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
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
		for _, statement := range program.Statements {
			fmt.Println(statement.String())
		}
	}
}
