package evaluator

import (
	"testing"

	"github.com/Jitesh117/brainrotLang-interpreter/lexer"
	"github.com/Jitesh117/brainrotLang-interpreter/object"
	"github.com/Jitesh117/brainrotLang-interpreter/parser"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"4", 4},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
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
	env := object.NewEnvironment()
	program := p.ParseProgram()

	return Eval(program, env)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
		return false
	}
	return true
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"based", true},
		{"cap", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"based == based", true},
		{"cap == cap", true},
		{"based == cap", false},
		{"based != cap", true},
		{"cap != based", true},
		{"(1 < 2) == based", true},
		{"(1 < 2) == cap", false},
		{"(1 > 2) == based", false},
		{"(1 > 2) == cap", true},
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
		t.Errorf("object has wrong value. got=%t, want=%t", result.Value, expected)
		return false
	}
	return true
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!based", false},
		{"!cap", true},
		{"!5", false},
		{"!!based", true},
		{"!!cap", false},
		{"!!5", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestSusFrExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"fr (based) { 10 }", 10},
		{"fr (cap) { 10 }", nil},
		{"fr (1) { 10 }", 10},
		{"fr (1 < 2) { 10 }", 10},
		{"fr (1 > 2) { 10 }", nil},
		{"fr (1 > 2) { 10 } sus { 20 }", 20},
		{"fr (1 < 2) { 10 } sus { 20 }", 10},
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

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}

func TestSlayStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"slay 10;", 10},
		{"slay 10; 9", 10},
		{"slay 2 * 5; 9;", 10},
		{"9; slay 2 * 5; 9;", 10},
		{`
			fr (10 > 1) {
				fr(10 > 1){
					slay 10;	
				}
				slay 1
			}
			slay 1
			`, 10},
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
			"5 + based;",
			"L + ratio + type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + based; 5;",
			"L + ratio + type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-based",
			"we don't do that here. unknown operator: -BOOLEAN",
		},
		{
			"based + cap;",
			"we don't do that here. unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; based + cap; 5",
			"we don't do that here. unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"fr (10 > 1) { based + based; }",
			"we don't do that here. unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`
				fr (10 > 1) {
				slay based + cap;
				}
				slay 1;
				}
			`, "we don't do that here. unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"foobar",
			"bruh moment! identifier not found: foobar",
		},
		{
			`"Hello" - "World"`,
			"we don't do that here. unknown operator: STRING - STRING",
		},
		{
			`{"name": "Monkey"}[vibe(x) { x }];`,
			"nah fam FUNCTION cannot be used as a hash key",
		},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)

		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned. got=%T(%+v)", evaluated, evaluated)
			continue
		}
		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q", tt.expectedMessage, errObj.Message)
		}

	}
}

func TestYeetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"yeet a = 5; a;", 5},
		{"yeet a = 5 * 5; a;", 25},
		{"yeet a = 5; yeet b = a; b;", 5},
		{"yeet a = 5; yeet b = a; yeet c = a + b + 5; c", 15},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestVibeObject(t *testing.T) {
	input := "vibe(x) { x + 2; };"
	evaluated := testEval(input)

	vb, ok := evaluated.(*object.Vibe)
	if !ok {
		t.Fatalf("object is not Vibe. got=%T (%+v)", evaluated, evaluated)
	}
	if len(vb.Parameters) != 1 {
		t.Fatalf("vibe has wrong parameters. Parameters=%+v", vb.Parameters)
	}

	if vb.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", vb.Parameters[0])
	}

	expectedBody := "(x + 2)"

	if vb.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, vb.Body.String())
	}
}

func TestVibeAppplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"yeet identity = vibe(x) { x; }; identity(5);", 5},
		{"yeet identity = vibe(x) { slay x; }; identity(5);", 5},
		{"yeet double = vibe(x) { x * 2; }; double(5);", 10},
		{"yeet add = vibe(x, y) { x + y; }; add(5, 5);", 10},
		{"yeet add = vibe(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"vibe(x) { x; }(5)", 5},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestClosures(t *testing.T) {
	input := `
	yeet newAdder = vibe(x) {
	vibe(y) { x + y };
	};

	yeet addTwo = newAdder(2);
	addTwo(2);
	`
	testIntegerObject(t, testEval(input), 4)
}

func TestStringLiteral(t *testing.T) {
	input := `"Hello World!"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestStringConcatenation(t *testing.T) {
	input := `"Hello" + " " + "World!"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestBuiltinVibes(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`rizzLevel("")`, 0},
		{`rizzLevel("four")`, 4},
		{`rizzLevel("hello world")`, 11},
		{`rizzLevel(1)`, "argument to `rizzLevel` not supported, got INTEGER"},
		{`rizzLevel("one", "two")`, "wrong number of arguments. got=2, want=1"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not Error. got=%T (%+v)", evaluated, evaluated)
				continue
			}
			if errObj.Message != expected {
				t.Errorf("wrong error message. expected=%q, got=%q", expected, errObj.Message)
			}
		}
	}
}

func TestArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"
	evaluated := testEval(input)
	result, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
	}
	if len(result.Elements) != 3 {
		t.Fatalf("array has wrong num of elements. got=%d",
			len(result.Elements))
	}
	testIntegerObject(t, result.Elements[0], 1)
	testIntegerObject(t, result.Elements[1], 4)
	testIntegerObject(t, result.Elements[2], 6)
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			"[1, 2, 3][0]",
			1,
		},
		{
			"[1, 2, 3][1]",
			2,
		},
		{
			"[1, 2, 3][2]",
			3,
		},
		{
			"yeet i = 0; [1][i];",
			1,
		},
		{
			"[1, 2, 3][1 + 1];",
			3,
		},
		{
			"yeet myArray = [1, 2, 3]; myArray[2];",
			3,
		},
		{
			"yeet myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];",
			6,
		},
		{
			"yeet myArray = [1, 2, 3]; yeet i = myArray[0]; myArray[i]",
			2,
		},
		{
			"[1, 2, 3][3]",
			nil,
		},
		{
			"[1, 2, 3][-1]",
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
	input := `yeet two = "two";
{
"one": 10 - 9,
two: 1 + 1,
"thr" + "ee": 6 / 2,
4: 4,
based: 5,
cap: 6
}`
	evaluated := testEval(input)
	result, ok := evaluated.(*object.Hash)
	if !ok {
		t.Fatalf("Eval didn't return Hash. got=%T (%+v)", evaluated, evaluated)
	}
	expected := map[object.HashKey]int64{
		(&object.String{Value: "one"}).HashKey():   1,
		(&object.String{Value: "two"}).HashKey():   2,
		(&object.String{Value: "three"}).HashKey(): 3,
		(&object.Integer{Value: 4}).HashKey():      4,
		BASED.HashKey():                            5,
		CAP.HashKey():                              6,
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

func TestHashIndexExpressions(t *testing.T) {
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
			`yeet key = "foo"; {"foo": 5}[key]`,
			5,
		},
		{
			`{}["foo"]`,
			nil,
		},

		{
			`{based: 5}[based]`,
			5,
		},
		{
			`{cap: 5}[cap]`,
			5,
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
