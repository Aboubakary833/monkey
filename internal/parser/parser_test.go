package parser

import (
	"fmt"
	"monkey/internal/ast"
	"monkey/internal/lexer"
	"slices"
	"testing"
)

func TestDeclarationStatement(t *testing.T) {
	input := `const XYZ = 255;
	let t = 10;
	let foobar = 838383;`

	lex := lexer.New(input)
	parser := New(lex)

	program := parser.ParseProgram()
	checkParserErrors(t, parser)

	fmt.Println(len(program.Statements))

	tests := []struct {
		expectedIdentifier string
	}{
		{"XYZ"},
		{"t"},
		{"foobar"},
	}

	for i, tt := range tests {
		s := program.Statements[i]
		name := tt.expectedIdentifier

		if !slices.Contains([]string{"const", "let"}, s.TokenLiteral()) {
			t.Fatalf(
				"[test #%d] - Expected statement to start with %q or %q, but got %q\n",
				i, "const", "let", s.TokenLiteral(),
			)
		}

		dStmt, ok := s.(*ast.DeclarationStatement)

		if !ok {
			t.Fatalf(
				"[test #%d] - Expected s to be ast.DeclarationStatement, but got %T\n",
				i, dStmt,
			)
		}

		if dStmt.Name.Value != name {
			t.Fatalf(
				"[test %d] - Expected statement Name field Value to be %q, but got %q\n",
				i, name, dStmt.Name.Value,
			)
		}

		if dStmt.Name.TokenLiteral() != name {
			t.Fatalf(
				"[test %d] - Expected statement.Name.TokenLiteral() method to return %q, but got %q\n",
				i, name, dStmt.TokenLiteral(),
			)
		}

	}
}

