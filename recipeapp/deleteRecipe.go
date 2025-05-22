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

func deleteRecipeLoop(args []string) error {
	instructions := "Please choose a recipe to delete\n" +
		"---------------------------------"
	menu, err := selectRecipeLoop(deleteRecipe, instructions)
	if err != nil {
		return err
	}
	// delMenu = menu
	_, err = fmt.Println(menu)
	if err != nil {
		fmt.Println(err.Error())
	}

	// 	menu.Instructions = "Please choose a recipe to delete" +
	// 		"---------------------------------"
	// 	if err != nil {
	// 		return err
	// 	}

	// 	err = menu.ShowMenu()
	// 	if err != nil {
	// 		return err
	// 	}

	// 	err = menu.MenuLoop()
	// 	if err != nil {
	// 		return err
	// 	}
	return nil
}

func deleteRecipe(args []string) error {

	// recipes, err := readRecipesJSON(jsonFileName)
	// if err != nil {
	// 	return err
	// }

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

// func deleteRecipeMenu(filename string) error {
// 	log.SetPrefix("deleteRecipeMenu: ")

// 	recipes, err := readRecipesJSON(filename)

// 	if err != nil {
// 		log.Fatal(err.Error())
// 	}

// 	var selection string
// 	for selection != "back" {
// 		selection, err = selectRecipeMenuNew(recipes)
// 		if err != nil {
// 			fmt.Println(err.Error())
// 		}

// 		if selection != "back" {
// 			deleteRecipe(selection)

// 			recipes, err = readRecipesJSON(filename)
// 			if err != nil {
// 				log.Fatal(err.Error())
// 			}
// 		}
// 	}

// 	return nil
// }

// // deletes a recipe from the recipes data using the selected index provided by user as a string
// func deleteRecipe(selection string) error {

// 	selectionInt, err := strconv.Atoi(selection)
// 	if err != nil {
// 		return err
// 	}

// 	err = removeRecipe(selectionInt-1, jsonFileName)

// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
