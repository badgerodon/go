package web

import (
	"fmt"
	"reflect"
	"strings"
)

type (
	RouteDefinition struct {
		Components []RouteComponent	
		Handler func(Context, []interface{}) error
		Template string
	}
)
func (this RouteDefinition) String() string {
	return fmt.Sprint(this.Components)
}
func (this RouteDefinition) Prefix() string {
	parts := []string{}
	for _, c := range this.Components {
		if static, ok := c.(StaticRouteComponent); ok {
			parts = append(parts, string(static))
		}
	}
	return strings.Join(parts, "/")
}
func (this RouteDefinition) Match(ctx Context, path string) ([]interface{}, bool) {
	ps := []interface{}{}
	
	var ok bool
	var val interface{}
	for _, c := range this.Components {
		path, ok, val = c.Match(path)
		if !ok {
			return ps, false
		}
		if val != nil {
			ps = append(ps, val)
		}
		if len(path) > 0 && path[0] == '/' {
			path = path[1:]
		}
	}
	
	return ps, len(path) == 0
}
func Parse(route string) []RouteComponent {
	if len(route) > 0 && route[0] == '/' {
		route = route[1:]
	}
	parts := strings.Split(route, "/")
	components := []RouteComponent{}
	for _, p := range parts {
		var component RouteComponent
		if p == "?" {
			component = &ParameterRouteComponent{reflect.String}
		} else {
			component = StaticRouteComponent(p)	
		}
		components = append(components, component)
	}
	return components
}

func (this *Application) addRoute(rd *RouteDefinition) {	
	rs := this.routes.Get(rd.Prefix())
	if rs == nil {
		rs = []*RouteDefinition{rd}
	} else {
		rs = append(rs.([]*RouteDefinition), rd)
	}
	this.routes.Insert(rd.Prefix(), rs)
}	

func (this *Application) mkTemplateName(name string) string {
	name = strings.Replace(name, "Controller/", "/", -1)
	name = strings.ToLower(name)
	return name
}
func (this *Application) routeController(prefix string, handler interface{}) {	
	typ := reflect.TypeOf(handler)
	nm := typ.Name()
	// De-pointer
	if typ.Kind() == reflect.Ptr {
		nm = typ.Elem().Name()
	}
	
	for i := 0; i < typ.NumMethod(); i++ {
		m := typ.Method(i)
		var rd *RouteDefinition
		switch m.Name {
		case "Index":
			rd = this.routeFunction(prefix, handler, m.Func.Interface())
		case "Show":
			rd = this.routeFunction(prefix + "/?", handler, m.Func.Interface())
		default:
			route := prefix + "/" + m.Name
			for j := 0; j < m.Type.NumIn()-2; j++ {
				route += "/?"
			}
			rd = this.routeFunction(route, handler, m.Func.Interface())
		}
		rd.Template = this.mkTemplateName(nm + "/" + m.Name)
		this.addRoute(rd)
	}
}

func (this *Application) routeFunction(route string, rcvr interface{}, handler interface{}) *RouteDefinition {
	var ctx Context = &defaultContext{}
	tctx := reflect.TypeOf(ctx)
	
	rd := &RouteDefinition{}
	
	typ := reflect.TypeOf(handler)
	// De-pointer
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	
	cs := Parse(route)
		
	i := 1
	if rcvr != nil {
		i++
	}
	
	if typ.NumIn() < i || !tctx.AssignableTo(typ.In(i-1)) {
		panic("Must take at least one parameter and it must be a Context")
	}
	
	for j := 0; j < len(cs) && i < typ.NumIn(); j++ {
		if prc, ok := cs[j].(*ParameterRouteComponent); ok {
			prc.Kind = typ.In(i).Kind()
			i++
		}
		if src, ok := cs[j].(*StaticRouteComponent); ok {
			rd.Template += "/" + string(*src)
		}
	}
	if rd.Template == "" || rd.Template == "/" {
		rd.Template += "index"
	}
	rd.Template = this.mkTemplateName(rd.Template)
	
	rd.Components = cs
	rd.Handler = func(ctx Context, parameters []interface{}) error {
		ctx.Layout(this.layout)
		ctx.Template(rd.Template)
		vs := make([]reflect.Value, typ.NumIn())
		i := 0
		if rcvr != nil {
			vs[i] = reflect.ValueOf(rcvr)
			i++
		}
		vs[i] = reflect.ValueOf(ctx)
		i++
		for j := 0; j < len(parameters) && (j+i) < len(vs); j++ {
			vs[j+i] = reflect.ValueOf(parameters[j])
		}
		reflect.ValueOf(handler).Call(vs)
		return nil
	}	
	
	return rd
}

func (this *Application) Route(route string, handler interface{}) {
	typ := reflect.TypeOf(handler)
	// De-pointer
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	
	switch typ.Kind() {
	case reflect.Func:
		this.addRoute(this.routeFunction(route, nil, handler))
	case reflect.Struct:
		this.routeController(route, handler)
	default:
		panic("Unknown handler type: `" + typ.String() + "`")
	}
}
func Route(route string, handler interface{}) {
	DefaultApplication.Route(route, handler)
}