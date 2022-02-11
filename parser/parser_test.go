package parser

import (
	"fmt"
	"testing"

	"github.com/tatsuya4559/monkey/ast"
	"github.com/tatsuya4559/monkey/lexer"
)

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{`let x = 5; // comment`, "x", 5},
		{`let y = true;`, "y", true},
		{`let foobar = y;`, "foobar", "y"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program, err := p.ParseProgram()
		if err != nil {
			t.Fatalf("parse error: %v", err)
		}

		if program == nil {
			t.Fatalf("ParseProgram() returned nil")
		}
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d",
				len(program.Statements))
		}

		stmt := program.Statements[0]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}

		val := stmt.(*ast.LetStatement).Value
		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got=%q", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not %s. got=%s", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral() not '%s'. got=%s",
			name, letStmt.Name.TokenLiteral())
		return false
	}

	return true
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{`return 5; // comment`, 5},
		{`return x;`, "x"},
		{`return false;`, false},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program, err := p.ParseProgram()
		if err != nil {
			t.Fatalf("parse error: %v", err)
		}

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d",
				len(program.Statements))
		}

		stmt := program.Statements[0]
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.ReturnStatement. got=%T", stmt)
			return
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral not 'return'. got=%q",
				returnStmt.TokenLiteral())
		}

		if !testLiteralExpression(t, returnStmt.ReturnValue, tt.expectedValue) {
			return
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := `foobar; // comment`

	l := lexer.New(input)
	p := New(l)
	program, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	if !testIdentifer(t, stmt.Expression, "foobar") {
		return
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := `5;`

	l := lexer.New(input)
	p := New(l)
	program, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	if !testIntegerLiteral(t, stmt.Expression, 5) {
		return
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue interface{}
	}{
		{"!5; //comment", "!", 5},
		{"-15;", "-", 15},
		{"!true;", "!", true},
		{"!false", "!", false},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program, err := p.ParseProgram()
		if err != nil {
			t.Fatalf("parse error: %v", err)
		}

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		expr, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("expr not *ast.PrefixExpression. got=%T", stmt.Expression)
		}
		if expr.Operator != tt.operator {
			t.Fatalf("expr.Operator not '%s'. got='%s'",
				tt.operator, expr.Operator)
		}
		if !testLiteralExpression(t, expr.Right, tt.integerValue) {
			return
		}
	}
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integer, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}
	if integer.Value != value {
		t.Errorf("integer.Value not %d. got=%d", value, integer.Value)
		return false
	}
	if integer.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integer.TokenLiteral not %d. got=%s", value,
			integer.TokenLiteral())
		return false
	}
	return true
}

func TestParsingInfixExpression(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5; // comment", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 % 5;", 5, "%", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program, err := p.ParseProgram()
		if err != nil {
			t.Fatalf("parse error: %v", err)
		}

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		testInfixExpression(t, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue)
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b // comment",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
		{
			"a * [1, 2, 3, 4][b * c] * d",
			"((a * ([1, 2, 3, 4][(b * c)])) * d)",
		},
		{
			"add(a * b[2], b[1], 2 * [1, 2][1])",
			"add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program, err := p.ParseProgram()
		if err != nil {
			t.Fatalf("parse error: %v", err)
		}

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func testIdentifer(t *testing.T, expr ast.Expression, value string) bool {
	ident, ok := expr.(*ast.Identifier)
	if !ok {
		t.Fatalf("expr not *ast.Identifier. got=%T", expr)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s. got=%s", value,
			ident.TokenLiteral())
		return false
	}

	return true
}

func testLiteralExpression(
	t *testing.T,
	expr ast.Expression,
	expected interface{},
) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, expr, int64(v))
	case int64:
		return testIntegerLiteral(t, expr, v)
	case string:
		return testIdentifer(t, expr, v)
	case bool:
		return testBooleanLiteral(t, expr, v)
	}
	t.Errorf("type of expr not handled. got=%T", expr)
	return false
}

func testInfixExpression(
	t *testing.T,
	expr ast.Expression,
	left interface{},
	operator string,
	right interface{},
) bool {
	opExpr, ok := expr.(*ast.InfixExpression)
	if !ok {
		t.Fatalf("expr not ast.InfixExpression. got=%T(%s)", expr, expr)
		return false
	}

	if !testLiteralExpression(t, opExpr.Left, left) {
		return false
	}

	if opExpr.Operator != operator {
		t.Errorf("expr.Operator not '%s'. got=%q", operator, opExpr.Operator)
		return false
	}

	if !testLiteralExpression(t, opExpr.Right, right) {
		return false
	}

	return true
}

