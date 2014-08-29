package main

import (
	"fmt"
	"io/ioutil"
)

func main() {
	buf, _ := ioutil.ReadFile("test.html")
	html := string(buf)
	buf, _ = ioutil.ReadFile("test.css")
	css := string(buf)
	root := parse_html(html)
	stylesheet := parse_css(css)
	style_root := style_tree(&root, &stylesheet)

	initial_containing_block := Dimensions{x: 0.0, y: 0.0, width: 800.0, height: 600.0}

	layout_root := layout(&style_root, initial_containing_block)
	fmt.Printf("root:\n%+v\n", root)
	fmt.Printf("style sheet:\n%+v\n", stylesheet)
	fmt.Println("style root:")
	fmt.Printf("%v\n", style_root)
	fmt.Println("debug:")
	fmt.Printf("%+v\n", layout_root.dimensions)
}
