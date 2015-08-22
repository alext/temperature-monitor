package webserver_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	. "github.com/alext/temperature-monitor/Godeps/_workspace/src/github.com/onsi/ginkgo"
	. "github.com/alext/temperature-monitor/Godeps/_workspace/src/github.com/onsi/gomega"

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

		It("returns the sensor data as JSON", func() {
			updateTime := time.Now().Add(-3 * time.Minute)
			sensor.SetTemperature(15643, updateTime)

			resp := doGetRequest(server, "/sensors/one")
			Expect(resp.Code).To(Equal(http.StatusOK))
			Expect(resp.Header().Get("Content-Type")).To(Equal("application/json"))

			data := decodeJsonResponse(resp)
			Expect(data["temperature"]).To(BeEquivalentTo(15643))

			updateTimeStr, _ := updateTime.MarshalText()
			Expect(data["updated_at"]).To(Equal(string(updateTimeStr)))
		})

		It("returns 404 for a non-existent sensor", func() {
			resp := doGetRequest(server, "/sensors/non-existent")
			Expect(resp.Code).To(Equal(http.StatusNotFound))
		})
	})

	Describe("sensors index", func() {
		var (
			s1 *dummySensor
			s2 *dummySensor
		)

		BeforeEach(func() {
			s1 = &dummySensor{}
			s1.SetTemperature(18345, time.Now())
			s2 = &dummySensor{}
			s2.SetTemperature(19542, time.Now())
			server.AddSensor("one", s1)
			server.AddSensor("two", s2)
		})

		It("returns details of all sensors", func() {
			resp := doGetRequest(server, "/sensors")
			Expect(resp.Code).To(Equal(http.StatusOK))
			Expect(resp.Header().Get("Content-Type")).To(Equal("application/json"))

			data := decodeJsonResponse(resp)
			Expect(data).To(HaveKey("one"))
			Expect(data).To(HaveKey("two"))
			data1 := data["one"].(map[string]interface{})
			Expect(data1["temperature"]).To(BeEquivalentTo(18345))
			data2 := data["two"].(map[string]interface{})
			Expect(data2["temperature"]).To(BeEquivalentTo(19542))
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
