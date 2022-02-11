package evaluator

import (
	"testing"

	"github.com/tatsuya4559/monkey/ast"
	"github.com/tatsuya4559/monkey/lexer"
	"github.com/tatsuya4559/monkey/object"
	"github.com/tatsuya4559/monkey/parser"
)

func TestDefineMacros(t *testing.T) {
	input := `
let number = 1;
let function = fn(x, y) { x + y };
let mymacro = macro(x, y) { x + y };
`
	env := object.NewEnvironment()
	program := testParseProgram(t, input)

	DefineMacros(program, env)

	if len(program.Statements) != 2 {
		t.Fatalf("program.Statements does not contain 2 statements. got=%d",
			len(program.Statements))
	}

	if _, ok := env.Get("number"); ok {
		t.Fatalf("number should not be defined")
	}
	if _, ok := env.Get("function"); ok {
		t.Fatalf("function should not be defined")
	}

	obj, ok := env.Get("mymacro")
	if !ok {
		t.Fatalf("macro not in environment")
	}
	macro, ok := obj.(*object.Macro)
	if !ok {
		t.Fatalf("object is not Macro. got=%T (%+v)", obj, obj)
	}

	if len(macro.Parameters) != 2 {
		t.Fatalf("wrong number of macro.Parameters. want=%d, got=%d",
			2, len(macro.Parameters))
	}
	if macro.Parameters[0].String() != "x" {
		t.Errorf("parameter[0] is not 'x'. got=%q", macro.Parameters[0])
	}
	if macro.Parameters[1].String() != "y" {
		t.Errorf("parameter[1] is not 'y'. got=%q", macro.Parameters[1])
	}

	expectedBody := `(x + y)`
	if macro.Body.String() != expectedBody {
		t.Errorf("body is not %q. got=%q", expectedBody, macro.Body.String())
	}
}

func testParseProgram(t *testing.T, input string) *ast.Program {
	t.Helper()
	l := lexer.New(input)
	p := parser.New(l)
	program, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	return program
}

func TestExpandMacros(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			`let infixExpr = macro() { quote(1 + 2); };
			infixExpr();
			`,
			`(1 + 2)`,
		},
		{
			// use unquote to access arguments
			// without unquote, returns just `b - a`
			`let reverse = macro(a, b) { quote(unquote(b) - unquote(a)); };
			reverse(2 + 2, 10 - 5);
			`,
			`(10 - 5) - (2 + 2)`,
		},
		{
			`
			let unless = macro(cond, consequence, alternative) {
				quote(
					if (!(unquote(cond))) {
						unquote(consequence)
					} else {
						unquote(alternative)
					}
				);
			};
			unless(10 > 5, puts("not greater"), puts("greater"));
			`,
			`if (!(10 > 5)) { puts("not greater") } else { puts("greater") }`,
		},
	}

	for _, tt := range tests {
		expected := testParseProgram(t, tt.expected)
		program := testParseProgram(t, tt.input)

		env := object.NewEnvironment()
		DefineMacros(program, env)
		expanded := ExpandMacros(program, env)

		if expanded.String() != expected.String() {
			t.Errorf("not equal. want=%q, got=%q",
				expected.String(), expanded.String())
		}
	}
}
