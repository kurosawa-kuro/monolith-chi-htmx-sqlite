package templates

import (
	"html/template"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
)

// Loader handles template loading and management
type Loader struct {
	templates *template.Template
}

// NewLoader creates a new template loader
func NewLoader() *Loader {
	return &Loader{}
}

// Load loads all templates from the templates directory
func (l *Loader) Load() error {
	// Create template with custom functions
	funcMap := template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
		"sub": func(a, b int) int {
			return a - b
		},
		"mul": func(a, b int) int {
			return a * b
		},
		"div": func(a, b int) int {
			if b == 0 {
				return 0
			}
			return a / b
		},
		"min": func(a, b int) int {
			return int(math.Min(float64(a), float64(b)))
		},
		"max": func(a, b int) int {
			return int(math.Max(float64(a), float64(b)))
		},
		"ge": func(a, b int) bool {
			return a >= b
		},
		"le": func(a, b int) bool {
			return a <= b
		},
		"gt": func(a, b int) bool {
			return a > b
		},
		"lt": func(a, b int) bool {
			return a < b
		},
		"sequence": func(n int) []int {
			result := make([]int, n)
			for i := range result {
				result[i] = i + 1
			}
			return result
		},
	}

	tmpl := template.New("").Funcs(funcMap)

	// Load all template files
	err := filepath.Walk("src/templates", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".html" {
			content, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			_, err = tmpl.New(filepath.Base(path)).Parse(string(content))
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	l.templates = tmpl
	return nil
}

// GetTemplates returns the loaded templates
func (l *Loader) GetTemplates() *template.Template {
	return l.templates
}