func TestBooleanExpression(t *testing.T) {
	input := `true; //comment`

	l := lexer.New(input)
	p := New(l)
	program, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	if !testBooleanLiteral(t, stmt.Expression, true) {
		return
	}
}

func testBooleanLiteral(t *testing.T, expr ast.Expression, value bool) bool {
	b, ok := expr.(*ast.Boolean)
	if !ok {
		t.Errorf("expr not *ast.Boolean. got=%T", expr)
		return false
	}
	if b.Value != value {
		t.Errorf("boolean.Value not %t. got=%t", value, b.Value)
		return false
	}
	if b.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("boolean.TokenLiteral not %t. got=%s", value,
			b.TokenLiteral())
		return false
	}
	return true
}

func TestIfExpression(t *testing.T) {
	input := `if x < y { x } // comment`

	l := lexer.New(input)
	p := New(l)
	program, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	expr, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression not *ast.IfExpression. got=%T",
			stmt.Expression)
	}

	if !testInfixExpression(t, expr.Condition, "x", "<", "y") {
		return
	}

	if len(expr.Consequence.Statements) != 1 {
		t.Errorf("consequence not 1 statements. got=%d",
			len(expr.Consequence.Statements))
	}

	consequence, ok := expr.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] not ast.ExpressionStatement. got=%T",
			expr.Consequence.Statements[0])
	}
	if !testIdentifer(t, consequence.Expression, "x") {
		return
	}

	if expr.Alternative != nil {
		t.Errorf("expr.Alternative was not nil. got=%+v", expr.Alternative)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := `if x < y { x } else { y } //comment`

	l := lexer.New(input)
	p := New(l)
	program, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	expr, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression not *ast.IfExpression. got=%T",
			stmt.Expression)
	}

	if !testInfixExpression(t, expr.Condition, "x", "<", "y") {
		return
	}

	if len(expr.Consequence.Statements) != 1 {
		t.Errorf("consequence not 1 statements. got=%d",
			len(expr.Consequence.Statements))
	}

	consequence, ok := expr.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] not ast.ExpressionStatement. got=%T",
			expr.Consequence.Statements[0])
	}
	if !testIdentifer(t, consequence.Expression, "x") {
		return
	}

	if len(expr.Alternative.Statements) != 1 {
		t.Errorf("alternative not 1 statements. got=%d",
			len(expr.Alternative.Statements))
	}

	alternative, ok := expr.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] not ast.ExpressionStatement. got=%T",
			expr.Alternative.Statements[0])
	}
	if !testIdentifer(t, alternative.Expression, "y") {
		return
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `fn(x, y) { x + y; } //comment`

	l := lexer.New(input)
	p := New(l)
	program, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	// FunctionLiteral
	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}
	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression not *ast.FunctionLiteral. got=%T",
			stmt.Expression)
	}

	// Parameters
	if len(function.Parameters) != 2 {
		t.Errorf("function literal parameters wrong. want 2, got=%d",
			len(function.Parameters))
	}
	testLiteralExpression(t, function.Parameters[0], "x")
	testLiteralExpression(t, function.Parameters[1], "y")

	// Body
	if len(function.Body.Statements) != 1 {
		t.Fatalf("function.Body.Statements has not 1 statements. got=%d",
			len(function.Body.Statements))
	}
	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("function body stmt is not ast.ExpressionStatement. got=%T",
			function.Body.Statements[0])
	}
	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "fn() {}", expectedParams: []string{}},
		{input: "fn(x) {}", expectedParams: []string{"x"}},
		{input: "fn(x, y, z) {}", expectedParams: []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program, err := p.ParseProgram()
		if err != nil {
			t.Fatalf("parse error: %v", err)
		}

		// do not check type conversion. it is tested above.
		stmt := program.Statements[0].(*ast.ExpressionStatement)
		function := stmt.Expression.(*ast.FunctionLiteral)

		if len(function.Parameters) != len(tt.expectedParams) {
			t.Errorf("length parameters wrong. want %d, got=%d",
				len(tt.expectedParams), len(function.Parameters))
		}

		for i, ident := range function.Parameters {
			testLiteralExpression(t, ident, tt.expectedParams[i])
		}
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := `add(1, 2 * 3, 4 + 5); //comment`

	l := lexer.New(input)
	p := New(l)
	program, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	expr, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression not *ast.CallExpression. got=%T",
			stmt.Expression)
	}

	if !testIdentifer(t, expr.Function, "add") {
		return
	}

	if len(expr.Arguments) != 3 {
		t.Fatalf("wrong length of arguments. got=%d", len(expr.Arguments))
	}

	testLiteralExpression(t, expr.Arguments[0], 1)
	testInfixExpression(t, expr.Arguments[1], 2, "*", 3)
	testInfixExpression(t, expr.Arguments[2], 4, "+", 5)
}

