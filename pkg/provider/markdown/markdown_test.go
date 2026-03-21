package markdown_test

import (
	"github.com/iamkirkbater/mealie-markdown-exporter/pkg/provider/markdown"
	"github.com/iamkirkbater/mealie-markdown-exporter/pkg/provider/mealie"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
)

var _ = Describe("WriteRecipes", func() {
	var fs afero.Fs

	BeforeEach(func() {
		fs = afero.NewMemMapFs()
		fs.MkdirAll("/output", 0755)
	})

	It("writes a markdown file for each recipe", func() {
		recipes := []mealie.Recipe{
			{Name: "Pancakes", Slug: "pancakes"},
			{Name: "Waffles", Slug: "waffles"},
		}

		err := markdown.WriteRecipes(fs, "/output", recipes)
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

		err := markdown.WriteRecipes(fs, "/output", recipes)
		Expect(err).NotTo(HaveOccurred())

		content, err := afero.ReadFile(fs, "/output/chocolate-cake.md")
		Expect(err).NotTo(HaveOccurred())
		Expect(string(content)).To(ContainSubstring(`title: "Chocolate Cake"`))
		Expect(string(content)).To(ContainSubstring("# Chocolate Cake"))
	})

	It("includes description when present", func() {
		recipes := []mealie.Recipe{
			{
				Name:        "Soup",
				Slug:        "soup",
				Description: "A warm bowl of soup",
			},
		}

		err := markdown.WriteRecipes(fs, "/output", recipes)
		Expect(err).NotTo(HaveOccurred())

		content, _ := afero.ReadFile(fs, "/output/soup.md")
		Expect(string(content)).To(ContainSubstring(`description: "A warm bowl of soup"`))
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

		err := markdown.WriteRecipes(fs, "/output", recipes)
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

		err := markdown.WriteRecipes(fs, "/output", recipes)
		Expect(err).NotTo(HaveOccurred())

		content, _ := afero.ReadFile(fs, "/output/bread.md")
		s := string(content)
		Expect(s).To(ContainSubstring(`prep_time: "20 minutes"`))
		Expect(s).To(ContainSubstring(`cook_time: "45 minutes"`))
		Expect(s).To(ContainSubstring(`total_time: "1 hour 5 minutes"`))
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

		err := markdown.WriteRecipes(fs, "/output", recipes)
		Expect(err).NotTo(HaveOccurred())

		content, _ := afero.ReadFile(fs, "/output/pizza.md")
		s := string(content)
		Expect(s).To(ContainSubstring("rating: 4.5"))
		Expect(s).To(ContainSubstring(`source_url: "https://example.com/pizza"`))
	})

	It("escapes quotes in the description", func() {
		recipes := []mealie.Recipe{
			{
				Name:        "Mom's \"Best\" Cookies",
				Slug:        "moms-best-cookies",
				Description: `A recipe for "the best" cookies`,
			},
		}

		err := markdown.WriteRecipes(fs, "/output", recipes)
		Expect(err).NotTo(HaveOccurred())

		content, _ := afero.ReadFile(fs, "/output/moms-best-cookies.md")
		s := string(content)
		Expect(s).To(ContainSubstring(`title: "Mom's \"Best\" Cookies"`))
		Expect(s).To(ContainSubstring(`description: "A recipe for \"the best\" cookies"`))
	})

	It("omits optional fields when not present", func() {
		recipes := []mealie.Recipe{
			{Name: "Simple", Slug: "simple"},
		}

		err := markdown.WriteRecipes(fs, "/output", recipes)
		Expect(err).NotTo(HaveOccurred())

		content, _ := afero.ReadFile(fs, "/output/simple.md")
		s := string(content)
		Expect(s).NotTo(ContainSubstring("description"))
		Expect(s).NotTo(ContainSubstring("categories"))
		Expect(s).NotTo(ContainSubstring("tags"))
		Expect(s).NotTo(ContainSubstring("prep_time"))
		Expect(s).NotTo(ContainSubstring("cook_time"))
		Expect(s).NotTo(ContainSubstring("rating"))
		Expect(s).NotTo(ContainSubstring("source_url"))
	})
})
