package webserver_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/alext/temperature-monitor/sensor"
	"github.com/alext/temperature-monitor/webserver"
)

func TestWebServer(t *testing.T) {
	RegisterFailHandler(Fail)
	log.SetOutput(ioutil.Discard)
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
			sensor.temp, sensor.updateTime = 15643, updateTime

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

	Describe("setting a sensor", func() {
		var (
			s1 sensor.SettableSensor
		)

		BeforeEach(func() {
			s1 = sensor.NewPushSensor("something")
			s1.Set(12345, time.Now().Add(-1*time.Hour))
			server.AddSensor("one", s1)
		})

		It("updates the sensor with the given details", func() {
			data := map[string]interface{}{
				"temperature": 15643,
			}
			resp := doJSONPutRequest(server, "/sensors/one", data)
			Expect(resp.Code).To(Equal(http.StatusOK))
			temp, updated := s1.Read()
			Expect(temp).To(Equal(15643))
			Expect(updated).To(BeTemporally("~", time.Now(), 100*time.Millisecond))

			respData := decodeJsonResponse(resp)
			Expect(respData["temperature"]).To(BeEquivalentTo(15643))
		})

		It("returns a 400 for invalid data", func() {
			data := map[string]interface{}{
				"foo": "bar",
			}
			resp := doJSONPutRequest(server, "/sensors/one", data)
			Expect(resp.Code).To(Equal(http.StatusBadRequest))
			temp, _ := s1.Read()
			Expect(temp).To(Equal(12345))
		})

		It("returns 405 for a non-writable sensor", func() {
			s2 := &dummySensor{}
			s2.temp, s2.updateTime = 12345, time.Now().Add(-1*time.Hour)
			server.AddSensor("two", s2)

			data := map[string]interface{}{
				"temperature": 15643,
			}
			resp := doJSONPutRequest(server, "/sensors/two", data)
			Expect(resp.Code).To(Equal(405))
			temp, _ := s1.Read()
			Expect(temp).To(Equal(12345))
		})
	})

	Describe("sensors index", func() {
		var (
			s1 *dummySensor
			s2 *dummySensor
		)

		BeforeEach(func() {
			s1 = &dummySensor{}
			s1.temp, s1.updateTime = 18345, time.Now()
			s2 = &dummySensor{}
			s2.temp, s2.updateTime = 19542, time.Now()
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
	req := httptest.NewRequest("GET", "http://example.com"+path, nil)
	w := httptest.NewRecorder()
	server.ServeHTTP(w, req)
	return w
}

func doJSONPutRequest(server http.Handler, path string, bodyData interface{}) *httptest.ResponseRecorder {
	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(bodyData)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	req := httptest.NewRequest("PUT", "http://example.com"+path, &body)
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
