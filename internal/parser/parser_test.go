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

	identifier, ok := stmt.Expression.(*ast.Identifier)

	if !ok {
		t.Fatalf(
			"Expecting stmt to be a ast.Indentifier, but got %T\n",
			stmt.Expression,
		)
	}

	if identifier.Value != input {
		t.Errorf(
			"Expecting identifier.Value to be %q, but got=%q",
			input, identifier.Value,
		)
	}

	if identifier.TokenLiteral() != input {
		t.Errorf(
			"Expecting identifier.TokenLiteral() to return %q, but got=%q",
			input, identifier.TokenLiteral(),
		)
	}
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

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)

	if !ok {
		t.Fatalf(
			"Expecting stmt.Expression to be of type *ast.IntegerLiteral, but got %T\n",
			literal,
		)
	}

	if literal.Value != 5 {
		t.Fatalf(
			"Expecting literal.Value to be '5', but got '%d'\n",
			literal.Value,
		)
	}

	if literal.TokenLiteral() != "5" {
		t.Fatalf(
			"Expecting literal.TokenLiteral to return '5', but got %q\n",
			literal.TokenLiteral(),
		)
	}
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
	literal, ok := stmt.Expression.(*ast.FloatLiteral)

	if !ok {
		t.Fatalf(
			"Expecting stmt.Expression to be of type *ast.FloatLiteral, but got %T\n",
			literal,
		)
	}

	if literal.Value != 10.5 {
		t.Fatalf(
			"Expecting literal.Value to be '10.5', but got '%.f'\n",
			literal.Value,
		)
	}

	if literal.TokenLiteral() != "10.5" {
		t.Fatalf(
			"Expecting literal.TokenLiteral to return '10.5', but got %q\n",
			literal.TokenLiteral(),
		)
	}
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

		if !testLiteral(t, expr.Right, tt.value) {
			return
		}

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
		expression, ok := stmt.Expression.(*ast.InfixExpression)

		if !ok {
			t.Fatalf(
				"Expecting stmt.Expression to be of type *ast.InfixExpression, but got %T\n",
				expression,
			)
		}

		if !testLiteral(t, expression.Left, tt.leftValue) {
			return
		}

		if expression.Operator != tt.operator {
			t.Fatalf(
				"Expecting expression.Operator to be %q but got %q\n",
				tt.operator, expression.Operator,
			)
		}

		if !testLiteral(t, expression.Right, tt.rightValue) {
			return
		}
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

	default:
		t.Fatalf(
			"Expecting value to be of type int or float, but got %T\n",
			_type,
		)
		return false
	}
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
