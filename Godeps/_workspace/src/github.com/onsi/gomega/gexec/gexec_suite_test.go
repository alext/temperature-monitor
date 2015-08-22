package gexec_test

import (
	. "github.com/alext/temperature-monitor/Godeps/_workspace/src/github.com/onsi/ginkgo"
	. "github.com/alext/temperature-monitor/Godeps/_workspace/src/github.com/onsi/gomega"
	"github.com/alext/temperature-monitor/Godeps/_workspace/src/github.com/onsi/gomega/gexec"

	"testing"
)

var fireflyPath string

func TestGexec(t *testing.T) {
	BeforeSuite(func() {
		var err error
		fireflyPath, err = gexec.Build("./_fixture/firefly")
		Ω(err).ShouldNot(HaveOccurred())
	})

	AfterSuite(func() {
		gexec.CleanupBuildArtifacts()
	})

	RegisterFailHandler(Fail)
	RunSpecs(t, "Gexec Suite")
}
