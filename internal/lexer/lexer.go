package lexer

import (
	"strings"
	"text/scanner"
)

// Token types
type TokenType int

const (
	ILLEGAL TokenType = iota
	EOF
	IDENT
	NUMBER
	VAR
	IF
	ELSE
	FOR
	ARRAY
	WHILE
	PRINT

	// Operators
	PLUS
	MINUS
	MULTIPLY
	DIVIDE
	ASSIGN
	DECREMENT
	INCREMENT

	// Delimiters
	LPAREN
	RPAREN
	LBRACE
	RBRACE
	LSBREC
	RSBREC
	SEMICOLON
	COMMA

	// Comparisons
	GT
	LT
	GT_EQ
	LT_EQ
	EQ
	NOT_EQ
	IN
)

type Token struct {
	Type    TokenType
	Literal string
}

// Lexer for tokenizing
type Lexer struct {
	scanner scanner.Scanner
	token   rune
}

func NewLexer(input string) *Lexer {
	var s scanner.Scanner
	s.Init(strings.NewReader(input))
	s.Mode = scanner.ScanIdents | scanner.ScanInts | scanner.ScanFloats
	return &Lexer{scanner: s, token: s.Scan()}
}

func (l *Lexer) NextToken() Token {
	var tok Token

	switch l.token {
	case scanner.EOF:
		tok = Token{EOF, ""}
	case scanner.Ident:
		literal := l.scanner.TokenText()
		switch literal {
		case "lazy":
			tok = Token{VAR, literal}
		case "lazyArray":
			tok = Token{ARRAY, literal}
		case "if":
			tok = Token{IF, literal}
		case "el":
			tok = Token{ELSE, literal}
		case "lazyPrint":
			tok = Token{PRINT, literal}
		case "for":
			tok = Token{FOR, literal}
		case "while":
			tok = Token{WHILE, literal}
		case "in":
			tok = Token{IN, literal}
		default:
			tok = Token{IDENT, literal}
		}
	case scanner.Int, scanner.Float:
		tok = Token{NUMBER, l.scanner.TokenText()}
	case '+':
		tok = Token{PLUS, "+"}
	case '-':
		tok = Token{MINUS, "-"}
	case '*':
		tok = Token{MULTIPLY, "*"}
	case '/':
		tok = Token{DIVIDE, "/"}
	case '=':
		next := l.scanner.Peek()
		if next == '=' {
			l.scanner.Next()      // consume the second '='
			tok = Token{EQ, "=="} // now we have a '==' token
		} else {
			tok = Token{ASSIGN, "="}
		}
	case '(':
		tok = Token{LPAREN, "("}
	case ')':
		tok = Token{RPAREN, ")"}
	case '{':
		tok = Token{LBRACE, "{"}
	case '[':
		tok = Token{LSBREC, "["}
	case ']':
		tok = Token{RSBREC, "]"}
	case ';':
		tok = Token{SEMICOLON, ";"}
	case ',':
		tok = Token{COMMA, ","}
	case '}':
		tok = Token{RBRACE, "}"}
	case '>':
		if l.scanner.Peek() == '=' {
			l.scanner.Next() // consume the '='
			tok = Token{GT_EQ, ">="}
		} else {
			tok = Token{GT, ">"}
		}
	case '<':
		if l.scanner.Peek() == '=' {
			l.scanner.Next() // consume the '='
			tok = Token{LT_EQ, "<="}
		} else {
			tok = Token{LT, "<"}
		}
	default:
		tok = Token{ILLEGAL, l.scanner.TokenText()}
	}

	l.token = l.scanner.Scan()
	return tok
}
