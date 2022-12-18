package url_problem

import (
	"bytes"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
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
			// Fails because of url3.Path == "ping" and urlParsed.Path == "/ping"
			Expect(url3).To(Equal(urlParsed))
		})
	})

	When("parse host with tailing / and without it", func() {
		url1, _ := url.Parse("http://localhost")
		url2, _ := url.Parse("http://localhost/")
		It("should be consistent", func() {
			// Fails because of url1.Path == "" and url2.Path == "/"
			Expect(url1).To(Equal(url2))
		})
	})

	When("check request is correct by conventional Gin methods", func() {
		url1, _ := url.Parse("http://localhost")
		url2, _ := url.Parse("ping")

		w := httptest.NewRecorder()

		req, _ := http.NewRequest("GET", "/", nil)
		req.URL = url1.JoinPath(url2.String())

		r := gin.New()
		r.GET("ping", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		r.ServeHTTP(w, req)

		It("should work perfectly", func() {
			// Fails because handler expects "/ping" and gets "ping" so 404 error will be here.
			Expect(w.Code).To(Equal(http.StatusOK))
		})
	})

	When("check request provides correct http request", func() {
		req, _ := http.NewRequest("GET", "/", nil)
		url1, _ := url.Parse("http://localhost")
		url2, _ := url.Parse("/ping")

		req.URL = url1.JoinPath(url2.String())

		buf := new(bytes.Buffer)

		err := req.Write(buf)
		It("should be correct request", func() {
			Expect(err).To(BeNil())
			/*
					Gets:

					GET ping HTTP/1.1
				    Host: localhost
				    User-Agent: Go-http-client/1.1

					instead of

					GET /ping HTTP/1.1
					Host: localhost
					User-Agent: Go-http-client/1.1

					Path "ping" causes server to answer 400 error because url is incorrect.
			*/
			Expect(buf.String()).To(ContainSubstring("/ping"))
		})
	})
})
