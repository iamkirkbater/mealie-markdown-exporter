package apitoken_test

import (
	"github.com/iamkirkbater/mealie-markdown-exporter/pkg/apitoken"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
)

var _ = Describe("Resolve", func() {
	var fs afero.Fs

	BeforeEach(func() {
		fs = afero.NewMemMapFs()
	})

	Context("when the token is provided directly", func() {
		It("returns the token as-is", func() {
			token, err := apitoken.Resolve(fs, "my-secret-token")
			Expect(err).NotTo(HaveOccurred())
			Expect(token).To(Equal("my-secret-token"))
		})

		It("returns an empty string as-is", func() {
			token, err := apitoken.Resolve(fs, "")
			Expect(err).NotTo(HaveOccurred())
			Expect(token).To(Equal(""))
		})
	})

	Context("when the token is a file://path/to/file", func() {
		It("reads the token from an absolute file path", func() {
			afero.WriteFile(fs, "/tmp/token.txt", []byte("token-from-file"), 0644)

			token, err := apitoken.Resolve(fs, "file:///tmp/token.txt")
			Expect(err).NotTo(HaveOccurred())
			Expect(token).To(Equal("token-from-file"))
		})

		It("reads the token from a relative file path in the current directory", func() {
			afero.WriteFile(fs, ".my-token", []byte("relative-token"), 0644)

			token, err := apitoken.Resolve(fs, "file://.my-token")
			Expect(err).NotTo(HaveOccurred())
			Expect(token).To(Equal("relative-token"))
		})

		It("reads the token from a relative file path in a subdirectory", func() {
			fs.MkdirAll("config", 0755)
			afero.WriteFile(fs, "config/my-token", []byte("nested-token"), 0644)

			token, err := apitoken.Resolve(fs, "file://config/my-token")
			Expect(err).NotTo(HaveOccurred())
			Expect(token).To(Equal("nested-token"))
		})

		It("trims whitespace from the file contents", func() {
			afero.WriteFile(fs, "/tmp/token.txt", []byte("  token-with-spaces  \n"), 0644)

			token, err := apitoken.Resolve(fs, "file:///tmp/token.txt")
			Expect(err).NotTo(HaveOccurred())
			Expect(token).To(Equal("token-with-spaces"))
		})

		It("returns an error when the file does not exist", func() {
			_, err := apitoken.Resolve(fs, "file:///tmp/nonexistent.txt")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to read api-token from file"))
		})
	})
})
