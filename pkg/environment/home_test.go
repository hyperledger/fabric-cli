/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package environment_test

import (
	"os"
	"path/filepath"

	"github.com/hyperledger/fabric-cli/pkg/environment"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Home", func() {
	var (
		home environment.Home
	)

	BeforeEach(func() {
		home = environment.Home(TestPath)
	})

	Describe("String", func() {
		var path string

		JustBeforeEach(func() {
			path = home.String()
		})

		It("should equal test path", func() {
			Expect(path).To(Equal(TestPath))
		})
	})

	Describe("Path", func() {
		var path string

		JustBeforeEach(func() {
			path = home.Path("test")
		})

		It("should append to home path", func() {
			Expect(path).To(Equal(filepath.Join(TestPath, "test")))
		})
	})

	Describe("Plugins", func() {
		var path string

		JustBeforeEach(func() {
			path = home.Plugins()
		})

		It("should return plugins path", func() {
			Expect(path).To(Equal(filepath.Join(TestPath, "plugins")))
		})
	})

	Describe("Init", func() {
		var err error

		JustBeforeEach(func() {
			err = home.Init()
		})

		JustAfterEach(func() {
			os.RemoveAll(TestPath)
		})

		It("should initialize home path", func() {
			Expect(err).To(BeNil())
		})
	})
})
