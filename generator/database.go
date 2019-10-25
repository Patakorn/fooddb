package generator

import (
	"errors"
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
	Name      map[string]string `yaml:"name"`
	ImagePath string
}

type Ingredient struct {
	Name       map[string]string `yaml:"name"`
	Categories []string          `yaml:"categories"`
	Units      []string          `yaml:"units"`
	ImagePath  string
}

func exploreYaml(location string, exploreFn func(basename string, content []byte) error) error {
	yamlPathes, err := filepath.Glob(filepath.Join(location, "*.yaml"))
	if err != nil {
		return fmt.Errorf("cannot glob path '%v': %w", location, err)
	}
	for _, path := range yamlPathes {
		ext := filepath.Ext(path)
		basename := strings.TrimSuffix(filepath.Base(path), ext)

		content, err := ioutil.ReadFile(path)
		if err != nil {
			return fmt.Errorf("cannot read file '%v': %w", path, err)
		}

		if err := exploreFn(basename, content); err != nil {
			return fmt.Errorf("cannot parse file '%v': %w", path, err)
		}
	}
	return nil
}

var ErrNoImageFound error = errors.New("No image found for related element")

func getImage(basepath string, name string) (string, error) {
	path := filepath.Join(basepath, name)
	for _, ext := range []string{"jpg", "jpeg", "png"} {
		imgPath := strings.Join([]string{path, ext}, ".")
		if _, err := os.Stat(imgPath); err == nil {
			return imgPath, nil
		}
	}
	return "", ErrNoImageFound
}

func LoadDatabase(location string) (*FoodDB, error) {
	database := FoodDB{
		Categories:  make(map[string]Category),
		Ingredients: make(map[string]Ingredient),
	}

	// Load category files
	categoriePath := filepath.Join(location, CategorieLocation)
	if err := exploreYaml(categoriePath, func(basename string, content []byte) error {
		var category Category
		if err := yaml.Unmarshal(content, &category); err != nil {
			return fmt.Errorf("cannot parse to YAML: %w", err)
		}
		if imagePath, err := getImage(categoriePath, basename); err == nil {
			category.ImagePath = imagePath
		}

		database.Categories[basename] = category

		return nil
	}); err != nil {
		return nil, fmt.Errorf("cannot read category files: %w", err)
	}

	// Load ingredient files
	ingredientPath := filepath.Join(location, IngredientLocation)
	if err := exploreYaml(ingredientPath, func(basename string, content []byte) error {
		var ingredient Ingredient
		if err := yaml.Unmarshal(content, &ingredient); err != nil {
			return fmt.Errorf("cannot parse to YAML: %w", err)
		}
		if imagePath, err := getImage(ingredientPath, basename); err == nil {
			ingredient.ImagePath = imagePath
		}

		database.Ingredients[basename] = ingredient

		return nil
	}); err != nil {
		return nil, fmt.Errorf("cannot read ingredient files: %w", err)
	}

	return &database, nil
}
