/*
Copyright State Street Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package templates_test

import (
	"html/template"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/hyperledger/fabric-cli/pkg/fabric/templates"
)

func TestTemplates(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Templates Suite")
}

var _ = Describe("Templates", func() {
	var (
		text string

		tmpl *template.Template
		err  error
	)

	JustBeforeEach(func() {
		tmpl, err = template.New("test").Parse(text)
	})

	Describe("Config", func() {
		BeforeEach(func() {
			text = templates.Config
		})

		It("should parse successfully", func() {
			Expect(err).To(BeNil())
			Expect(tmpl).NotTo(BeNil())
		})
	})
})
