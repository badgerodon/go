package dom

import (
	"bytes"
	"encoding/xml"
	//"fmt"
	"io"
)

type (
	Node interface {
		Append(Node)
		Children() []Node
		Clone() Node
		Export(io.Writer)
		Get(string) (string, bool)
		Insert(int, Node)
		Parent() Node
		Remove(Node) bool
		Set(string, string)
		SetParent(Node)
	}
	Attribute struct {
		Name, Value string
	}
	Container struct {
		children []Node
	}
	Element struct {
		parent Node
		Tag string
		attributes map[string]string
		Container
	}
	Text struct {
		parent Node
		Content string
	}
	Html struct { Text }
	Fragment struct {
		parent Node
		Container
	}

	Visitor func(node Node, visit func())
)

var SelfClosingTags = map[string]bool{
	"base": true,
	"basefont": true,
	"frame": true,
	"link": true,
	"meta": true,
	"area": true,
	"br": true,
	"col": true,
	"hr": true,
	"img": true,
	"input": true,
	"param": true,
}

func Replace(n1, n2 Node) {
	p1 := n1.Parent()
	p2 := n2.Parent()
	if p2 != nil {
		p2.Remove(n2)
	}
	if p1 != nil {
		i := 0
		cs := p1.Children()
		for ; i < len(cs); i++ {
			if cs[i] == n1 {
				break
			}
		}
		p1.Insert(i, n2)
		p1.Remove(n1)
	}
}
func TextContent(root Node) string {
	if t, ok := root.(*Text); ok {
		return t.Content
	}
	str := ""
	for _, c := range root.Children() {
		str += TextContent(c)
	}
	return str
}
func Visit(root Node, handler func(Node,func())) {
	visit := func() {
		for _, c := range root.Children() {
			Visit(c, handler)
		}
	}
	handler(root, visit)
}
func FromXml(decoder *xml.Decoder) (*Fragment, error) {
	var cur Node = &Fragment{nil, Container{make([]Node, 0)}}
	stack := []Node{cur}
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		switch o := token.(type) {
		case xml.StartElement:
			n := NewElement(o.Name.Local)
			for _, a := range o.Attr {
				n.Set(a.Name.Local, a.Value)
			}
			cur.Append(n)
			stack = append(stack, n)
			cur = n
		case xml.EndElement:
			stack = stack[:len(stack)-1]
			if len(stack) > 0 {
				cur = stack[len(stack)-1]
			} else {
				cur = nil
			}
		case xml.CharData:
			if cur != nil {
				cur.Append(&Text{nil, string(o)})
			}
		case xml.Comment:
		case xml.ProcInst:
		case xml.Directive:
		}
	}
	return cur.(*Fragment), nil
}
func ParseHtml(html []byte) *Fragment {
	return (&parser{html,0,make([]int,0)}).parse()
}

