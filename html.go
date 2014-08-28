package main

import (
	"fmt"
	"strings"
	"unicode"
)

type Parser struct {
	pos   int
	input []rune
}

// Read the next rune without consuming it.
func (p Parser) next_rune() rune {
	return p.input[p.pos]
}

// Do the next runes start with the given string?
func (p Parser) starts_with(s string) bool {
	return strings.HasPrefix(string(p.input[p.pos:]), s)
}

// Return true if all input is consumed.
func (p Parser) eof() bool {
	return p.pos >= len(p.input)
}

func (p *Parser) consume_rune() rune {
	r := p.input[p.pos]
	p.pos = p.pos + 1
	return r
}

func (p *Parser) consume_while(test func(rune) bool) []rune {
	result := make([]rune, 0)
	for !p.eof() && test(p.next_rune()) {
		result = append(result, p.consume_rune())
	}
	return result
}

func (p *Parser) consume_whitespace() {
	p.consume_while(unicode.IsSpace)
}

func (p *Parser) parse_tag_name() []rune {
	return p.consume_while(func(r rune) bool {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			return true
		} else {
			return false
		}
	})
}

func (p *Parser) parse_node() (n Node) {
	switch p.next_rune() {
	case rune('<'):
		n = p.parse_element()
	default:
		n = p.parse_text()
	}
	return n
}

func (p *Parser) parse_text() Node {
	return text(string(p.consume_while(func(r rune) bool {
		return r != rune('<')
	})))
}

func (p *Parser) parse_element() Node {
	var tag_name string
	var attrs AttrMap
	if p.consume_rune() == rune('<') {
		tag_name = string(p.parse_tag_name())
		attrs = p.parse_attributes()
		if p.consume_rune() != rune('>') {
			panic("did not find rune >")
		}
	}

	children := p.parse_nodes()
	if p.consume_rune() != rune('<') {
		panic("didn't find closing tag <")
	}
	if p.consume_rune() != rune('/') {
		panic("didn't find closing tag /")
	}
	closing_tag := p.parse_tag_name()
	if string(closing_tag) != string(tag_name) {
		panic(fmt.Sprintf("closing tag name %s didn't match opening tag %s", string(closing_tag), string(tag_name)))
	}
	if p.consume_rune() != rune('>') {
		panic("didn't find closing tag >")
	}
	return elem(string(tag_name), attrs, children)
}

func (p *Parser) parse_attr() (string, string) {
	name := string(p.parse_tag_name())
	// fix debug
	r := p.consume_rune()
	if r != rune('=') {
		panic(fmt.Sprintf("parse_attr expected =, got %s after %s", string(r), name))
	}
	value := string(p.parse_attr_value())
	return name, value
}

func (p *Parser) parse_attr_value() string {
	open_quote := p.consume_rune()
	if !(open_quote == rune('"') || open_quote == rune('\'')) {
		panic(fmt.Sprintf("parse_attr_value open_quote not found: %s", string(open_quote)))
	}
	// this can miss case where there are mismatched quotes
	// class='name" werps!
	// when return r != open_quote
	value := p.consume_while(func(r rune) bool {
		return r != open_quote
		//return (r != rune('"') || r != rune('\''))
	})
	if p.consume_rune() != open_quote {
		panic("parse_attr_value no open quote found at end of attribute")
	}
	return string(value)
}

func (p *Parser) parse_attributes() AttrMap {
	attributes := make(map[string]string)
	for {
		p.consume_whitespace()
		if p.next_rune() == rune('>') {
			break
		}
		name, value := p.parse_attr()
		attributes[name] = value
	}
	return attributes
}

func (p *Parser) parse_nodes() []Node {
	nodes := make([]Node, 0)
	for {
		p.consume_whitespace()
		if p.eof() || p.starts_with("</") {
			break
		}
		nodes = append(nodes, p.parse_node())
	}
	return nodes
}

func parse_html(source string) Node {
	np := Parser{pos: 0, input: []rune(source)}
	nodes := np.parse_nodes()
	if len(nodes) == 1 {
		return nodes[0]
	} else {
		return elem("html", nil, nodes)
	}
}
