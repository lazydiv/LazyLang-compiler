package parser

import (
	"github.com/lazydiv/lazyLang-compiler/internal/lexer"
	"strconv"
)

// Parser builds the AST
type Parser struct {
	lexer        *lexer.Lexer
	currentToken lexer.Token
	peekToken    lexer.Token
}

func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{lexer: l}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) expectCurrent(t lexer.TokenType) bool {
	return p.currentToken.Type == t
}

func (p *Parser) expectPeek(t lexer.TokenType) bool {
	if p.peekToken.Type == t {
		p.nextToken()
		return true
	}
	return false
}

func (p *Parser) ParseProgram() *Program {
	program := &Program{Statements: []Statement{}}

	for p.currentToken.Type != lexer.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() Statement {
	switch p.currentToken.Type {
	case lexer.VAR:
		return p.parseVarStatement()
	case lexer.ARRAY:
		return p.parseArray()

	case lexer.FOR:
		return p.parseForStatement()
	case lexer.IF:
		return p.parseIfStatement()
	case lexer.PRINT:
		return p.parsePrintStatement()
	default:
		return nil
	}
}

func (p *Parser) parseArray() *ArrayStatement {
	stmt := &ArrayStatement{}

	// After 'lazyArray', expect an identifier
	if !p.expectPeek(lexer.IDENT) {
		return nil
	}

	stmt.Name = p.currentToken.Literal

	// Expect assignment operator
	if !p.expectPeek(lexer.ASSIGN) {
		return nil
	}

	// Expect opening bracket
	if !p.expectPeek(lexer.LSBREC) {
		return nil
	}

	// Parse array elements
	values := []Expression{}

	// Check if the array is empty
	if p.peekToken.Type != lexer.RSBREC {
		p.nextToken()
		values = append(values, p.parseExpression())

		// Parse remaining elements
		for p.peekToken.Type == lexer.COMMA {
			p.nextToken() // consume comma
			p.nextToken() // move to next expression
			values = append(values, p.parseExpression())
		}
	}

	// Expect closing bracket
	if !p.expectPeek(lexer.RSBREC) {
		return nil
	}

	stmt.Values = values
	return stmt
}

func (p *Parser) parseIndexExpression(left Expression) Expression {
	expr := &IndexExpression{Array: left}

	// Consume the opening bracket
	p.nextToken()

	// Parse the index expression
	expr.Index = p.parseExpression()

	// Expect closing bracket
	if !p.expectPeek(lexer.RSBREC) {
		return nil
	}

	return expr
}

func (p *Parser) parseVarStatement() *VarStatement {
	stmt := &VarStatement{}

	if !p.expectPeek(lexer.IDENT) {
		return nil
	}

	stmt.Name = p.currentToken.Literal

	if !p.expectPeek(lexer.ASSIGN) {
		return nil
	}

	p.nextToken()
	stmt.Value = p.parseExpression()

	return stmt
}

func (p *Parser) parseForStatement() *ForStatement {
	stmt := &ForStatement{}

	if !p.expectPeek(lexer.LPAREN) {
		return nil
	}

	p.nextToken() // Move to the first token of the init statement
	if p.currentToken.Type != lexer.SEMICOLON {
		if p.currentToken.Type == lexer.IDENT {
			name := p.currentToken.Literal

			if !p.expectPeek(lexer.ASSIGN) {
				return nil
			}

			p.nextToken()
			value := p.parseExpression()

			stmt.Init = &VarStatement{
				Name:  name,
				Value: value,
			}
		} else {
			return nil
		}

		if !p.expectPeek(lexer.SEMICOLON) {
			return nil
		}
	} else {
		p.nextToken()
	}

	if p.currentToken.Type != lexer.SEMICOLON {
		stmt.Condition = p.parseExpression()
	}

	if !p.expectPeek(lexer.SEMICOLON) {
		return nil
	}

	if p.peekToken.Type != lexer.RPAREN {
		p.nextToken() // Move to the first token of the post statement

		if p.currentToken.Type == lexer.IDENT {
			name := p.currentToken.Literal

			if !p.expectPeek(lexer.ASSIGN) {
				return nil
			}

			p.nextToken()
			value := p.parseExpression()

			stmt.Post = &VarStatement{
				Name:  name,
				Value: value,
			}
		} else {
			return nil
		}
	}

	if !p.expectPeek(lexer.RPAREN) {
		return nil
	}

	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}

	p.nextToken() // Move to the first token in the body
	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseIfStatement() *IfStatement {
	stmt := &IfStatement{}

	p.nextToken()
	stmt.Condition = p.parseExpression()

	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}

	p.nextToken()
	stmt.Consequence = p.parseBlockStatement()

	if p.peekToken.Type == lexer.ELSE {
		p.nextToken()

		if !p.expectPeek(lexer.LBRACE) {
			return nil
		}

		p.nextToken()
		stmt.Alternative = p.parseBlockStatement()
	}

	return stmt
}

