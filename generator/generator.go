package generator

import (
	"errors"
	"fmt"
	"io"
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

func copyFile(src, dest string) error {
	from, err := os.Open(src)
	if err != nil {
		return err
	}
	defer from.Close()

	to, err := os.OpenFile(dest, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer to.Close()

	_, err = io.Copy(to, from)
	if err != nil {
		return err
	}
	return nil
}

func copyImageTypeScript(src, outputDir string) error {
	if src == "" {
		return nil
	}
	dest := filepath.Join(outputDir, filepath.Base(src))
	if err := copyFile(src, dest); err != nil {
		return fmt.Errorf("cannot copy file '%v' to '%v': %w", src, dest, err)
	}
	return nil
}

func Generate(database *FoodDB, language string, outputDir string) error {
	switch language {
	case "typescript":
		// Used to import images in TypeScript easily (as module)
		err := generateTpl(database, TypeScriptModuleTemplate, filepath.Join(outputDir, "fooddb.d.ts"))
		if err != nil {
			return fmt.Errorf("cannot generate Typescript module file: %w", err)
		}

		err = generateTpl(database, TypeScriptTemplate, filepath.Join(outputDir, "fooddb.gen.ts"))
		if err != nil {
			return fmt.Errorf("cannot generate Typescript database: %w", err)
		}

		// Copy images to a dedicated folder that can be imported from TypeScript code
		imageDir := filepath.Join(outputDir, TypeScriptImageDir)
		err = os.MkdirAll(imageDir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("cannot create image dir '%v': %w", imageDir, err)
		}
		for _, c := range database.Categories {
			if err := copyImageTypeScript(c.ImagePath, imageDir); err != nil {
				return fmt.Errorf("cannot copy image file: %w", err)
			}
		}
		for _, i := range database.Ingredients {
			if err := copyImageTypeScript(i.ImagePath, imageDir); err != nil {
				return fmt.Errorf("cannot copy image file: %w", err)
			}
		}
	default:
		return fmt.Errorf("%q: %w", language, ErrInvalidLanguage)
	}

	return nil
}
