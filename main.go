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
	root := parse(html)
	stylesheet := parsecss(css)
	fmt.Printf("%+v\n", root)
	fmt.Printf("%+v\n", stylesheet)
}
