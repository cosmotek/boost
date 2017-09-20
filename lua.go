package boost

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"reflect"

	"github.com/BixData/gluasocket"
	"github.com/ailncode/gluaxmlpath"
	"github.com/cjoudrey/gluahttp"
	"github.com/cjoudrey/gluaurl"
	"github.com/kohkimakimoto/gluayaml"
	"github.com/rs/zerolog"
	log "github.com/rs/zerolog/log"
	"github.com/serenize/snaker"
	"github.com/yuin/gluamapper"
	"github.com/yuin/gluare"
	lua "github.com/yuin/gopher-lua"
	gjson "layeh.com/gopher-json"
	lfs "layeh.com/gopher-lfs"
)

// App is used to maintain the underlying
// state for the lua vm
type App struct {
	state *lua.LState
	name  string
	print func(data string)
}

// Global is a simple object used to create
// a lua object (const, table etc) that can
// be injected into the app config for use
// within the lua script
type Global interface {
	GetKey() string
	GetValue() lua.LValue
}

// NewAppConfig creates a lua vm session and app config
// using an application name (used to create the root table),
// filename/path and some globals (optional). The passed file
// may have any extension as long as it is a valid lua script.
// You may prefer to use a .lua extension for automatic syntax
// highlighting in editors.
func NewAppConfig(name, filename string, enableLogging bool, globals ...Global) (*App, error) {
	print := func(data string) {

	}

	if enableLogging {
		logger := log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		print = func(data string) {
			logger.Debug().Msg(data)
		}
	}

	print("creating new VM state")
	l := lua.NewState()

	print("adding module: 'lfs'")
	lfs.Preload(l)

	print("adding module: 'socket'")
	gluasocket.Preload(l)

	print("adding module: 're (regex)'")
	l.PreloadModule("re", gluare.Loader)

	print("adding module: 'yaml'")
	l.PreloadModule("yaml", gluayaml.Loader)

	print("adding module: 'url'")
	l.PreloadModule("url", gluaurl.Loader)

	print("adding module: 'xmlpath'")
	gluaxmlpath.Preload(l)

	print("adding module: 'gjson'")
	gjson.Preload(l)

	print("adding module: 'http'")
	l.PreloadModule("http", gluahttp.NewHttpModule(&http.Client{}).Loader)

	print("creating root table: '" + name + "'")
	l.SetGlobal(name, l.NewTable())

	print("parsing globals")
	for _, global := range globals {
		print(fmt.Sprintf("adding global %s=>%s", global.GetKey(), global.GetValue()))
		l.SetGlobal(global.GetKey(), global.GetValue())
	}

	print(fmt.Sprintf("running file: %s", filename))
	if err := l.DoFile(filename); err != nil {
		return nil, err
	}

	return &App{
		state: l,
		name:  name,
		print: print,
	}, nil
}

// Cleanup closes and cleans up the lua VM, this must be
// called when all interaction with the config app is complete
func (a *App) Cleanup() {
	a.print("cleaning up")
	a.state.Close()
}

// GetGlobal retrieves a global object from the lua vm and
// maps it to the provided mapping pointer. This seems to only
// work when the object matching the provided key is a table.
func (a *App) GetGlobal(key string, mapping interface{}) error {
	if reflect.ValueOf(mapping).Kind() != reflect.Ptr {
		return errors.New("input mapping must be a pointer")
	}

	a.print(fmt.Sprintf("mapping global: %s to provided mapping", key))
	if err := gluamapper.Map(a.state.GetGlobal(key).(*lua.LTable), mapping); err != nil {
		return err
	}
	return nil
}

func (a *App) digg(parent string, code *string, object reflect.Type) {
	for i := 0; i < object.NumField(); i++ {
		typ := object.Field(i).Type
		kind := typ.Kind()

		// need to check for objects inside this object
		if kind == reflect.Struct || kind == reflect.Map {
			name := snaker.CamelToSnake(object.Field(i).Name)
			*code += fmt.Sprintf("%s.%s = {}\n", parent, name)
			a.digg(parent+"."+name, code, typ)
		}
	}
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
	a.print(fmt.Sprintf("creating config passthru object: %s", configObj))

	bootstrapStr := fmt.Sprintf("%s = {}\n", configObj)
	a.digg(configObj, &bootstrapStr, reflect.ValueOf(mapping).Type().Elem())

	a.print(fmt.Sprintf("adding mapping bootstrap tables: \n%s", bootstrapStr))
	err := l.DoString(bootstrapStr)
	if err != nil {
		return err
	}

	a.print(fmt.Sprintf("calling method: %s.%s(%s)", a.name, method, configObj))
	if err := l.CallByParam(lua.P{
		Fn:      l.GetGlobal(a.name).(*lua.LTable).RawGetString(method),
		NRet:    0,
		Protect: true,
	}, l.GetGlobal(configObj)); err != nil {
		return err
	}

	a.print("mapping passthru object to mapping")
	if err := gluamapper.Map(l.GetGlobal(configObj).(*lua.LTable), mapping); err != nil {
		return err
	}

	return nil
}
