package evaluator

import (
	"testing"

	"github.com/tatsuya4559/monkey/object"
)

func TestQuote(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			`quote(5);`,
			`5`,
		},
		{
			`quote(5 + 8);`,
			`(5 + 8)`,
		},
		{
			`quote(foobar);`,
			`foobar`,
		},
		{
			`quote(foobar + barfoo);`,
			`(foobar + barfoo)`,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		quote, ok := evaluated.(*object.Quote)
		if !ok {
			t.Fatalf("object is not *object.Quote. got=%T (%+v)",
				evaluated, evaluated)
		}

		if quote.Node == nil {
			t.Fatalf("quote.Node is nil")
		}

		if quote.Node.String() != tt.expected {
			t.Errorf("not equal. want=%q, got=%q",
				tt.expected, quote.Node.String())
		}
	}
}

func TestQuoteUnquote(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			`quote(unquote(4));`,
			`4`,
		},
		{
			`quote(unquote(4 + 3));`,
			`7`,
		},
		{
			`quote(5 + unquote(4 + 3));`,
			`(5 + 7)`,
		},
		{
			`quote(unquote(4 + 3) + 5);`,
			`(7 + 5)`,
		},
		{
			`let foo = 8;
			quote(foo)`,
			`foo`,
		},
		{
			`let foo = 8;
			quote(unquote(foo))`,
			`8`,
		},
		{
			`quote(unquote(true));`,
			`true`,
		},
		{
			`quote(unquote(true == false));`,
			`false`,
		},
		{
			`quote(unquote(quote(4 + 4)));`,
			`(4 + 4)`,
		},
		{
			`let quotedInfixExpression = quote(4 + 4);
			quote(unquote(4 + 4) + unquote(quotedInfixExpression));`,
			`(8 + (4 + 4))`,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		quote, ok := evaluated.(*object.Quote)
		if !ok {
			t.Fatalf("object is not *object.Quote. got=%T (%+v)",
				evaluated, evaluated)
		}

		if quote.Node == nil {
			t.Fatalf("quote.Node is nil")
		}

		if quote.Node.String() != tt.expected {
			t.Errorf("not equal. want=%q, got=%q",
				tt.expected, quote.Node.String())
		}
	}
}
