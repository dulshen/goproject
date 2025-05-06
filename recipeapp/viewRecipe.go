package main

import (
	"fmt"
	"log"
	"strconv"
)

func viewRecipeMenu(filename string) error {
	log.SetPrefix("viewRecipeMenu: ")

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
			err = viewRecipe(recipes, selection)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}
	return nil
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
func viewRecipe(recipes *[]Recipe, selection string) error {
	selectionInt, err := strconv.Atoi(selection)

	if err != nil {
		return err
	}

	// recipe := (*recipes)[selectionInt-1]
	recipe, err := getRecipe(selectionInt-1, jsonFileName)
	if err != nil {
		return err
	}

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
