package web

import (
	"fmt"
	"net/http"
)

type (
	Context interface {
		Layout(layout string)
		Render(data interface{})
		Request() *http.Request
		Response() http.ResponseWriter
		Template(tpl string)
		Write(data interface{})
	}
	
	defaultContext struct {
		application *Application
		response http.ResponseWriter
		request *http.Request
		template string
		layout string
	}
)

func (this *defaultContext) Render(data interface{}) {
	bs, err := this.application.Render(this.layout, this.template, data)
	if err != nil {
		this.application.ErrorHandler(this, err)
		return
	}
	_, err = this.response.Write(bs)
	if err != nil {
		this.application.ErrorHandler(this, err)
		return
	}
}
func (this *defaultContext) Request() *http.Request {
	return this.request
}
func (this *defaultContext) Response() http.ResponseWriter {
	return this.response
}
func (this *defaultContext) Layout(layout string) {
	this.layout = layout
}
func (this *defaultContext) Template(template string) {
	this.template = template
}
func (this *defaultContext) Write(data interface{}) {
	_, err := this.response.Write([]byte(fmt.Sprint(data)))
	if err != nil {
		this.application.ErrorHandler(this, err)
	}
}