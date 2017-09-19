package main

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/ailncode/gluaxmlpath"
	"github.com/cjoudrey/gluahttp"
	"github.com/cjoudrey/gluaurl"
	"github.com/kohkimakimoto/gluayaml"

	"github.com/yuin/gluamapper"
	"github.com/yuin/gluare"
	lua "github.com/yuin/gopher-lua"
	gjson "layeh.com/gopher-json"
)

type App struct {
}

func NewAppConfig(name, filename string) (*App, error) {
	return nil, nil
}

func (a *App) SetGlobal(key string, val interface{}) error {
	return nil
}

func (a *App) ParseFunction(method string, mapping interface{}) error {
	return nil
}

func ParseByApp(appname, methodname, filename string, mapping interface{}) error {
	if reflect.ValueOf(mapping).Kind() != reflect.Ptr {
		return errors.New("input mapping must be a pointer")
	}

	l := lua.NewState()
	defer l.Close()

	l.PreloadModule("re", gluare.Loader)
	l.PreloadModule("yaml", gluayaml.Loader)

	l.PreloadModule("url", gluaurl.Loader)
	gluaxmlpath.Preload(l)

	gjson.Preload(l)
	l.PreloadModule("http", gluahttp.NewHttpModule(&http.Client{}).Loader)

	l.SetGlobal(appname, l.NewTable())
	configObj := fmt.Sprintf("__boost__%s_%s__boost__", appname, methodname)
	l.SetGlobal(configObj, l.NewTable())

	mappingType := reflect.ValueOf(mapping).Type().Elem()
	for i := 0; i < mappingType.NumField(); i++ {
		typ := mappingType.Field(i).Type.Kind()

		// need to check for objects inside this object
		if typ == reflect.Struct || typ == reflect.Map {
			name := mappingType.Field(i).Name
			l.GetGlobal(configObj).(*lua.LTable).RawSetString(strings.ToLower(name), l.NewTable())
		}
	}

	if err := l.DoFile(filename); err != nil {
		panic(err)
	}

	if err := l.CallByParam(lua.P{
		Fn:      l.GetGlobal(appname).(*lua.LTable).RawGetString(methodname),
		NRet:    0,
		Protect: true,
	}, l.GetGlobal(configObj)); err != nil {
		panic(err)
	}

	if err := gluamapper.Map(l.GetGlobal(configObj).(*lua.LTable), mapping); err != nil {
		panic(err)
	}

	return nil
}

type Example struct {
	Favorites struct {
		Fruits  []string
		Color   string
		Animals bool
	}
}

func main() {
	var example Example
	err := ParseByApp("MyApp", "configure", "example.conf", &example)
	if err != nil {
		panic(err)
	}

	fmt.Println("returned", example)
}
