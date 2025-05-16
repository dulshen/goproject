package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/dulshen/goproject/climenus"
)

// struct describing a recipe
type Recipe struct {
	Name        string       // name of the recipe
	Ingredients []Ingredient // ingredients list
}

// struct describing an Ingredient
type Ingredient struct {
	Name     string // name of the ingredient
	Quantity int    // quantity of the ingredient
	Unit     string // unit of quantity
}

func main() {

	initializeJSONFile(jsonFileName, jsonDirectoryName, false)

	log.SetPrefix("climenu: ")
	log.SetFlags(0)

	menuData, err := climenus.BuildMenu(mainMenuOptions)

	if err != nil {
		log.Fatal("failed to build menu")
	}

	menuString, err := climenus.PrintMenu(menuData)

	if err != nil {
		log.Fatal("failed to print menu")
	}

	var selection string

	for selection != "exit" {
		fmt.Print(menuString)
		isInput := false
		for !isInput {
			_, err = fmt.Scan(&selection)

			if err == nil {
				selection, err = climenus.MakeSelection(menuData, selection)
			}

			if err != nil {
				fmt.Println("Selection does not exist. Please enter a valid selection.")
			} else {
				isInput = true
			}
		}

		// fmt.Print(selection)

		if selection == "add" {
			AddRecipeLoop()
		}
		if selection == "view" {
			viewRecipeMenu(jsonFileName)
			// isDeleteMode := false
			// SelectRecipeMenu(jsonFileName, isDeleteMode)
		}
		if selection == "edit" {
			editRecipesMenu(jsonFileName)
		}
		if selection == "del" {
			deleteRecipeMenu(jsonFileName)
			// isDeleteMode := true
			// SelectRecipeMenu(jsonFileName, isDeleteMode)
		}
		if selection == "clearAll" {
			clearAllRecipes()
		}
	}
}

func clearAllRecipes() error {

	countConfirms := 0
	fmt.Println("This will delete all recipes saved so far. Are you sure? (Y/N)")

	for countConfirms < 2 {
		var input string
		_, err := fmt.Scan(&input)

		if err != nil {
			log.Fatal(err)
		}

		if strings.ToLower(input) == "n" {
			return nil
		} else {
			countConfirms++
			if countConfirms == 1 {
				fmt.Println("Enter Y once more to proceed with deleting all recipes. (Y/N)")
			}
		}
	}

	initializeJSONFile(jsonFileName, jsonDirectoryName, true)

	return nil
}
