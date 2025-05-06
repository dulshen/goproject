package main

import (
	"fmt"
	"log"
	"strconv"
)

func deleteRecipeMenu(filename string) error {
	log.SetPrefix("deleteRecipeMenu: ")

	recipes, err := readRecipesJSON(filename)

	if err != nil {
		log.Fatal(err.Error())
	}

	var selection string
	for selection != "back" {
		selection, err = selectRecipeMenuNew(recipes)
		if err != nil {
			fmt.Println(err.Error())
		}

		if selection != "back" {
			deleteRecipe(selection)

			recipes, err = readRecipesJSON(filename)
			if err != nil {
				log.Fatal(err.Error())
			}
		}
	}

	return nil
}

// deletes a recipe from the recipes data using the selected index provided by user as a string
func deleteRecipe(selection string) error {

	selectionInt, err := strconv.Atoi(selection)
	if err != nil {
		return err
	}

	err = removeRecipe(selectionInt-1, jsonFileName)

	if err != nil {
		return err
	}

	return nil
}
