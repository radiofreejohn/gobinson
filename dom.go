package main

import "strings"

const (
	Element = iota
	Text
)

type Node struct {
	node_type  int
	depth      int
	children   []Node
	text       string
	attributes AttrMap
}

/*
func (n Node) String() string {
	types := []string{"Element", "Text"}
	s := ""
	children := ""
	if len(n.children) > 0 {
		children = "children: ["
		result := ""
		for _, c := range n.children {
			result = result + fmt.Sprintf("%s ", c)
		}
		children = children + result + "]"
	}

	for k := range n.attributes {
		s = s + fmt.Sprintf("%s=%s ", k, n.attributes[k])
	}
	return strings.TrimSpace(fmt.Sprintf("%s: %s %s%s", types[n.node_type], n.text, s, children))
}
*/

type AttrMap map[string]string

func text(data string) Node {
	return Node{node_type: Text, children: nil, text: data}
}

func elem(name string, attrs AttrMap, children []Node) Node {
	return Node{children: children, text: name, attributes: attrs}
}

func (e Node) get_attribute(key string) string {
	return e.attributes[key]
}

func (e Node) id() string {
	return e.get_attribute("id")
}

func (e Node) classes() []string {
	return strings.Fields(e.get_attribute("class"))
}
