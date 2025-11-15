package parser

import (
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

	if program == nil {
		t.Fatal("ParseProgam() return nil")
	}

	tests := []struct{
		expectedIdentifier	string
	}{
		{ "XYZ" },
		{ "t" },
		{ "foobar" },
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

func checkParserErrors(t *testing.T, parser *Parser) {
	errors := parser.errors;

	if len(errors) == 0 { return }

	t.Errorf("Parser has %d error(s).\n", len(errors))

	for msg := range errors {
		t.Errorf("Parser error: %q\n", msg)
	}
	t.FailNow()
}
