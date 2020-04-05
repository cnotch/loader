// Copyright (c) 2019,CAO HONGJU. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package loader

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
)

// JSONLoader implements the Loader interface.
// It loads the value from a json source.
type JSONLoader struct {
	// Path is json file source's path
	Path string
	// Reader is json reader source, priority is higher than Path
	Reader io.Reader
	// CreatedIfNonExsit indicates whether the json file is created automatically if it does not exist
	CreatedIfNonExsit bool
}

// Load loads the value from json and stores it in the value pointed to by v.
func (j *JSONLoader) Load(v interface{}) error {
	if j.Reader != nil {
		return json.NewDecoder(j.Reader).Decode(v)
	}

	path := j.Path

	// if it does not exist, created automatically
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if j.CreatedIfNonExsit {
			if err := j.createConfigFile(path, v); err != nil {
				return err
			}
		}
	} else {
		b, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(b, v); err != nil {
			return err
		}
	}

	return nil
}

func (j *JSONLoader) createConfigFile(path string, v interface{}) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}

	defer f.Close()

	var formatted bytes.Buffer
	body, err := json.Marshal(v)
	if err != nil {
		return err
	}

	if err := json.Indent(&formatted, body, "", "\t"); err != nil {
		return err
	}

	if _, err := f.Write(formatted.Bytes()); err != nil {
		return err
	}
	if err := f.Sync(); err != nil {
		return err
	}

	return nil
}
