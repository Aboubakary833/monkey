package ast

import (
	"monkey/internal/token"
	"testing"
)

func TestString(t *testing.T) {

	program := &Program{
		[]Statement{
			&DeclarationStatement{
				Token: token.Token{ Type: token.LET, Literal: "let" },
				Name: &Identifier{
					Token: token.Token{ Type: token.IDENTIFIER, Literal: "myVar" },
					Value: "myVar",
				},
				Value: &Identifier{
					Token: token.Token{ Type: token.IDENTIFIER, Literal: "anotherVar" },
					Value: "anotherVar",
				},
			},
		},
	}

	expected := "let myVar = anotherVar;"
	got := program.String()

	if got != expected {
		t.Fatalf(
			"Expecting program.String() to return `%s` but got `%s`\n",
			expected, got,
		)
	}
}
