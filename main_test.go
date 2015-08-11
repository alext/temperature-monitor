package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/alext/afero"
)

func TestMain(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Main")
}

var _ = BeforeSuite(func() {
	log.SetOutput(ioutil.Discard)
})

const testLogFile = "/var/log/temperature-monitor"

var _ = Describe("setting up logging", func() {

	BeforeEach(func() {
		fs = &afero.MemMapFs{}
	})

	Describe("logging to stdout/stderr", func() {
		var (
			realStdout, realStderr *os.File
		)
		BeforeEach(func() {
			realStdout = os.Stdout
			realStderr = os.Stderr
		})
		AfterEach(func() {
			os.Stdout = realStdout
			os.Stderr = realStderr
		})

		It("logs to stdout when requested", func(done Done) {
			writer, outputC := setupPipe()
			os.Stdout = writer

			Expect(setupLogging("STDOUT")).To(Succeed())

			log.Print("test log entry")
			writer.Close()
			Expect(<-outputC).To(ContainSubstring("test log entry"))
			close(done)
		})

		It("logs to stderr when requested", func(done Done) {
			writer, outputC := setupPipe()
			os.Stderr = writer

			Expect(setupLogging("STDERR")).To(Succeed())

			log.Print("test log entry")
			writer.Close()
			Expect(<-outputC).To(ContainSubstring("test log entry"))
			close(done)
		})
	})

	Describe("logging to a file", func() {
		It("creates the given logfile if necessary", func() {
			Expect(setupLogging(testLogFile)).To(Succeed())

			log.Print("test log entry")
			Expect(testLogContents()).To(ContainSubstring("test log entry"))

		})

		It("appends to an existing logfile", func() {
			file, err := fs.Create(testLogFile)
			Expect(err).NotTo(HaveOccurred())
			_, err = fmt.Fprintln(file, "existing log line")
			Expect(err).NotTo(HaveOccurred())
			err = file.Close()
			Expect(err).NotTo(HaveOccurred())

			Expect(setupLogging(testLogFile)).To(Succeed())
			log.Print("test log entry")
			logContents := testLogContents()
			Expect(logContents).To(ContainSubstring("existing log line"))
			Expect(logContents).To(ContainSubstring("test log entry"))
		})
	})
})

func testLogContents() string {
	file, err := fs.Open(testLogFile)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	contents, err := ioutil.ReadAll(file)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	return string(contents)
}

func setupPipe() (*os.File, <-chan string) {
	r, w, err := os.Pipe()
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	outC := make(chan string)
	go func() {
		defer GinkgoRecover()
		buf := bytes.Buffer{}
		_, err := io.Copy(&buf, r)
		ExpectWithOffset(1, err).NotTo(HaveOccurred())
		outC <- buf.String()
	}()
	return w, outC
}