// Container
func (this *Container) Append(n Node) {
	p := n.Parent()
	if p != nil {
		p.Remove(n)
	}
	this.children = append(this.children, n)
}
func (this *Container) Children() []Node {
	cs := make([]Node, len(this.children))
	copy(cs, this.children)
	return cs
}
func (this *Container) Insert(pos int, node Node) {
	if pos < len(this.children) {
		cs := make([]Node, len(this.children) + 1)
		copy(cs, this.children[:pos])
		cs[pos] = node
		copy(cs[pos+1:], this.children[pos:])
		this.children = cs
	} else {
		this.Append(node)
	}
}
func (this *Container) Remove(n Node) bool {
	for i, c := range this.children {
		if c == n {
			c.SetParent(nil)
			cs := this.children[:i]
			cs = append(cs, this.children[i+1:]...)
			this.children = cs
			return true
		}
	}
	return false
}
// Element
func NewElement(tag string) *Element {
	return &Element{nil, tag, make(map[string]string), Container{make([]Node, 0)}}
}
func (this *Element) Append(node Node) {
	if _, ok := node.(*Fragment); ok {
		for _, c := range node.Children() {
			this.Append(c)
		}
		return
	}
	this.Container.Append(node)
	node.SetParent(this)
}
func (this *Element) Clone() Node {
	n := NewElement(this.Tag)
	for _, c := range this.Children() {
		n.Append(c.Clone())
	}
	for k, v := range this.attributes {
		n.attributes[k] = v
	}
	return n
}
func (this *Element) Export(w io.Writer) {
	io.WriteString(w, "<")
	io.WriteString(w, this.Tag)
	keys := make([]string, 0, len(this.attributes))
	for k, _ := range this.attributes {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for k := range keys {
		v := this.attributes[k]
		io.WriteString(w, " ")
		xml.Escape(w, []byte(k))
		io.WriteString(w, "=\"")
		xml.Escape(w, []byte(v))
		io.WriteString(w, "\"")
	}
	if this.IsSelfClosing() {
		io.WriteString(w, " />")
	} else {
		io.WriteString(w, ">")
		for _, c := range this.children {
			c.Export(w)
		}
		io.WriteString(w, "</")
		io.WriteString(w, this.Tag)
		io.WriteString(w, ">")
	}
}
func (this *Element) Get(name string) (string, bool) {
	v, ok := this.attributes[name]
	return v, ok
}
func (this *Element) Insert(pos int, node Node) {
	this.Container.Insert(pos, node)
	node.SetParent(this)
}
func (this *Element) IsSelfClosing() bool {
	_, ok := SelfClosingTags[this.Tag]
	return ok
}
func (this *Element) Parent() Node {
	return this.parent
}
func (this *Element) Set(name, value string) {
	this.attributes[name] = value
}
func (this *Element) SetParent(parent Node) {
	this.parent = parent
}
func (this *Element) String() string {
	var buf bytes.Buffer
	this.Export(&buf)
	return buf.String()
}

// Text
func NewText(content string) *Text {
	return &Text{nil, content}
}
func (this *Text) Append(n Node) {}
func (this *Text) Children() []Node {
	return make([]Node, 0)
}
func (this *Text) Clone() Node {
	return NewText(this.Content)
}
func (this *Text) Export(w io.Writer) {
	xml.Escape(w, []byte(this.Content))
}
func (this *Text) Get(name string) (string, bool) {
	return "", false
}
func (this *Text) Insert(pos int, node Node) {}
func (this *Text) Parent() Node {
	return this.parent
}
func (this *Text) Remove(n Node) bool {
	return false
}
func (this *Text) Set(name, value string) {}
func (this *Text) SetParent(parent Node) {
	this.parent = parent
}
func (this *Text) String() string {
	return this.Content
}

// HTML
func NewHtml(content string) *Html {
	return &Html{Text{nil, content}}
}
func (this *Html) Clone() Node {
	return NewHtml(this.Text.Content)
}
func (this *Html) Export(w io.Writer) {
	io.WriteString(w, this.Text.Content)
}

// Fragment
func NewFragment() *Fragment {
	return &Fragment{nil, Container{make([]Node, 0)}}
}
func (this *Fragment) Append(n Node) {
	if _, ok := n.(*Fragment); ok {
		for _, c := range n.Children() {
			this.Append(c)
		}
		return
	}
	this.Container.Append(n)
	n.SetParent(this)
}
func (this *Fragment) Clone() Node {
	n := NewFragment()
	for _, c := range this.Children() {
		n.Append(c.Clone())
	}
	return n
}
func (this *Fragment) Export(w io.Writer) {
	for _, c := range this.children {
		c.Export(w)
	}
}
func (this *Fragment) Get(name string) (string, bool) {
	return "", false
}
func (this *Fragment) Insert(pos int, node Node) {
	this.Container.Insert(pos, node)
	node.SetParent(this)
}
func (this *Fragment) Parent() Node {
	return this.parent
}
func (this *Fragment) Set(name, value string) {

}
func (this *Fragment) SetParent(parent Node) {
	this.parent = parent
}
func (this *Fragment) String() string {
	var buf bytes.Buffer
	this.Export(&buf)
	return buf.String()
}
