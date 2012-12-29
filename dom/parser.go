package dom

import (
	//"strings"
	"unicode"
	"unicode/utf8"
)

type parser struct {
	buffer []byte
	pos int
	stack []int
}

func (p *parser) save() {
	p.stack = append(p.stack, p.pos)
}
func (p *parser) restore() {
	p.pos = p.stack[len(p.stack)-1]
}
func (p *parser) pop() {
	p.stack = p.stack[:len(p.stack)-1]
}
func (p *parser) peek() rune {
	if p.pos < len(p.buffer) {
		r, _ := utf8.DecodeRune(p.buffer[p.pos:])
		return r
	}
	return 0
}
func (p *parser) next() rune {
	if p.pos < len(p.buffer) {
		r, sz := utf8.DecodeRune(p.buffer[p.pos:])
		p.pos += sz
		return r
	}
	return 0
}
func (p *parser) whitespace() (string, bool) {
	p.save()
	defer p.pop()

	str := ""
	for {
		c := p.peek()
		if unicode.IsSpace(c) {
			str += string(c)
			p.next()
		} else {
			break
		}
	}

	if len(str) == 0 {
		p.restore()
		return "", false
	}

	return str, true
}
func (p *parser) identifier() (string, bool) {
	p.save()
	defer p.pop()

	c := p.next()
	if !unicode.IsLetter(c) {
		p.restore()
		return "", false
	}

	str := string(c)
	for {
		c := p.peek()
		if unicode.IsLetter(c) || unicode.IsDigit(c) || c == '-' || c == '_' {
			str += string(c)
			p.next()
		} else {
			break
		}
	}
	return str, true
}
func (p *parser) string() (string, bool) {
	p.save()
	defer p.pop()

	sc := p.next()
	if sc != '"' && sc != '\'' {
		p.restore()
		return "", false
	}

	str := ""
	for {
		c := p.next()
		if c == 0 {
			p.restore()
			return "", false
		}
		if c == sc {
			break
		}
		str += string(c)
	}

	return str, true
}
func (p *parser) startElement() (*Element, bool) {
	p.save()
	defer p.pop()

	var el *Element

	if p.next() != '<' {
		p.restore()
		return nil, false
	}

	tag, ok := p.identifier()
	if !ok {
		p.restore()
		return nil, false
	}

	el = NewElement(tag)

	p.whitespace()

	for {
		attributeName, ok := p.identifier()
		if !ok {
			break
		}
		p.whitespace()
		if p.peek() == '=' {
			p.whitespace()
			attributeValue, ok := p.string()
			if !ok {
				attributeValue, ok = p.identifier()
			}
			el.Set(attributeName, attributeValue)
			if !ok {
				break
			}
		} else {
			el.Set(attributeName, "")
		}
	}

	p.whitespace()

	if p.next() != '>' {
		p.restore()
		return nil, false
	}

	return el, true
}
func (p *parser) endElement() (string, bool) {
	p.save()
	defer p.pop()

	if p.next() != '<' {
		p.restore()
		return "", false
	}

	if p.next() != '/' {
		p.restore()
		return "", false
	}

	tag, ok := p.identifier()
	if !ok {
		p.restore()
		return "", false
	}

	p.whitespace()

	if p.next() != '>' {
		p.restore()
		return "", false
	}

	return tag, true
}
func (p *parser) comment() (string, bool) {
	return "", false
}
func (p *parser) parse() *Fragment {
	var root *Fragment = NewFragment()
	var cur Node = root
	txt := ""
	for {
		el, ok := p.startElement()
		if ok {
			cur.Append(NewText(txt))
			txt = ""
			cur.Append(el)
			cur = el
			continue
		}
		_, ok = p.comment()
		if ok {
			cur.Append(NewText(txt))
			txt = ""
			continue
		}
		_, ok = p.endElement()
		if ok {
			cur.Append(NewText(txt))
			if cur.Parent() != nil {
				cur = cur.Parent()
			}
			continue
		}
		c := p.next()
		if c == 0 {
			break
		}
		txt += string(c)
	}
	cur.Append(NewText(txt))
	return root
}