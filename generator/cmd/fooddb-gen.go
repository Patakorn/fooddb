package main

import (
	"fmt"
	"os"

	"github.com/fuhrmannb/fooddb/generator"
	"github.com/spf13/cobra"
)

var dbLocation, outputDir string

var rootCmd = &cobra.Command{
	Use:   "fooddb-gen",
	Short: "Generate a food static database for a specific programming language",
	Long: `fooddb-gen convert the YAML food database to static objects that can be used directly into a program
Doing so avoid your program to use a dedicated YAML library to load raw food database
Supported languages are: typescript`,
	RunE: func(cmd *cobra.Command, languages []string) error {
		db, err := generator.LoadDatabase(dbLocation)
		if err != nil {
			return err
		}

		err = os.MkdirAll(outputDir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("cannot create output directory: %w", err)
		}

		for _, l := range languages {
			err = generator.Generate(db, l, outputDir)
			if err != nil {
				return err
			}
		}
		return nil
	},
}

func main() {
	rootCmd.Flags().StringVarP(&dbLocation, "db-location", "l", "./database",
		`food database location`)
	rootCmd.Flags().StringVarP(&outputDir, "output-dir", "o", "./output",
		`output directory`)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
