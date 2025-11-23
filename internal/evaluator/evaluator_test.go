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
		{ "23.1", 23.1 },
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		if !testFloatObject(t, evaluated, tt.expected) {
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
			obj,
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
			obj,
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


func testEval(input string) object.Object {
	lex := lexer.New(input)
	parser := parser.New(lex)
	program := parser.ParseProgram()

	return Eval(program)
}

