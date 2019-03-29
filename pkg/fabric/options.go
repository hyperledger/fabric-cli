/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package fabric

// Option provides the ability to override client defaults
type Option func(*client)

// WithFactory overrides the default client factory
func WithFactory(f Factory) Option {
	return func(c *client) {
		c.factory = f
	}
}
