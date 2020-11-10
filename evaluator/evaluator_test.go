package evaluator

import (
	"github.com/tatsuya4559/monkey/lexer"
	"github.com/tatsuya4559/monkey/object"
	"github.com/tatsuya4559/monkey/parser"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10;", 10},
		{"-5", -5},
		{"-10;", -10},
		{"5 + 5 - 8", 2},
		{"2 * 4 * 1", 8},
		{"-50 + 10 - 1", -41},
		{"3 + 5 * 2", 13},
		{"10 + 5 * -2", 0},
		{"50 / 2 * 3 + 10", 85},
		{"2 * (5 + 10)", 30},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	return Eval(program, env)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. want=%d, got=%d",
			expected, result.Value)
		return false
	}

	return true
}

func TestEvalStringExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"Hello world"`, "Hello world"},
		{`"foo";`, "foo"},
		{`"foo" + " " + "bar";`, "foo bar"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testStringObject(t, evaluated, tt.expected)
	}
}

func testStringObject(t *testing.T, obj object.Object, expected string) bool {
	str, ok := obj.(*object.String)
	if !ok {
		t.Errorf("object is not String. got=%T (%+v)", obj, obj)
		return false
	}

	if str.Value != expected {
		t.Errorf("String has wrong value. want=%q, got=%q",
			expected, str.Value)
		return false
	}

	return true
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false;", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
		{`"foo" == "foo"`, true},
		{`"foobar" == "foo" + "bar"`, true},
		{`"foo" != "bar"`, true},
		{`[1, 2, 3] == [1, 2, 3]`, true},
		{`[1, 1 * 2, 3] == [1, 2]`, false},
		{`[1, 2, 1 + 2] == [1, 2, 4]`, false},
		{`["foo", "bar"] == ["foo", "bar"]`, true},
		{`["foo", "bar"] == ["foo"]`, false},
		{`["foo", "bar"] == ["foo", "baz"]`, false},
		{`["foo", "bar"] == []`, false},
		{`[true, false] == [true, false]`, true},
		{`[true, false] == [true]`, false},
		{`[true, false] == [true, true]`, false},
		{`[] == []`, true},
		{`[[1, 2], ["foo", "bar"]] == [[1, 2], ["foo", "bar"]]`, true},
		{`{"foo": 1, "bar": 2} == {"foo": 1, "bar": 2}`, true},
		{`{} == {}`, true},
		{`{"foo": 1, "bar": 2} == {"foo": 0, "bar": 2}`, false},
		{`{"foo": 1, "bar": 2} == {"baz": 1}`, false},
		{`{"foo": 1, "bar": 2} == {}`, false},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. want=%t, got=%t",
			expected, result.Value)
		return false
	}

	return true
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestIfElseExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 < 2) { 10 } else { 20 }", 10},
		{"if (1 > 2) { 10 } else { 20 }", 20},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		default:
			testNullObject(t, evaluated)
		}
	}
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}

func TestReturnStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{`
if (10 > 1) {
	if (10 > 1) {
		return 10;
	}

	return 1;
}
`,
			10,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"5 + true;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + true; 5",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-true",
			"unknown operator: -BOOLEAN",
		},
		{
			"true + false",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; true + false; 5",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if (10 > 1) { true + false }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`
if (10 > 1) {
	if (10 > 1) {
		return true + false;
	}
	return 1;
}
`,
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"foobar;",
			"identifier not found: foobar",
		},
		{
			`"hello" - "world"`,
			"unknown operator: STRING - STRING",
		},
		{
			`{[1,2]: "Monkey"};`,
			"unusable as hash key: ARRAY",
		},
		{
			`{"name": "Monkey"}[fn(x){x}];`,
			"unusable as hash key: FUNCTION",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned. got=%T (%+v)",
				evaluated, evaluated)
			continue
		}

		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. want %q, got=%q",
				tt.expectedMessage, errObj.Message)
		}
	}
}

func TestLetStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a =  5; let b = a; b;", 5},
		{"let a =  5; let b = a; let c = a + b + 5; c;", 15},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	input := `fn(x) { x + 2; };`
	expectedBody := "(x + 2)"

	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. got=%+v", fn.Parameters)
	}
	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fn.Parameters[0])
	}
	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let identity = fn(x) { x; }; identity(5);", 5},
		{"let identity = fn(x) { return x; }; identity(5);", 5},
		{"let double = fn(x) { x * 2; }; double(5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5, 5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"fn(x) { x; }(5)", 5},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestClosures(t *testing.T) {
	input := `
let newAdder = fn(x) {
	fn(y) { x + y };
};

