package generator

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

const CategorieLocation = "/categories"
const IngredientLocation = "/ingredients"

type FoodDB struct {
	Categories  map[string]Category
	Ingredients map[string]Ingredient
}

type Category struct {
	Name map[string]string `yaml:"name"`
}

type Ingredient struct {
	Name       map[string]string `yaml:"name"`
	Image      string            `yaml:"image"`
	Categories []string          `yaml:"categories"`
	Units      []string          `yaml:"units"`
}

func exploreFolder(location string, exploreFn func(basename string, content []byte) error) error {
	return filepath.Walk(location, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("cannot walk into path '%v': %w", location, err)
		}

		ext := filepath.Ext(path)
		basename := strings.TrimSuffix(filepath.Base(path), ext)

		if ext != ".yaml" {
			return nil
		}

		content, err := ioutil.ReadFile(path)
		if err != nil {
			return fmt.Errorf("cannot read file '%v': %w", path, err)
		}

		if err := exploreFn(basename, content); err != nil {
			return fmt.Errorf("cannot parse file '%v': %w", path, err)
		}

		return nil
	})
}

func LoadDatabase(location string) (*FoodDB, error) {
	database := FoodDB{
		Categories:  make(map[string]Category),
		Ingredients: make(map[string]Ingredient),
	}

	// Load category files
	categoriePath := filepath.Join(location, CategorieLocation)
	if err := exploreFolder(categoriePath, func(basename string, content []byte) error {
		var category Category
		if err := yaml.Unmarshal(content, &category); err != nil {
			return fmt.Errorf("cannot parse to YAML: %w", err)
		}
		database.Categories[basename] = category
		return nil
	}); err != nil {
		return nil, fmt.Errorf("cannot read category files: %w", err)
	}

	// Load ingredient files
	ingredientPath := filepath.Join(location, IngredientLocation)
	if err := exploreFolder(ingredientPath, func(basename string, content []byte) error {
		var ingredient Ingredient
		if err := yaml.Unmarshal(content, &ingredient); err != nil {
			return fmt.Errorf("cannot parse to YAML: %w", err)
		}
		database.Ingredients[basename] = ingredient
		return nil
	}); err != nil {
		return nil, fmt.Errorf("cannot read ingredient files: %w", err)
	}

	return &database, nil
}
