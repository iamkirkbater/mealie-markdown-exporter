package markdown

import (
	"bytes"
	_ "embed"
	"fmt"
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
	Recipe mealie.Recipe
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
}

// WriteRecipes renders each recipe as a markdown file and writes it to the output directory.
func WriteRecipes(fs afero.Fs, outputDir string, recipes []mealie.Recipe) error {
	tmpl, err := template.New("recipe").Funcs(funcMap).Parse(defaultTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	for _, recipe := range recipes {
		if err := writeRecipe(fs, tmpl, outputDir, recipe); err != nil {
			log.Errorf("Failed to write recipe %q: %v", recipe.Slug, err)
			return err
		}
	}

	return nil
}

func writeRecipe(fs afero.Fs, tmpl *template.Template, outputDir string, recipe mealie.Recipe) error {
	data := TemplateData{Recipe: recipe}

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
