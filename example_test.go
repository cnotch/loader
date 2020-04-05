package loader_test

import (
	"fmt"

	 "github.com/cnotch/loader"
)

func Example_initLoader() {
	// Our struct which is used for configuration
	type ServerConfig struct {
		Name    string
		Port    int
		Enabled bool
		Users   []string
	}
	s := &ServerConfig{}

	if err := loader.Load(s, loader.Func(func(cfg interface{}) error {
		s := cfg.(*ServerConfig)
		s.Name = "notch"
		s.Port = 554
		return nil
	})); err != nil {
		panic(err)
	}

	fmt.Println("Host-->", s.Name)
	fmt.Println("Port-->", s.Port)

	// Output:
	// Host--> notch
	// Port--> 554
}
