package css

import (
	//"log"
	"strings"
	. "github.com/badgerodon/go/dom"
)

type Filter func(Node)bool

func allNodes(root Node) []Node {
	ns := []Node{root}
	for _, c := range root.Children() {
		ns = append(ns, allNodes(c)...)
	}
	return ns
}

func filterAdjacentSiblings(left, right Filter) Filter {
	return func(node Node) bool {
		if !right(node) {
			return false
		}
		p := node.Parent()
		if p == nil {
			return false
		}
		var prev Node
		for _, next := range p.Children() {
			if next == node {
				if prev == nil {
					return false
				} else {
					return left(prev)
				}
			}
			prev = next
		}
		return false
	}
}
func filterAttribute(name string) Filter {
	return func(node Node) bool {
		_, ok := node.Get(name)
		return ok
	}
}
func filterAttributeContains(name, value string) Filter {
	return func(node Node) bool {
		v, ok := node.Get(name)
		if !ok {
			return false
		}
		return strings.Contains(v, value)
	}
}
func filterAttributeEquals(name, value string) Filter {
	return func(node Node) bool {
		v, ok := node.Get(name)
		if !ok {
			return false
		}
		return v == value
	}
}
func filterAttributePrefix(name, value string) Filter {
	return func(node Node) bool {
		v, ok := node.Get(name)
		if !ok {
			return false
		}
		return strings.HasPrefix(value, v)
	}
}
func filterAttributePrefixDash(name, value string) Filter {
	return func(node Node) bool {
		v, ok := node.Get(name)
		if !ok {
			return false
		}
		return strings.HasPrefix(value, v + "-")
	}
}
func filterAttributeSpaceContains(name, value string) Filter {
	return func(node Node) bool {
		v, ok := node.Get(name)
		if !ok {
			return false
		}
		vs := strings.Fields(v)
		for _, vf := range vs {
			if vf == value {
				return true
			}
		}
		return false
	}
}
func filterAttributeSuffix(name, value string) Filter {
	return func(node Node) bool {
		v, ok := node.Get(name)
		if !ok {
			return false
		}
		return strings.HasSuffix(value, v)
	}
}
func filterClass(cls string) Filter {
	return func(node Node) bool {
		v, ok := node.Get("class")
		if !ok {
			return false
		}
		return v == cls
	}
}
func filterDescendants(left, right Filter) Filter {
	return func(node Node) bool {
		if !right(node) {
			return false
		}
		p := node
		for {
			p = p.Parent()
			if p == nil {
				break
			}
			if left(p) {
				return true
			}
		}
		return false
	}
}
func filterEmpty() Filter {
	return func(node Node) bool {
		return len(node.Children()) == 0
	}
}
func filterFirstChild() Filter {
	return filterNthChild(1)
}
func filterFirstOfType() Filter {
	return filterNthOfType(1)
}
func filterId(id string) Filter {
	return func(node Node) bool {
		v, ok := node.Get("id")
		return ok && v == id
	}
}
func filterIntersection(filters ... Filter) Filter {
	return func(node Node) bool {
		all := true
		for _, f := range filters {
			all = all && f(node)
		}
		return all
	}
}
func filterLastChild() Filter {
	return filterNthLastChild(1)
}
func filterLastOfType() Filter {
	return filterNthLastOfType(1)
}
func filterNone() Filter {
	return func(node Node) bool {
		return false
	}
}
func filterNthChild(n int) Filter {
	return func(node Node) bool {
		p := node.Parent()
		if p == nil {
			return false
		}
		cs := p.Children()
		if (n-1) < len(cs) {
			return cs[n-1] == node
		}
		return false
	}
}
func filterNthLastChild(n int) Filter {
	return func(node Node) bool {
		p := node.Parent()
		if p == nil {
			return false
		}
		cs := p.Children()
		n = len(cs) - n
		if (n-1) < len(cs) {
			return cs[n-1] == node
		}
		return false
	}
}
func filterNthLastOfType(n int) Filter {
	return func(node Node) bool {
		return false
	}
}
func filterNthOfType(n int) Filter {
	return func(node Node) bool {
		return false
	}
}
func filterOnlyChild() Filter {
	return filterIntersection(filterFirstChild(), filterLastChild())
}
func filterTag(tag string) Filter {
	return func(node Node) bool {
		if el, ok := node.(*Element); ok {
			if el.Tag == tag {
				return true
			}
		}
		return false
	}
}
func filterRoot() Filter {
	return func(node Node) bool {
		return node.Parent() == nil
	}
}
func filterUnion(filters ... Filter) Filter {
	return func(node Node) bool {
		for _, f := range filters {
			if f(node) {
				return true
			}
		}
		return false
	}
}
func filterUniversal() Filter {
	return func(node Node) bool {
		return true
	}
}
func apply(filter Filter, nodes []Node) []Node {
	ns := make([]Node, 0)
	for _, n := range nodes {
		if filter(n) {
			ns = append(ns, n)
		}
	}
	return ns
}

func Find(root Node, exp string) []Node {
	f, ok := Parse(exp)
	if !ok {
		return []Node{}
	}
	return apply(f, allNodes(root))
}
func First(root Node, exp string) Node {
	ns := Find(root, exp)
	if len(ns) > 0 {
		return ns[0]
	}
	return nil
}
func NextTill(root Node, exp string) []Node {
	nodes := []Node{}

	f, ok := Parse(exp)
	if !ok {
		return nodes
	}

	p := root.Parent()
	if p != nil {
		seen := false
		for _, n := range p.Children() {
			if n == root {
				seen = true
				continue
			}
			if seen {
				if f(n) {
					break
				} else {
					nodes = append(nodes, n)
				}
			}
		}
	}

	return nodes
}
func PrevAll(root Node) []Node {
	nodes := []Node{}
	p := root.Parent()
	if p != nil {
		for _, n := range p.Children() {
			if n == root {
				break
			}
			nodes = append(nodes, n)
		}
	}
	return nodes
}
