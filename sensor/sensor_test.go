package sensor

import (
	"os"
	"testing"
	"time"

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

type dummyTicker struct {
	duration time.Duration
	C        chan time.Time
	notify   chan struct{}
	stopped  bool
}

func (t dummyTicker) Channel() <-chan time.Time {
	select {
	case t.notify <- struct{}{}:
	default:
	}
	return t.C
}
func (t *dummyTicker) Stop() {
	t.stopped = true
}

var _ = Describe("a sensor", func() {
	var (
		tkr *dummyTicker
	)

	BeforeEach(func() {
		fs = &afero.MemMapFs{}

		newTicker = func(d time.Duration) ticker {
			tkr = &dummyTicker{
				duration: d,
				C:        make(chan time.Time, 1),
				notify:   make(chan struct{}, 1),
			}
			return tkr
		}
	})

	Describe("constructing a sensor", func() {
		var (
			sensor Sensor
			err    error
		)

		BeforeEach(func() {
			populateValueFile(testDeviceID, sampleData1)
			sensor, err = New(testDeviceID)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should read the initial temperature", func() {
			Expect(sensor.Temperature()).To(Equal(19437))
		})

		It("should start a ticker to poll the temperature every minute", func(done Done) {
			<-tkr.notify
			Expect(tkr.duration).To(Equal(time.Minute))

			close(done)
		})

		It("should update the temperature on each tick", func(done Done) {
			<-tkr.notify
			Expect(sensor.Temperature()).To(Equal(19437))

			populateValueFile(testDeviceID, sampleData2)
			tkr.C <- time.Now()
			<-tkr.notify
			Expect(sensor.Temperature()).To(Equal(18062))

			populateValueFile(testDeviceID, sampleData1)
			tkr.C <- time.Now()
			<-tkr.notify
			Expect(sensor.Temperature()).To(Equal(19437))

			close(done)
		})
	})

	Describe("closing a sensor", func() {
		var (
			sensor Sensor
			err    error
		)

		BeforeEach(func() {
			populateValueFile(testDeviceID, sampleData1)
			sensor, err = New(testDeviceID)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should stop the ticker", func(done Done) {
			<-tkr.notify

			sensor.Close()
			Expect(tkr.stopped).To(BeTrue())

			close(done)
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
