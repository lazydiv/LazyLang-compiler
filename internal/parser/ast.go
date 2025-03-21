package parser

import (
	"fmt"
	"strconv"
	"strings"
)

type Node interface {
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

// Program is the root node of our AST
type Program struct {
	Statements []Statement
}

func (p *Program) String() string {
	var out strings.Builder
	for _, s := range p.Statements {
		out.WriteString(s.String() + "\n")
	}
	return out.String()
}

type VarStatement struct {
	Name  string
	Value Expression
}

type ArrayStatement struct {
	Name   string
	Values []Expression
}

func (vs *VarStatement) statementNode() {}

func (vs *VarStatement) String() string {
	return fmt.Sprintf("var %s = %s", vs.Name, vs.Value.String())
}

// IndexExpression represents accessing an array element by index: array[index]
type IndexExpression struct {
	Array Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode() {}
func (ie *IndexExpression) String() string {
	return fmt.Sprintf("%s[%s]", ie.Array.String(), ie.Index.String())

}

func (as *ArrayStatement) statementNode() {}
func (as *ArrayStatement) String() string {
	var out strings.Builder
	out.WriteString("lazyArray")
	out.WriteString(as.Name)
	out.WriteString(" = [")
	for i, v := range as.Values {
		if i != 0 {
			out.WriteString(", ")
		}
		out.WriteString(v.String())
	}
	out.WriteString("]")
	return out.String()
}

type IfStatement struct {
	Condition   Expression
	Consequence []Statement
	Alternative []Statement
}

type ForStatement struct {
	Init      Statement   // Initialization statement
	Condition Expression  // Loop condition
	Post      Statement   // Post iteration statement
	Body      []Statement // Loop body
}

func (is *IfStatement) statementNode() {}
func (is *IfStatement) String() string {
	var out strings.Builder
	out.WriteString("if ")
	out.WriteString(is.Condition.String())
	out.WriteString(" { ")
	for _, stmt := range is.Consequence {
		out.WriteString(stmt.String() + "; ")
	}
	out.WriteString(" }")

	if len(is.Alternative) > 0 {
		out.WriteString(" else { ")
		for _, stmt := range is.Alternative {
			out.WriteString(stmt.String() + "; ")
		}
		out.WriteString(" }")
	}

	return out.String()
}

func (fs *ForStatement) statementNode() {}
func (fs *ForStatement) String() string {
	var out strings.Builder
	out.WriteString("for ")

	if fs.Init != nil {
		out.WriteString(fs.Init.String() + "; ")
	}
	if fs.Condition != nil {
		out.WriteString(fs.Condition.String() + "; ")
	}
	if fs.Post != nil {
		out.WriteString(fs.Post.String())
	}
	out.WriteString(" { ")
	for _, stmt := range fs.Body {
		out.WriteString(stmt.String() + "; ")
	}
	out.WriteString(" }")

	return out.String()
}

type PrintStatement struct {
	Value Expression
}

func (ps *PrintStatement) statementNode() {}
func (ps *PrintStatement) String() string {
	return fmt.Sprintf("print(%s)", ps.Value.String())
}

type Identifier struct {
	Value string
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) String() string  { return i.Value }

type NumberLiteral struct {
	Value float64
}

func (nl *NumberLiteral) expressionNode() {}
func (nl *NumberLiteral) String() string  { return strconv.FormatFloat(nl.Value, 'f', -1, 64) }

type InfixExpression struct {
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode() {}
func (ie *InfixExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", ie.Left.String(), ie.Operator, ie.Right.String())
}
