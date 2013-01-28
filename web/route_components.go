package web

import (
	"reflect"
	"strconv"
	"strings"
)

type (
	RouteComponent interface{
		Match(string) (string, bool, interface{})
	} 
	StaticRouteComponent string	
	ParameterRouteComponent struct {
		Kind reflect.Kind
	}
)

func (this *ParameterRouteComponent) Match(path string) (string, bool, interface{}) {
	var val interface{}
	
	str := strings.SplitN(path, "/", 1)[0]
	switch this.Kind {
	case reflect.Bool:
		val = str == "true" || str == "on" || str == "1"
	case reflect.Int:
		var err error
		val, err = strconv.Atoi(str)
		if err != nil {
			return path, false, val
		}
	case reflect.Float32:
		f, err := strconv.ParseFloat(str, 32)
		if err != nil {
			return path, false, val
		}
		val = float32(f)
	case reflect.Float64:
		f, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return path, false, val
		}
		val = float64(f)
	case reflect.String:
		val = str
	default:
		panic("Unknown Kind")
	}
	
	path = path[len(str):]
	
	return path, true, val
}

func (this StaticRouteComponent) Match(path string) (string, bool, interface{}) {
	if strings.HasPrefix(path, string(this)) {
		return path[len(string(this)):], true, nil
	}
	return path, false, nil
}