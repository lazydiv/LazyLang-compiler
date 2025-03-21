package codegen

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/lazydiv/lazyLang-compiler/internal/parser"
)

type CodeGen struct {
	variables map[string]bool
}

func NewCodeGen() *CodeGen {
	return &CodeGen{
		variables: make(map[string]bool),
	}
}
func (cg *CodeGen) Generate(program *parser.Program) string {
	var out strings.Builder

	out.WriteString("package main\n\n")
	out.WriteString("import \"fmt\"\n\n")
	out.WriteString("func main() {\n")

	for _, stmt := range program.Statements {
		out.WriteString("\t" + cg.generateStatement(stmt) + "\n")
	}

	out.WriteString("}\n")
	return out.String()
}

// compiler transalotr from ast to go CodeGen
// compiler to asm
func (cg *CodeGen) generateStatement(stmt parser.Statement) string {
	switch s := stmt.(type) {
	case *parser.VarStatement:
		varName := s.Name
		varExpr := cg.generateExpression(s.Value)

		if cg.variables[varName] {
			return fmt.Sprintf("%s = %s", varName, varExpr)
		} else {
			cg.variables[varName] = true
			return fmt.Sprintf("%s := %s", varName, varExpr)
		}
	case *parser.ForStatement:
		var out strings.Builder
		out.WriteString("for ")

		// Handle initialization
		if s.Init != nil {
			initStmt := cg.generateStatement(s.Init)
			if varStmt, ok := s.Init.(*parser.VarStatement); ok {
				varName := varStmt.Name
				if cg.variables[varName] {
					// Variable already declared, use "=" instead of ":="
					if parts := strings.Split(initStmt, ":="); len(parts) == 2 {
						initStmt = parts[0] + "= " + parts[1]
					}
				}
			}
			out.WriteString(initStmt)
		}
		out.WriteString("; ")

		// Handle condition
		if s.Condition != nil {
			out.WriteString(cg.generateExpression(s.Condition))
		}
		out.WriteString("; ")

		// Handle post statement
		if s.Post != nil {
			postStmt := cg.generateStatement(s.Post)
			if varStmt, ok := s.Post.(*parser.VarStatement); ok {
				varName := varStmt.Name
				if cg.variables[varName] {
					// Variable already declared, use "=" instead of ":="
					if parts := strings.Split(postStmt, ":="); len(parts) == 2 {
						postStmt = parts[0] + "= " + parts[1]
					}
				}
			}
			out.WriteString(postStmt)
		}

		out.WriteString(" {\n")

		for _, stmt := range s.Body {
			out.WriteString("\t" + cg.generateStatement(stmt) + "\n")
		}

		out.WriteString("}")

		return out.String()

	case *parser.ArrayStatement:
		arrName := s.Name
		var out strings.Builder

		if cg.variables[arrName] {
			out.WriteString(fmt.Sprintf("%s = []interface{}{", arrName))
		} else {
			cg.variables[arrName] = true
			out.WriteString(fmt.Sprintf("%s := []interface{}{", arrName))
		}

		for i, v := range s.Values {
			if i != 0 {
				out.WriteString(", ")
			}
			if v != nil {
				out.WriteString(cg.generateExpression(v))
			} else {
				out.WriteString("nil")
			}
		}
		out.WriteString("}")

		return out.String()
	case *parser.IfStatement:
		var out strings.Builder
		condition := cg.generateExpression(s.Condition)

		out.WriteString(fmt.Sprintf("if %s {\n", condition))
		for _, stmt := range s.Consequence {
			out.WriteString("\t\t" + cg.generateStatement(stmt) + "\n")
		}

		if len(s.Alternative) > 0 {
			out.WriteString("\t} else {\n")
			for _, stmt := range s.Alternative {
				out.WriteString("\t\t" + cg.generateStatement(stmt) + "\n")
			}
		}

		out.WriteString("\t}")
		return out.String()

	case *parser.PrintStatement:
		expr := cg.generateExpression(s.Value)
		return fmt.Sprintf("fmt.Println(%s)", expr)
	default:
		return ""
	}
}

func (cg *CodeGen) generateExpression(expr parser.Expression) string {
	switch e := expr.(type) {
	case *parser.Identifier:
		return e.Value
	case *parser.NumberLiteral:
		return strconv.FormatFloat(e.Value, 'f', -1, 64)
	case *parser.InfixExpression:
		left := cg.generateExpression(e.Left)
		right := cg.generateExpression(e.Right)
		return fmt.Sprintf("(%s %s %s)", left, e.Operator, right)
	case *parser.IndexExpression:
		array := cg.generateExpression(e.Array)
		index := cg.generateExpression(e.Index)
		return fmt.Sprintf("%s[%s]", array, index)
	default:
		return ""
	}
}
