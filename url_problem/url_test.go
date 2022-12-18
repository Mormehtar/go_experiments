package url_problem

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"net/url"
	"testing"
)

func TestUrlBug(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Url bug suite")
}

var _ = Describe("Check url problem", func() {
	When("join two urls", func() {
		url1, _ := url.Parse("http://localhost")
		url2, _ := url.Parse("ping")
		url3 := url1.JoinPath(url2.String())
		It("should be consistent", func() {
			urlParsed, _ := url.Parse(url3.String())
			Expect(url3).To(Equal(urlParsed))
		})
	})

	When("parse host with tailing / and without it", func() {
		url1, _ := url.Parse("http://localhost")
		url2, _ := url.Parse("http://localhost/")
		It("should be consistent", func() {
			Expect(url1).To(Equal(url2))
		})
	})
})
