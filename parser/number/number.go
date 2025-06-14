package number

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

// Number-related regexp as in <number> and children (7.1.1. Lexical structure).
const (
	prefix2  = `#b(#[ie])?|(#[ie])?#b`
	prefix8  = `#o(#[ie])?|(#[ie])?#o`
	prefix10 = `(#d)?(#[ie])?|(#[ie])?(#d)?`
	prefix16 = `#x(#[ie])?|(#[ie])?#x`

	uinteger2  = `[0-1]+#*`
	uinteger8  = `[0-7]+#*`
	uinteger10 = `[0-9]+#*`
	uinteger16 = `[0-9a-f]+#*`

	decimal10 = `(\.?` + uinteger10 + `|[0-9]+\.[0-9]*#*` + `|[0-9]+#+\.#*)` + `([esfdl][+-]?[0-9]+)?`

	ureal2  = `(?P<dividend>` + uinteger2 + `)(/(?P<divisor>` + uinteger2 + `))?`
	ureal8  = `(?P<dividend>` + uinteger8 + `)(/(?P<divisor>` + uinteger8 + `))?`
	ureal10 = `(?P<dividend>` + uinteger10 + `)(/(?P<divisor>` + uinteger10 + `))?|(?P<decimal>` + decimal10 + `)`
	ureal16 = `(?P<dividend>` + uinteger16 + `)(/(?P<divisor>` + uinteger16 + `))?`

	real2  = `(?P<realSign>[+-]?)(?P<realUreal>` + ureal2 + `)`
	real8  = `(?P<realSign>[+-]?)(?P<realUreal>` + ureal8 + `)`
	real10 = `(?P<realSign>[+-]?)(?P<realUreal>` + ureal10 + `)`
	real16 = `(?P<realSign>[+-]?)(?P<realUreal>` + ureal16 + `)`

	complex2  = `(?P<complexReal>` + real2 + `)(@(?P<complexImag>` + real2 + `)|(?P<complexImagSign>[+-])(?P<complexImag>` + ureal2 + `)?i` + `)?|(?P<complexImagSign>[+-])(?P<complexImag>` + ureal2 + `)?i`
	complex8  = `(?P<complexReal>` + real8 + `)(@(?P<complexImag>` + real8 + `)|(?P<complexImagSign>[+-])(?P<complexImag>` + ureal8 + `)?i` + `)?|(?P<complexImagSign>[+-])(?P<complexImag>` + ureal8 + `)?i`
	complex10 = `(?P<complexReal>` + real10 + `)(@(?P<complexImag>` + real10 + `)|(?P<complexImagSign>[+-])(?P<complexImag>` + ureal10 + `)?i` + `)?|(?P<complexImagSign>[+-])(?P<complexImag>` + ureal10 + `)?i`
	complex16 = `(?P<complexReal>` + real16 + `)(@(?P<complexImag>` + real16 + `)|(?P<complexImagSign>[+-])(?P<complexImag>` + ureal16 + `)?i` + `)?|(?P<complexImagSign>[+-])(?P<complexImag>` + ureal16 + `)?i`

	number2  = `(?P<prefix>` + prefix2 + `)(?P<complex>` + complex2 + `)`
	number8  = `(?P<prefix>` + prefix8 + `)(?P<complex>` + complex8 + `)`
	number10 = `(?P<prefix>` + prefix10 + `)(?P<complex>` + complex10 + `)`
	number16 = `(?P<prefix>` + prefix16 + `)(?P<complex>` + complex16 + `)`

	number = number2 + `|` + number8 + `|` + number10 + `|` + number16
)

const (
	typeNumber = iota
	typeComplex
	typeReal
	typeUreal

	baseN  = 0
	base2  = 2
	base8  = 8
	base10 = 10
	base16 = 16
)

var regexps map[int]map[int]*struct {
	Regexp *regexp.Regexp
	Groups []string
}

type Number struct {
	literal string

	complex complex128
	inexact bool

	isNumber bool
	radixVal int
}

func (n *Number) String() string {
	return fmt.Sprintf("%e (i=%t)", n.complex, n.inexact)
}

func init() {
	regexps = map[int]map[int]*struct {
		Regexp *regexp.Regexp
		Groups []string
	}{
		typeNumber: {
			baseN: {Regexp: compileRegexp(number)},
		},
		typeComplex: {
			base2:  {Regexp: compileRegexp(complex2)},
			base8:  {Regexp: compileRegexp(complex8)},
			base10: {Regexp: compileRegexp(complex10)},
			base16: {Regexp: compileRegexp(complex16)},
		},
		typeReal: {
			base2:  {Regexp: compileRegexp(real2)},
			base8:  {Regexp: compileRegexp(real8)},
			base10: {Regexp: compileRegexp(real10)},
			base16: {Regexp: compileRegexp(real16)},
		},
		typeUreal: {
			base2:  {Regexp: compileRegexp(ureal2)},
			base8:  {Regexp: compileRegexp(ureal8)},
			base10: {Regexp: compileRegexp(ureal10)},
			base16: {Regexp: compileRegexp(ureal16)},
		},
	}

	for t := range regexps {
		for b := range regexps[t] {
			regexps[t][b].Groups = regexps[t][b].Regexp.SubexpNames()
		}
	}
}

