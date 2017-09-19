# Boost
[![godoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/rucuriousyet/boost)

Scripted Lua configuration for Go!

Boost is a simple (and somewhat fast) configuration/bootstrapping file engine designed to make script-based configurations dead simple. Using Boost, Applications can interface with a simple Lua program, mapping data from said program to the map or struct of choice. Boost was somewhat inspired the Vagrantfile used by Hashicorp's Vagrant as well as the makefile and Caddyfile. Boost was created to fill the need for a dynamic configuration system that could be used in the Stackmesh daemon, sandbox, cli tool etc.

## Install
This package is go get-able so I'm sure you know what to do...
(but just in case you don't ↲)

`go get github.com/rucuriousyet/boost`

## Example
**example.conf**
```lua
random = {
  one = 1,
  two = 2,
  strawberries = "red"
}

MyApp.configure = function(config)
  config.favorites.fruits = {
    "apple",
    "orange"
  }

  config.favorites.color = "red"
  config.favorites.animals = true
end

-- config functions can only take a
-- single parameter, the name of that
-- parameter doesn't matter
MyApp.sandbox = function(controls)
  controls.lighting = true
  controls.sfx = 200
end

```

**main.go**
```go
package main

import (
  "fmt"

  "github.com/rucuriousyet/boost"
)

type Example struct {
  Favorites struct {
    Fruits []string
    Color string
    Animals bool
  }
}

type AnotherExample struct {
  Lighting bool
  SFX uint
}

func main() {
  config, err := boost.NewAppConfig("MyApp", "example.conf")
  defer config.Close()
  if err != nil {
    panic(err)
  }

  example := &Example{}
  err = config.ParseFunction("configure", example)
  if err != nil {
    panic(err)
  }

  anotherExample := &AnotherExample{}
  err = config.ParseFunction("sandbox", anotherExample)
  if err != nil {
    panic(err)
  }

  globals := map[string]interface{}{}
  err = config.GetGlobal("random", &globals)
  if err != nil {
    panic(err)
  }

  fmt.Println(globals["one"]) // => 1
  fmt.Println(example.Fruits) // => [apple, orange]
  fmt.Println(anotherExample.SFX) // => 200
}
```
## Shout-outs

Special Thanks to the creator of https://github.com/yuin/gopher-lua which this library is almost fully dependent on.

## Contribution and maintenance

This library is currently under active development alongside Stackmesh. Please note that Boost has not yet been fully tested, be careful when using this library for anything in sensitive environments. If you would like to make a suggestion or report a bug, please feel free to submit a PR or issue. I really hope you enjoy using Boost! Thanks!
