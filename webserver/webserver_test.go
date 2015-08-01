package webserver_test

import (
	"encoding/json"
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
		It("resturns a 200", func() {
			resp := doGetRequest(server, "/")
			Expect(resp.Code).To(Equal(http.StatusOK))

		})
	})

	Describe("reading a sensor", func() {
		var (
			sensor *dummySensor
		)

		BeforeEach(func() {
			sensor = &dummySensor{}
			server.AddSensor("one", sensor)
		})

		It("returns the sensor value as JSON", func() {
			sensor.SetTemperature(15.643)
			resp := doGetRequest(server, "/sensors/one")
			Expect(resp.Code).To(Equal(http.StatusOK))
			Expect(resp.Header().Get("Content-Type")).To(Equal("application/json"))

			data := decodeJsonResponse(resp)
			Expect(data["temperature"]).To(BeNumerically("~", 15.643))
		})

		It("returns 404 for a non-existent sensor", func() {
			resp := doGetRequest(server, "/sensors/non-existent")
			Expect(resp.Code).To(Equal(http.StatusNotFound))
		})
	})
})

func doGetRequest(server http.Handler, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("GET", "http://example.com"+path, nil)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	return w
}

func decodeJsonResponse(w *httptest.ResponseRecorder) map[string]interface{} {
	var data map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &data)
	ExpectWithOffset(1, err).To(BeNil())
	return data
}
