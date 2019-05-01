/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package environment

import (
	"io"
	"os"
)

// DefaultStreams contains the default io streams
var DefaultStreams = Streams{
	In:  os.Stdin,
	Out: os.Stdout,
	Err: os.Stderr,
}

// Streams contains io stream configuration
type Streams struct {
	In  io.Reader
	Out io.Writer
	Err io.Writer
}
