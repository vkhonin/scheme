package lexer

import (
	"errors"
	"github.com/vkhonin/scheme/parser/number"
	"strings"
	"text/scanner"
)

// Type of token as in <token> (7.1.1. Lexical structure).
const (
	LPAREN TokenType = iota // Literal: (
	RPAREN                  // Literal: )
	HPAREN                  // Literal: #(
	SQUOTE                  // Literal: '
	BQUOTE                  // Literal: `
	COMMA                   // Literal: ,
	COMMAT                  // Literal: ,@
	DOT                     // Literal: .
	BOOL                    // Literal example: #t
	CHAR                    // Literal example: #\t
	IDENT                   // Literal example: t
	STRING                  // Literal example: "t"
	NUMBER                  // Literal example: 1
)

var (
	EOF            = errors.New("EOF")
	INVALID_DOT    = errors.New("invalid dot token")
	INVALID_HASH   = errors.New("invalid hash prefixed token")
	INVALID_IDENT  = errors.New("invalid identifier")
	INVALID_NUMBER = errors.New("invalid number")
	UNEXPECTED_EOF = errors.New("unexpected EOF")
	UNKNOWN_NCHAR  = errors.New("unknown character name")
)

type Lexer struct {
	Scanner scanner.Scanner
}

type Token struct {
	Type    TokenType
	Literal string
}

type TokenType uint8

func (l *Lexer) NextToken() (Token, error) {
	l.skipAtmosphere()

	switch r := l.Scanner.Next(); r {
	case scanner.EOF:
		return Token{}, EOF
	case '(':
		return Token{Type: LPAREN, Literal: "("}, nil
	case ')':
		return Token{Type: RPAREN, Literal: ")"}, nil
	case '\'':
		return Token{Type: SQUOTE, Literal: "'"}, nil
	case '`':
		return Token{Type: BQUOTE, Literal: "`"}, nil
	case ',':
		if l.Scanner.Peek() == '@' {
			l.Scanner.Next()
			return Token{Type: COMMAT, Literal: ",@"}, nil
		}
		return Token{Type: COMMA, Literal: ","}, nil
	case '.':
		if l.isDelimiter(l.Scanner.Peek()) {
			return Token{Type: DOT, Literal: "."}, nil
		} else if '0' <= l.Scanner.Peek() && l.Scanner.Peek() <= '9' {
			return l.scanNumber(r)
		} else if l.Scanner.Next() == '.' && l.Scanner.Next() == '.' {
			return Token{Type: IDENT, Literal: "..."}, nil
		}
		return Token{}, INVALID_DOT
	case '"':
		return l.scanString()
	case '#':
		switch l.Scanner.Peek() {
		case '(':
			return Token{Type: HPAREN, Literal: "#" + string(l.Scanner.Next())}, nil
		case 't', 'f':
			return Token{Type: BOOL, Literal: "#" + string(l.Scanner.Next())}, nil
		case '\\':
			l.Scanner.Next()
			char := l.Scanner.Next()
			if l.isDelimiter(l.Scanner.Peek()) {
				return Token{Type: CHAR, Literal: "#\\" + string(char)}, nil
			}
			return l.scanNchar(char)
		case 'i', 'e', 'b', 'o', 'd', 'x':
			return l.scanNumber(r)
		default:
			return Token{}, INVALID_HASH
		}
	case '+', '-':
		if l.isDelimiter(l.Scanner.Peek()) {
			return Token{Type: IDENT, Literal: string(r)}, nil
		}
		return l.scanNumber(r)
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return l.scanNumber(r)
	default:
		return l.scanIdentifier(r)
	}
}

func (l *Lexer) skipAtmosphere() {
	for l.isAtmosphere(l.Scanner.Peek()) {
		if l.isComment(l.Scanner.Peek()) {
			for !l.isNewline(l.Scanner.Peek()) {
				l.Scanner.Next()
			}
		}
		l.Scanner.Next()
	}
}

func (l *Lexer) isAtmosphere(r rune) bool {
	return l.isWhitespace(r) || l.isComment(r)
}

func (l *Lexer) isWhitespace(r rune) bool {
	return r == ' ' || l.isNewline(r)
}

func (l *Lexer) isNewline(r rune) bool {
	return r == '\n'
}

func (l *Lexer) isComment(r rune) bool {
	return r == ';'
}

func (l *Lexer) isDelimiter(r rune) bool {
	return l.isWhitespace(r) || strings.ContainsRune("();\"", r)
}

func (l *Lexer) scanNchar(prefix rune) (Token, error) {
	var sb strings.Builder

	sb.WriteRune(prefix)

	for r := l.Scanner.Peek(); !l.isDelimiter(r) && r != scanner.EOF; r = l.Scanner.Peek() {
		sb.WriteRune(l.Scanner.Next())
	}

	if sb.String() != "space" && sb.String() != "newline" {
		return Token{}, UNKNOWN_NCHAR
	}

	return Token{Type: CHAR, Literal: "#\\" + sb.String()}, nil
}

func (l *Lexer) scanNumber(prefix rune) (Token, error) {
	var sb strings.Builder

	sb.WriteRune(prefix)

	for r := l.Scanner.Peek(); !l.isDelimiter(r) && r != scanner.EOF; r = l.Scanner.Peek() {
		sb.WriteRune(l.Scanner.Next())
	}

	if !number.NewFromLiteral(sb.String()).IsNumber() {
		return Token{}, INVALID_NUMBER
	}

	return Token{Type: NUMBER, Literal: sb.String()}, nil
}

func (l *Lexer) scanString() (Token, error) {
	var sb strings.Builder

	for p, c := '"', l.Scanner.Next(); !(p != '\\' && c == '"'); p, c = c, l.Scanner.Next() {
		if c == scanner.EOF {
			return Token{}, UNEXPECTED_EOF
		}

		sb.WriteRune(c)
	}

	return Token{Type: STRING, Literal: sb.String()}, nil
}

func (l *Lexer) scanIdentifier(initial rune) (Token, error) {
	if !l.isIdentifierInitial(initial) {
		return Token{}, INVALID_IDENT
	}

	var sb strings.Builder

	sb.WriteRune(initial)

	for r := l.Scanner.Peek(); !l.isDelimiter(r) && r != scanner.EOF; r = l.Scanner.Peek() {
		if !l.isIdentifierSubsequent(r) {
			return Token{}, INVALID_IDENT
		}

		sb.WriteRune(l.Scanner.Next())
	}

	return Token{Type: IDENT, Literal: sb.String()}, nil
}

func (l *Lexer) isIdentifierInitial(r rune) bool {
	return ('a' <= r && r <= 'z') ||
		('A' <= r && r <= 'Z') ||
		strings.ContainsRune("!$%&*/:<=>?^_~", r)
}

func (l *Lexer) isIdentifierSubsequent(r rune) bool {
	return l.isIdentifierInitial(r) ||
		('0' <= r && r <= '9') ||
		strings.ContainsRune("+-.@", r)
}