func TestReturnStatement(t *testing.T) {
	input := `
	return 5;
	return 10;
	return 993322;
	`
	lex := lexer.New(input)
	parser := New(lex)

	program := parser.ParseProgram()
	checkParserErrors(t, parser)

	if len(program.Statements) != 3 {
		t.Fatalf(
			"Expected program.Statements to contains 3 statements, but got %d\n",
			len(program.Statements),
		)
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)

		if !ok {
			t.Errorf(
				"Expecting stmt to be a ast.ReturnStatement, but got %T\n",
				stmt,
			)
			continue
		}

		if returnStmt.TokenLiteral() != "return" {
			t.Errorf(
				"Expected returnStmt.TokenLiteral() to return \"return\", but got %q\n",
				returnStmt.TokenLiteral(),
			)
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar"
	lex := lexer.New(input)
	parser := New(lex)

	program := parser.ParseProgram()
	checkParserErrors(t, parser)

	if len(program.Statements) != 1 {
		t.Fatalf(
			"Expected program.Statements to contains 1 statements, but got %d\n",
			len(program.Statements),
		)
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf(
			"Expecting program.Statements[0] to be a ast.ExpressionStatement, but got %T\n",
			stmt,
		)
	}

	testIdentifier(t, stmt.Expression, "foobar")
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"
	lex := lexer.New(input)
	parser := New(lex)

	program := parser.ParseProgram()
	checkParserErrors(t, parser)

	if len(program.Statements) != 1 {
		t.Fatalf(
			"Expecting program.Statements to contains 1 statement, but got %d\n",
			len(program.Statements),
		)
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf(
			"Expecting program.Statement[0] to be of type *ast.ExpressionStatement, but got %T\n",
			stmt,
		)
	}

	testIntegerLiteral(t, stmt.Expression, 5)
}

func TestFloatLiteralExpression(t *testing.T) {
	input := "10.5;"
	lex := lexer.New(input)
	parser := New(lex)

	program := parser.ParseProgram()
	checkParserErrors(t, parser)

	if len(program.Statements) != 1 {
		t.Fatalf(
			"Expecting program.Statements to contains 1 Statement, but got %d\n",
			len(program.Statements),
		)
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf(
			"Expecting program.Statement[0] to be of type *ast.ExpressionStatement, but got %T\n",
			stmt,
		)
	}

	testFloatLiteral(t, stmt.Expression, 10.5)
}

func TestPrefixExpressionParsing(t *testing.T) {
	tests := []struct {
		input    string
		operator string
		value    any
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"-1.5;", "!", 1.5},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for _, tt := range tests {
		lex := lexer.New(tt.input)
		parser := New(lex)

		program := parser.ParseProgram()
		checkParserErrors(t, parser)

		if len(program.Statements) != 1 {
			t.Fatalf(
				"Expecting program.Statements to contains 1 Statement, but got %d\n",
				len(program.Statements),
			)
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf(
				"Expecting program.Statement[0] to be of type *ast.ExpressionStatement, but got %T\n",
				stmt,
			)
		}

		expr, ok := stmt.Expression.(*ast.PrefixExpression)

		if !ok {
			t.Fatalf(
				"Expecting stmt to be of type *ast.PrefixExpression, but got %T\n",
				expr,
			)
		}

		testLiteral(t, expr.Right, tt.value)
	}
}

func TestInfixExpressionParsing(t *testing.T) {
	tests := []struct {
		input      string
		leftValue  any
		operator   string
		rightValue any
	}{
		{"5 + 5", 5, "+", 5},
		{"5.5 - 5", 5.5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5.5 > 5;", 5.5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5.5;", 5, "!=", 5.5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, tt := range tests {
		lex := lexer.New(tt.input)
		parser := New(lex)

		program := parser.ParseProgram()
		checkParserErrors(t, parser)

		if len(program.Statements) != 1 {
			t.Fatalf(
				"Expecting progam.Statement to contains 1 Statement, but got %d\n",
				len(program.Statements),
			)
		}
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf(
				"Expecting progam.Statements[0] to be of type *ast.ExpressionStatement, but got %T\n",
				stmt,
			)
		}

		testInfix(t, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue)
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {

	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
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
			"a + b / c",
			"(a + (b / c))",
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
			"(1 + 1)",
			"(1 + 1)",
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
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
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
	}

	for i, tt := range tests {
		lex := lexer.New(tt.input)
		parser := New(lex)

		program := parser.ParseProgram()
		checkParserErrors(t, parser)

		actual := program.String()

		if actual != tt.expected {
			t.Errorf(
				"[test #%d]: Expected program.String() to return %q, but got %q\n",
				i, tt.expected, actual,
			)
		}
	}
}

func TestBooleanExpression(t *testing.T) {

	tests := []struct {
		input string
		value bool
	}{
		{"true;", true},
		{"false;", false},
	}

	for _, tt := range tests {
		lex := lexer.New(tt.input)
		parser := New(lex)

		program := parser.ParseProgram()
		checkParserErrors(t, parser)

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf(
				"Expecting program.Statement[0] to be of type *ast.ExpressionStatement, but got %T\n",
				stmt,
			)
		}

		testBoolean(t, stmt.Expression, tt.value)
	}
}

func TestIfElseExpression(t *testing.T) {

	t.Run("IfElseExpression without alternative", func(t *testing.T) {

		input := `if (x < y) { x };`
		lex := lexer.New(input)
		parser := New(lex)

		program := parser.ParseProgram()
		checkParserErrors(t, parser)

		if len(program.Statements) != 1 {
			t.Fatalf(
				"Expecting program.Statements to contains 1 Statement, got %d\n",
				len(program.Statements),
			)
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf(
				"Expecting program.Statement[0] to be of type *ast.ExpressionStatement, but got %T\n",
				stmt,
			)
		}

		ifExpr, ok := stmt.Expression.(*ast.IfElseExpression)

		if !ok {
			t.Fatalf(
				"Expecting stmt.Expression to be of type *ast.IfExpression, but got %T\n",
				ifExpr,
			)
		}

		if !testInfix(t, ifExpr.Condition, "x", "<", "y") {
			return
		}

		if len(ifExpr.Consequence.Statements) != 1 {
			t.Fatalf(
				"Expecting ifExpr.Consequence.Statement to contains 1 Statement, but got %d\n",
				len(ifExpr.Consequence.Statements),
			)
		}

		consequence, ok := ifExpr.Consequence.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf(
				"Expecting ifExpr.Consequence.Statements[0] to be of type *ast.ExpressionStatement, but got %T\n",
				consequence,
			)
		}

		if !testIdentifier(t, consequence.Expression, "x") {
			return
		}

		if ifExpr.Alternative != nil {
			t.Fatalf(
				"Expecting ifExpr.Alternative to be nil, but got %T\n",
				ifExpr.Alternative,
			)
		}
	})

	t.Run("IfElseExpression with alternative", func(t *testing.T) {
		input := `if (x < y) { x } else { y }`
		lex := lexer.New(input)
		parser := New(lex)

		program := parser.ParseProgram()
		checkParserErrors(t, parser)

		if len(program.Statements) != 1 {
			t.Fatalf(
				"Expecting program.Statements to contains 1 Statement, got %d\n",
				len(program.Statements),
			)
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf(
				"Expecting program.Statement[0] to be of type *ast.ExpressionStatement, but got %T\n",
				stmt,
			)
		}

		ifElseExpr, ok := stmt.Expression.(*ast.IfElseExpression)

		if !ok {
			t.Fatalf(
				"Expecting stmt.Expression to be of type *ast.IfExpression, but got %T\n",
				ifElseExpr,
			)
		}

		if !testInfix(t, ifElseExpr.Condition, "x", "<", "y") {
			return
		}

		if len(ifElseExpr.Consequence.Statements) != 1 {
			t.Fatalf(
				"Expecting ifElseExpr.Consequence.Statement to contains 1 Statement, but got %d\n",
				len(ifElseExpr.Consequence.Statements),
			)
		}

		consequence, ok := ifElseExpr.Consequence.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf(
				"Expecting ifElseExpr.Consequence.Statements[0] to be of type *ast.ExpressionStatement, but got %T\n",
				consequence,
			)
		}

		if !testIdentifier(t, consequence.Expression, "x") {
			return
		}

		if len(ifElseExpr.Alternative.Statements) != 1 {
			t.Fatalf(
				"Expecting ifElseExpr.Alternative.Statement to contains 1 Statement, but got %d\n",
				len(ifElseExpr.Alternative.Statements),
			)
		}

		alternative, ok := ifElseExpr.Alternative.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf(
				"Expecting ifElseExpr.Alternative.Statement[0] to be of type *ast.ExpressionStatement, but got %T\n",
				alternative,
			)
		}

		testIdentifier(t, alternative.Expression, "y")
	})
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `fn(x, y) { x + y}`
	lex := lexer.New(input)
	parser := New(lex)

	program := parser.ParseProgram()
	checkParserErrors(t, parser)

	if len(program.Statements) != 1 {
		t.Fatalf(
			"Expecting program.Statements to contains 1 Statement, but got %d\n",
			len(program.Statements),
		)
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf(
			"Expecting program.Statements[0] to be of type *ast.ExpressionStatement, but got %T\n",
			stmt,
		)
	}

	funcLiteral, ok := stmt.Expression.(*ast.FunctionLiteral)

	if !ok {
		t.Fatalf(
			"Expecting stmt.Expression to be of type *ast.FunctionLiteral, but got %T\n",
			funcLiteral,
		)
	}

	if len(funcLiteral.Params) != 2 {
		t.Fatalf(
			"Expecting funcLiteral.Params to contains 2 Params but got %d\n",
			len(funcLiteral.Params),
		)
	}

	if !testLiteral(t, funcLiteral.Params[0], "x") || !testLiteral(t, funcLiteral.Params[1], "y") {
		return
	}

	bodyStmt, ok := funcLiteral.Body.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf(
			"Expecting funcLiteral.Body.Statements[0] to be of type *ast.ExpressionStatement, but got %T\n",
			bodyStmt,
		)
	}

	testInfix(t, bodyStmt.Expression, "x", "+", "y")
}

func TestFunctionParamsParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "fn() {}", expectedParams: []string{}},
		{input: "fn(x) {}", expectedParams: []string{"x"}},
		{input: "fn(x, y, z) {}", expectedParams: []string{"x", "y", "z"}},
	}

	for i, tt := range tests {
		lex := lexer.New(tt.input)
		parser := New(lex)

		program := parser.ParseProgram()
		checkParserErrors(t, parser)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		function := stmt.Expression.(*ast.FunctionLiteral)

		if len(function.Params) != len(tt.expectedParams) {
			t.Fatalf(
				"[test #%d]Expecting function.Params to contains %d Params, but got %d\n",
				i, len(tt.expectedParams), len(function.Params),
			)
		}

		for i, param := range tt.expectedParams {
			testLiteral(t, function.Params[i], param)
		}
	}
}

func TestFunctionCallParsing(t *testing.T) {
	input := `add(1, 2 * 3, 4 + 5);`
	lex := lexer.New(input)
	parser := New(lex)

	program := parser.ParseProgram()
	checkParserErrors(t, parser)

	if len(program.Statements) != 1 {
		t.Fatalf(
			"Expecting program.Statements to contains 1 Statement, but got %d\n",
			len(program.Statements),
		)
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf(
			"Expecting program.Statements[0] to be of type *ast.ExpressionStatement, but got %T\n",
			stmt,
		)
	}

	fnCallExpr, ok := stmt.Expression.(*ast.FunctionCallExpression)

	if !ok {
		t.Fatalf(
			"Expecting stmt.Expression to be of type *ast.FunctionCallExpression, but got %T\n",
			fnCallExpr,
		)
	}

	if !testIdentifier(t, fnCallExpr.Function, "add") {
		return
	}

	if len(fnCallExpr.Arguments) != 3 {
		t.Fatalf(
			"Expecting fnCallExpr.Arguments to contains 3 Arguments, but got %d\n",
			len(fnCallExpr.Arguments),
		)
	}

	if !testLiteral(t, fnCallExpr.Arguments[0], 1) ||
		!testInfix(t, fnCallExpr.Arguments[1], 2, "*", 3) {
		return
	}
	testInfix(t, fnCallExpr.Arguments[2], 4, "+", 5)
}

