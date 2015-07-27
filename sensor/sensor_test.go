package sensor

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
)

func TestSensor(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Sensor")
}

const testDeviceId = "28-0123456789ab"
const testValueFile = "/sys/bus/w1/devices/" + testDeviceId + "/w1_slave"
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
		It("should return the temperature", func() {
			file, err := fs.Create(testValueFile)
			Expect(err).NotTo(HaveOccurred())
			_, err = file.Write([]byte(sampleData1))
			Expect(err).NotTo(HaveOccurred())

			tfile, err := fs.Open(testValueFile)
			Expect(err).NotTo(HaveOccurred())
			b := make([]byte, 100)
			_, err = tfile.Read(b)
			Expect(err).NotTo(HaveOccurred())

			Expect(b).To(ContainSubstring("t=19437"))

		})
	})

	It("should pass", func() {
		Expect(true).To(BeTrue())
	})
})
