package dom

import (
	"bytes"
	"fmt"
	"testing"
	"encoding/xml"
)

func Test(t *testing.T) {
}

func TestXml(t *testing.T) {
	var buf bytes.Buffer
	buf.WriteString(`<?xml version="1.0" encoding="UTF-8"?><test><child></child></test>`)
	d := xml.NewDecoder(&buf)
	fmt.Println(FromXml(d))
}