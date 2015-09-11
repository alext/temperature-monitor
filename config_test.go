package main

import (
	"encoding/json"

	"github.com/alext/afero"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Parsing the config file", func() {
	BeforeEach(func() {
		fs = &afero.MemMapFs{}
	})

	Context("with a config file", func() {

		It("should set the port", func() {
			createConfigFile("/etc/config.json", configData{
				"port": 1234,
			})

			config, err := loadConfig("/etc/config.json")
			Expect(err).NotTo(HaveOccurred())
			Expect(config.Port).To(Equal(1234))
		})

		It("should set a default port if none is specified", func() {
			createConfigFile("/etc/config.json", configData{})

			config, err := loadConfig("/etc/config.json")
			Expect(err).NotTo(HaveOccurred())
			Expect(config.Port).To(Equal(defaultPort))
		})

		It("should setup the sensor details", func() {
			createConfigFile("/etc/config.json", configData{
				"sensors": map[string]map[string]interface{}{
					"foo": map[string]interface{}{
						"id": "28-12345678",
					},
					"bar": map[string]interface{}{
						"id": "28-87654321",
					},
				},
			})

			config, err := loadConfig("/etc/config.json")
			Expect(err).NotTo(HaveOccurred())
			Expect(config.Sensors).To(HaveLen(2))

			Expect(config.Sensors["foo"].ID).To(Equal("28-12345678"))
			Expect(config.Sensors["bar"].ID).To(Equal("28-87654321"))
		})

		It("should have an empty list of sensors if none given", func() {
			createConfigFile("/etc/config.json", configData{})

			config, err := loadConfig("/etc/config.json")
			Expect(err).NotTo(HaveOccurred())
			Expect(config.Sensors).To(HaveLen(0))
		})
	})

	Context("when the config file doesn't exist", func() {
		It("should set a default port", func() {
			config, err := loadConfig("/etc/config.json")
			Expect(err).NotTo(HaveOccurred())
			Expect(config.Port).To(Equal(defaultPort))
		})

		It("should set an empty list of sensors", func() {
			config, err := loadConfig("/etc/config.json")
			Expect(err).NotTo(HaveOccurred())
			Expect(config.Sensors).To(HaveLen(0))
		})
	})
})

type configData map[string]interface{}

func createConfigFile(filename string, data configData) {
	file, err := fs.Create(filename)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	defer file.Close()
	err = json.NewEncoder(file).Encode(data)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
}
