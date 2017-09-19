# boost
[![godoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/rucuriousyet/boost)

scripted Lua configuration for Go!

# Install
`go get github.com/rucuriousyet/boost`

# Example
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
