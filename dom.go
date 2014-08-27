package main

import (
	"strings"
)

const (
	Element = iota
	Text
)

type NodeData struct {
	text       string
	attributes AttrMap
}

type Node struct {
	node_type int
	depth     int
	children  []Node
	data      NodeData
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
	return Node{node_type: Text, children: nil, data: NodeData{text: data}}
}

func elem(name string, attrs AttrMap, children []Node) Node {
	return Node{children: children, data: NodeData{text: name, attributes: attrs}}
}

func (n NodeData) get_attribute(key string) string {
	return n.attributes[key]
}

func (n NodeData) id() string {
	return n.get_attribute("id")
}

func (n NodeData) classes() []string {
	return strings.Fields(n.get_attribute("class"))
}
