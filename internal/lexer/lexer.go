package lexer

import (
	"monkey/internal/helper"
	"monkey/internal/token"
	"slices"
)

type Lexer struct {
	input      string
	currentPos int  // current char position in input (current char)
	nextPos    int  // next char position (after current char)
	char       byte // current char under examination
}


func (lex *Lexer) NextToken() (_token token.Token) {

	lex.skipWhitespace()

	switch true {

	case lex.char == 0:
		_token = lex.newToken(token.EOF, "")

	case helper.IsDigit(lex.char):
		literal, _type := lex.readNumber()
		_token.Literal = literal
		_token.Type = _type
		return
	
	case lex.isStartOfTwoCharToken():
		_token = lex.getTwoCharToken()

	case slices.Contains(token.SPECIAL_CHARS_KEYS, lex.char):
		_token = lex.newSpecialCharToken(lex.char)

	case helper.IsCharAllowedInKeyOrVar(lex.char):
		_token.Literal = lex.readWord()
		_token.Type = token.LookupWord(_token.Literal)
		return

	default:
		_token = lex.newToken(token.ILLEGAL, "")
	}

	lex.readChar()

	return
}

// readChar move the reading position and set the next char
// as the current char. This action also update the `nexPos`.
// If the end of the input is reached, it set char to 0 which
// is the equivalent of NUL, in our case an EOF.
func (lex *Lexer) readChar() {

	if lex.nextPos >= len(lex.input) {
		lex.char = 0
	} else {
		lex.char = lex.input[lex.nextPos]
	}
	lex.currentPos = lex.nextPos
	lex.nextPos++
}

// peekChar return the next character to be read.
func (lex *Lexer) peekChar() byte {
	if lex.nextPos >= len(lex.input) {
		return 0
	}

	return lex.input[lex.nextPos]
}

// skipWhitespace skip all space, tab, newline or carriage return
func (lex *Lexer) skipWhitespace() {
	if slices.Contains([]byte{' ', '\t', '\n', '\r'}, lex.char) {
		lex.readChar()
	}
}

// readWord read and return a keyword like "let", "const", "fn"
// or return an identifier
func (lex *Lexer) readWord() string {
	currentPos := lex.currentPos

	for helper.IsCharAllowedInKeyOrVar(lex.char) {
		lex.readChar()
	}

	return lex.input[currentPos:lex.currentPos]
}

// readNumber read and return a number and a boolean that specify
// the type of the number. If true is return the number is an
// integer, otherwise it's a float.
func (lex *Lexer) readNumber() (string, token.TokenType) {
	var number string

	_type := token.INTEGER

	for helper.IsDigit(lex.char) || (_type == token.INTEGER && lex.char == '.') {
		
		if lex.char == '.' && helper.IsDigit(lex.peekChar()) {
			_type = token.FLOAT
			number += string(lex.char)
		} else {
			number += string(lex.char)
		}
		
		lex.readChar()
	}

	return number, _type
}

func (lex *Lexer) newToken(_type token.TokenType, literal string) token.Token {
	return token.Token{
		Type:    _type,
		Literal: literal,
	}
}

func (lex *Lexer) getTwoCharToken() token.Token {
	_literal := string(lex.char)
	lex.readChar()
	_literal += string(lex.char)

	if _type, ok := token.TWO_CHARS[_literal]; ok {
		return token.Token{
			Type: _type,
			Literal: _literal,
		}
	}

	return token.Token{
		Type: token.ILLEGAL,
		Literal: "",
	}
}

// isStartOfTwoCharToken check if the current token is a start
// of a two character token like '==' or '!='
func (lex *Lexer) isStartOfTwoCharToken() bool {
	return slices.Contains([]byte{ '=', '!', '<', '>' }, lex.char) && lex.peekChar() == '='
}

// newSpecialCharToken return new special character like '+', '=' token
func (lex *Lexer) newSpecialCharToken(char byte) token.Token {
	literal := string(char)
	_type := token.SPECIAL_CHARS[char]

	return token.Token{Type: _type, Literal: literal}
}

func New(input string) *Lexer {
	lex := &Lexer{input: input}
	lex.readChar()

	return lex
}
