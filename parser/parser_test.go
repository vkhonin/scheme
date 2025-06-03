package parser_test

import (
	"github.com/vkhonin/scheme/lexer"
	"github.com/vkhonin/scheme/parser"
	"testing"
	"unsafe"
)

type testCase struct {
	Description string
	Input       []lexer.Token
	Output      []parser.Sexpr
}

func TestParser_Parse(t *testing.T) {
	p := parser.Parser{}

	var (
		trueValue        = true
		falseValue       = false
		spaceCharValue   = ' '
		newlineCharValue = '\n'
		aCharValue       = 'a'
		stringValue      = "string"
		symbolValue      = "symbol"

		quoteValue           = "quote"
		quasiquoteValue      = "quasiquote"
		unquoteValue         = "unquote"
		unquoteSplicingValue = "unquote-splicing"
	)

	testCases := []testCase{
		{
			Description: "Bool",
			Input: []lexer.Token{
				{Type: lexer.BOOL, Literal: "#t"},
				{Type: lexer.BOOL, Literal: "#f"},
			},
			Output: []parser.Sexpr{
				&parser.Atom{Type: parser.BOOL, Value: unsafe.Pointer(&trueValue)},
				&parser.Atom{Type: parser.BOOL, Value: unsafe.Pointer(&falseValue)},
			},
		},
		{
			Description: "Character",
			Input: []lexer.Token{
				{Type: lexer.CHAR, Literal: "#\\space"},
				{Type: lexer.CHAR, Literal: "#\\newline"},
				{Type: lexer.CHAR, Literal: "#\\a"},
			},
			Output: []parser.Sexpr{
				&parser.Atom{Type: parser.CHAR, Value: unsafe.Pointer(&spaceCharValue)},
				&parser.Atom{Type: parser.CHAR, Value: unsafe.Pointer(&newlineCharValue)},
				&parser.Atom{Type: parser.CHAR, Value: unsafe.Pointer(&aCharValue)},
			},
		},
		{
			Description: "String",
			Input: []lexer.Token{
				{Type: lexer.STRING, Literal: "string"},
			},
			Output: []parser.Sexpr{
				&parser.Atom{Type: parser.STRING, Value: unsafe.Pointer(&stringValue)},
			},
		},
		{
			Description: "Symbol",
			Input: []lexer.Token{
				{Type: lexer.IDENT, Literal: "symbol"},
			},
			Output: []parser.Sexpr{
				&parser.Atom{Type: parser.SYMBOL, Value: unsafe.Pointer(&symbolValue)},
			},
		},
		{
			Description: "Vector",
			Input: []lexer.Token{
				{Type: lexer.HPAREN, Literal: "#("},
				{Type: lexer.RPAREN, Literal: ")"},
				{Type: lexer.HPAREN, Literal: "#("},
				{Type: lexer.STRING, Literal: "string"},
				{Type: lexer.RPAREN, Literal: ")"},
				{Type: lexer.HPAREN, Literal: "#("},
				{Type: lexer.STRING, Literal: "string"},
				{Type: lexer.IDENT, Literal: "symbol"},
				{Type: lexer.RPAREN, Literal: ")"},
			},
			Output: []parser.Sexpr{
				&parser.Atom{Type: parser.VECTOR, Value: unsafe.Pointer(&[]parser.Sexpr{})},
				&parser.Atom{Type: parser.VECTOR, Value: unsafe.Pointer(&[]parser.Sexpr{
					&parser.Atom{Type: parser.STRING, Value: unsafe.Pointer(&stringValue)},
				})},
				&parser.Atom{Type: parser.VECTOR, Value: unsafe.Pointer(&[]parser.Sexpr{
					&parser.Atom{Type: parser.STRING, Value: unsafe.Pointer(&stringValue)},
					&parser.Atom{Type: parser.SYMBOL, Value: unsafe.Pointer(&symbolValue)},
				})},
			},
		},
		{
			Description: "Abbreviation",
			Input: []lexer.Token{
				{Type: lexer.SQUOTE, Literal: "'"},
				{Type: lexer.IDENT, Literal: "symbol"},
				{Type: lexer.BQUOTE, Literal: "`"},
				{Type: lexer.IDENT, Literal: "symbol"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.IDENT, Literal: "symbol"},
				{Type: lexer.COMMAT, Literal: ",@"},
				{Type: lexer.IDENT, Literal: "symbol"},
			},
			Output: []parser.Sexpr{
				&parser.Expr{
					Car: &parser.Atom{Type: parser.SYMBOL, Value: unsafe.Pointer(&quoteValue)},
					Cdr: &parser.Expr{
						Car: &parser.Atom{Type: parser.SYMBOL, Value: unsafe.Pointer(&symbolValue)},
						Cdr: &parser.Expr{Car: nil, Cdr: nil},
					},
				},
				&parser.Expr{
					Car: &parser.Atom{Type: parser.SYMBOL, Value: unsafe.Pointer(&quasiquoteValue)},
					Cdr: &parser.Expr{
						Car: &parser.Atom{Type: parser.SYMBOL, Value: unsafe.Pointer(&symbolValue)},
						Cdr: &parser.Expr{Car: nil, Cdr: nil},
					},
				},
				&parser.Expr{
					Car: &parser.Atom{Type: parser.SYMBOL, Value: unsafe.Pointer(&unquoteValue)},
					Cdr: &parser.Expr{
						Car: &parser.Atom{Type: parser.SYMBOL, Value: unsafe.Pointer(&symbolValue)},
						Cdr: &parser.Expr{Car: nil, Cdr: nil},
					},
				},
				&parser.Expr{
					Car: &parser.Atom{Type: parser.SYMBOL, Value: unsafe.Pointer(&unquoteSplicingValue)},
					Cdr: &parser.Expr{
						Car: &parser.Atom{Type: parser.SYMBOL, Value: unsafe.Pointer(&symbolValue)},
						Cdr: &parser.Expr{Car: nil, Cdr: nil},
					},
				},
			},
		},
	}

	for _, c := range testCases {
		p.Tokens = c.Input

		result := p.Parse()

		if len(result) != len(c.Output) {
			t.Errorf("expected %v got %v", c.Output, result)
			continue
		}

		for i, r := range result {
			if !r.Equals(c.Output[i]) {
				t.Errorf("expected %v got %v", c.Output, result)
			}
		}
	}
}
