package main

import (
	"io"
	"io/ioutil"
	"os"

	"github.com/tatsuya4559/monkey/evaluator"
	"github.com/tatsuya4559/monkey/lexer"
	"github.com/tatsuya4559/monkey/object"
	"github.com/tatsuya4559/monkey/parser"
)

func main() {
	if len(os.Args) < 2 {
		panic("no src file")
	}

	src, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	l := lexer.New(string(src))
	p := parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		printParseErrors(os.Stderr, p.Errors())
		os.Exit(1)
	}

	env := object.NewEnvironment()
	macroEnv := object.NewEnvironment()

	evaluator.DefineMacros(program, macroEnv)
	expanded := evaluator.ExpandMacros(program, macroEnv)

	evaluated := evaluator.Eval(expanded, env)
	if evaluated != nil {
		io.WriteString(os.Stdout, evaluated.Inspect())
		io.WriteString(os.Stdout, "\n")
	}
}

func printParseErrors(out io.Writer, errors []string) {
	io.WriteString(out, "Woops! We ran into some monkey business here!\n")
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
