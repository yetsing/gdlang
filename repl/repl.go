package repl

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"weilang/lexer"
	"weilang/token"
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
		for {
			tok := l.NextToken()
			if tok.TypeIs(token.EOF) {
				break
			}
			if tok.TypeIs(token.ILLEGAL) {
				fmt.Println("illegal token", tok)
				break
			}
			fmt.Println(tok)
		}
	}
}
