package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/user"

	"github.com/tatsuya4559/monkey/evaluator"
	"github.com/tatsuya4559/monkey/lexer"
	"github.com/tatsuya4559/monkey/object"
	"github.com/tatsuya4559/monkey/parser"
	"github.com/tatsuya4559/monkey/repl"
)

func main() {
	if len(os.Args) < 2 {
		startREPL()
	} else {
		interpretFile(os.Args[1])
	}
}

func startREPL() {
	user, err := user.Current()
	if err != nil {
		log.Fatalf("cannot get current user: %v", err)
	}
	fmt.Printf("Hello %s! This is the Monkey programming language!\n",
		user.Username)
	fmt.Printf("Feel free to type in commands\n")
	repl.Start(os.Stdin, os.Stdout)
}

func interpretFile(filename string) {
	src, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("cannot read %s: %v", filename, err)
	}

	l := lexer.New(string(src))
	p := parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		io.WriteString(os.Stderr, "Woops! We ran into some monkey business here!\n")
		io.WriteString(os.Stderr, " parser errors:\n")
		for _, msg := range p.Errors() {
			io.WriteString(os.Stderr, "\t"+msg+"\n")
		}
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
