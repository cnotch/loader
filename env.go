// Copyright (c) 2019,CAO HONGJU. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package loader

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// EnvLoader implements the Loader interface.
// It loads the value from the environment variables.
type EnvLoader struct {
	// Prefix  prepends given string to every environment variable
	Prefix string
	// Exclude excluded names in FlagSet
	Exclude []string
	// FlagSet it defines the environment variables name source.
	// If FlagSet is nil, use flag.CommandLine
	FlagSet *flag.FlagSet
}

// Load loads the value from the environment variables and stores it in the value pointed to by v.
func (env *EnvLoader) Load(v interface{}) (err error) {
	fs := env.FlagSet
	if fs == nil {
		fs = flag.CommandLine
	}

	envKeys := []string{}
	prefix := env.Prefix

	fs.VisitAll(func(f *flag.Flag) {
		if env.isExclude(f.Name) {
			return
		}

		// get environment variable name
		key := strings.ToUpper(strings.Replace(prefix+"_"+f.Name, "-", "_", -1))
		envKeys = append(envKeys, key)
		if value, ok := os.LookupEnv(key); ok {
			if err2 := f.Value.Set(value); err2 != nil {
				err = err2
			}
		}
	})

	// set the Usage method to add the environment variables section
	fs.Usage = func() {
		// output default
		_, proc := filepath.Split(os.Args[0])
		fmt.Fprintf(fs.Output(), "Usage of %s:\n", proc)
		fs.PrintDefaults()

		// output environment variables
		fmt.Fprintf(fs.Output(), "\nEnvironment variables:\n")
		for _, key := range envKeys {
			fmt.Fprintf(fs.Output(), "  %s\n", key)
		}
	}

	return
}

func (env *EnvLoader) isExclude(name string) bool {
	for _, v := range env.Exclude {
		if v == name {
			return true
		}
	}
	return false
}
