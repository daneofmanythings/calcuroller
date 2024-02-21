package ast

import (
	"bytes"

	"github.com/daneofmanythings/diceroni/pkg/interpreter/token"
)

type Node interface {
	TokenLiteral() string
	String() string
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
	if len(p.Statements) == 0 {
		return ""
	}
	return p.Statements[0].TokenLiteral()
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

type ExpressionStatement struct {
	Expression Expression
	Token      token.Token
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression == nil {
		return ""
	}
	return es.Expression.String()
}

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

type DiceLiteral struct {
	Token         token.Token
	Value         IntegerLiteral
	QuantModifier IntegerLiteral
	MaxModifier   IntegerLiteral
	MinModifier   IntegerLiteral
	HighModifier  IntegerLiteral
	LowModifier   IntegerLiteral
}

func (dl *DiceLiteral) expressionNode()      {}
func (dl *DiceLiteral) TokenLiteral() string { return dl.Token.Literal }
func (dl *DiceLiteral) String() string {
	var out bytes.Buffer

	if dl.QuantModifier.Value > 0 {
		out.WriteString(dl.QuantModifier.String())
	}

	out.WriteString("d")
	out.WriteString(dl.Value.String())

	if dl.MinModifier.Value > 0 {
		out.WriteString("mi" + dl.MinModifier.String())
	}
	if dl.MaxModifier.Value > 0 {
		out.WriteString("ma" + dl.MaxModifier.String())
	}
	if dl.LowModifier.Value > 0 {
		out.WriteString("ml" + dl.LowModifier.String())
	}
	if dl.HighModifier.Value > 0 {
		out.WriteString("mh" + dl.HighModifier.String())
	}

	return out.String()
}

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}
