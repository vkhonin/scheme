package parser_test

import (
	"github.com/vkhonin/scheme/lexer"
	"github.com/vkhonin/scheme/parser"
	"github.com/vkhonin/scheme/parser/number"
	"testing"
)

type testCase struct {
	Description string
	Input       []lexer.Token
	Output      []parser.Sexpr
}

func TestParser_Parse(t *testing.T) {
	p := parser.Parser{}

	testCases := []testCase{
		{
			Description: "Bool",
			Input: []lexer.Token{
				{Type: lexer.BOOL, Literal: "#t"},
				{Type: lexer.BOOL, Literal: "#f"},
			},
			Output: []parser.Sexpr{
				&parser.Atom{Type: parser.BOOL, Value: true},
				&parser.Atom{Type: parser.BOOL, Value: false},
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
				&parser.Atom{Type: parser.CHAR, Value: ' '},
				&parser.Atom{Type: parser.CHAR, Value: '\n'},
				&parser.Atom{Type: parser.CHAR, Value: 'a'},
			},
		},
		{
			Description: "String",
			Input: []lexer.Token{
				{Type: lexer.STRING, Literal: "string"},
			},
			Output: []parser.Sexpr{
				&parser.Atom{Type: parser.STRING, Value: "string"},
			},
		},
		{
			Description: "Symbol",
			Input: []lexer.Token{
				{Type: lexer.IDENT, Literal: "symbol"},
			},
			Output: []parser.Sexpr{
				&parser.Atom{Type: parser.SYMBOL, Value: "symbol"},
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
				&parser.Atom{Type: parser.VECTOR, Value: []parser.Sexpr{}},
				&parser.Atom{Type: parser.VECTOR, Value: []parser.Sexpr{
					&parser.Atom{Type: parser.STRING, Value: "string"},
				}},
				&parser.Atom{Type: parser.VECTOR, Value: []parser.Sexpr{
					&parser.Atom{Type: parser.STRING, Value: "string"},
					&parser.Atom{Type: parser.SYMBOL, Value: "symbol"},
				}},
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
					Car: &parser.Atom{Type: parser.SYMBOL, Value: "quote"},
					Cdr: &parser.Expr{
						Car: &parser.Atom{Type: parser.SYMBOL, Value: "symbol"},
						Cdr: &parser.Expr{Car: nil, Cdr: nil},
					},
				},
				&parser.Expr{
					Car: &parser.Atom{Type: parser.SYMBOL, Value: "quasiquote"},
					Cdr: &parser.Expr{
						Car: &parser.Atom{Type: parser.SYMBOL, Value: "symbol"},
						Cdr: &parser.Expr{Car: nil, Cdr: nil},
					},
				},
				&parser.Expr{
					Car: &parser.Atom{Type: parser.SYMBOL, Value: "unquote"},
					Cdr: &parser.Expr{
						Car: &parser.Atom{Type: parser.SYMBOL, Value: "symbol"},
						Cdr: &parser.Expr{Car: nil, Cdr: nil},
					},
				},
				&parser.Expr{
					Car: &parser.Atom{Type: parser.SYMBOL, Value: "unquote-splicing"},
					Cdr: &parser.Expr{
						Car: &parser.Atom{Type: parser.SYMBOL, Value: "symbol"},
						Cdr: &parser.Expr{Car: nil, Cdr: nil},
					},
				},
			},
		},
		{
			Description: "Number",
			Input: []lexer.Token{
				{Type: lexer.NUMBER, Literal: "#b10"},
				{Type: lexer.NUMBER, Literal: "#b#e0#/10"},
				{Type: lexer.NUMBER, Literal: "#b#i+10/1#"},
				{Type: lexer.NUMBER, Literal: "#e#o-70/1+i"},
				{Type: lexer.NUMBER, Literal: "#i#x-fi"},
				{Type: lexer.NUMBER, Literal: "1#e-1"},
				{Type: lexer.NUMBER, Literal: "2s+2"},
				{Type: lexer.NUMBER, Literal: ".3#f+33"},
				{Type: lexer.NUMBER, Literal: "4.4#d+4"},
				{Type: lexer.NUMBER, Literal: "55#.l-5"},
			},
			Output: []parser.Sexpr{
				&parser.Atom{Type: parser.NUMBER, Value: number.NewFromValue(complex(2, 0), false)},
				&parser.Atom{Type: parser.NUMBER, Value: number.NewFromValue(complex(0, 0), true)},
				&parser.Atom{Type: parser.NUMBER, Value: number.NewFromValue(complex(1, 0), true)},
				&parser.Atom{Type: parser.NUMBER, Value: number.NewFromValue(complex(-56, 1), false)},
				&parser.Atom{Type: parser.NUMBER, Value: number.NewFromValue(complex(0, -15), true)},
				&parser.Atom{Type: parser.NUMBER, Value: number.NewFromValue(complex(1, 0), true)},
				&parser.Atom{Type: parser.NUMBER, Value: number.NewFromValue(complex(200, 0), false)},
				&parser.Atom{Type: parser.NUMBER, Value: number.NewFromValue(complex(3e32, 0), true)},
				&parser.Atom{Type: parser.NUMBER, Value: number.NewFromValue(complex(44000, 0), true)},
				&parser.Atom{Type: parser.NUMBER, Value: number.NewFromValue(complex(0.0055, 0), true)},
			},
		},
		{
			Description: "List",
			Input: []lexer.Token{
				{Type: lexer.LPAREN, Literal: "("},
				{Type: lexer.RPAREN, Literal: ")"},
				{Type: lexer.LPAREN, Literal: "("},
				{Type: lexer.STRING, Literal: "string"},
				{Type: lexer.RPAREN, Literal: ")"},
				{Type: lexer.LPAREN, Literal: "("},
				{Type: lexer.STRING, Literal: "string"},
				{Type: lexer.STRING, Literal: "string"},
				{Type: lexer.RPAREN, Literal: ")"},
				{Type: lexer.LPAREN, Literal: "("},
				{Type: lexer.STRING, Literal: "string"},
				{Type: lexer.DOT, Literal: "."},
				{Type: lexer.STRING, Literal: "string"},
				{Type: lexer.RPAREN, Literal: ")"},
			},
			Output: []parser.Sexpr{
				&parser.Expr{Car: nil, Cdr: nil},
				&parser.Expr{
					Car: &parser.Atom{Type: parser.STRING, Value: "string"},
					Cdr: &parser.Expr{Car: nil, Cdr: nil},
				},
				&parser.Expr{
					Car: &parser.Atom{Type: parser.STRING, Value: "string"},
					Cdr: &parser.Expr{
						Car: &parser.Atom{Type: parser.STRING, Value: "string"},
						Cdr: &parser.Expr{Car: nil, Cdr: nil},
					},
				},
				&parser.Expr{
					Car: &parser.Atom{Type: parser.STRING, Value: "string"},
					Cdr: &parser.Atom{Type: parser.STRING, Value: "string"},
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
				t.Errorf("expected %v got %v", c.Output[i], result[i])
			}
		}
	}
}
