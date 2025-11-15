package token

import (
	"maps"
	"slices"
)

type TokenType int

const (
	EOF TokenType = iota
	ILLEGAL

	// identifiers, ...etc
	IDENTIFIER
	NUMBER

	// Operators
	ASSIGN
	PLUS
	MINUS
	ASTERISK
	SLASH
	MODULO

	BANG
	LESSER_THAN
	GREATER_THAN

	EQUAL
	NOT_EQUAL

	LESSER_OR_EQUAL_TO
	GREATER_OR_EQUAL_TO
	// Delimiters
	COMMA
	SEMICOLON

	LPAREN   // (
	RPAREN   // )
	LBRACE   // {
	RBRACE   // }
	LBRACKET // [
	RBRACKET // ]

	// Keywords
	FUNCTION
	LET
	CONST
	TRUE
	FALSE
	IF
	ELSE
	RETURN
)

var SPECIAL_CHARS = map[byte]TokenType{
	'=': ASSIGN,
	'+': PLUS,
	'-': MINUS,
	'*': ASTERISK,
	'/': SLASH,
	'%': MODULO,
	'!': BANG,
	'<': LESSER_THAN,
	'>': GREATER_THAN,
	',': COMMA,
	';': SEMICOLON,
	'(': LPAREN,
	')': RPAREN,
	'{': LBRACE,
	'}': RBRACE,
	'[': LBRACKET,
	']': RBRACKET,
}

var SPECIAL_CHARS_KEYS = slices.AppendSeq([]byte{}, maps.Keys(SPECIAL_CHARS))

var TWO_CHARS = map[string]TokenType{
	"==": EQUAL,
	"!=": NOT_EQUAL,
	"<=": LESSER_OR_EQUAL_TO,
	">=": GREATER_OR_EQUAL_TO,
}

var KEYWORDS = map[string]TokenType{
	"fn": FUNCTION,
	"let": LET,
	"const": CONST,
	"true": TRUE,
	"false": FALSE,
	"if": IF,
	"else": ELSE,
	"return": RETURN,
}

type Token struct {
	Type    TokenType
	Literal string
}

// LookupWord serach for the type of a given word.
// If the word is not a keyword, it return IDENTIFIER TokenType
func LookupWord(w string) TokenType {
	if _type, ok := KEYWORDS[w]; ok {
		return _type
	}

	return IDENTIFIER
}
