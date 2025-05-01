package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

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
