package main

import (
	"fmt"
	"sort"
)

type PropertyMap map[string]Value

type Display int

const (
	Inline Display = iota
	Block
	None
)

type StyledNode struct {
	node             Node
	specified_values PropertyMap
	children         []StyledNode
}

// not sure I need the bool here
func (s StyledNode) value(name string) (Value, bool) {
	// Rust returns an explicitly copied value here
	// I am assuming this isn't needed here
	v, ok := s.specified_values[name]
	return v, ok
}

func (s StyledNode) lookup(name, fallback_name string, default_val Value) Value {
	v, ok := s.value(name)
	if !ok {
		v, ok = s.value(fallback_name)
		if !ok {
			v = default_val
		}
	}
	return v
}

func (s StyledNode) display() Display {
	// not sure I need the bool here
	v, _ := s.value("display")
	switch v {
	case KeywordValue("inline"):
		return Inline
	case KeywordValue("none"):
		return None
	default:
		return Block
	}
}

func style_tree(root Node, stylesheet Stylesheet) StyledNode {
	sn := StyledNode{node: root, children: make([]StyledNode, 0)}
	var sv PropertyMap
	// if node is an element, return the specified_values of that element's
	// data against provided stylesheet
	// else make specified_values an empty PropertyMap
	// does it need to be empty or can it be nil?
	if root.node_type == Element {
		sv = specified_values(root.data, stylesheet)
		// something about Nodes
	} else {
		sv = make(PropertyMap)
	}
	sn.specified_values = sv
	for _, child := range root.children {
		sn.children = append(sn.children, style_tree(child, stylesheet))
	}
	return sn
}

type MRBySpecificity []MatchedRule

func (s MRBySpecificity) Len() int      { return len(s) }
func (s MRBySpecificity) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// sort highest first
func (s MRBySpecificity) Less(i, j int) bool {
	return s[i].selector.Specificity() > s[j].selector.Specificity()
}

func specified_values(elem NodeData, stylesheet Stylesheet) PropertyMap {
	// use of rules here is confusing since it's matchedrule which has
	// selectors and rules within
	values := make(PropertyMap)
	rules := matching_rules(elem, stylesheet)
	sort.Sort(MRBySpecificity(rules))
	for _, rule := range rules {
		for _, declaration := range rule.rule.declarations {
			values[declaration.name] = declaration.value
		}
	}
	return values
}

type MatchedRule struct {
	selector Selector
	rule     Rule
}

func matching_rules(elem NodeData, stylesheet Stylesheet) []MatchedRule {
	// todo
	matchedrules := make([]MatchedRule, 0)
	for _, rule := range stylesheet.rules {
		for _, selector := range rule.selectors {
			// rust uses .find – return first element satisfying
			if matches(elem, selector) {
				matchedrules = append(matchedrules, MatchedRule{selector, rule})
				// make sure this breaks out of the loop correctly
				break
			}
		}
	}
	return matchedrules
}

func matches(elem NodeData, selector Selector) bool {
	// todo
	// god this is convoluted
	// i assume this exists to extend beyond SimpleSelector?
	return matches_simple_selector(elem, selector)
}

// I use Selector and not a sum type so just use Selector
func matches_simple_selector(elem NodeData, selector Selector) bool {
	// if selector.tag_name.iter().any(|name| elem.tag_name != *name) {
	// return false }
	// what does .any do in this context in rust?
	// looks like it does what I'd expect, so this means if any tag name in this array
	// isn't equal to elem.tag_name return false
	// i was confusing myself because I thought this meant .. nevermind
	// selector isn't an array, does rust just use .iter() to get .any for free?

	if selector.tag_name != elem.text {
		return false
	}
	if selector.id != elem.id() {
		return false
	}
	// find more efficient way to do this
	for _, selector_class := range selector.class {
		for _, elem_class := range elem.classes() {
			if selector_class != elem_class {
				return false
			}
		}
	}

	return true
}
