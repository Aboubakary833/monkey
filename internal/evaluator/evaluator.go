package evaluator

import (
	"monkey/internal/ast"
	"monkey/internal/object"
)


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
