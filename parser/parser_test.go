package parser_test

import (
	"github.com/vkhonin/scheme/lexer"
	"github.com/vkhonin/scheme/parser"
	"github.com/vkhonin/scheme/parser/number"
	"math"
	"testing"
)

type testCase struct {
	Description string
	Input       []lexer.Token
	Output      []parser.Sexpr
}

type numberTestCase struct {
	Literal string
	Value   complex128
	Inexact bool
}

func TestParser_Parse(t *testing.T) {
	p := parser.Parser{}

	numberTestCases := []numberTestCase{
		{"#b0", complex(0, 0), false},
		{"#b1", complex(1, 0), false},
		{"#b10", complex(2, 0), false},
		{"#b1#", complex(2, 0), true},
		{"#b0/1", complex(0, 0), false},
		{"#b1/1", complex(1, 0), false},
		{"#b1/10", complex(0.5, 0), false},
		{"#b#i0", complex(0, 0), true},
		{"#b#i1", complex(1, 0), true},
		{"#b#e0", complex(0, 0), false},
		{"#b#e1", complex(1, 0), false},
		{"#i#b0", complex(0, 0), true},
		{"#i#b1", complex(1, 0), true},
		{"#e#b0", complex(0, 0), false},
		{"#e#b1", complex(1, 0), false},
		{"#b-0", complex(0, 0), false},
		{"#b-1", complex(-1, 0), false},
		{"#b-10", complex(-2, 0), false},
		{"#b-0/1", complex(0, 0), false},
		{"#b-1/1", complex(-1, 0), false},
		{"#b-1/10", complex(-0.5, 0), false},
		{"#b0@0", complex(0*math.Cos(0), 0*math.Sin(0)), false},
		{"#b1@1", complex(1*math.Cos(1), 1*math.Sin(1)), true},
		{"#b0+0i", complex(0, 0), false},
		{"#b0+1i", complex(0, 1), false},
		{"#b1+0i", complex(1, 0), false},
		{"#b1+1i", complex(1, 1), false},
		{"#b0-0i", complex(0, 0), false},
		{"#b0-1i", complex(0, -1), false},
		{"#b1-0i", complex(1, 0), false},
		{"#b1-1i", complex(1, -1), false},
		{"#b-0+0i", complex(0, 0), false},
		{"#b-0+1i", complex(0, 1), false},
		{"#b-1+0i", complex(-1, 0), false},
		{"#b-1+1i", complex(-1, 1), false},
		{"#b-0-0i", complex(0, 0), false},
		{"#b-0-1i", complex(0, -1), false},
		{"#b-1-0i", complex(-1, 0), false},
		{"#b-1-1i", complex(-1, -1), false},
		{"#b+i", complex(0, 1), false},
		{"#b+0i", complex(0, 0), false},
		{"#b+1i", complex(0, 1), false},
		{"#b-i", complex(0, -1), false},
		{"#b-0i", complex(0, 0), false},
		{"#b-1i", complex(0, -1), false},
		{"#o0", complex(0, 0), false},
		{"#o1", complex(1, 0), false},
		{"#o7", complex(7, 0), false},
		{"#o10", complex(8, 0), false},
		{"#o1#", complex(8, 0), true},
		{"#o0/7", complex(0, 0), false},
		{"#o1/7", complex(1.0/7, 0), false},
		{"#o10/1", complex(8, 0), false},
		{"#o#i0", complex(0, 0), true},
		{"#o#i7", complex(7, 0), true},
		{"#o#e0", complex(0, 0), false},
		{"#o#e10", complex(8, 0), false},
		{"#i#o0", complex(0, 0), true},
		{"#i#o7", complex(7, 0), true},
		{"#e#o0", complex(0, 0), false},
		{"#e#o10", complex(8, 0), false},
		{"#o-0", complex(0, 0), false},
		{"#o-1", complex(-1, 0), false},
		{"#o-7", complex(-7, 0), false},
		{"#o-10", complex(-8, 0), false},
		{"#o-0/7", complex(0, 0), false},
		{"#o-1/7", complex(-1.0/7, 0), false},
		{"#o-10/1", complex(-8, 0), false},
		{"#o0@0", complex(0*math.Cos(0), 0*math.Sin(0)), false},
		{"#o7@7", complex(7*math.Cos(7), 7*math.Sin(7)), true},
		{"#o0+0i", complex(0, 0), false},
		{"#o0+7i", complex(0, 7), false},
		{"#o7+0i", complex(7, 0), false},
		{"#o7+7i", complex(7, 7), false},
		{"#o0-0i", complex(0, 0), false},
		{"#o0-7i", complex(0, -7), false},
		{"#o7-0i", complex(7, 0), false},
		{"#o7-7i", complex(7, -7), false},
		{"#o-0+0i", complex(0, 0), false},
		{"#o-0+7i", complex(0, 7), false},
		{"#o-7+0i", complex(-7, 0), false},
		{"#o-7+7i", complex(-7, 7), false},
		{"#o-0-0i", complex(0, 0), false},
		{"#o-0-7i", complex(0, -7), false},
		{"#o-7-0i", complex(-7, 0), false},
		{"#o-7-7i", complex(-7, -7), false},
		{"#o+i", complex(0, 1), false},
		{"#o+0i", complex(0, 0), false},
		{"#o+7i", complex(0, 7), false},
		{"#o-i", complex(0, -1), false},
		{"#o-0i", complex(0, 0), false},
		{"#o-7i", complex(0, -7), false},
		{"0", complex(0, 0), false},
		{"1", complex(1, 0), false},
		{"12", complex(12, 0), false},
		{"1#", complex(10, 0), true},
		{"0/1", complex(0, 0), false},
		{"1/2", complex(0.5, 0), false},
		{"3/4", complex(0.75, 0), false},
		{"#i0", complex(0, 0), true},
		{"#i1", complex(1, 0), true},
		{"#e0", complex(0, 0), false},
		{"#e12", complex(12, 0), false},
		{"#d0", complex(0, 0), false},
		{"#d1", complex(1, 0), false},
		{"#d#i0", complex(0, 0), true},
		{"#d#i1", complex(1, 0), true},
		{"#d#e0", complex(0, 0), false},
		{"#d#e1", complex(1, 0), false},
		{"#i#d0", complex(0, 0), true},
		{"#i#d1", complex(1, 0), true},
		{"#e#d0", complex(0, 0), false},
		{"#e#d1", complex(1, 0), false},
		{"-0", complex(0, 0), false},
		{"-1", complex(-1, 0), false},
		{"-12", complex(-12, 0), false},
		{"-0/1", complex(0, 0), false},
		{"-1/2", complex(-0.5, 0), false},
		{"-3/4", complex(-0.75, 0), false},
		{"0.0", complex(0.0, 0), true},
		{"1.2", complex(1.2, 0), true},
		{".1", complex(0.1, 0), true},
		{"1.", complex(1.0, 0), true},
		{"0e0", complex(0.0, 0), true},
		{"1e1", complex(10.0, 0), true},
		{"1e+1", complex(10.0, 0), true},
		{"1e-1", complex(0.1, 0), true},
		{"1s1", complex(10.0, 0), true},
		{"1f1", complex(10.0, 0), true},
		{"1d1", complex(10.0, 0), true},
		{"1l1", complex(10.0, 0), true},
		{".1e1", complex(1.0, 0), true},
		{"1.2e1", complex(12.0, 0), true},
		{"1##.", complex(100.0, 0), true},
		{"1##.e1", complex(1000.0, 0), true},
		{"1##.e+1", complex(1000.0, 0), true},
		{"1##.e-1", complex(10.0, 0), true},
		{"1##.s1", complex(1000.0, 0), true},
		{"1#.#", complex(10.0, 0), true},
		{"1##.##", complex(100.0, 0), true},
		{"0@0", complex(0*math.Cos(0), 0*math.Sin(0)), false},
		{"1@1", complex(1*math.Cos(1), 1*math.Sin(1)), true},
		{"1@-1", complex(1*math.Cos(-1), 1*math.Sin(-1)), true},
		{"-1@1", complex(-1*math.Cos(1), -1*math.Sin(1)), true},
		{"-1@-1", complex(-1*math.Cos(-1), -1*math.Sin(-1)), true},
		{"1.2@3.4", complex(1.2*math.Cos(3.4), 1.2*math.Sin(3.4)), true},
		{"-1.2@3.4", complex(-1.2*math.Cos(3.4), -1.2*math.Sin(3.4)), true},
		{"1.2@-3.4", complex(1.2*math.Cos(-3.4), 1.2*math.Sin(-3.4)), true},
		{"-1.2@-3.4", complex(-1.2*math.Cos(-3.4), -1.2*math.Sin(-3.4)), true},
		{"0+0i", complex(0, 0), false},
		{"1+2i", complex(1, 2), false},
		{"1-2i", complex(1, -2), false},
		{"-1+2i", complex(-1, 2), false},
		{"-1-2i", complex(-1, -2), false},
		{"1.2+3.4i", complex(1.2, 3.4), true},
		{"1.2-3.4i", complex(1.2, -3.4), true},
		{"+i", complex(0, 1), false},
		{"+1i", complex(0, 1), false},
		{"+1.2i", complex(0, 1.2), true},
		{"+.1i", complex(0, 0.1), true},
		{"+1e1i", complex(0, 10.0), true},
		{"+1##.e1i", complex(0, 1000.0), true},
		{"-i", complex(0, -1), false},
		{"-1i", complex(0, -1), false},
		{"-1.2i", complex(0, -1.2), true},
		{"-.1i", complex(0, -0.1), true},
		{"-1e1i", complex(0, -10.0), true},
		{"-1##.e1i", complex(0, -1000.0), true},
		{"#x0", complex(0, 0), false},
		{"#x1", complex(1, 0), false},
		{"#x9", complex(9, 0), false},
		{"#xa", complex(10, 0), false},
		{"#xf", complex(15, 0), false},
		{"#x10", complex(16, 0), false},
		{"#x1#", complex(16, 0), true},
		{"#x1a", complex(26, 0), false},
		{"#x0/1", complex(0, 0), false},
		{"#x1/f", complex(1.0/15, 0), false},
		{"#xa/f", complex(10.0/15, 0), false},
		{"#x#i0", complex(0, 0), true},
		{"#x#i1", complex(1, 0), true},
		{"#x#i9", complex(9, 0), true},
		{"#x#ia", complex(10, 0), true},
		{"#x#e0", complex(0, 0), false},
		{"#x#ef", complex(15, 0), false},
		{"#i#x0", complex(0, 0), true},
		{"#i#x1", complex(1, 0), true},
		{"#i#xa", complex(10, 0), true},
		{"#e#x0", complex(0, 0), false},
		{"#e#xf", complex(15, 0), false},
		{"#x-0", complex(0, 0), false},
		{"#x-1", complex(-1, 0), false},
		{"#x-9", complex(-9, 0), false},
		{"#x-a", complex(-10, 0), false},
		{"#x-f", complex(-15, 0), false},
		{"#x-10", complex(-16, 0), false},
		{"#x-0/1", complex(0, 0), false},
		{"#x-1/f", complex(-1.0/15, 0), false},
		{"#x-a/f", complex(-10.0/15, 0), false},
		{"#x0@0", complex(0*math.Cos(0), 0*math.Sin(0)), false},
		{"#x1@1", complex(1*math.Cos(1), 1*math.Sin(1)), true},
		{"#x0+0i", complex(0, 0), false},
		{"#x0+fi", complex(0, 15), false},
		{"#xa+0i", complex(10, 0), false},
		{"#xa+fi", complex(10, 15), false},
		{"#x0-0i", complex(0, 0), false},
		{"#x0-fi", complex(0, -15), false},
		{"#xa-0i", complex(10, 0), false},
		{"#xa-fi", complex(10, -15), false},
		{"#x-0+0i", complex(0, 0), false},
		{"#x-0+fi", complex(0, 15), false},
		{"#x-a+0i", complex(-10, 0), false},
		{"#x-a+fi", complex(-10, 15), false},
		{"#x-0-0i", complex(0, 0), false},
		{"#x-0-fi", complex(0, -15), false},
		{"#x-a-0i", complex(-10, 0), false},
		{"#x-a-fi", complex(-10, -15), false},
		{"#x+i", complex(0, 1), false},
		{"#x+0i", complex(0, 0), false},
		{"#x+ai", complex(0, 10), false},
		{"#x-i", complex(0, -1), false},
		{"#x-0i", complex(0, 0), false},
		{"#x-ai", complex(0, -10), false},
	}

	testCases := []testCase{
		{
			Description: "Number",
			Input:       make([]lexer.Token, len(numberTestCases), len(numberTestCases)),
			Output:      make([]parser.Sexpr, len(numberTestCases), len(numberTestCases)),
		},
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

	for i, c := range numberTestCases {
		testCases[0].Input[i] = lexer.Token{Type: lexer.NUMBER, Literal: c.Literal}
		testCases[0].Output[i] = &parser.Atom{Type: parser.NUMBER, Value: number.NewFromValue(c.Value, c.Inexact)}
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
				t.Errorf("expected %v got %v for %s", c.Output[i], result[i], c.Input[i].Literal)
			}
		}
	}
}
