package token

import (
	"maps"
	"monkey/internal/helper"
	"slices"
)

type TokenType int

const (
	EOF TokenType = iota
	ILLEGAL

	// identifiers, ...etc
	IDENTIFIER
	INTEGER
	FLOAT

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

var FLIPPED_SPECIAL_CHARS = helper.FlipMap(SPECIAL_CHARS)

var SPECIAL_CHARS_KEYS = slices.AppendSeq([]byte{}, maps.Keys(SPECIAL_CHARS))


var TWO_CHARS = map[string]TokenType{
	"==": EQUAL,
	"!=": NOT_EQUAL,
	"<=": LESSER_OR_EQUAL_TO,
	">=": GREATER_OR_EQUAL_TO,
}

var FLIPPED_TWO_CHARS = helper.FlipMap(TWO_CHARS)


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

var FLIPPED_KEYWORDS = helper.FlipMap(KEYWORDS)


var OTHERS = map[string]TokenType{
	"identifier": IDENTIFIER,
	"integer": INTEGER,
	"float": FLOAT,
	"eof": EOF,
	"illegal": ILLEGAL,
}

var FLIPPED_OTHERS = helper.FlipMap(OTHERS)


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

func GetLiteralByType(_type TokenType) string {

	strValueMaps := []map[TokenType]string{
		FLIPPED_TWO_CHARS,
		FLIPPED_KEYWORDS,
		FLIPPED_OTHERS,
	}

	if l, ok := FLIPPED_SPECIAL_CHARS[_type]; ok {
		return string(l)
	}
	
	for _, m := range strValueMaps {
		if l, ok := m[_type]; ok {
			return l
		}
	}
	
	return ""
}
