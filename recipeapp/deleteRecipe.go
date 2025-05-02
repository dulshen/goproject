package main

import (
	"strconv"

	"github.com/dulshen/goproject/climenus"
)

// deletes a recipe from the recipes data using the selected index provided by user as a string
func deleteRecipe(recipes *[]Recipe, viewMenuData *[]climenus.MenuOptionData, selection string) error {

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
