package eval

import (
	"testing"

	"goto/lexer"
	"goto/object"
	"goto/parser"
)

func evalInput(inp string) object.Object {
	l := lexer.New(inp)
	p := parser.New(l)

	program := p.ParseProgram()

	return Eval(program)
}

func testIntegerObject(t *testing.T, obj object.Object, exp int64) bool {
	intobj, ok := obj.(*object.Integer)

	if !ok {
		t.Errorf("Expected object type to be integer. got=%T", obj)
		return false
	}

	if intobj.Value != exp {
		t.Errorf("Expected %d but got %d instead", exp, intobj.Value)
		return false
	}

	return true
}

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input string
		exp   int64
	}{
		{"5", 5},
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
		out := evalInput(tt.input)
		testIntegerObject(t, out, tt.exp)
	}
}

func testBooleanObject(t *testing.T, obj object.Object, exp bool) bool {
	boolobj, ok := obj.(*object.Boolean)

	if !ok {
		t.Errorf("Expected object type to be boolean. got=%T", obj)
		return false
	}

	if boolobj.Value != exp {
		t.Errorf("Expected %v but got %v instead", exp, boolobj.Value)
		return false
	}

	return true
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input string
		exp   bool
	}{
		{"true", true},
		{"false", false},
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
	}

	for _, tt := range tests {
		out := evalInput(tt.input)
		testBooleanObject(t, out, tt.exp)
	}
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("Expected object to be NULL. got=%T", obj)
		return false
	}
	return true
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input string
		exp   interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
		{"if (1 > 2) { 10 } else if ( 3 > 4 ) { 20 } else { 30 }", 30},
		{"if (1 > 2) { 10 } else if ( 3 < 4 ) { 20 } else { 30 }", 20},
	}

	for _, tt := range tests {
		out := evalInput(tt.input)
		intg, ok := tt.exp.(int)

		if ok {
			testIntegerObject(t, out, int64(intg))
		} else {
			testNullObject(t, out)
		}
	}
}
