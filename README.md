# loader
loader provides easy to use object loading capability.

## Motivation

Sometimes, we need to load the same object from multiple sources in order of priority.

For example, the configuration data of the application:
 1. Built-in default value
 2. Load from the configuration file
 3. Load from environment variables
 4. Parse from the command line parameter
 5. Other sources (centralized configuration management, etc.)


loader provides a surprisingly simple pattern.

## Installing

1. Get package:

	```Shell
	go get -u github.com/cnotch/loader
	```

2. Import it in your code:

	```Go
	import "github.com/cnotch/loader"
	```

## Usage

For example, load configuration, which typically consists of the following steps:

1. Define configuration structure

``` go
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
```

2. Load by priority

```go
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
```

When used to load configurations, loader takes care of simplicity and compatibility, and take full advantage of the built-in flag package with little additional constraints.