let addTwo = newAdder(2);
addTwo(5);`
	testIntegerObject(t, testEval(input), 7)
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len(1)`, "argument to `len` not supported, got INTEGER"},
		{`len("one", "two")`, "wrong number of arguments. want=1, got=2"},
		{`len([])`, 0},
		{`len(["foo", "bar"])`, 2},
		{`len([1, 2 * 3, 4 + 5 + 6, "foo"])`, 4},
		{`first([])`, nil},
		{`first([1, 2])`, 1},
		{`first([2 * 3, 4 + 5 + 6])`, 6},
		{`first(1)`, "argument to `first` must be ARRAY, got INTEGER"},
		{`first("one", "two")`, "wrong number of arguments. want=1, got=2"},
		{`last([])`, nil},
		{`last([1, 2])`, 2},
		{`last([2 * 3, 4 + 5 + 6, 7 + 8])`, 15},
		{`last(1)`, "argument to `last` must be ARRAY, got INTEGER"},
		{`last("one", "two")`, "wrong number of arguments. want=1, got=2"},
		{`rest([])`, nil},
		{`rest([1])`, []int{}},
		{`rest([1, 2])`, []int{2}},
		{`rest([2 * 3, 4 + 5 + 6, 7 + 8])`, []int{15, 15}},
		{`rest(1)`, "argument to `rest` must be ARRAY, got INTEGER"},
		{`rest("one", "two")`, "wrong number of arguments. want=1, got=2"},
		{`push([], 1)`, []int{1}},
		{`push([1], 2)`, []int{1, 2}},
		{`push([1, 2], 5 - 2)`, []int{1, 2, 3}},
		{`push(1)`, "wrong number of arguments. want=2, got=1"},
		{`push(1, 2)`, "first argument to `push` must be ARRAY, got INTEGER"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Fatalf("object is not Error. got=%T (%+v)",
					evaluated, evaluated)
			}
			if errObj.Message != expected {
				t.Errorf("wrong error message. want=%q, got=%q",
					expected, errObj.Message)
			}
		case []int:
			testArrayObject(t, evaluated, expected)
		default:
			testNullObject(t, evaluated)
		}
	}
}

func testArrayObject(t *testing.T, obj object.Object, expected []int) bool {
	arr, ok := obj.(*object.Array)
	if !ok {
		t.Errorf("object is not ARRAY. got=%T (%+v)", obj, obj)
		return false
	}
	if len(arr.Elements) != len(expected) {
		t.Errorf("wrong number of elements. want=%d, got=%d",
			len(expected), len(arr.Elements))
		return false
	}
	for i, e := range arr.Elements {
		if !testIntegerObject(t, e, int64(expected[i])) {
			return false
		}
	}
	return true
}

func TestArrayLiteral(t *testing.T) {
	input := `[1, 2 * 2, 3 + 3];`

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("object is not array. got=%T (%+v)", evaluated, evaluated)
	}

	if len(result.Elements) != 3 {
		t.Fatalf("array has wrong num of elements. got=%d",
			len(result.Elements))
	}

	testIntegerObject(t, result.Elements[0], 1)
	testIntegerObject(t, result.Elements[1], 4)
	testIntegerObject(t, result.Elements[2], 6)
}

func TestArrayIndexExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`[1, 2, 3][1]`,
			2,
		},
		{
			`[1, 2, 3][2]`,
			3,
		},
		{
			`let i = 0; [1][i]`,
			1,
		},
		{
			`[1, 2, 3][1 + 1]`,
			3,
		},
		{
			`let myArray = [1, 2, 3]; myArray[2]`,
			3,
		},
		{
			`let myArray = [1, 2, 3]; myArray[0] + myArray[1];`,
			3,
		},
		{
			`let myArray = [1, 2, 3]; let i = myArray[0]; myArray[i];`,
			2,
		},
		{
			`[1, 2, 3][3]`,
			nil,
		},
		{
			`[1, 2, 3][-1]`,
			nil,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestHashLiterals(t *testing.T) {
	input := `let two = "two";
{
	"one": 10 - 9,
	two: 1 + 1,
	"thr" + "ee": 6 / 2,
	4: 4,
	true: 5,
	false: 6
}`
	expected := map[object.HashKey]int64{
		(&object.String{Value: "one"}).HashKey():   1,
		(&object.String{Value: "two"}).HashKey():   2,
		(&object.String{Value: "three"}).HashKey(): 3,
		(&object.Integer{Value: 4}).HashKey():      4,
		TRUE.HashKey():                             5,
		FALSE.HashKey():                            6,
	}

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Hash)
	if !ok {
		t.Fatalf("object is not hash. got=%T (%+v)", evaluated, evaluated)
	}

	if len(result.Pairs) != len(expected) {
		t.Fatalf("Hash has wrong num of pairs. got=%d", len(result.Pairs))
	}

	for expectedKey, expectedValue := range expected {
		pair, ok := result.Pairs[expectedKey]
		if !ok {
			t.Errorf("no pair for given key in Pairs")
		}
		testIntegerObject(t, pair.Value, expectedValue)
	}
}

func TestHashIndexExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`{"foo": 5}["foo"]`,
			5,
		},
		{
			`{"foo": 5}["bar"]`,
			nil,
		},
		{
			`let key = "foo"; {"foo": 5}[key]`,
			5,
		},
		{
			`{}["foo"]`,
			nil,
		},
		{
			`{5: 5}[5]`,
			5,
		},
		{
			`{true: 5}[true]`,
			5,
		},
		{
			`{false: 5}[false]`,
			5,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		default:
			testNullObject(t, evaluated)
		}

	}
}
