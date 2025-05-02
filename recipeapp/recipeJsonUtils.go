package main

import (
	"errors"
	"slices"
)

var errRecipeAlreadyExists = errors.New("recipe to add already exists in dataset")

// util function to get a single recipe from the json file
func getRecipe(index int, filename string) (Recipe, error) {
	recipes, err := readRecipesJSON(filename)

	if err != nil {
		return Recipe{}, err
	}

	if !isValidIndex(index, recipes) {
		return Recipe{}, errors.New("getRecipe: index out of range")
	}

	return (*recipes)[index], nil
}

// util function to remove a single recipe from the json file
func removeRecipe(index int, filename string) error {
	recipes, err := readRecipesJSON(filename)

	if err != nil {
		return err
	}

	if !isValidIndex(index, recipes) {
		return errors.New("removeRecipe: index out of range")
	}

	updatedRecipes := slices.Delete((*recipes), index, index+1)

	writeRecipesJSON(filename, &updatedRecipes)

	return nil
}

// adds a provided recipe to the json recipe list data
// if the recipe name already exists in dataset, returns an error
// unless overwrite argument is set to true
func addRecipe(recipe Recipe, filename string, overwrite bool) error {
	recipes, err := readRecipesJSON(jsonFileName)

	if err != nil {
		return err
	}
	// first check if recipe already exists, return error if yes and overwrite=false
	// or overwrite the recipe in list with the new recipe if overwrite=true
	for i, oldRecipe := range *recipes {
		if oldRecipe.Name == recipe.Name {
			if !overwrite {
				return errRecipeAlreadyExists
			} else {
				(*recipes)[i].Ingredients = recipe.Ingredients
				return nil
			}
		}
	}

	// if got here then didn't find the recipe so add it
	*recipes = append((*recipes), recipe)

	err = writeRecipesJSON(jsonFileName, recipes)
	if err != nil {
		return err
	}

	return nil
}

// checks if a provided index for recipe list is valid
func isValidIndex(index int, recipes *[]Recipe) bool {
	if index < 0 || index > len(*recipes)-1 {
		return false
	}

	return true
}
