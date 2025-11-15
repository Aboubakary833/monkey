package ast

import (
	"bytes"
	"monkey/internal/token"
)

type Node interface {
	TokenLiteral() 	string
	String()		string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}

	return ""
}

func (p *Program) String() string {
	var output bytes.Buffer

	for _, stmt := range p.Statements {
		output.WriteString(stmt.String())
	}

	return output.String()
}

type DeclarationStatement struct {
	Token 	token.Token
	Name  	*Identifier
	Value 	Expression
}

func (ds *DeclarationStatement) statementNode() {}
func (ds *DeclarationStatement) TokenLiteral() string { return ds.Token.Literal }

func (ds *DeclarationStatement) String() string {
	var output bytes.Buffer

	output.WriteString(ds.TokenLiteral() + " ")
	output.WriteString(ds.Name.String() + " = ")
	
	if ds.Value != nil {
		output.WriteString(ds.Value.String())
	}
	output.WriteString(";")

	return output.String()
}



type Identifier struct {
	Token 	token.Token
	Value 	string
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string { return i.Value }



type ReturnStatement struct {
	Token		token.Token
	ReturnValue	Expression
}
func (rs *ReturnStatement) statementNode() {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }

func (rs *ReturnStatement) String() string {
	var output bytes.Buffer

	output.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		output.WriteString(rs.ReturnValue.String())
	}
	output.WriteString(";")

	return output.String()
}



type ExpressionStatement struct {
	Token		token.Token
	Expression	Expression
}
func (es *ExpressionStatement) statementNode() {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }

func (es *ExpressionStatement) String() string {

	if es.Expression != nil {
		return es.Expression.String()
	}

	return ""
}

