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
