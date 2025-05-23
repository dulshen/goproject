package main

import (
	"fmt"
	"strconv"

	"github.com/dulshen/goproject/climenus"
)

const delName = "del"
const delDescr = "Delete Recipe"

func registerDeleteRecipeCommand(menu *climenus.Menu) {
	menu.AddCommand(&climenus.Command{Name: delName, Description: delDescr, Execute: deleteRecipeLoop})
}

func deleteRecipeLoop(args []string, menu *climenus.Menu) error {
	instructions := "Please choose a recipe to delete\n" +
		"---------------------------------"
	menu, err := selectRecipeLoop(deleteRecipe, instructions)
	if err != nil {
		return err
	}

	_, err = fmt.Println(menu)
	if err != nil {
		fmt.Println(err.Error())
	}

	return nil
}

func deleteRecipe(args []string) error {

	selectionInt, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}

	index := selectionInt - 1
	err = removeRecipe(index, jsonFileName)
	if err != nil {
		return err
	}
	recipes, err := readRecipesJSON(jsonFileName)
	if err != nil {
		return err
	}
	// re-initialize the select recipe menu commands list
	InitializeSelectRecipeCommands(selectRecipeMenu, recipes, deleteRecipe)

	return nil
}
