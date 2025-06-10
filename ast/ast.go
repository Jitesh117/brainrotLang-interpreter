package ast

import (
	"bytes"

	"github.com/Jitesh117/brainrotLang-interpreter/token"
)

type Node interface {
	TokenLiteral() string
	String() string
}
type Statement interface {
	Node
	StatementNode()
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
	} else {
		return ""
	}
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

type YeetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (ls *YeetStatement) StatementNode()       {}
func (ls *YeetStatement) TokenLiteral() string { return ls.Token.Literal }

func (ls *YeetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";")

	return out.String()
}

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

type SlayStatement struct {
	Token     token.Token
	Slayvalue Expression
}

func (rs *SlayStatement) StatementNode()       {}
func (rs *SlayStatement) TokenLiteral() string { return rs.Token.Literal }

func (rs *SlayStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.Slayvalue != nil {
		out.WriteString(rs.Slayvalue.String())
	}

	out.WriteString(";")

	return out.String()
}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) StatementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }

func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }
