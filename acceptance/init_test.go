package acceptance_test

import (
	"os"
	"strings"
	"testing"

	"github.com/onsi/gomega/gexec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	boshioCheck string
	boshioIn    string
)

var _ = BeforeSuite(func() {
	var err error

	Expect(
		strings.Contains(os.Getenv("GOFLAGS"), "-mod=vendor"),
	).To(BeTrue(), "GOFLAGS must include `-mod-vendor`")

	if _, err = os.Stat("/opt/resource/check"); err == nil {
		boshioCheck = "/opt/resource/check"
	} else {
		boshioCheck, err = gexec.Build("../cmd/check")
		Expect(err).NotTo(HaveOccurred())
	}

	if _, err = os.Stat("/opt/resource/in"); err == nil {
		boshioIn = "/opt/resource/in"
	} else {
		boshioIn, err = gexec.Build("../cmd/in")
		Expect(err).NotTo(HaveOccurred())
	}
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})

func TestAcceptance(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Acceptance Suite")
}
