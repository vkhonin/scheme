package parser

import (
	"github.com/vkhonin/scheme/lexer"
	"github.com/vkhonin/scheme/parser/number"
)

const (
	BOOL AtomType = iota
	NUMBER
	CHAR
	STRING
	SYMBOL
	VECTOR
)

var (
	abbrevToIdent = map[string]string{
		"'":  "quote",
		"`":  "quasiquote",
		",":  "unquote",
		",@": "unquote-splicing",
	}
)

type Parser struct {
	Tokens []lexer.Token
	index  int
}

type Sexpr interface {
	Equals(s Sexpr) bool
}

type Atom struct {
	Type  AtomType
	Value interface{}
}

func (a *Atom) Equals(s Sexpr) bool {
	a2, ok := s.(*Atom)
	if !ok {
		return false
	}

	if a.Type != a2.Type {
		return false
	}

	switch a.Type {
	case BOOL:
		return (a.Value).(bool) == (a2.Value).(bool)
	case CHAR:
		return (a.Value).(rune) == (a2.Value).(rune)
	case STRING, SYMBOL:
		return (a.Value).(string) == (a2.Value).(string)
	case VECTOR:
		aVector := (a.Value).([]Sexpr)
		a2Vector := (a2.Value).([]Sexpr)
		la := len(aVector)
		la2 := len(a2Vector)
		if la != la2 {
			return false
		}
		for i := range la {
			if !aVector[i].Equals(a2Vector[i]) {
				return false
			}
		}
		return true
	case NUMBER:
		aNum := (a.Value).(*number.Number)
		a2Num := (a2.Value).(*number.Number)
		return aNum.IsNumber() && a2Num.IsNumber() && aNum.Inexact() == a2Num.Inexact() && aNum.Value() == a2Num.Value()
	default:
		panic("type comparison not implemented")
	}
}

type AtomType uint8

type Expr struct {
	Car Sexpr
	Cdr Sexpr
}

func (e *Expr) Equals(s Sexpr) bool {
	e2, ok := s.(*Expr)
	if !ok {
		return false
	}

	if e == nil || e2 == nil {
		return e == e2
	}

	var isCarsEqual, isCdrsEqual bool

	if e.Car == nil || e2.Car == nil {
		isCarsEqual = e.Car == e2.Car
	} else {
		isCarsEqual = e.Car.Equals(e2.Car)
	}

	if e.Cdr == nil || e2.Cdr == nil {
		isCdrsEqual = e.Cdr == e2.Cdr
	} else {
		isCdrsEqual = e.Cdr.Equals(e2.Cdr)
	}

	return isCarsEqual && isCdrsEqual
}

func (p *Parser) Parse() []Sexpr {
	p.index = 0

	var program []Sexpr

	for p.index < len(p.Tokens) {
		program = append(program, p.ParseNextNode())
	}

	return program
}

func (p *Parser) ParseNextNode() Sexpr {
	currentToken := &p.Tokens[p.index]
	var sexpr Sexpr

	switch currentToken.Type {
	case lexer.BOOL:
		sexpr = &Atom{Type: BOOL, Value: p.parseBool(currentToken.Literal)}
	case lexer.NUMBER:
		sexpr = &Atom{Type: NUMBER, Value: p.parseNumber(currentToken.Literal)}
	case lexer.CHAR:
		sexpr = &Atom{Type: CHAR, Value: p.parseChar(currentToken.Literal)}
	case lexer.STRING:
		sexpr = &Atom{Type: STRING, Value: currentToken.Literal}
	case lexer.IDENT:
		sexpr = &Atom{Type: SYMBOL, Value: currentToken.Literal}
	case lexer.HPAREN:
		sexpr = &Atom{Type: VECTOR, Value: p.parseVector()}
	case lexer.SQUOTE, lexer.BQUOTE, lexer.COMMA, lexer.COMMAT:
		sexpr = p.parseAbbrev()
	case lexer.LPAREN:
		sexpr = p.parseList()
	}
	p.index++

	return sexpr
}

func (*Parser) parseBool(literal string) bool {
	return literal[1] == 't'
}

func (p *Parser) parseNumber(literal string) any {
	return number.NewFromLiteral(literal).Parse()
}

func (*Parser) parseChar(literal string) rune {
	var char rune
	switch literal[2:] {
	case "space":
		char = ' '
	case "newline":
		char = '\n'
	default:
		for i, c := range literal {
			if i == 2 {
				char = c
				break
			}
		}
	}
	return char
}

func (p *Parser) parseVector() []Sexpr {
	value := make([]Sexpr, 0)

	p.index++
	node := &p.Tokens[p.index]

	for node.Type != lexer.RPAREN {
		value = append(value, p.ParseNextNode())
		node = &p.Tokens[p.index]
	}

	return value
}

func (p *Parser) parseAbbrev() *Expr {
	node := &p.Tokens[p.index]

	value := Expr{
		Car: &Atom{Type: SYMBOL, Value: abbrevToIdent[node.Literal]},
	}

	p.index++

	value.Cdr = &Expr{
		Car: p.ParseNextNode(),
		Cdr: &Expr{Car: nil, Cdr: nil},
	}

	p.index--

	return &value
}

func (p *Parser) parseList() *Expr {
	var value Expr
	var previousNode *Expr
	currentNode := &value

	p.index++
	node := &p.Tokens[p.index]

	for node.Type != lexer.RPAREN {
		if node.Type == lexer.DOT {
			p.index++
			previousNode.Cdr = p.ParseNextNode()

			node = &p.Tokens[p.index]
			if node.Type != lexer.RPAREN {
				// TODO: replace with error
				panic("list end expected")
			}
		}

		currentNode.Car = p.ParseNextNode()
		currentNode.Cdr = &Expr{}
		previousNode = currentNode
		currentNode = currentNode.Cdr.(*Expr)

		node = &p.Tokens[p.index]
	}

	return &value
}
