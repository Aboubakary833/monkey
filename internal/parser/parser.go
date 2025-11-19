package parser

import (
	"fmt"
	"monkey/internal/ast"
	"monkey/internal/lexer"
	"monkey/internal/token"
	"strconv"
)

const (
	_ 				int = iota
	LOWEST
	EQUALS // comparision(==)
	LESS_OR_GREATER // < or >
	SUM // addition(+)
	PRODUCT // *
	PREFIX // -x or !x
	CALL // myFunc(x)
)

var precedences = map[token.TokenType]int{
	token.EQUAL: EQUALS,
	token.NOT_EQUAL: EQUALS,
	token.LESSER_THAN: LESS_OR_GREATER,
	token.GREATER_THAN: LESS_OR_GREATER,
	token.PLUS: SUM,
	token.MINUS: SUM,
	token.SLASH: PRODUCT,
	token.ASTERISK: PRODUCT,
}

type (
	prefixParseFn 	func() ast.Expression
	infixParseFn 	func(ast.Expression) ast.Expression
)

type Parser struct {
	lex				*lexer.Lexer

	currentToken	token.Token
	peekToken		token.Token

	prefixParseFns	map[token.TokenType]prefixParseFn
	infixParseFns	map[token.TokenType]infixParseFn

	errors			[]string
}

func New(lex *lexer.Lexer) *Parser {
	parser := &Parser{
		lex: lex,
		errors: []string{},
	}

	parser.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	parser.infixParseFns = make(map[token.TokenType]infixParseFn)

	parser.registerPrefixesAndInfixes()

	parser.nextToken()
	parser.nextToken()

	return parser
}

// registerPrefixesAndInfixes register all prefix/infix
// parsers functions.
func (p *Parser) registerPrefixesAndInfixes() {

	// Prefixes
	p.registerPrefix(token.IDENTIFIER, p.parseIdentifier)
	p.registerPrefix(token.INTEGER, p.parseInteger)
	p.registerPrefix(token.FLOAT, p.parseFloat)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)

	// Infixes
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQUAL, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQUAL, p.parseInfixExpression)
	p.registerInfix(token.LESSER_THAN, p.parseInfixExpression)
	p.registerInfix(token.GREATER_THAN, p.parseInfixExpression)

}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) addError(err string) {
	p.errors = append(p.errors, err)
}

func (p *Parser) peekError(_type token.TokenType) {
	msg := fmt.Sprintf(
		"Expected next token to be '%s', but got '%s' instead.",
		token.GetLiteralByType(_type),
		token.GetLiteralByType(p.peekToken.Type),
	)

	p.addError(msg)
}

func (p *Parser) noPrefixParseFnError(_type token.TokenType) {
	msg := fmt.Sprintf(
		"No prefix parse function for '%s' found",
		token.GetLiteralByType(_type),
	)
	p.addError(msg)
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lex.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.currentTokenIs(token.EOF) {
		stmt := p.parseStatement()

		program.Statements = append(program.Statements, stmt)

		p.nextToken()
	}

	return program
}


func (p *Parser) parseStatement() ast.Statement {

	switch p.currentToken.Type {

	case token.CONST, token.LET:
		return p.parseDeclarationStatement()

	case token.RETURN:
		return p.parseReturnStatement()

	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseDeclarationStatement() *ast.DeclarationStatement {
	
	stmt := &ast.DeclarationStatement{Token: p.currentToken}
	
	if !p.expectPeekTokenToBe(token.IDENTIFIER) {
		return nil
	}
	
	stmt.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
	
	if !p.expectPeekTokenToBe(token.ASSIGN) {
		return nil
	}
	// TODO: We're skipping the expressions until we
	// encounter a semicolon
	for !p.currentTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {

	stmt := &ast.ReturnStatement{Token: p.currentToken}

	p.nextToken()

	for !p.currentTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{ Token: p.currentToken }

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix, ok := p.prefixParseFns[p.currentToken.Type]
	
	if !ok {
		p.noPrefixParseFnError(p.currentToken.Type)
		return nil
	}
	leftExpression := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix, ok := p.infixParseFns[p.peekToken.Type]

		if !ok {
			return leftExpression
		}
		p.nextToken()

		leftExpression = infix(leftExpression)
	}

	return leftExpression
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{
		Token: p.currentToken,
		Value: p.currentToken.Literal,
	}
}

func (p *Parser) parseInteger() ast.Expression {
	intLiteral := &ast.IntegerLiteral{ Token: p.currentToken }

	value, err := strconv.ParseInt(p.currentToken.Literal, 0, 64)
	
	if err != nil {
		msg := fmt.Sprintf("Could not parse %q as integer\n", p.currentToken.Literal)
		p.addError(msg)
	}
	intLiteral.Value = value

	return intLiteral
}

func (p *Parser) parseFloat() ast.Expression {
	flotLiteral := &ast.FloatLiteral{ Token: p.currentToken }

	value, err := strconv.ParseFloat(p.currentToken.Literal, 64)

	if err != nil {
		msg := fmt.Sprintf("Could not parse %q as float\n", p.currentToken.Literal)
		p.addError(msg)
	}
	flotLiteral.Value = value

	return flotLiteral
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token: p.currentToken,
		Operator: p.currentToken.Literal,
	}

	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token: p.currentToken,
		Left: left,
		Operator: p.currentToken.Literal,
	}
	precedence := precedences[p.currentToken.Type]
	
	p.nextToken()

	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) registerPrefix(_type token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[_type] = fn
}

func (p *Parser) registerInfix(_type token.TokenType, fn infixParseFn) {
	p.infixParseFns[_type] = fn
}

func (p *Parser) currentTokenIs(_type token.TokenType) bool {
	return p.currentToken.Type == _type
}

func (p *Parser) peekTokenIs(_type token.TokenType) bool {
	return p.peekToken.Type == _type
}

func (p *Parser) expectPeekTokenToBe(_type token.TokenType) bool {

	if p.peekTokenIs(_type) {
		p.nextToken()
		return true
	}
	p.peekError(_type)
	
	return false
}

func (p *Parser) currentPrecedence() int {
	if prec, ok := precedences[p.currentToken.Type]; ok {
		return prec
	}

	return LOWEST
}

func (p *Parser) peekPrecedence() int {
	if prec, ok := precedences[p.peekToken.Type]; ok {
		return prec
	}

	return LOWEST
}
