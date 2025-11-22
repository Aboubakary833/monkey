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

type IntegerLiteral struct {
	Token		token.Token
	Value		int64
}
func (il *IntegerLiteral) expressionNode() {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string { return il.TokenLiteral() }



type FloatLiteral struct {
	Token		token.Token
	Value		float64
}
func (fl *FloatLiteral) expressionNode() {}
func (fl *FloatLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FloatLiteral) String() string { return fl.TokenLiteral() }



type Boolean struct {
	Token 		token.Token
	Value		bool
}
func (b *Boolean) expressionNode() {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string { return b.TokenLiteral() }



type PrefixExpression struct {
	Token		token.Token
	Operator	string
	Right		Expression
}
func (pxe *PrefixExpression) expressionNode() {}
func (pxe *PrefixExpression) TokenLiteral() string { return pxe.Token.Literal }
func (pxe *PrefixExpression) String() string {
	var output bytes.Buffer

	output.WriteString("(")
	output.WriteString(pxe.Operator)
	output.WriteString(pxe.Right.String())
	output.WriteString(")")

	return output.String()
}



type InfixExpression struct {
	Token		token.Token
	Left		Expression
	Operator	string
	Right		Expression
}
func (ixf *InfixExpression) expressionNode() {}
func (ixf *InfixExpression) TokenLiteral() string { return ixf.Token.Literal }
func (ixf *InfixExpression) String() string {
	var output bytes.Buffer

	output.WriteString("(")
	output.WriteString(ixf.Left.String())
	output.WriteString(" " + ixf.Operator + " ")
	output.WriteString(ixf.Right.String())
	output.WriteString(")")

	return output.String()
}


type BlockStatement struct {
	Token			token.Token
	Statements		[]Statement
}
func (block *BlockStatement) statementNode() {}
func (block *BlockStatement) TokenLiteral() string { return block.Token.Literal }
func (block *BlockStatement) String() string {
	var output bytes.Buffer

	for _, stmt := range block.Statements {
		output.WriteString(stmt.String())
	}

	return output.String()
}


type IfElseExpression struct {
	Token			token.Token
	Condition		Expression
	Consequence		*BlockStatement
	Alternative		*BlockStatement
}
func (ieExpr *IfElseExpression) expressionNode() {}
func (ieExpr *IfElseExpression) TokenLiteral() string { return ieExpr.Token.Literal }
func (ieExpr *IfElseExpression) String() string {
	var output bytes.Buffer

	output.WriteString("if")
	output.WriteString(ieExpr.Condition.String())
	output.WriteString(ieExpr.Consequence.String())

	if ieExpr.Alternative != nil {
		output.WriteString("else")
		output.WriteString(ieExpr.Alternative.String())
	}

	return output.String()
}



type FunctionLiteral struct {
	Token		token.Token
	Params		[]*Identifier
	Body		*BlockStatement
}
func (fn *FunctionLiteral) expressionNode() {}
func (fn *FunctionLiteral) TokenLiteral() string { return fn.Token.Literal }
func (fn *FunctionLiteral) String() string {
	var output bytes.Buffer
	var lastIdx = len(fn.Params) - 1

	output.WriteString(fn.TokenLiteral())
	output.WriteString("(")

	for i, param := range fn.Params {
		output.WriteString(param.String())
		
		if i < lastIdx {
			output.WriteString(", ")
		}
	}

	output.WriteString(") ")
	output.WriteString(fn.Body.String())

	return output.String()
}



type FunctionCallExpression struct {
	Token		token.Token
	Function	Expression // Identifier or FunctionLiteral
	Arguments	[]Expression
}
func (fnCall *FunctionCallExpression) expressionNode() {}
func (fnCall *FunctionCallExpression) TokenLiteral() string { return fnCall.Token.Literal }
func (fnCall *FunctionCallExpression) String() string {
	var output bytes.Buffer
	var lastIdx = len(fnCall.Arguments) - 1

	output.WriteString(fnCall.Function.String())
	output.WriteString("(")

	for i, arg := range fnCall.Arguments {
		output.WriteString(arg.String())
		
		if i < lastIdx {
			output.WriteString(", ")
		}
	}

	output.WriteString(")")

	return output.String()
}