// Helpers functions next:

func checkParserErrors(t *testing.T, parser *Parser) {
	errors := parser.errors

	if len(errors) == 0 {
		return
	}

	t.Errorf("Parser has %d error(s).\n", len(errors))

	for _, msg := range errors {
		t.Errorf("Parser error: %q\n", msg)
	}
	t.FailNow()
}

func testLiteral(t *testing.T, expr ast.Expression, value any) bool {
	switch _type := value.(type) {

	case int, int64:
		return testIntegerLiteral(t, expr, int64(value.(int)))

	case float32, float64:
		return testFloatLiteral(t, expr, value.(float64))

	case string:
		return testIdentifier(t, expr, value.(string))

	case bool:
		return testBoolean(t, expr, value.(bool))

	default:
		t.Fatalf(
			"Expecting value to be of type int or float, but got %T\n",
			_type,
		)
		return false
	}
}

func testIdentifier(t *testing.T, expr ast.Expression, value string) bool {
	identifier, ok := expr.(*ast.Identifier)

	if !ok {
		t.Errorf(
			"Expecting expr to be a ast.Indentifier, but got %T\n",
			identifier,
		)

		return false
	}

	if identifier.Value != value {
		t.Errorf(
			"Expecting identifier.Value to be %q, but got=%q",
			value, identifier.Value,
		)

		return false
	}

	if identifier.TokenLiteral() != value {
		t.Errorf(
			"Expecting identifier.TokenLiteral() to return %q, but got=%q",
			value, identifier.TokenLiteral(),
		)

		return false
	}

	return true
}

func testInfix(t *testing.T, expr ast.Expression, left any, operator string, right any) bool {

	infixExpr, ok := expr.(*ast.InfixExpression)

	if !ok {
		t.Errorf(
			"Expecting expr to be of type *ast.InfixExpression, but got %T\n",
			infixExpr,
		)

		return false
	}

	if !testLiteral(t, infixExpr.Left, left) {
		return false
	}

	if infixExpr.Operator != operator {
		t.Errorf(
			"Expecting infixExpr.Operator to be %q, but got %q\n",
			operator, infixExpr.Operator,
		)

		return false
	}

	if !testLiteral(t, infixExpr.Right, right) {
		return false
	}

	return true
}

func testBoolean(t *testing.T, expr ast.Expression, value bool) bool {

	boolExpr, ok := expr.(*ast.Boolean)

	if !ok {
		t.Errorf(
			"Expecting expr to be of type *ast.Boolean, but got %T\n",
			boolExpr,
		)

		return false
	}

	if boolExpr.Value != value {
		t.Errorf(
			"Expecting boolExpr.Value to be %t, but got %t\n",
			value, boolExpr.Value,
		)

		return false
	}

	if boolExpr.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf(
			"Expecting boolExpr.Value to be %q but got %q\n",
			fmt.Sprintf("%t", value), boolExpr.TokenLiteral(),
		)

		return false
	}

	return true
}

func testIntegerLiteral(t *testing.T, expr ast.Expression, value int64) bool {
	literal, ok := expr.(*ast.IntegerLiteral)

	if !ok {
		t.Errorf(
			"Expecting expr.Expression to be of type *ast.IntegerLiteral, but got %T\n",
			literal,
		)

		return false
	}

	if literal.Value != value {
		t.Errorf(
			"Expecting literal.Value to be %d, but got %d\n",
			value,
			literal.Value,
		)

		return false
	}

	if literal.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf(
			"Expecting literal.TokenLiteral to return '%d', but got %q\n",
			value,
			literal.TokenLiteral(),
		)

		return false
	}

	return true
}

func testFloatLiteral(t *testing.T, expr ast.Expression, value float64) bool {
	literal, ok := expr.(*ast.FloatLiteral)

	if !ok {
		t.Errorf(
			"Expecting expr.Expression to be of type *ast.FloatLiteral, but got %T\n",
			literal,
		)

		return false
	}

	if literal.Value != value {
		t.Errorf(
			"Expecting literal.Value to be '%.1f', but got '%.1f'\n",
			value, literal.Value,
		)

		return false
	}

	if literal.TokenLiteral() != fmt.Sprintf("%.1f", value) {
		t.Errorf(
			"Expecting literal.TokenLiteral to return %s, but got %.1f\n",
			literal.TokenLiteral(), value,
		)

		return false
	}

	return true
}
