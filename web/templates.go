package web

import (
	"bytes"
	"errors"
	"html/template"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type (
	LayoutModel struct {
		Body template.HTML
		Scripts []string
		Styles []string
		Model interface{}
	}
)

func getTemplates() (*template.Template, error) {
	t := template.New("")
	root := "app/tpl"
	err := filepath.Walk(root, func(p string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		
		if path.Ext(p) != ".tpl" {
			return nil
		}
		
		if err != nil {
			return err
		}
		
		contents, err := ioutil.ReadFile(p)
		if err != nil {
			return err
		}
		
		nm := strings.Replace(p, "\\", "/", -1)
		nm = nm[len(root):]
		nm = nm[:len(nm) - len(path.Ext(nm))]
		
		child := t.New(nm)
		_, err = child.Parse(string(contents))
		
		return err
	})
	return t, err
}

func (this *Application) Render(layout, tpl string, model interface{}) ([]byte, error) {
	// Get the templates
	t, err := getTemplates()
	if err != nil {
		return nil, err
	}
	
	tt := t.Lookup(this.route + tpl)
	lt := t.Lookup(this.route + layout)
	
	if tt == nil {
		return nil, errors.New("Unknown template `" + tpl + "`")
	}
	
	// Render the template
	var buf bytes.Buffer
	err = tt.Execute(&buf, model)
	if err != nil {
		return nil, err
	}
	// Render the layout
	if layout != "" && layout != "/" && lt != nil {
		scripts, err := getScriptURLs(this.route)
		if err != nil {
			return nil, err
		}
		styles, err := getStyleURLs(this.route)
		if err != nil {
			return nil, err
		}
	
		body := template.HTML(buf.String())
		buf.Reset()
		err = lt.Execute(&buf, LayoutModel{
			Body: body,
			Scripts: scripts,
			Styles: styles,
			Model: model,
		})
		if err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}