package parser

import (
	"monkey/internal/ast"
	"monkey/internal/lexer"
	"monkey/internal/token"
)

type Parser struct {
	lex				*lexer.Lexer

	currentToken	token.Token
	peekToken		token.Token
}

func New(lex *lexer.Lexer) *Parser {
	parser := &Parser{
		lex: lex,
	}

	parser.nextToken()
	parser.nextToken()

	return parser
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lex.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.currentToken.Type != token.EOF {
		stmt := p.parseStatement()

		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}


func (p *Parser) parseStatement() ast.Statement {

	switch p.currentToken.Type {

	case token.CONST, token.LET:
		return p.parseDeclarationStatement()

	default:
		return nil
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

	return false
}
