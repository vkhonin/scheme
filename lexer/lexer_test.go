package lexer_test

import (
	"errors"
	"github.com/vkhonin/scheme/lexer"
	"reflect"
	"strings"
	"testing"
)

type testCase struct {
	Description string
	Input       string
	Output      []lexer.Token
}

func TestLexer_NextToken(t *testing.T) {
	l := lexer.Lexer{}

	testCases := []testCase{
		{
			Description: "Identifiers",
			Input:       "+ - ... !$%&*/:<=>?^_~1qQ+-.@",
			Output: []lexer.Token{
				{Type: lexer.IDENT, Literal: "+"},
				{Type: lexer.IDENT, Literal: "-"},
				{Type: lexer.IDENT, Literal: "..."},
				{Type: lexer.IDENT, Literal: "!$%&*/:<=>?^_~1qQ+-.@"},
			},
		},
		{
			Description: "Booleans",
			Input:       "#t#f",
			Output: []lexer.Token{
				{Type: lexer.BOOL, Literal: "#t"},
				{Type: lexer.BOOL, Literal: "#f"},
			},
		},
		{
			Description: "Numbers",
			Input:       "#b10 #b#e0#/10 #b#i+10/1# #e#o-70/1+i #i#x-fi 1#e-1 2s+2 .3#f+33 4.4#d+4 55#.l-5",
			Output: []lexer.Token{
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
		},
		{
			Description: "Characters",
			Input:       "#\\a #\\space #\\newline",
			Output: []lexer.Token{
				{Type: lexer.CHAR, Literal: "#\\a"},
				{Type: lexer.CHAR, Literal: "#\\space"},
				{Type: lexer.CHAR, Literal: "#\\newline"},
			},
		},
		{
			Description: "Strings",
			Input:       "\"\" \"a\" \"\n\"",
			Output: []lexer.Token{
				{Type: lexer.STRING, Literal: ""},
				{Type: lexer.STRING, Literal: "a"},
				{Type: lexer.STRING, Literal: "\n"},
			},
		},
		{
			Description: "Special tokens",
			Input:       "()#('`,,@. ",
			Output: []lexer.Token{
				{Type: lexer.LPAREN, Literal: "("},
				{Type: lexer.RPAREN, Literal: ")"},
				{Type: lexer.HPAREN, Literal: "#("},
				{Type: lexer.SQUOTE, Literal: "'"},
				{Type: lexer.BQUOTE, Literal: "`"},
				{Type: lexer.COMMA, Literal: ","},
				{Type: lexer.COMMAT, Literal: ",@"},
				{Type: lexer.DOT, Literal: "."},
			},
		},
	}

	for _, c := range testCases {
		l.Scanner.Init(strings.NewReader(c.Input))

		tokens := make([]lexer.Token, 0, len(c.Output))

		for token, err := l.NextToken(); ; token, err = l.NextToken() {
			if err != nil {
				if errors.Is(err, lexer.EOF) {
					break
				}

				t.Error(err)

				continue
			}

			tokens = append(tokens, token)
		}

		if !reflect.DeepEqual(c.Output, tokens) {
			t.Errorf("expected %v got %v", c.Output, tokens)
		}
	}
}
