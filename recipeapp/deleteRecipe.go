package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/dulshen/goproject/climenus"
)

// menu name for delete recipe
const delName = "del"

// menu description for delete recipe
const delDescr = "Delete Recipe"

// Registers the delete recipe command in the main menu
func registerDeleteRecipeCommand(menu *climenus.Menu) {
	menu.AddCommand(&climenus.Command{Name: delName, Description: delDescr, Execute: deleteRecipeLoop})
}

// Main loop for deleting a recipe. Uses the select recipe loop to get a recipe selection
// from the user, and then executes the deleteRecipe function
func deleteRecipeLoop(args []string, menu *climenus.Menu) error {
	instructions := "Please choose a recipe to delete\n" +
		"---------------------------------"
	err := selectRecipeLoop(deleteRecipe, instructions)
	if err != nil {
		return err
	}

	return nil
}

// Removes the recipe indicated by the index provided in args from the stored recipe data
func deleteRecipe(args []string, menu *climenus.Menu) error {

	selectionInt, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}

	index := selectionInt - 1
	err = removeRecipe(index, jsonFileName)
	if err != nil {
		return err
	}

	fmt.Printf("Succesfully deleted recipe %s\n", args[0])
	time.Sleep(1 * time.Second)

	recipes, err := readRecipesJSON(jsonFileName)
	if err != nil {
		return err
	}
	// re-initialize the select recipe menu commands list
	InitializeSelectRecipeCommands(menu, recipes, deleteRecipe)

	return nil
}
