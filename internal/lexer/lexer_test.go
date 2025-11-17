package lexer

import (
	"monkey/internal/token"
	"testing"
)

func TestReadChar(t *testing.T) {
	t.Run("readChar should set char to 'l'", func(t *testing.T) {
		input := "let x = 5"
		lex := Lexer{input: input}

		lex.readChar()

		if lex.currentPos != 0 {
			t.Fatalf(
				"Expected lexer current position to be 0, but got %d\n",
				lex.currentPos,
			)
		}

		if lex.nextPos != 1 {
			t.Fatalf(
				"Expected lexer next position to be 1, but got %d\n",
				lex.nextPos,
			)
		}

		if lex.char != byte('l') {
			t.Fatalf(
				"Expected lexer current char to be 'l', but got '%c'\n",
				lex.char,
			)
		}
	})

	t.Run("readChar should set char to 'n' and update positions", func(t *testing.T) {
		input := "const INVARIABLE = 0"
		lex := Lexer{input: input, currentPos: 1, nextPos: 2, char: 'o'}

		lex.readChar()

		if lex.currentPos != 2 {
			t.Fatalf(
				"Expected lexer current position to be 2, but got %d\n",
				lex.currentPos,
			)
		}

		if lex.nextPos != 3 {
			t.Fatalf(
				"Expected lexer next position to be 3, but got %d\n",
				lex.nextPos,
			)
		}

		if lex.char != byte('n') {
			t.Fatalf(
				"Expected lexer current char to be 'n', but got '%c'\n",
				lex.char,
			)
		}
	})

	t.Run("readChar should set char to 0 if it reached the end of line", func(t *testing.T) {
		input := "let func = fn(x, y) { x + y }"
		lex := Lexer{input: input, currentPos: 30, nextPos: 31, char: '}'}

		lex.readChar()

		if lex.currentPos != 31 {
			t.Fatalf(
				"Expected lexer current position to be 31, but got %d\n",
				lex.currentPos,
			)
		}

		if lex.nextPos != 32 {
			t.Fatalf(
				"Expected lexer next position to be 32, but got %d\n",
				lex.nextPos,
			)
		}

		if lex.char != 0 {
			t.Fatalf(
				"Expected lexer current char to be 0, but got '%c'\n",
				lex.char,
			)
		}
	})
}

func TestSkipWhitespace(t *testing.T) {
	t.Run("skipWhitespace should skip space", func(t *testing.T) {
		input := "let x = 5"
		lex := Lexer{input: input, currentPos: 3, nextPos: 4, char: ' '}
		lex.skipWhitespace()

		if lex.char != 'x' {
			t.Fatalf(
				"Expected skipWhitespace to skip space and set char to 'm', but go '%c'\n",
				lex.char,
			)
		}
	})

	t.Run("skipWhitespace should skip LN & TAB", func(t *testing.T) {
		input := `let multiplyByTwo = fn(x) {
	return x * 2
}`
		lex := Lexer{input: input}

		for lex.char != '{' {
			lex.readChar()
		}
		lex.readChar()

		for range 2 {
			lex.skipWhitespace()
		}

		if lex.char != 'r' {
			t.Fatalf(
				"Expected skipWhitespace to skip LN & TAB and set char to 'r', but go '%c'\n",
				lex.char,
			)
		}
	})
}

func TestIsStartOfTwoCharToken(t *testing.T) {

	t.Run("it should return false", func(t *testing.T) {
		lex := Lexer{ input: "const X = 10", currentPos: 8, nextPos: 9, char: '=' }

		if lex.isStartOfTwoCharToken() {
			t.Fatal("Expected isStartOfTwoCharToken() to return false, but got true.")
		}
	})

	t.Run("it should return true", func(t *testing.T) {
		var lex Lexer
		lexDataSlices := []struct{
			input 		string
			char		byte
		}{
			{ "x == 0", '=' },
			{ "x != 0", '!' },
			{ "x <= 0", '<' },
			{ "x >= 0", '>' },
		}

		for i, s := range lexDataSlices {
			lex = Lexer{ input: s.input, currentPos: 2, nextPos: 3, char: s.char }

			if !lex.isStartOfTwoCharToken() {
				t.Fatalf(
					"[tests #%d]: Expected isStartOfTwoCharToken() to return true, but got false.",
					i,
				)
			}
		}
	})
}

func TestNextToken(t *testing.T) {

	input := `const five = 5;
let ten = 10;
let add = fn(x, y) {
x + y;
};
let result = add(five, ten);
!-/*5;
5 < 10 > 5;
if (5 < 10) {
return true;
} else {
return false;
}
10 == 10;
10.5 != 9;
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.CONST, "const"},
		{token.IDENTIFIER, "five"},
		{token.ASSIGN, "="},
		{token.INTEGER, "5"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENTIFIER, "ten"},
		{token.ASSIGN, "="},
		{token.INTEGER, "10"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENTIFIER, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENTIFIER, "x"},
		{token.COMMA, ","},
		{token.IDENTIFIER, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENTIFIER, "x"},
		{token.PLUS, "+"},
		{token.IDENTIFIER, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENTIFIER, "result"},
		{token.ASSIGN, "="},
		{token.IDENTIFIER, "add"},
		{token.LPAREN, "("},
		{token.IDENTIFIER, "five"},
		{token.COMMA, ","},
		{token.IDENTIFIER, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.ASTERISK, "*"},
		{token.INTEGER, "5"},
		{token.SEMICOLON, ";"},
		{token.INTEGER, "5"},
		{token.LESSER_THAN, "<"},
		{token.INTEGER, "10"},
		{token.GREATER_THAN, ">"},
		{token.INTEGER, "5"},
		{token.SEMICOLON, ";"},
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.INTEGER, "5"},
		{token.LESSER_THAN, "<"},
		{token.INTEGER, "10"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		
		{token.INTEGER, "10"},
		{token.EQUAL, "=="},
		{token.INTEGER, "10"},
		{token.SEMICOLON, ";"},
		{token.FLOAT, "10.5"},
		{token.NOT_EQUAL, "!="},
		{token.INTEGER, "9"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	lex := New(input)

	for i, tt := range tests {

		_token := lex.NextToken()
		
		if _token.Literal != tt.expectedLiteral {
			t.Fatalf(
				"[test #%d] - Wrong token literal. Expected literal %q, got %q\n",
				i, tt.expectedLiteral, _token.Literal,
			)
		}
		if _token.Type != tt.expectedType {
			t.Fatalf(
				"[test #%d] - Wrong token type. Expected %d, got %d\n",
				i, tt.expectedType, _token.Type,
			)
		}


	}
}
