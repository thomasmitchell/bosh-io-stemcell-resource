package acceptance_test

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/onsi/gomega/gbytes"
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

const bothTypesStemcellRequest = `
{
	"source": {
		"name": "bosh-aws-xen-ubuntu-trusty-go_agent"
	},
	"version": {
		"version": "3262.4"
	}
}`

const bothTypesForceRegularStemcellRequest = `
{
	"source": {
		"name": "bosh-aws-xen-ubuntu-trusty-go_agent",
		"force_regular": true
	},
	"version": {
		"version": "3262.4"
	}
}`

const stemcellRequestWithFileName = `
{
	"source": {
		"name": "bosh-aws-xen-hvm-ubuntu-trusty-go_agent"
	},
	"params": {
		"preserve_filename": true
	},
	"version": {
		"version": "3262.12"
	}
}`

const invalidRequestVersion = `
{
	"source": {
		"name": "bosh-aws-xen-hvm-ubuntu-trusty-go_agent"
	},
	"params": {
		"preserve_filename": true
	},
	"version": {
		"version": "AAAAA"
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

				<-session.Exited
				Expect(session.ExitCode()).To(Equal(0))
				Expect(session.Out).To(gbytes.Say(`{"version":{"version":"3262.4"},"metadata":\[{"name":"url","value":"https://d26ekeud912fhb.cloudfront.net/bosh-stemcell/aws/light-bosh-stemcell-3262.4-aws-xen-hvm-ubuntu-trusty-go_agent.tgz"},{"name":"sha1","value":"58b80c916ad523defea9e661045b7fc700a9ec4f"}\]}`))

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

				<-session.Exited
				Expect(session.ExitCode()).To(Equal(0))

				tarballBytes, err := ioutil.ReadFile(filepath.Join(contentDir, "stemcell.tgz"))
				Expect(err).NotTo(HaveOccurred())

				checksum, err := ioutil.ReadFile(filepath.Join(contentDir, "sha1"))
				Expect(err).NotTo(HaveOccurred())

				Expect(session.Out).To(gbytes.Say(fmt.Sprintf(`{"version":{"version":"3262.9"},"metadata":\[{"name":"url","value":"https://s3.amazonaws.com/bosh-azure-stemcells/bosh-stemcell-3262.9-azure-hyperv-ubuntu-trusty-go_agent.tgz"},{"name":"sha1","value":"%s"}\]}`, string(checksum))))
				Expect(string(checksum)).To(Equal(fmt.Sprintf("%x", sha1.Sum(tarballBytes))))
			})
		})
	})

	Context("when a stemcell is requested that supports both light and regular", func() {
		var (
			command    *exec.Cmd
			contentDir string
		)

		BeforeEach(func() {
			var err error
			contentDir, err = ioutil.TempDir("", "")
			Expect(err).NotTo(HaveOccurred())

			command = exec.Command(boshioIn, contentDir)
			command.Stdin = bytes.NewBufferString(bothTypesStemcellRequest)
		})

		AfterEach(func() {
			err := os.RemoveAll(contentDir)
			Expect(err).NotTo(HaveOccurred())
		})

		It("downloads the light stemcell with metadata", func() {
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			<-session.Exited
			Expect(session.ExitCode()).To(Equal(0))

			tarballBytes, err := ioutil.ReadFile(filepath.Join(contentDir, "stemcell.tgz"))
			Expect(err).NotTo(HaveOccurred())

			checksum, err := ioutil.ReadFile(filepath.Join(contentDir, "sha1"))
			Expect(err).NotTo(HaveOccurred())
			Expect(string(checksum)).To(Equal(fmt.Sprintf("%x", sha1.Sum(tarballBytes))))

			urlBytes, err := ioutil.ReadFile(filepath.Join(contentDir, "url"))
			Expect(err).NotTo(HaveOccurred())
			Expect(string(urlBytes)).To(ContainSubstring("light"))
		})
	})

	Context("when a stemcell is requested that supports both light and regular and force_regular is true", func() {
		var (
			command    *exec.Cmd
			contentDir string
		)

		BeforeEach(func() {
			var err error
			contentDir, err = ioutil.TempDir("", "")
			Expect(err).NotTo(HaveOccurred())

			command = exec.Command(boshioIn, contentDir)
			command.Stdin = bytes.NewBufferString(bothTypesForceRegularStemcellRequest)
		})

		AfterEach(func() {
			err := os.RemoveAll(contentDir)
			Expect(err).NotTo(HaveOccurred())
		})

		It("downloads the regular stemcell with metadata", func() {
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			<-session.Exited
			Expect(session.ExitCode()).To(Equal(0))

			tarballBytes, err := ioutil.ReadFile(filepath.Join(contentDir, "stemcell.tgz"))
			Expect(err).NotTo(HaveOccurred())

			checksum, err := ioutil.ReadFile(filepath.Join(contentDir, "sha1"))
			Expect(err).NotTo(HaveOccurred())
			Expect(string(checksum)).To(Equal(fmt.Sprintf("%x", sha1.Sum(tarballBytes))))

			urlBytes, err := ioutil.ReadFile(filepath.Join(contentDir, "url"))
			Expect(err).NotTo(HaveOccurred())
			Expect(string(urlBytes)).NotTo(ContainSubstring("light"))
		})
	})

	Context("when a stemcell is requested with the original filename", func() {
		var (
			command    *exec.Cmd
			contentDir string
		)

		BeforeEach(func() {
			var err error
			contentDir, err = ioutil.TempDir("", "")
			Expect(err).NotTo(HaveOccurred())

			command = exec.Command(boshioIn, contentDir)
			command.Stdin = bytes.NewBufferString(stemcellRequestWithFileName)
		})

		AfterEach(func() {
			err := os.RemoveAll(contentDir)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when the tarball is requested", func() {
			It("saves the stemcell to the correct filename", func() {
				session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				<-session.Exited
				Expect(session.ExitCode()).To(Equal(0))

				tarballBytes, err := ioutil.ReadFile(filepath.Join(contentDir, "light-bosh-stemcell-3262.12-aws-xen-hvm-ubuntu-trusty-go_agent.tgz"))
				Expect(err).NotTo(HaveOccurred())

				checksum, err := ioutil.ReadFile(filepath.Join(contentDir, "sha1"))
				Expect(session.Out).To(gbytes.Say(fmt.Sprintf(`{"version":{"version":"3262.12"},"metadata":\[{"name":"url","value":"https://d26ekeud912fhb.cloudfront.net/bosh-stemcell/aws/light-bosh-stemcell-3262.12-aws-xen-hvm-ubuntu-trusty-go_agent.tgz"},{"name":"sha1","value":"%s"}\]}`, string(checksum))))

				Expect(err).NotTo(HaveOccurred())
				Expect(string(checksum)).To(Equal(fmt.Sprintf("%x", sha1.Sum(tarballBytes))))
			})
		})
	})

	Context("when an error occurs", func() {
		var (
			command    *exec.Cmd
			contentDir string
		)

		BeforeEach(func() {
			var err error
			contentDir, err = ioutil.TempDir("", "")
			Expect(err).NotTo(HaveOccurred())

			command = exec.Command(boshioIn, contentDir)
			command.Stdin = bytes.NewBufferString(stemcellRequestWithFileName)
		})

		AfterEach(func() {
			err := os.RemoveAll(contentDir)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when the provided location is not writeable", func() {
			BeforeEach(func() {
				err := os.RemoveAll(contentDir)
				Expect(err).NotTo(HaveOccurred())

				contentFile, err := os.Create(contentDir)
				Expect(err).NotTo(HaveOccurred())

				err = contentFile.Close()
				Expect(err).NotTo(HaveOccurred())

				// make it a file instead
				contentDir = contentFile.Name()
			})

			It("returns an error", func() {
				session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				<-session.Exited
				Expect(session.ExitCode()).To(Equal(1))
				Expect(session.Err).To(gbytes.Say("not a directory"))
			})
		})

		Context("when the request version does not exist", func() {
			It("returns an error", func() {
				command.Stdin = bytes.NewBufferString(invalidRequestVersion)
				session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				<-session.Exited
				Expect(session.ExitCode()).To(Equal(1))
				Expect(session.Err).To(gbytes.Say("failed to find stemcell matching version:"))
			})
		})

		Context("when the json provided is malformed", func() {
			It("returns an error", func() {
				command.Stdin = bytes.NewBufferString("%%%%%%%")
				session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				<-session.Exited
				Expect(session.ExitCode()).To(Equal(1))
				Expect(session.Err).To(gbytes.Say("invalid character"))
			})
		})
	})
})
