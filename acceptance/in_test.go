package acceptance_test

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/onsi/gomega/gexec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const lightStemcellRequest = `
{
	"source": {
		"name": "bosh-aws-xen-hvm-ubuntu-trusty-go_agent"
	},
	"params": {
		"tarball": false
	},
	"version": {
		"version": "3262.4"
	}
}`

const regularStemcellRequest = `
{
	"source": {
		"name": "bosh-azure-hyperv-ubuntu-trusty-go_agent"
	},
	"version": {
		"version": "3262.9"
	}
}`

var _ = Describe("in", func() {
	Context("when a light stemcell is requested", func() {
		var (
			command    *exec.Cmd
			contentDir string
		)

		BeforeEach(func() {
			var err error
			contentDir, err = ioutil.TempDir("", "")
			Expect(err).NotTo(HaveOccurred())

			command = exec.Command(boshioIn, contentDir)
			command.Stdin = bytes.NewBufferString(lightStemcellRequest)
		})

		AfterEach(func() {
			err := os.RemoveAll(contentDir)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when no tarball is requested", func() {
			It("writes just the metadata", func() {
				session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session, "30s").Should(gexec.Exit(0))

				version, err := ioutil.ReadFile(filepath.Join(contentDir, "version"))
				Expect(err).NotTo(HaveOccurred())
				Expect(string(version)).To(Equal("3262.4"))

				url, err := ioutil.ReadFile(filepath.Join(contentDir, "url"))
				Expect(err).NotTo(HaveOccurred())
				Expect(string(url)).To(Equal("https://d26ekeud912fhb.cloudfront.net/bosh-stemcell/aws/light-bosh-stemcell-3262.4-aws-xen-hvm-ubuntu-trusty-go_agent.tgz"))

				checksum, err := ioutil.ReadFile(filepath.Join(contentDir, "sha1"))
				Expect(err).NotTo(HaveOccurred())
				Expect(string(checksum)).To(Equal("58b80c916ad523defea9e661045b7fc700a9ec4f"))
			})
		})
	})

	Context("when a regular stemcell is requested", func() {
		var (
			command    *exec.Cmd
			contentDir string
		)

		BeforeEach(func() {
			var err error
			contentDir, err = ioutil.TempDir("", "")
			Expect(err).NotTo(HaveOccurred())

			command = exec.Command(boshioIn, contentDir)
			command.Stdin = bytes.NewBufferString(regularStemcellRequest)
		})

		AfterEach(func() {
			err := os.RemoveAll(contentDir)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when the tarball is requested", func() {
			It("downloads the stemcell with metadata", func() {
				session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session, "30s").Should(gexec.Exit(0))
				tarballBytes, err := ioutil.ReadFile(filepath.Join(contentDir, "stemcell.tgz"))
				Expect(err).NotTo(HaveOccurred())

				checksum, err := ioutil.ReadFile(filepath.Join(contentDir, "sha1"))
				Expect(err).NotTo(HaveOccurred())
				Expect(string(checksum)).To(Equal(fmt.Sprintf("%x", sha1.Sum(tarballBytes))))
			})
		})
	})
})
