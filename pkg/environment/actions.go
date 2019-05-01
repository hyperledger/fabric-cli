/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package environment

// Action is used to modify config
type Action func(c *Config)

// SetCurrentContext updates the current context
func SetCurrentContext(context string) Action {
	return func(c *Config) {
		c.CurrentContext = context
	}
}

// SetContext adds or updates the specified context
func SetContext(name string, context *Context) Action {
	return func(c *Config) {
		current, ok := c.Contexts[name]
		if !ok {
			c.Contexts[name] = context
			return
		}

		// override existing keys if context already exists

		if len(context.Network) > 0 {
			current.Network = context.Network
		}

		if len(context.Organization) > 0 {
			current.Organization = context.Organization
		}

		if len(context.User) > 0 {
			current.User = context.User
		}

		if len(context.Channel) > 0 {
			current.Channel = context.Channel
		}

		if len(context.Orderers) > 0 {
			current.Orderers = context.Orderers
		}

		if len(context.Peers) > 0 {
			current.Peers = context.Peers
		}
	}
}

// DeleteContext deletes a specified context
func DeleteContext(name string) Action {
	return func(c *Config) {
		delete(c.Contexts, name)
	}
}

// SetNetwork adds or updates the specified network
func SetNetwork(name string, network *Network) Action {
	return func(c *Config) {
		c.Networks[name] = network
	}
}

// DeleteNetwork deletes the specified network
func DeleteNetwork(name string) Action {
	return func(c *Config) {
		delete(c.Networks, name)
	}
}
