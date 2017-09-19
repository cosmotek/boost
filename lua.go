package boost

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

// App is used to maintain the underlying
// state for the lua vm
type App struct {
	state *lua.LState
	name  string
}

// GlobalType maps some lua vm types to
// golang for use in creating global objects
type GlobalType lua.LValue

// Global is a simple object used to create
// a lua object (const, table etc) that can
// be injected into the app config for use
// within the lua script
type Global struct {
	Key   string
	Type  GlobalType
	Value interface{}
}

// NewAppConfig creates a lua vm session and app config
// using an application name (used to create the root table),
// filename/path and some globals (optional). The passed file
// may have any extension as long as it is a valid lua script.
// You may prefer to use a .lua extension for automatic syntax
// highlighting in editors.
func NewAppConfig(name, filename string, globals ...Global) (*App, error) {
	l := lua.NewState()
	l.PreloadModule("re", gluare.Loader)
	l.PreloadModule("yaml", gluayaml.Loader)

	l.PreloadModule("url", gluaurl.Loader)
	gluaxmlpath.Preload(l)

	gjson.Preload(l)
	l.PreloadModule("http", gluahttp.NewHttpModule(&http.Client{}).Loader)

	l.SetGlobal(name, l.NewTable())
	// for _, global := range globals {
	// 	// iterate and create globals
	// }

	if err := l.DoFile(filename); err != nil {
		panic(err)
	}

	return &App{
		state: l,
		name:  name,
	}, nil
}

// Cleanup closes and cleans up the lua VM, this must be
// called when all interaction with the config app is complete
func (a *App) Cleanup() {
	a.state.Close()
}

// GetGlobal retrieves a global object from the lua vm and
// maps it to the provided mapping pointer. This seems to only
// work when the object matching the provided key is a table.
func (a *App) GetGlobal(key string, mapping interface{}) error {
	if reflect.ValueOf(mapping).Kind() != reflect.Ptr {
		return errors.New("input mapping must be a pointer")
	}

	if err := gluamapper.Map(a.state.GetGlobal(key).(*lua.LTable), mapping); err != nil {
		panic(err)
	}
	return nil
}

// ParseFunction runs the method by provided method name, app name and maps
// the result to the provided 'mapping' pointer. The mapping must be a
// ptr of a struct or map type.
func (a *App) ParseFunction(method string, mapping interface{}) error {
	l := a.state

	if reflect.ValueOf(mapping).Kind() != reflect.Ptr {
		return errors.New("input mapping must be a pointer")
	}

	configObj := fmt.Sprintf("__boost__%s_%s__boost__", a.name, method)
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

	if err := l.CallByParam(lua.P{
		Fn:      l.GetGlobal(a.name).(*lua.LTable).RawGetString(method),
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
