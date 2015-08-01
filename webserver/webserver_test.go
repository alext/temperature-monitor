package webserver_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/alext/temperature-monitor/webserver"
)

func TestWebServer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Webserver")
}

var _ = Describe("the webserver", func() {
	var (
		server *webserver.Webserver
	)

	BeforeEach(func() {
		server = webserver.New(8080)
	})

	Describe("the root page", func() {
		It("should return a 200", func() {
			resp := doGetRequest(server, "/")
			Expect(resp.Code).To(Equal(http.StatusOK))

		})
	})
})

func doGetRequest(server http.Handler, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("GET", "http://example.com"+path, nil)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	return w
}
