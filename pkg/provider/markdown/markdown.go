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
	Title       string
	Description string
	DateAdded   string
	DateUpdated string
	Categories  []string
	Tags        []string
	PrepTime    string
	CookTime    string
	TotalTime   string
	Servings    float64
	Rating      float64
	OrgURL      string
}

func templateDataFromRecipe(recipe mealie.Recipe) TemplateData {
	data := TemplateData{}

	if recipe.Name != nil {
		data.Title = *recipe.Name
	}
	if recipe.Description != nil {
		data.Description = *recipe.Description
	}
	if recipe.DateAdded != nil {
		data.DateAdded = *recipe.DateAdded
	}
	if recipe.DateUpdated != nil {
		data.DateUpdated = *recipe.DateUpdated
	}
	if recipe.PrepTime != nil {
		data.PrepTime = *recipe.PrepTime
	}
	if recipe.CookTime != nil {
		data.CookTime = *recipe.CookTime
	}
	if recipe.TotalTime != nil {
		data.TotalTime = *recipe.TotalTime
	}
	if recipe.Rating != nil {
		data.Rating = *recipe.Rating
	}
	if recipe.OrgURL != nil {
		data.OrgURL = *recipe.OrgURL
	}

	data.Servings = recipe.RecipeServings

	for _, cat := range recipe.RecipeCategory {
		data.Categories = append(data.Categories, cat.Name)
	}
	for _, tag := range recipe.Tags {
		data.Tags = append(data.Tags, tag.Name)
	}

	return data
}

// WriteRecipes renders each recipe as a markdown file and writes it to the output directory.
func WriteRecipes(fs afero.Fs, outputDir string, recipes []mealie.Recipe) error {
	funcMap := template.FuncMap{
		"escapeQuotes": func(s string) string {
			return strings.ReplaceAll(s, `"`, `\"`)
		},
	}
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
	data := templateDataFromRecipe(recipe)

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
