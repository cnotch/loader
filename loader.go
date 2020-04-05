// Copyright (c) 2019,CAO HONGJU. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

/*
Package loader provides easy to use object loading capability.

Sometimes, we need to load the same object from multiple sources
in order of priority.

For example, the configuration data of the application:
 1. Built-in default value
 2. Load from the configuration file
 3. Load from environment variables
 4. Parse from the command line parameter
 5. Other sources (centralized configuration management, etc.)


loader provides a surprisingly simple pattern.
For example, load configuration, which typically consists of the following steps:
1. Define configuration structure

	type Config struct {
		Log  LogConfig  `json:"log,omitempty"`
		Rtsp RtspConfig `json:"rtsp"`
	}

	type RtspConfig struct {
		Port int `json:"port"`
	}

	type LogConfig struct {
		Debug bool `json:"debug"`
	}

2. Load by priority

	import (
		"flag"

		"github.com/cnotch/loader"
	)

	const (
		envPrefix = "APP"
	)

	func LoadConfig(configPath string) *Config {
		conf := new(Config)

		// Define the initialization loader
		initLoader := func(cfg interface{}) error {
			flag.BoolVar(&conf.Log.Debug, "log-debug", false, "Enable debug level log")
			flag.IntVar(&conf.Rtsp.Port, "rtsp-port", 554, "Set RTSP listen port")
			return nil
		}

		// load the config in the following order:
		// load from init then
		// load from json then
		// load from environment vars then
		// load from command-line flag.
		if err := loader.Load(conf,
			loader.Func(initLoader),
			&loader.JSONLoader{Path: configPath, CreatedIfNonExsit: true},
			&loader.EnvLoader{Prefix: envPrefix},
			&loader.FlagLoader{}); err != nil {
			panic(err)
		}
		return conf
	}


When used to load configurations, loader takes care of simplicity
and compatibility, and take full advantage of the built-in flag package
with little additional constraints.
*/
package loader

// Loader loads the value from a source.
type Loader interface {
	// Load loads the value and stores it in the value pointed to by v.
	Load(v interface{}) error
}

// Func is an adapter to allow the use of ordinary functions as Loader.
type Func func(v interface{}) error

// Load loads the value and stores it in the value pointed to by v.
func (l Func) Load(v interface{}) error {
	return l(v)
}

// Load loads the values from the loader in order
// and store them in the values pointed to by v.
func Load(v interface{}, loaders ...Loader) error {
	for _, loader := range loaders {
		if err := loader.Load(v); err != nil {
			return err
		}
	}
	return nil
}
