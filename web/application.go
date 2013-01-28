package web

import (
	"errors"
	"fmt"
	"net/http"
	
	"github.com/badgerodon/collections/tst"
)

type (
	Application struct {
		ErrorHandler func(Context, error)
		MissingHandler func(Context)
		route string
		layout string
		routes *tst.TernarySearchTree
	}
)

var (
	Applications map[string]*Application
	DefaultApplication *Application
)

func init() {
	Applications = make(map[string]*Application)
	DefaultApplication = NewApplication(
		"/", 
		"layouts/application", 
	)
}

func NewApplication(route, layout string) *Application {
	app := &Application{
		ErrorHandler: DefaultErrorHandler, 
		MissingHandler: DefaultMissingHandler,
		route: route,
		layout: layout,
		routes: tst.New(),
	}
	Applications[route] = app
	return app
}

func DefaultErrorHandler(ctx Context, err error) {
	res := ctx.Response()
	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(500)
	res.Write([]byte(fmt.Sprint(err)))
}
func DefaultMissingHandler(ctx Context) {
	res := ctx.Response()
	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(404)
	res.Write([]byte("Unknown page"))
}
func (this *Application) Handle(w http.ResponseWriter, req *http.Request) {
	ctx := &defaultContext{this, w, req, "", ""}
	
	path := req.URL.Path[len(this.route):]
	v := this.routes.GetLongestPrefix(path)
	if v == nil {
		this.MissingHandler(ctx)
		return
	}
	// Multiple routes can have the same prefix
	routes, ok := v.([]*RouteDefinition)
	if !ok {
		this.ErrorHandler(ctx, errors.New("Invalid handler function"))
		return
	}
	for _, route := range routes {
		parameters, ok := route.Match(ctx, path)
		if ok {
			err := route.Handler(ctx, parameters)
			if err != nil {
				this.ErrorHandler(ctx, err)
			}
			return
		}
	}
	
	// Guess nothing matched
	this.MissingHandler(ctx)
}