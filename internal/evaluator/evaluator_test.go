package evaluator

import (
	"monkey/internal/lexer"
	"monkey/internal/object"
	"monkey/internal/parser"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct{
		input		string
		expected	int64
	}{
		{ "5", 5 },
		{ "10", 10 },
		{ "187", 187 },
		{ "-20", -20 },
		{ "-3", -3 },
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{ "12 % 10", 2 },
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		if !testIntegerObject(t, evaluated, tt.expected) {
			return
		}
	}
}

func TestEvalFloatExpression(t *testing.T) {
		tests := []struct{
		input		string
		expected	float64
	}{
		{ "3.14", 3.14 },
		{ ".25", 0.25 },
		{ "-23.1", -23.1 },
		{ "-18.5", -18.5 },
		{"5.25 + 5.25 + 5.25 + 5.50 - 10", 11.25},
		{"2.5 * 2.5 * 2.5 * 2", 31.25},
		{"-50 + 100 + -50.50", -0.5},
		{"5 * 2.5 + 10", 22.5},
		{"5 + 2 * 10.25", 25.5},
		{"50 / 2 * 2 + 10.5", 60.5},
		{ "10.5 % 10", 0.5 },
		{"2.5 * (5.5 + 10)", 38.75},
		{"3 * 3 * 3 + 10.75", 37.75},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		if !testFloatObject(t, evaluated, tt.expected) {
			return
		}
	}
}


func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct{
		input		string
		expected	bool
	}{
		{ "true", true },
		{ "false", false },
		{ "1 < 2", true },
		{ "1 > 2", false },
		{ "1 < 1", false },
		{ "1 > 1", false },
		{ "1 == 1", true },
		{ "1 != 1", false },
		{ "1 == 2", false },
		{ "1 != 2", true },
		{ "1 >= 0", true },
		{ "2 <= 1", false },
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}

		for _, tt := range tests {
		evaluated := testEval(tt.input)
		if !testBooleanObject(t, evaluated, tt.expected) {
			return
		}
	}
}


func TestEvalBangOperator(t *testing.T) {
	tests := []struct{
		input		string
		expected	bool
	}{
		{ "!true", false },
		{ "!false", true },
		{ "!5", false },
		{ "!0", true },
		{ "!!true", true },
		{ "!!false", false },
		{ "!!5", true },
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		if !testBooleanObject(t, evaluated, tt.expected) {
			return
		}
	}
}



// Helpers functions:


func testIntegerObject(t *testing.T, got object.Object, expected int64) bool {
	obj, ok := got.(*object.Integer)

	if !ok {
		t.Errorf(
			"Expecting obj to be of type object.Integer, but got %T\n",
			got,
		)

		return false
	}

	if expected != obj.Value {
		t.Errorf(
			"Expecting obj.Value to be %d, but got %d\n",
			expected, obj.Value,
		)

		return false
	}

	return true
}


func testFloatObject(t *testing.T, got object.Object, expected float64) bool {
	obj, ok := got.(*object.Float)

	if !ok {
		t.Errorf(
			"Expecting obj to be of type object.Float, but got %T\n",
			got,
		)

		return false
	}

	if expected != obj.Value {
		t.Errorf(
			"Expecting obj.Value to be %g, but got %g\n",
			expected, obj.Value,
		)

		return false
	}

	return true
}


func testBooleanObject(t *testing.T, got object.Object, expected bool) bool {
	obj, ok := got.(*object.Boolean)

	if !ok {
		t.Errorf(
			"Expecting obj to be of type object.Boolean, but got %T\n",
			got,
		)

		return false
	}

	if expected != obj.Value {
		t.Errorf(
			"Expecting obj.Value to be %t, but got %t\n",
			expected, obj.Value,
		)

		return false
	}

	return true
}


func testEval(input string) object.Object {
	lex := lexer.New(input)
	parser := parser.New(lex)
	program := parser.ParseProgram()

	return Eval(program)
}