func (p *Parser) parseBlockStatement() []Statement {
	statements := []Statement{}

	for p.currentToken.Type != lexer.RBRACE && p.currentToken.Type != lexer.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			statements = append(statements, stmt)
		}
		p.nextToken()
	}

	return statements
}

func (p *Parser) parsePrintStatement() *PrintStatement {
	stmt := &PrintStatement{}

	if !p.expectPeek(lexer.LPAREN) {
		return nil
	}

	p.nextToken()
	stmt.Value = p.parseExpression()

	if !p.expectPeek(lexer.RPAREN) {
		return nil
	}

	return stmt
}

func (p *Parser) precedence(tokenType lexer.TokenType) int {

	switch tokenType {
	case lexer.LSBREC:
		return 4
	case lexer.MULTIPLY, lexer.DIVIDE:
		return 3
	case lexer.PLUS, lexer.MINUS:
		return 2
	case lexer.GT, lexer.LT, lexer.EQ, lexer.NOT_EQ, lexer.GT_EQ, lexer.LT_EQ:
		return 1
	default:
		return 0
	}
}

func (p *Parser) parseExpression() Expression {
	return p.parseExpressionWithPrecedence(0)
}

func (p *Parser) parseExpressionWithPrecedence(precedence int) Expression {
	left := p.parsePrimary()

	if left == nil {
		return nil
	}

	for p.peekToken.Type != lexer.EOF &&
		p.peekToken.Type != lexer.SEMICOLON &&
		p.peekToken.Type != lexer.RPAREN &&
		p.peekToken.Type != lexer.RBRACE &&
		precedence < p.precedence(p.peekToken.Type) {

		switch p.peekToken.Type {
		case lexer.PLUS, lexer.MINUS, lexer.MULTIPLY, lexer.DIVIDE, lexer.GT, lexer.LT, lexer.EQ, lexer.GT_EQ, lexer.LT_EQ, lexer.NOT_EQ:
			p.nextToken()
			left = p.parseInfixExpression(left)
		case lexer.LSBREC:
			p.nextToken()
			left = p.parseIndexExpression(left)

		default:
			return left
		}
	}

	return left
}

func (p *Parser) parseInfixExpression(left Expression) Expression {
	expression := &InfixExpression{
		Left:     left,
		Operator: p.currentToken.Literal,
	}

	precedence := p.precedence(p.currentToken.Type)
	p.nextToken()
	expression.Right = p.parseExpressionWithPrecedence(precedence)

	return expression
}

func (p *Parser) parsePrimary() Expression {
	switch p.currentToken.Type {
	case lexer.IDENT:
		return &Identifier{Value: p.currentToken.Literal}
	case lexer.NUMBER:
		value, _ := strconv.ParseFloat(p.currentToken.Literal, 64)
		return &NumberLiteral{Value: value}
	default:
		return nil
	}
}