func TestCallExpressionArgumentParsing(t *testing.T) {
	tests := []struct {
		input        string
		expectedArgs []string
	}{
		{input: `exit(); //comment`, expectedArgs: []string{}},
		{input: `print(x)`, expectedArgs: []string{"x"}},
		{input: `add(x, y, z)`, expectedArgs: []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program, err := p.ParseProgram()
		if err != nil {
			t.Fatalf("parse error: %v", err)
		}

		// do not check type conversion. it is tested above.
		stmt := program.Statements[0].(*ast.ExpressionStatement)
		expr := stmt.Expression.(*ast.CallExpression)

		if len(expr.Arguments) != len(tt.expectedArgs) {
			t.Errorf("wrong length of arguments. want %d, got=%d",
				len(tt.expectedArgs), len(expr.Arguments))
		}

		for i, ident := range expr.Arguments {
			testLiteralExpression(t, ident, tt.expectedArgs[i])
		}
	}
}

func TestStringLiteralExpression(t *testing.T) {
	input := `"hello world"; //comment`

	l := lexer.New(input)
	p := New(l)
	program, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	str, ok := stmt.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("expr not *ast.StringLiteral. got=%T", stmt.Expression)
	}
	if str.Value != "hello world" {
		t.Errorf("str.Value not %s. got=%s", "hello world", str.Value)
	}
	if str.TokenLiteral() != "hello world" {
		t.Errorf("str.TokenLiteral not %s. got=%s", "hello world",
			str.TokenLiteral())
	}
}

func TestParsingArrayLiteral(t *testing.T) {
	input := `[1, 2 * 2, 3 + 3] //comment`

	l := lexer.New(input)
	p := New(l)
	program, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	array, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("expr not *ast.ArrayLiteral. got=%T", stmt.Expression)
	}
	if len(array.Elements) != 3 {
		t.Fatalf("len(array.Elements) not 3. got=%d", len(array.Elements))
	}

	testIntegerLiteral(t, array.Elements[0], 1)
	testInfixExpression(t, array.Elements[1], 2, "*", 2)
	testInfixExpression(t, array.Elements[2], 3, "+", 3)
}

func TestParsingEmptyArray(t *testing.T) {
	input := `[]`

	l := lexer.New(input)
	p := New(l)
	program, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	array, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("expr not *ast.ArrayLiteral. got=%T", stmt.Expression)
	}
	if len(array.Elements) != 0 {
		t.Fatalf("len(array.Elements) not 0. got=%d", len(array.Elements))
	}
}

func TestParsingIndexExpression(t *testing.T) {
	input := `myArray[1 + 1];`

	l := lexer.New(input)
	p := New(l)
	program, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	indexExpr, ok := stmt.Expression.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("expr not *ast.IndexExpression. got=%T", stmt.Expression)
	}
	if !testIdentifer(t, indexExpr.Left, "myArray") {
		return
	}
	if !testInfixExpression(t, indexExpr.Index, 1, "+", 1) {
		return
	}
}

func TestParsingHashLiteralStringKey(t *testing.T) {
	input := `{"one": 1, "two": 2, "three": 3} //comment`
	expected := map[string]int64{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	l := lexer.New(input)
	p := New(l)
	program, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("expr not *ast.HashLiteral. got=%T", stmt.Expression)
	}

	if len(hash.Pairs) != 3 {
		t.Fatalf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
	}

	for key, value := range hash.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key %s is not *ast.StringLiteral. got=%T",
				key.String(), key)
		}

		expectedValue := expected[literal.Value]
		testIntegerLiteral(t, value, expectedValue)
	}
}

func TestParsingEmptyHashLiteral(t *testing.T) {
	input := `{}`

	l := lexer.New(input)
	p := New(l)
	program, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("expr not *ast.HashLiteral. got=%T", stmt.Expression)
	}

	if len(hash.Pairs) != 0 {
		t.Errorf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
	}
}

