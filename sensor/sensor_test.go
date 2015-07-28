package sensor

import (
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
)

func TestSensor(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Sensor")
}

const testDeviceID = "28-0123456789ab"

const sampleData1 = `37 01 4b 46 7f ff 09 10 26 : crc=26 YES
37 01 4b 46 7f ff 09 10 26 t=19437
`
const sampleData2 = `21 01 4b 46 7f ff 0f 10 4b : crc=4b YES
21 01 4b 46 7f ff 0f 10 4b t=18062
`

var _ = Describe("a sensor", func() {

	BeforeEach(func() {
		fs = &afero.MemMapFs{}
	})

	Describe("reading the temperature", func() {
		var sensor Sensor

		BeforeEach(func() {
			populateValueFile(testDeviceID, sampleData1)
			var err error
			sensor, err = New(testDeviceID)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should return the temperature", func() {
			Expect(sensor.Temperature()).To(Equal(19.437))
		})

		It("should allow multiple reads", func() {
			sensor.Temperature()
			Expect(sensor.Temperature()).To(Equal(19.437))
		})

		It("should handle changed file contents", func() {
			sensor.Temperature()
			populateValueFile(testDeviceID, sampleData2)
			Expect(sensor.Temperature()).To(Equal(18.062))
		})
	})
})

func populateValueFile(deviceID, contents string) {
	valueFilePath := w1DevicesPath + deviceID + "/w1_slave"
	file, err := fs.OpenFile(valueFilePath, os.O_RDWR, 0644)
	if err != nil {
		if os.IsNotExist(err) {
			file, err = fs.Create(valueFilePath)
		}
		ExpectWithOffset(1, err).NotTo(HaveOccurred())
	}
	_, err = file.Write([]byte(contents))
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
}
