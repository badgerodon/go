package css

import (
	//"fmt"
	"testing"

	. "github.com/badgerodon/go/dom"
	. "github.com/badgerodon/go/dom/dsl"
)

type (
	Case struct {
		Selector string
		Tree Node
		Result []Node
	}
)

var cases []Case

func init() {
	h1 := E("h1", "Hello World")
	p1 := E("p", "paragraph")
	d1 := E("div", "another div")
	d2 := E("div", A("id", "main"), A("class", "container"), h1, p1, d1)

	t1 := E("html",
		E("head"),
		E("body",
			d2,
		),
	)

	cases = []Case{
		{"",t1,[]Node{}},
		{" ",t1,[]Node{}},
		{"div",t1,[]Node{d2,d1}},
		{"body div",t1,[]Node{d2,d1}},
		{"div div",t1,[]Node{d1}},
		{"div.container",t1,[]Node{d2}},
		{".container",t1,[]Node{d2}},
		{"div#main",t1,[]Node{d2}},
		{"#main",t1,[]Node{d2}},
	}
}

func Test(t *testing.T) {
	for _, c := range cases {
		ns := Find(c.Tree, c.Selector)
		if len(ns) != len(c.Result) {
			t.Error("Expected", len(c.Result), "Got", len(ns))
		} else {
			for i, n := range ns {
				if n != c.Result[i] {
					t.Error("Expected", c.Result[i], "Got", n)
				}
			}
		}
	}
}
