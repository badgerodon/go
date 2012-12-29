package css

import (
	//"strings"
	"unicode"
	"unicode/utf8"
)

type Parser struct {
	buffer []byte
	pos int
	stack []int
}

func (p *Parser) save() {
	p.stack = append(p.stack, p.pos)
}
func (p *Parser) restore() {
	p.pos = p.stack[len(p.stack)-1]
}
func (p *Parser) pop() {
	p.stack = p.stack[:len(p.stack)-1]
}
func (p *Parser) peek() rune {
	if p.pos < len(p.buffer) {
		r, _ := utf8.DecodeRune(p.buffer[p.pos:])
		return r
	}
	return 0
}
func (p *Parser) next() rune {
	if p.pos < len(p.buffer) {
		r, sz := utf8.DecodeRune(p.buffer[p.pos:])
		p.pos += sz
		return r
	}
	return 0
}
func (p *Parser) whitespace() (string, bool) {
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
func (p *Parser) name() (string, bool) {
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
func (p *Parser) universal() bool {
	p.save()
	defer p.pop()

	if p.next() == '*' {
		return true
	}

	p.restore()
	return false
}
func (p *Parser) id() (string, bool) {
	p.save()
	defer p.pop()

	if p.next() == '#' {
		if name, ok := p.name(); ok {
			return name, true
		}
	}

	p.restore()
	return "", false
}
func (p *Parser) class() (string, bool) {
	p.save()
	defer p.pop()

	if p.next() == '.' {
		if name, ok := p.name(); ok {
			return name, true
		}
	}

	p.restore()
	return "", false
}
func (p *Parser) simple() (Filter, bool) {
	filters := []Filter{}
	if name, ok := p.name(); ok {
		filters = append(filters, filterTag(name))
	}
	if ok := p.universal(); ok {
		filters = append(filters, filterUniversal())
	}
	for {
		if cls, ok := p.class(); ok {
			filters = append(filters, filterClass(cls))
		}	else if id, ok := p.id(); ok {
			filters = append(filters, filterId(id))
		} else {
			break
		}
	}
	if len(filters) > 0 {
		return filterIntersection(filters...), true
	}
	return nil, false
}
func (p *Parser) expression() (Filter, bool) {
	p.save()
	defer p.pop()

	f, ok := p.simple()
	if !ok {
		p.restore()
		return nil, false
	}
	if _, ok := p.whitespace(); ok {
		f2, ok := p.expression()
		if ok {
			f = filterDescendants(f, f2)
		}
	}
	if p.peek() == '+' {
		p.next()
		p.whitespace()
		f2, ok := p.expression()
		if ok {
			f = filterAdjacentSiblings(f, f2)
		}
	}
	if p.peek() == ',' {
		p.next()
		p.whitespace()
		f2, ok := p.expression()
		if ok {
			f = filterUnion(f, f2)
		}
	}

	return f, p.pos == len(p.buffer)
}

func Parse(exp string) (Filter, bool) {
	p := &Parser{[]byte(exp),0,make([]int, 0)}
	return p.expression()
}