func TestParsingHashLiteralIntegerKey(t *testing.T) {
	input := `{100: 1, 200: 2, 300: 3}`
	expected := map[int64]int64{
		100: 1,
		200: 2,
		300: 3,
	}

	l := lexer.New(input)
	p := New(l)
	program, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("expr not *ast.HashLiteral. got=%T", stmt.Expression)
	}

	if len(hash.Pairs) != 3 {
		t.Fatalf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
	}

	for key, value := range hash.Pairs {
		literal, ok := key.(*ast.IntegerLiteral)
		if !ok {
			t.Errorf("key %s is not *ast.IntegerLiteral. got=%T",
				key.String(), key)
		}

		expectedValue := expected[literal.Value]
		testIntegerLiteral(t, value, expectedValue)
	}
}

func TestParsingHashLiteralBooleanKey(t *testing.T) {
	input := `{true: 1, false: 0}`
	expected := map[bool]int64{
		true:  1,
		false: 0,
	}

	l := lexer.New(input)
	p := New(l)
	program, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("expr not *ast.HashLiteral. got=%T", stmt.Expression)
	}

	if len(hash.Pairs) != 2 {
		t.Fatalf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
	}

	for key, value := range hash.Pairs {
		literal, ok := key.(*ast.Boolean)
		if !ok {
			t.Errorf("key %s is not *ast.Boolean. got=%T",
				key.String(), key)
		}

		expectedValue := expected[literal.Value]
		testIntegerLiteral(t, value, expectedValue)
	}
}

func TestParsingHashLiteralWithExpressions(t *testing.T) {
	input := `{"one": 0+1, "two": 10-8, "three": 15/5}`
	tests := map[string]func(ast.Expression){
		"one": func(e ast.Expression) {
			testInfixExpression(t, e, 0, "+", 1)
		},
		"two": func(e ast.Expression) {
			testInfixExpression(t, e, 10, "-", 8)
		},
		"three": func(e ast.Expression) {
			testInfixExpression(t, e, 15, "/", 5)
		},
	}

	l := lexer.New(input)
	p := New(l)
	program, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("expr not *ast.HashLiteral. got=%T", stmt.Expression)
	}

	if len(hash.Pairs) != 3 {
		t.Fatalf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
	}

	for key, value := range hash.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key %s is not *ast.StringLiteral. got=%T",
				key.String(), key)
		}

		testFunc, ok := tests[literal.Value]
		if !ok {
			t.Errorf("No test function for key %q found", literal.Value)
			continue
		}
		testFunc(value)
	}
}

func TestMacroLiteralParsing(t *testing.T) {
	input := `macro(x, y) { x + y; } //comment`

	l := lexer.New(input)
	p := New(l)
	program, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	// FunctionLiteral
	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}
	macro, ok := stmt.Expression.(*ast.MacroLiteral)
	if !ok {
		t.Fatalf("stmt.Expression not *ast.MacroLiteral. got=%T",
			stmt.Expression)
	}

	// Parameters
	if len(macro.Parameters) != 2 {
		t.Errorf("macro literal parameters wrong. want 2, got=%d",
			len(macro.Parameters))
	}
	testLiteralExpression(t, macro.Parameters[0], "x")
	testLiteralExpression(t, macro.Parameters[1], "y")

	// Body
	if len(macro.Body.Statements) != 1 {
		t.Fatalf("macro.Body.Statements has not 1 statements. got=%d",
			len(macro.Body.Statements))
	}
	bodyStmt, ok := macro.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("macro body stmt is not ast.ExpressionStatement. got=%T",
			macro.Body.Statements[0])
	}
	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestWhileStatement(t *testing.T) {
	input := `
while (i < 10) {
	let i = i * 2;
	// comment
}
	`

	l := lexer.New(input)
	p := New(l)
	program, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.WhileStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.WhileStatement. got=%T",
			program.Statements[0])
	}

	if !testInfixExpression(t, stmt.Condition, "i", "<", 10) {
		return
	}

	if len(stmt.Body.Statements) != 1 {
		t.Errorf("body not 1 statements. got=%d",
			len(stmt.Body.Statements))
	}

	if !testLetStatement(t, stmt.Body.Statements[0], "i") {
		return
	}

	val := stmt.Body.Statements[0].(*ast.LetStatement).Value
	if !testInfixExpression(t, val, "i", "*", 2) {
		return
	}
}
