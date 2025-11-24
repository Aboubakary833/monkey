package evaluator

import (
	"fmt"
	"math"
	"monkey/internal/ast"
	"monkey/internal/object"
	"slices"
	"strconv"
)


var (
	NULL	= &object.Null{}
	TRUE 	= &object.Boolean{ Value: true }
	FALSE	= &object.Boolean{ Value: false }
)

var booleanOperators = []string{ "==", "!=", "<", ">", "<=", ">=" }


func Eval(node ast.Node) object.Object {

	switch node := node.(type) {
	
	case *ast.Program:
		return evalStatement(node.Statements)

	case *ast.ExpressionStatement:
		return Eval(node.Expression)

	case *ast.IntegerLiteral:
		return &object.Integer{ Value: node.Value }

	case *ast.FloatLiteral:
		return &object.Float{ Value: node.Value }

	case *ast.Boolean:
		return evalToNativeBool(node.Value)

	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evaluatePrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)
		return evaluateInfixExpression(node.Operator, left, right)
	}

	return nil
}


func evalStatement(statements []ast.Statement) object.Object {
	var result object.Object

	for _, stmt := range statements {
		result = Eval(stmt)
	}

	return result
}

func evalToNativeBool(val bool) *object.Boolean {
	if val {
		return TRUE
	}
	return FALSE
}

func evaluatePrefixExpression(operator string, right object.Object) object.Object {
	if operator == "!" {
		return evalBangOperatorExpression(right)
	}

	if operator == "-" {
		return evalMinusOperatorExpression(right)
	}

	return NULL
}

func evalBangOperatorExpression(right object.Object) object.Object {

	switch right.Type() {

	case object.BOOLEAN_OBJ:
		value := right.(*object.Boolean).Value
		return evalToNativeBool(!value)

	case object.INTEGER_OBJ, object.FLOAT_OBJ:
		value := getObjectNumberValue(right)
		return evalToNativeBool(value == 0.0)

	default:
		return FALSE
	}

}

func evalMinusOperatorExpression(right object.Object) object.Object {
	
	switch right.Type() {

	case object.BOOLEAN_OBJ:
		if boolean, ok := right.(*object.Boolean); ok && boolean.Value {
			return &object.Integer{ Value: -1 }
		}
		return &object.Integer{ Value: 0 }

	case object.INTEGER_OBJ:
		value := right.(*object.Integer).Value
		return &object.Integer{ Value: -value }

	case object.FLOAT_OBJ:
		value := right.(*object.Float).Value
		return &object.Float{ Value: -value }

	default:
		return NULL
	}
}


func evaluateInfixExpression(operator string, left, right object.Object) object.Object {

	// In case we're dealing with booleans

	if left.Type() == object.BOOLEAN_OBJ && right.Type() == object.BOOLEAN_OBJ {
		switch operator {
		case "!=":
			return evalToNativeBool(left != right)

		case "<", ">", "<=", ">=":
			leftVal := getObjectNumberValue(left)
			rightVal := getObjectNumberValue(right)

			return evaluateArithmeticOperatorExpression(operator, leftVal, rightVal)

		default:
			return evalToNativeBool(left == right)
		}
	}

	// Otherwise

	leftValue := getObjectNumberValue(left)
	rightValue := getObjectNumberValue(right)

	if slices.Contains(booleanOperators, operator) {
		return evaluateLogicalOperatorExpression(operator, leftValue, rightValue)
	}

	return evaluateArithmeticOperatorExpression(operator, leftValue, rightValue)
}


func evaluateArithmeticOperatorExpression(operator string, leftValue, rightValue float64) object.Object {

	var result float64

	switch operator {
	
	case "+":
		result = leftValue + rightValue

	case "-":
		result = leftValue - rightValue

	case "*":
		result = leftValue * rightValue

	case "/":
		result = leftValue / rightValue

	case "%":
		result = math.Mod(leftValue, rightValue)
	
	default:
		return NULL

	}
	
	value, err := strconv.Atoi(fmt.Sprintf("%g", result))
	if err != nil {
		return &object.Float{ Value: result }
	}

	return &object.Integer{ Value: int64(value) }
}

func evaluateLogicalOperatorExpression(operator string, leftValue, rightValue float64) object.Object {
	switch operator {
	
	case "==":
		return evalToNativeBool(leftValue == rightValue)

	case "!=":
		return evalToNativeBool(leftValue != rightValue)

	case "<":
		return evalToNativeBool(leftValue < rightValue)

	case ">":
		return evalToNativeBool(leftValue > rightValue)

	case "<=":
		return evalToNativeBool(leftValue <= rightValue)

	case ">=":
		return evalToNativeBool(leftValue >= rightValue)

	default:
		return NULL
	}
}


// getObjectNumberValue return the object float64 representation value.
// For integer, for example, the value will get converted to float64
// and then return. Also, NULL & FALSE objects return 0 while TRUE return 1.
func getObjectNumberValue(obj object.Object) float64 {
	
	switch obj.Type() {

	case object.INTEGER_OBJ:
		return float64(obj.(*object.Integer).Value)

	case object.FLOAT_OBJ:
		return obj.(*object.Float).Value

	case object.BOOLEAN_OBJ:
		if boolean, ok := obj.(*object.Boolean); ok && boolean.Value {
			return 1
		}
		return 0
	}

	return 0
}
