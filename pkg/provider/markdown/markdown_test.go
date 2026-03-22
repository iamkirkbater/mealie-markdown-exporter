package markdown_test

import (
	"os"

	"github.com/iamkirkbater/mealie-markdown-exporter/pkg/provider/markdown"
	"github.com/iamkirkbater/mealie-markdown-exporter/pkg/provider/mealie"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
)

var _ = Describe("WriteRecipes", func() {
	var (
		fs       afero.Fs
		provider *markdown.Provider
	)

	BeforeEach(func() {
		fs = afero.NewMemMapFs()
		fs.MkdirAll("/output", 0755)

		var err error
		provider, err = markdown.NewProvider()
		Expect(err).NotTo(HaveOccurred())
	})

	It("writes a markdown file for each recipe", func() {
		recipes := []mealie.Recipe{
			{Name: "Pancakes", Slug: "pancakes"},
			{Name: "Waffles", Slug: "waffles"},
		}

		err := provider.WriteRecipes(fs, "/output", recipes, nil)
		Expect(err).NotTo(HaveOccurred())

		exists, _ := afero.Exists(fs, "/output/pancakes.md")
		Expect(exists).To(BeTrue())
		exists, _ = afero.Exists(fs, "/output/waffles.md")
		Expect(exists).To(BeTrue())
	})

	It("includes the title in the front matter and body", func() {
		recipes := []mealie.Recipe{
			{Name: "Chocolate Cake", Slug: "chocolate-cake"},
		}

		err := provider.WriteRecipes(fs, "/output", recipes, nil)
		Expect(err).NotTo(HaveOccurred())

		content, err := afero.ReadFile(fs, "/output/chocolate-cake.md")
		Expect(err).NotTo(HaveOccurred())
		Expect(string(content)).To(ContainSubstring(`title: "Chocolate Cake"`))
	})

	It("includes description when present", func() {
		recipes := []mealie.Recipe{
			{
				Name:        "Soup",
				Slug:        "soup",
				Description: "A warm bowl of soup",
			},
		}

		err := provider.WriteRecipes(fs, "/output", recipes, nil)
		Expect(err).NotTo(HaveOccurred())

		content, _ := afero.ReadFile(fs, "/output/soup.md")
		Expect(string(content)).To(ContainSubstring(`summary: "A warm bowl of soup"`))
		Expect(string(content)).To(ContainSubstring("A warm bowl of soup"))
	})

	It("includes categories and tags", func() {
		recipes := []mealie.Recipe{
			{
				Name: "Tacos",
				Slug: "tacos",
				RecipeCategory: []mealie.RecipeCategory{
					{Name: "Dinner", Slug: "dinner"},
					{Name: "Mexican", Slug: "mexican"},
				},
				Tags: []mealie.RecipeTag{
					{Name: "quick", Slug: "quick"},
				},
			},
		}

		err := provider.WriteRecipes(fs, "/output", recipes, nil)
		Expect(err).NotTo(HaveOccurred())

		content, _ := afero.ReadFile(fs, "/output/tacos.md")
		s := string(content)
		Expect(s).To(ContainSubstring(`- "Dinner"`))
		Expect(s).To(ContainSubstring(`- "Mexican"`))
		Expect(s).To(ContainSubstring(`- "quick"`))
	})

	It("includes time fields when present", func() {
		recipes := []mealie.Recipe{
			{
				Name:      "Bread",
				Slug:      "bread",
				PrepTime:  "20 minutes",
				CookTime:  "45 minutes",
				TotalTime: "1 hour 5 minutes",
			},
		}

		err := provider.WriteRecipes(fs, "/output", recipes, nil)
		Expect(err).NotTo(HaveOccurred())

		content, _ := afero.ReadFile(fs, "/output/bread.md")
		s := string(content)
		Expect(s).To(ContainSubstring("**Prep Time**: 20 minutes"))
		Expect(s).To(ContainSubstring("**Cook Time**: 45 minutes"))
		Expect(s).To(ContainSubstring("**Total Time**: 1 hour 5 minutes"))
	})

	It("includes rating and source URL when present", func() {
		recipes := []mealie.Recipe{
			{
				Name:   "Pizza",
				Slug:   "pizza",
				Rating: 4.5,
				OrgURL: "https://example.com/pizza",
			},
		}

		err := provider.WriteRecipes(fs, "/output", recipes, nil)
		Expect(err).NotTo(HaveOccurred())

		content, _ := afero.ReadFile(fs, "/output/pizza.md")
		s := string(content)
		Expect(s).To(ContainSubstring("**Rating**: 4.5"))
		Expect(s).To(ContainSubstring("https://example.com/pizza"))
	})

	It("escapes quotes in the description", func() {
		recipes := []mealie.Recipe{
			{
				Name:        "Mom's \"Best\" Cookies",
				Slug:        "moms-best-cookies",
				Description: `A recipe for "the best" cookies`,
			},
		}

		err := provider.WriteRecipes(fs, "/output", recipes, nil)
		Expect(err).NotTo(HaveOccurred())

		content, _ := afero.ReadFile(fs, "/output/moms-best-cookies.md")
		s := string(content)
		Expect(s).To(ContainSubstring(`title: "Mom's \"Best\" Cookies"`))
		Expect(s).To(ContainSubstring(`summary: "A recipe for \"the best\" cookies"`))
	})

	It("omits optional fields when not present", func() {
		recipes := []mealie.Recipe{
			{Name: "Simple", Slug: "simple"},
		}

		err := provider.WriteRecipes(fs, "/output", recipes, nil)
		Expect(err).NotTo(HaveOccurred())

		content, _ := afero.ReadFile(fs, "/output/simple.md")
		s := string(content)
		Expect(s).NotTo(ContainSubstring("summary"))
		Expect(s).NotTo(ContainSubstring("categories"))
		Expect(s).NotTo(ContainSubstring("tags"))
		Expect(s).NotTo(ContainSubstring("Prep Time"))
		Expect(s).NotTo(ContainSubstring("Cook Time"))
		Expect(s).NotTo(ContainSubstring("Rating"))
		Expect(s).NotTo(ContainSubstring("source_url"))
	})

	It("includes image path in front matter when provided", func() {
		recipes := []mealie.Recipe{
			{Name: "Pancakes", Slug: "pancakes"},
		}
		images := map[string]string{"pancakes": "pancakes.webp"}

		err := provider.WriteRecipes(fs, "/output", recipes, images)
		Expect(err).NotTo(HaveOccurred())

		content, _ := afero.ReadFile(fs, "/output/pancakes.md")
		Expect(string(content)).To(ContainSubstring(`image: "pancakes.webp"`))
	})

	It("omits image from front matter when not provided", func() {
		recipes := []mealie.Recipe{
			{Name: "Pancakes", Slug: "pancakes"},
		}

		err := provider.WriteRecipes(fs, "/output", recipes, nil)
		Expect(err).NotTo(HaveOccurred())

		content, _ := afero.ReadFile(fs, "/output/pancakes.md")
		Expect(string(content)).NotTo(ContainSubstring("image:"))
	})

	Context("with a custom template file", func() {
		It("uses the provided template", func() {
			tmpFile, err := os.CreateTemp("", "custom-*.tmpl")
			Expect(err).NotTo(HaveOccurred())
			defer os.Remove(tmpFile.Name())

			_, err = tmpFile.WriteString(`# {{ .Recipe.Name }}`)
			Expect(err).NotTo(HaveOccurred())
			tmpFile.Close()

			customProvider, err := markdown.NewProvider(markdown.WithTemplateFilePath(tmpFile.Name()))
			Expect(err).NotTo(HaveOccurred())

			recipes := []mealie.Recipe{
				{Name: "Custom Recipe", Slug: "custom-recipe"},
			}

			err = customProvider.WriteRecipes(fs, "/output", recipes, nil)
			Expect(err).NotTo(HaveOccurred())

			content, _ := afero.ReadFile(fs, "/output/custom-recipe.md")
			Expect(string(content)).To(Equal("# Custom Recipe"))
		})

		It("returns an error if the template file does not exist", func() {
			_, err := markdown.NewProvider(markdown.WithTemplateFilePath("/nonexistent/template.tmpl"))
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to read template file"))
		})
	})
})