func compileRegexp(s string) *regexp.Regexp {
	return regexp.MustCompile(`^(` + s + `)$`)
}

func NewFromLiteral(literal string) *Number {
	return &Number{
		literal:  literal,
		isNumber: regexps[typeNumber][baseN].Regexp.MatchString(literal),
	}
}

func NewFromValue(value complex128, inexact bool) *Number {
	return &Number{
		complex:  value,
		inexact:  inexact,
		isNumber: true,
		radixVal: 10,
	}
}

func (n *Number) IsNumber() bool {
	return n.isNumber
}

func (n *Number) Inexact() bool {
	return n.inexact
}

func (n *Number) Value() complex128 {
	return n.complex
}

func (n *Number) Parse() *Number {
	groupVals := n.getGroupVals(n.literal, typeNumber, baseN)

	n.parsePrefix(groupVals["prefix"])
	n.parseComplex(groupVals["complex"])

	return n
}

func (n *Number) parseComplex(literal string) {
	groupVals := n.getGroupVals(literal, typeComplex, n.radixVal)

	var (
		rVal = n.parseReal(groupVals["complexReal"])
		iVal float64
	)

	if strings.ContainsRune(literal, '@') {
		iRaw := n.parseReal(groupVals["complexImag"])
		sin := math.Sin(iRaw)
		if math.Abs(sin) > 1e-52 {
			n.inexact = true
		}
		iVal = rVal * sin
		rVal = rVal * math.Cos(iRaw)
	} else if strings.ContainsRune(literal, 'i') {
		iRaw := 1.0
		if groupVals["complexImag"] != "" {
			iRaw = n.parseUreal(groupVals["complexImag"])
		}
		iVal = n.getSign(groupVals["complexImagSign"]) * iRaw
	}

	n.complex = complex(rVal, iVal)
}

func (n *Number) parseReal(literal string) float64 {
	if literal == "" {
		return 0
	}

	groupVals := n.getGroupVals(literal, typeReal, n.radixVal)

	ureal := n.parseUreal(groupVals["realUreal"])

	return n.getSign(groupVals["realSign"]) * ureal
}

func (n *Number) parseUreal(literal string) float64 {
	groupVals := n.getGroupVals(literal, typeUreal, n.radixVal)

	if groupVals["decimal"] != "" {
		return n.parseDecimal(groupVals["decimal"])
	}

	dividend := n.parseUint(groupVals["dividend"])

	divisor := 1.0
	if strings.ContainsRune(literal, '/') {
		divisor = n.parseUint(groupVals["divisor"])
	}

	return dividend / divisor
}

func (n *Number) parseDecimal(literal string) float64 {
	literal = strings.Map(func(r rune) rune {
		switch r {
		case 's', 'f', 'd', 'l':
			return 'e'
		case '#':
			n.inexact = true
			return '0'
		}
		return r
	}, literal)

	if strings.ContainsRune(literal, '.') || strings.ContainsRune(literal, 'e') {
		n.inexact = true
	}

	value, err := strconv.ParseFloat(literal, 0)
	if err != nil {
		panic(err)
	}

	return value
}

func (n *Number) parseUint(literal string) float64 {
	if strings.ContainsRune(literal, '#') {
		literal = strings.ReplaceAll(literal, "#", "0")
		n.inexact = true
	}

	value, err := strconv.ParseInt(literal, n.radixVal, 0)
	if err != nil {
		panic(err)
	}

	return float64(value)
}

func (n *Number) parsePrefix(literal string) {
	switch {
	case strings.ContainsRune(literal, 'b'):
		n.radixVal = base2
	case strings.ContainsRune(literal, 'o'):
		n.radixVal = base8
	case strings.ContainsRune(literal, 'x'):
		n.radixVal = base16
	default:
		n.radixVal = base10
	}

	if strings.ContainsRune(literal, 'i') {
		n.inexact = true
	}
}

func (n *Number) getGroupVals(l string, t, b int) map[string]string {
	regex := regexps[t][b].Regexp
	groups := regexps[t][b].Groups

	matches := regex.FindStringSubmatch(l)

	vals := make(map[string]string, len(matches))

	for i, match := range matches {
		if match != "" && groups[i] != "" {
			vals[groups[i]] = match
		}
	}

	return vals
}

func (n *Number) getSign(l string) float64 {
	if l == "-" {
		return -1
	}

	return 1
}
