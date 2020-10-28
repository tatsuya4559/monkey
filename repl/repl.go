package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/tatsuya4559/monkey/lexer"
	"github.com/tatsuya4559/monkey/parser"
)

const PROMPT = ">> "
const MONKEY_FACE = `
　　　　／三ヽ
　　　 /(( ‥|)
　　 ／　( ┴)
　～(　|| ノ||
　　(_(∪)_)∪
`

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParseErrors(out, p.Errors())
			continue
		}

		io.WriteString(out, program.String())
		io.WriteString(out, "\n")
	}
}

func printParseErrors(out io.Writer, errors []string) {
	io.WriteString(out, MONKEY_FACE)
	io.WriteString(out, "Woops! We ran into some monkey business here!\n")
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
