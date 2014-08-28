package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

type Stylesheet struct {
	rules []Rule
}

type Unit int

const (
	Px Unit = iota
)

type Rule struct {
	selectors    []Selector
	declarations []Declaration
}

type Selector struct {
	tag_name string
	id       string
	class    []string
}

type BySpecificity []Selector

func (s BySpecificity) Len() int      { return len(s) }
func (s BySpecificity) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// sort highest first
func (s BySpecificity) Less(i, j int) bool {
	return s[i].Specificity() > s[j].Specificity()
}

func (s Selector) Specificity() uint64 {
	var result uint64
	result ^= uint64(len(s.tag_name))
	result ^= uint64(len(s.class) << 8)
	result ^= uint64(len(s.id) << 16)
	return result
}

func (p *Parser) parse_simple_selector() Selector {
	result := Selector{class: make([]string, 0)}

LOOP:
	for !p.eof() {
		//r := p.next_rune()
		switch r := p.next_rune(); {
		case r == rune('#'):
			p.consume_rune()
			result.id = p.parse_identifier()
		case r == rune('.'):
			p.consume_rune()
			result.class = append(result.class, p.parse_identifier())
		case r == rune('*'):
			p.consume_rune()
		case valid_identifier_rune(r):
			result.tag_name = p.parse_identifier()
		default:
			break LOOP
		}
	}
	return result
}

type Declaration struct {
	name  string
	value Value
}

// Value types
type Value interface{}
type KeywordValue string
type LengthValue struct {
	length float32
	unit   Unit
}
type ColorValue [4]uint8

func (v LengthValue) to_px() float32 {
	if v.unit == Px {
		return v.length
	} else {
		return 0.0
	}
}

// adding this since I can't make an interface method
func (v KeywordValue) to_px() float32 {
	return 0.0
}

func (p *Parser) parse_rule() Rule {
	return Rule{selectors: p.parse_selectors(),
		declarations: p.parse_declarations()}
}

func (p *Parser) parse_identifier() string {
	return string(p.consume_while(valid_identifier_rune))
}

func (p *Parser) parse_color() Value {
	if p.consume_rune() != rune('#') {
		panic("parse_color expected #, but got something else")
	}
	return Value(ColorValue{p.parse_hex_pair(), p.parse_hex_pair(), p.parse_hex_pair(), 255})
}

func (p *Parser) parse_unit() {
	s := strings.ToLower(p.parse_identifier())
	switch s {
	case "px":
		return
	default:
		panic(fmt.Sprintf("parse_unit expected 'px' but got %s", s))
	}
}

func (p *Parser) parse_value() Value {
	//r := p.next_rune()
	var v Value
	switch r := p.next_rune(); {
	case unicode.IsNumber(r):
		v = p.parse_length()
	case r == rune('#'):
		v = p.parse_color()
	default:
		v = p.parse_identifier()
	}
	return v
}

func (p *Parser) parse_length() Value {
	l := p.parse_float()
	p.parse_unit()
	return Value(l)
}

func (p *Parser) parse_float() float32 {
	s := string(p.consume_while(func(r rune) bool {
		if unicode.IsNumber(r) || r == rune('.') {
			return true
		} else {
			return false
		}
	}))
	f, _ := strconv.ParseFloat(s, 32)
	return float32(f)
}

func (p *Parser) parse_hex_pair() uint8 {
	s := string(p.input[p.pos : p.pos+2])
	p.pos = p.pos + 2
	i, _ := strconv.ParseUint(s, 0x10, 8)
	ui := uint8(i)
	return ui
}

func (p *Parser) parse_selectors() []Selector {
	selectors := make([]Selector, 0)

LOOP:
	for {
		selectors = append(selectors, p.parse_simple_selector())
		p.consume_whitespace()
		r := p.next_rune()
		switch r {
		case rune(','):
			p.consume_rune()
		case rune('{'):
			break LOOP
		default:
			panic(fmt.Sprintf("Unexpected character %s in selector list", string(r)))
		}
	}
	sort.Sort(BySpecificity(selectors))
	return selectors
}

func (p *Parser) parse_declarations() []Declaration {
	if p.consume_rune() != rune('{') {
		panic("parse_declarations expected { but got something else")
	}
	declarations := make([]Declaration, 0)
	for {
		p.consume_whitespace()
		if p.next_rune() == rune('}') {
			p.consume_rune()
			break
		}
		declarations = append(declarations, p.parse_declaration())
	}
	return declarations
}

func (p *Parser) parse_declaration() Declaration {
	property_name := p.parse_identifier()
	p.consume_whitespace()
	if p.consume_rune() != rune(':') {
		panic("parse_declaration expected : got something else")
	}
	p.consume_whitespace()
	value := p.parse_value()
	p.consume_whitespace()
	if p.consume_rune() != rune(';') {
		panic("parse declaration expected ; at end of value, but got something else")
	}
	return Declaration{name: property_name, value: value}
}

func (p *Parser) parse_rules() []Rule {
	rules := make([]Rule, 0)
	for {
		p.consume_whitespace()
		if p.eof() {
			break
		}
		rules = append(rules, p.parse_rule())
	}
	return rules
}

func parse_css(source string) Stylesheet {
	np := Parser{pos: 0, input: []rune(source)}
	return Stylesheet{rules: np.parse_rules()}
}

func valid_identifier_rune(r rune) bool {
	if unicode.IsLetter(r) || unicode.IsNumber(r) || r == rune('-') || r == rune('_') {
		return true
	} else {
		return false
	}
}
