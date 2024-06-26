package ast

import (
	"bytes"
	"fmt"

	"github.com/daneofmanythings/calcuroller/pkg/interpreter/token"
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
	Tags  []string
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

type IllegalLiteral struct {
	Token   token.Token
	Literal string
	pos     int
}

func (il *IllegalLiteral) expressionNode()      {}
func (il *IllegalLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IllegalLiteral) String() string       { return il.Literal }

type DiceLiteral struct {
	Token       token.Token
	Tags        []string
	Size        uint32
	Quantity    uint32
	MaxValue    uint32
	MinValue    uint32
	KeepHighest uint32
	KeepLowest  uint32
}

func (dl *DiceLiteral) expressionNode()      {}
func (dl *DiceLiteral) TokenLiteral() string { return dl.Token.Literal }
func (dl *DiceLiteral) String() string {
	var out bytes.Buffer

	if dl.Quantity > 0 {
		out.WriteString(fmt.Sprintf("%d", dl.Quantity))
	}

	out.WriteString("d")
	out.WriteString(fmt.Sprintf("%d", dl.Size))

	if dl.MinValue > 0 {
		out.WriteString("mi" + fmt.Sprintf("%d", dl.MinValue))
	}
	if dl.MaxValue > 0 {
		out.WriteString("ma" + fmt.Sprintf("%d", dl.MaxValue))
	}
	if dl.KeepLowest > 0 {
		out.WriteString("kl" + fmt.Sprintf("%d", dl.KeepLowest))
	}
	if dl.KeepHighest > 0 {
		out.WriteString("kh" + fmt.Sprintf("%d", dl.KeepHighest))
	}
	for _, tag := range dl.Tags {
		out.WriteString("[")
		out.WriteString(tag)
		out.WriteString("]")
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
