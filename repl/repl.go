package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/tatsuya4559/monkey/evaluator"
	"github.com/tatsuya4559/monkey/lexer"
	"github.com/tatsuya4559/monkey/object"
	"github.com/tatsuya4559/monkey/parser"
)

const PROMPT = ">> "
const MONKEY_FACE = `
        ／三ヽ
       /(( ‥|)
     ／  ( ┴)
  ～(  || ノ||
    (_(∪)_)∪
`

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()
	macroEnv := object.NewEnvironment()

	for {
		fmt.Print(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program, err := p.ParseProgram()
		if err != nil {
			printParseError(out, err)
			continue
		}

		evaluator.DefineMacros(program, macroEnv)
		expanded := evaluator.ExpandMacros(program, macroEnv)

		evaluated := evaluator.Eval(expanded, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParseError(out io.Writer, err error) {
	io.WriteString(out, MONKEY_FACE)
	io.WriteString(out, "Woops! We ran into some monkey business here!\n")
	io.WriteString(out, " parser error:\n")
	io.WriteString(out, "\t"+err.Error()+"\n")
}
