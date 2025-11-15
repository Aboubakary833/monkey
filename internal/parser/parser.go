package parser

import (
	"fmt"
	"monkey/internal/ast"
	"monkey/internal/lexer"
	"monkey/internal/token"
)

const (
	_ 		int = iota
	LOWEST
	EQUALS // comparision(==)
	LESS_OR_GREATER // < or >
	SUM // addition(+)
	PRODUCT // *
	PREFIX // -x or !x
	CALL // myFunc(x)
)

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

	//Prefix/Infix parser functions maps registrations
	parser.registerPrefix(token.IDENTIFIER, parser.parseIdentifier)

	parser.nextToken()
	parser.nextToken()

	return parser
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(_type token.TokenType) {
	msg := fmt.Sprintf(
		"Expected next token to be %q, but got %q instead.",
		_type, p.peekToken.Type,
	)

	p.errors = append(p.errors, msg)
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
	prefix := p.prefixParseFns[p.currentToken.Type]
	
	if prefix == nil {
		return nil
	}
	leftExpression := prefix()

	return leftExpression
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{
		Token: p.currentToken,
		Value: p.currentToken.Literal,
	}
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
