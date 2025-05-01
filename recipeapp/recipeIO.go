package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
)

// string representing the filename to store json data
const jsonFileName = "../data/recipes.json"

// string representing the directory json data is stored in
const jsonDirectoryName = "../data"

// checks if the directory and json file for storing data are set up yet
// and creates the directory and file if needed, to prevent errors later
// on when writing data to json file
func initializeJSONFile(filename string, directory string, overwrite bool) error {
	log.SetPrefix("initializeJSONFile: ")
	log.SetFlags(0)

	var emptyList []Recipe

	// first check if directory exists, and create it if not
	stat, err := os.Stat(directory)

	if err != nil {
		err := os.Mkdir(directory, 0755)
		if err != nil {
			log.Fatal(err)
		}
	} else if !stat.IsDir() {
		// err was nil but still need to checkc if path is a directory here
		log.Fatalf("%s is an existing path but is not a directory", directory)
	}

	// next check if file exists
	_, err = os.Stat(filename)

	// if file doesn't exist then create it (or if overwrite flag was set)
	if overwrite || errors.Is(err, fs.ErrNotExist) {
		jsonData, err := json.MarshalIndent(emptyList, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		// write the empty list to the file to initialize
		os.WriteFile(filename, jsonData, 0644)
	}

	return nil
}

// reads in recipes from a JSON file
func readRecipesJSON(filename string) (*[]Recipe, error) {
	var existingRecipes []Recipe

	jsonDataFromFile, err := os.ReadFile(filename)

	if err == nil {
		err = json.Unmarshal(jsonDataFromFile, &existingRecipes)
	}

	if err != nil {
		fmt.Println(err.Error())
		return nil, errors.New("failed to serialize recipe ingredient list")
	}

	return &existingRecipes, nil
}

// writes out recipe list to a JSON file
func writeRecipesJSON(filename string, recipes *[]Recipe) error {

	updatedJson, err := json.MarshalIndent(*recipes, "", "  ")

	if err != nil {
		return errors.New("failed to serialize updated recipes")
	}

	err = os.WriteFile(filename, updatedJson, 0644)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}
