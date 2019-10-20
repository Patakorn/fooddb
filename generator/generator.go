package generator

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

var ErrInvalidLanguage = errors.New("invalid language")

func generateTpl(database *FoodDB, languageTpl string, outFile string) error {
	tpl, err := template.New(outFile).Parse(languageTpl)
	if err != nil {
		return fmt.Errorf("cannot parse template: %w", err)
	}
	f, err := os.Create(outFile)
	if err != nil {
		return fmt.Errorf("cannot create generated file: %w", err)
	}

	err = tpl.Execute(f, database)
	if err != nil {
		return fmt.Errorf("cannot render template: %w", err)
	}
	return nil
}

func Generate(database *FoodDB, language string, outputDir string) error {
	switch language {
	case "typescript":
		err := generateTpl(database, TypeScriptTemplate, filepath.Join(outputDir, "fooddb.ts"))
		if err != nil {
			return fmt.Errorf("cannot generate Typescript database: %w", err)
		}
	default:
		return fmt.Errorf("%q: %w", language, ErrInvalidLanguage)
	}

	return nil
}
