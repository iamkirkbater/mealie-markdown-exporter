package markdown

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/iamkirkbater/mealie-markdown-exporter/pkg/provider/mealie"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

//go:embed default.tmpl
var defaultTemplate string

type TemplateData struct {
	Recipe    mealie.Recipe
	ImagePath string
}

var funcMap = template.FuncMap{
	"escapeQuotes": func(s string) string {
		return strings.ReplaceAll(s, `"`, `\"`)
	},
	"categoryNames": func(cats []mealie.RecipeCategory) []string {
		names := make([]string, len(cats))
		for i, c := range cats {
			names[i] = c.Name
		}
		return names
	},
	"tagNames": func(tags []mealie.RecipeTag) []string {
		names := make([]string, len(tags))
		for i, t := range tags {
			names[i] = t.Name
		}
		return names
	},
	"add": func(a, b int) int {
		return a + b
	},
	"hasNutrition": func(n mealie.Nutrition) bool {
		return n.Calories != "" || n.CarbohydrateContent != "" || n.ProteinContent != "" ||
			n.FatContent != "" || n.SaturatedFatContent != "" || n.UnsaturatedFatContent != "" ||
			n.TransFatContent != "" || n.CholesterolContent != "" || n.SodiumContent != "" ||
			n.FiberContent != "" || n.SugarContent != ""
	},
}

type Provider struct {
	templateContent string
}

type Option func(*Provider) error

func WithTemplateFilePath(path string) Option {
	return func(p *Provider) error {
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read template file %q: %w", path, err)
		}
		p.templateContent = string(data)
		return nil
	}
}

func NewProvider(opts ...Option) (*Provider, error) {
	p := &Provider{
		templateContent: defaultTemplate,
	}
	for _, opt := range opts {
		if err := opt(p); err != nil {
			return nil, err
		}
	}
	return p, nil
}

// WriteRecipes renders each recipe as a markdown file and writes it to the output directory.
// The images map contains slug -> image filename for recipes that have downloaded images.
func (p *Provider) WriteRecipes(fs afero.Fs, outputDir string, recipes []mealie.Recipe, images map[string]string) error {
	tmpl, err := template.New("recipe").Funcs(funcMap).Parse(p.templateContent)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	for _, recipe := range recipes {
		if err := writeRecipe(fs, tmpl, outputDir, recipe, images[recipe.Slug]); err != nil {
			log.Errorf("Failed to write recipe %q: %v", recipe.Slug, err)
			return err
		}
	}

	return nil
}

func writeRecipe(fs afero.Fs, tmpl *template.Template, outputDir string, recipe mealie.Recipe, imagePath string) error {
	data := TemplateData{Recipe: recipe, ImagePath: imagePath}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("failed to render template for %q: %w", recipe.Slug, err)
	}

	filename := filepath.Join(outputDir, recipe.Slug+".md")
	if err := afero.WriteFile(fs, filename, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write file %q: %w", filename, err)
	}

	log.Debugf("Wrote %s", filename)
	return nil
}
