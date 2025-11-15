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
