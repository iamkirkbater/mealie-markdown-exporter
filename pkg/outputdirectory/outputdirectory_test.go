package outputdirectory_test

import (
	"github.com/iamkirkbater/mealie-markdown-exporter/pkg/outputdirectory"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
)

var _ = Describe("Resolve", func() {
	var fs afero.Fs

	BeforeEach(func() {
		fs = afero.NewMemMapFs()
	})

	Context("when the output directory is empty", func() {
		It("returns an error", func() {
			err := outputdirectory.Resolve(fs, "")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("output-dir is required"))
		})
	})

	Context("when using the default output directory", func() {
		It("creates the directory if it does not exist", func() {
			err := outputdirectory.Resolve(fs, outputdirectory.DefaultOutputDir)
			Expect(err).NotTo(HaveOccurred())

			exists, _ := afero.DirExists(fs, outputdirectory.DefaultOutputDir)
			Expect(exists).To(BeTrue())
		})

		It("does not error if the directory already exists", func() {
			fs.MkdirAll(outputdirectory.DefaultOutputDir, 0755)

			err := outputdirectory.Resolve(fs, outputdirectory.DefaultOutputDir)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("when using a user-provided output directory", func() {
		It("succeeds if the directory exists", func() {
			fs.MkdirAll("/custom/output", 0755)

			err := outputdirectory.Resolve(fs, "/custom/output")
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns an error if the directory does not exist", func() {
			err := outputdirectory.Resolve(fs, "/custom/nonexistent")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("output directory \"/custom/nonexistent\" does not exist"))
		})
	})
})
