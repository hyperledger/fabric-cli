/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package environment_test

import (
	"os"

	"github.com/hyperledger/fabric-cli/pkg/environment"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

//go:generate counterfeiter -o mocks/config.go --fake-name DefaultConfig . Config

var _ = Describe("Config", func() {
	var (
		settings, file, result *environment.Settings
		err                    error
	)

	BeforeEach(func() {
		settings, err = environment.GetSettings()

		Expect(err).NotTo(HaveOccurred())

		settings.Home = environment.Home(os.TempDir())
	})

	JustBeforeEach(func() {
		err = settings.Home.Init()

		Expect(err).NotTo(HaveOccurred())
	})

	Context("when loading from file", func() {
		JustBeforeEach(func() {
			file, err = settings.FromFile()

			Expect(err).NotTo(HaveOccurred())
		})

		It("should not be nil", func() {
			Expect(file).NotTo(BeNil())
		})

		It("should be able to save", func() {
			file.ActiveProfile = "foo"

			Expect(file.Save()).Should(Succeed())

			result, err = settings.FromFile()

			Expect(result.ActiveProfile).To(Equal("foo"))
		})
	})
})
