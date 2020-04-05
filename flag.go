// Copyright (c) 2019,CAO HONGJU. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package loader

import (
	"flag"
	"os"
)

// FlagLoader implements the Loader interface.
// It parses the command-line flags from os.Args[1:].
// Must be called after all flags are defined and
// before flags are accessed by the program.
type FlagLoader struct {
	// FlagSet is a FlagSet used to parse the command line
	// If FlagSet is nil, use flag.CommandLine
	FlagSet *flag.FlagSet
}

// Load loads the value from command line and stores it in the value pointed to by v.
func (fl *FlagLoader) Load(v interface{}) (err error) {
	if fl.FlagSet == nil {
		return flag.CommandLine.Parse(os.Args[1:])
	}
	return fl.FlagSet.Parse(os.Args[1:])
}
