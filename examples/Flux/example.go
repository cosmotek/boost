package main

import (
	"fmt"

	"github.com/rucuriousyet/boost"
	"github.com/rucuriousyet/boost/types"
)

type Configuration struct {
	Version     string
	Hash        string
	Name        string
	Destination string
	Entrypoint  string
	DryRun      struct {
		Type   uint
		Config string
	}
	Upload bool
}

func main() {
	conf, err := boost.NewAppConfig(
		"Flux",
		"Fluxfile",
		true,
		types.NewString("host_os", "linux"),
		types.NewNumber("docker_compose", 0),
		types.NewString("destination", "release"),
		types.NewBool("manual_exec", false),
		types.NewBool("immutable_env", false),
	)
	defer conf.Cleanup()
	if err != nil {
		panic(err)
	}

	config := &Configuration{}
	err = conf.ParseFunction("buildconf", config)
	if err != nil {
		panic(err)
	}

	fmt.Println(config)
	mapping := &map[string]interface{}{}

	conf.GetGlobal("job", mapping)
	fmt.Println(mapping)
}
