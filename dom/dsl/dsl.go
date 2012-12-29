package dsl

import (
	"fmt"
	. "github.com/badgerodon/go/dom"
)

func add(node Node, nodes ... interface{}) {
	for _, n := range nodes {
		switch o := n.(type) {
		case *Html:
			node.Append(o)
		case *Text:
			node.Append(o)
		case *Element:
			node.Append(o)
		case *Fragment:
			node.Append(o)
		case *Attribute:
			node.Set(o.Name, o.Value)
		case string:
			node.Append(T(o))
		default:
			node.Append(T(fmt.Sprint(o)))
		}
	}
}

func E(tag string, nodes ... interface{}) *Element {
	el := NewElement(tag)
	add(el, nodes...)
	return el
}
func A(name, value string) *Attribute {
	return &Attribute{name, value}
}
func T(content string) *Text {
	return NewText(content)
}
func H(content string) *Html {
	return NewHtml(content)
}
func F(nodes ... interface{}) *Fragment {
	f := NewFragment()
	add(f, nodes...)
	return f
}
