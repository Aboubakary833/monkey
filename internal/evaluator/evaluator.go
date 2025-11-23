package evaluator

import (
	"monkey/internal/ast"
	"monkey/internal/object"
)


var (
	NULL	= &object.Null{}
	TRUE 	= &object.Boolean{ Value: true }
	FALSE	= &object.Boolean{ Value: false }
)


var boolAndNullBangOpMap = map[string]*object.Boolean{
	"null": TRUE,
	"false": TRUE,
	"true": FALSE,
}

var boolAndNullMinusMap = map[string]*object.Integer{
	"true": { Value: -1 },
	"false": { Value: -0 },
	"null": { Value: -0 },
}


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

	if boolObj, ok := boolAndNullBangOpMap[right.Inspect()]; ok {
		return boolObj
	}

	if integer, ok := right.(*object.Integer); ok && integer.Value == 0 {
		return TRUE
	}

	return FALSE
}

func evalMinusOperatorExpression(right object.Object) object.Object {

	if result, ok := boolAndNullMinusMap[right.Inspect()]; ok {
		return result
	}
	
	switch right.Type() {

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
