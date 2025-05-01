package main

import (
	"fmt"
	"strconv"

	"github.com/dulshen/goproject/climenus"
)

// Presents a CLI menu for user to view a recipe, selected from the
// recipes located in the specified file
func viewRecipe(filename string) error {

	recipes, err := readRecipesJSON(filename)
	if err != nil {
		return err
	}

	viewMenuData, err := buildViewRecipeMenu(recipes)
	if err != nil {
		return err
	}

	menuString, err := climenus.PrintMenu(*viewMenuData)
	if err != nil {
		return err
	}

	var selection string
	// accept user input until user selects to back out of this menu
	for selection != "back" {
		fmt.Print(menuString)
		isInput := false
		// loop until valid user input
		for !isInput {
			_, err = fmt.Scan(&selection)
			// check if valid selection
			if err == nil {
				selection, err = climenus.MakeSelection(*viewMenuData, selection)
			}

			if err != nil {
				fmt.Println("Selection does not exist. Please enter a valid selection.")
			} else {
				isInput = true
			}
		}
		// skip printing recipe if user selected to go back, just exit loop instead
		if selection != "back" {
			printRecipe(recipes, selection)
		}
	}

	return nil
}

// Build the view recipe menu from the list of recipes.
// Returns a pointer to a slice of MenuOptionData.
func buildViewRecipeMenu(recipes *[]Recipe) (*[]climenus.MenuOptionData, error) {
	var viewMenuOptions []map[string]string

	for i, recipe := range *recipes {
		option := map[string]string{
			"description": recipe.Name,
			"menuKey":     strconv.Itoa(i + 1), // no support for menuKey for view recipe at the moment, just use number
		}
		viewMenuOptions = append(viewMenuOptions, option)
	}
	// add one more option to menu for going back
	viewMenuOptions = append(viewMenuOptions, map[string]string{
		"description": "Return to main menu",
		"menuKey":     "back",
	})

	viewMenuData, err := climenus.BuildMenu(viewMenuOptions)
	if err != nil {
		return nil, err
	}
	return &viewMenuData, nil
}

// Used to scale a recipe.
// Requests a float value from the user, then multiplies each ingredient quantity
// by this amount and prints the resulting recipe
func scaleRecipe(recipe Recipe) error {
	fmt.Println("enter amount to multiply by")
	var mult float32
	_, err := fmt.Scan(&mult)
	if err != nil {
		return err
	}

	for _, ingredient := range recipe.Ingredients {
		fmt.Printf("%s: %.2f %s\n", ingredient.Name, (float32(ingredient.Quantity) * mult), ingredient.Unit)
	}

	return nil
}

// Prints the recipe that was selected to be viewed by the user.
// Selection argument is a string representing the option number of the recipe.
func printRecipe(recipes *[]Recipe, selection string) error {
	selectionInt, err := strconv.Atoi(selection)

	if err != nil {
		return err
	}

	recipe := (*recipes)[selectionInt-1]

	fmt.Printf("\nRecipe: %s\n", recipe.Name)
	fmt.Println("----------------------------------")

	for _, ingredient := range recipe.Ingredients {
		fmt.Printf("%s: %d %s\n", ingredient.Name, ingredient.Quantity, ingredient.Unit)
	}

	var input string

	for input != "back" {
		fmt.Println("\nenter 'back' to return to previous menu, enter 'scale' to multiply the recipe.")
		_, err = fmt.Scan(&input)

		if err != nil {
			return err
		}
		if input == "scale" {
			scaleRecipe(recipe)
		}
	}

	return nil
}
