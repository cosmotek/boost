package main

import (
	"fmt"

	"github.com/rucuriousyet/boost"
	"github.com/rucuriousyet/boost/types"
)

type Configuration struct {
	GpuProcessing struct {
		Enabled bool
		Device  string
	}
	Networks struct {
		Npipe struct {
			Address  string
			MaxConns uint
		}

		TCP struct {
			Address      string
			MaxConns     uint
			ReadDeadline uint16
		}

		UDP struct {
			Address      string
			MaxConns     uint
			ReadDeadline uint16
		}
	}
	Resources struct {
		MaxMem uint32
	}
}

func main() {
	conf, err := boost.NewAppConfig(
		"Stackmesh",
		"stackmesh.conf",
		false,
		types.NewString("host_os", "linux"),
		types.NewBool("is_desktop", true),
		types.NewNumber("mem", 1024),
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
