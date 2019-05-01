/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package environment

import (
	"os"
	"path/filepath"
)

// DefaultHome is the default cli home directory
const DefaultHome = Home("${HOME}/.fabric")

// Home is the location of the configuration files
// By default, the files are stored in ~/.fabric
type Home string

// String resolves variables and returns home as a string
func (h Home) String() string {
	return os.ExpandEnv(string(h))
}

// Path appends compoenents to home and returns it as a string
func (h Home) Path(components ...string) string {
	return filepath.Join(append([]string{h.String()}, components...)...)
}

// Plugins returns the path to the plugins directory
func (h Home) Plugins() string {
	return h.Path("plugins")
}

// Init creates a home directory if it does not already exist
func (h Home) Init() error {
	// continue if the home directory already exists
	if _, err := os.Stat(h.String()); os.IsNotExist(err) {
		// create a new home directory
		if err = os.MkdirAll(h.String(), 0755); err != nil {
			return err
		}
	}

	return nil
}
