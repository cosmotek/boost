# Boost
[![godoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/rucuriousyet/boost)

Scripted Lua configuration for Go!

![boosters](https://upload.wikimedia.org/wikipedia/commons/thumb/4/45/020408_STS110_Atlantis_launch.jpg/1158px-020408_STS110_Atlantis_launch.jpg)

Boost is a simple (and somewhat fast) configuration/bootstrapping file engine designed to make script-based configurations dead simple. Using Boost, Applications can interface with a simple Lua program, mapping data from said program to the map or struct of choice. Boost was somewhat inspired the Vagrantfile used by Hashicorp's Vagrant as well as the makefile and Caddyfile. Boost was created to fill the need for a dynamic configuration system that could be used in the Stackmesh daemon, sandbox, cli tool etc. Additionally, I came up with the early concept of Boost when working on a app packaging tool I call Shuttle.

## Install
This package is go get-able so I'm sure you know what to do...
(but just in case you don't ↲)

`go get github.com/rucuriousyet/boost`

## Example
In addition to the example below, I've included an much more detailed example in the `examples` folder.

**example.conf**
```lua
-- some global values
random = {
  one = 1,
  two = 2,
  strawberries = "red"
}

-- configuration function
MyApp.configure = function(config)
  config.favorites.fruits = {
    "apple",
    "orange"
  }

  config.favorites.color = "red"
  config.favorites.animals = true
end
```

**main.go**
```go
package main

import (
	"fmt"

	"github.com/rucuriousyet/boost"
	"github.com/rucuriousyet/boost/types"
)

type Configuration struct {
	Favorites struct {
    Color string
    Animals bool
    Fruits []string
  }
}

func main() {
	conf, err := boost.NewAppConfig(
		"MyApp",
		"example.conf",
    // enabled debug logging
		true,
		types.NewString("name", "seth"),
	)

	defer conf.Cleanup()
	if err != nil {
		panic(err)
	}

	config := &Configuration{}
	err = conf.ParseFunction("configure", config)
	if err != nil {
		panic(err)
	}

	fmt.Println(config)
}
```

## API Coverage
At the moment all Lua stdlib packages appear to be fully functional (os, io, file, math etc..) as well as a few Gopher-Lua modules for serialization, http and utility. Native/Pure Lua libraries and C-Wrapping modules are yet to be tested. Most pure Lua code should work, anything interfacing with C should not be expected to work.

Included Gopher-Lua Modules:
+ http (http client wrapping net/http)
+ re (regex)
+ json
+ yaml
+ xmlpath
+ url (parser)
+ lfs (Lua Filesystem, not complete)

## Shout-outs

Special Thanks to the creator of https://github.com/yuin/gopher-lua which this library is almost fully dependent on.

## Contribution and maintenance

This library is currently under active development alongside Stackmesh. Please note that Boost has not yet been fully tested, be careful when using this library for anything in sensitive environments. If you would like to make a suggestion or report a bug, please feel free to submit a PR or issue. I really hope you enjoy using Boost! Thanks!
