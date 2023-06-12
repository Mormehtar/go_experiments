package utils

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestUtils(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Test utils test")
}

var _ = Describe("ChangeDBName", func() {
	When("there is no database in the url", func() {
		It("should add database to the url", func() {
			Expect(ChangeDBName("blah=blah port=port", "new_name")).
				To(Equal("blah=blah port=port database=new_name"))
		})
	})
	When("dsn is empty", func() {
		It("should add database to the url", func() {
			Expect(ChangeDBName("", "new_name")).To(Equal("database=new_name"))
		})
	})
	When("there is a database in the url", func() {
		It("should change database in the url", func() {
			Expect(ChangeDBName("blah=blah port=port database=old_name", "new_name")).
				To(Equal("blah=blah port=port database=new_name"))
		})
	})
	When("there is a database in the url and it is not the last parameter", func() {
		It("should change database in the url", func() {
			Expect(ChangeDBName("blah=blah database=old_name port=port", "new_name")).
				To(Equal("blah=blah database=new_name port=port"))
		})
	})
})
