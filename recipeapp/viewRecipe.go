package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/dulshen/goproject/climenus"
)

// menu name for view recipe
const viewName = "view"

// menu description for view recipe
const viewDescr = "View a Recipe"

// max size for a row of text for recipe steps
const maxStepLineSize = 50

// Function used to register the view recipe command in the main menu
func registerViewRecipeCommand(menu *climenus.Menu) {
	c := climenus.Command{Name: viewName, Description: viewDescr, Execute: viewRecipeLoop}
	menu.AddCommand(&c)
}

// Main loop for the view recipe menu
// Makes use of the select recipe loop, which prompts the user
// to select a recipe to view.
func viewRecipeLoop(args []string, menu *climenus.Menu) error {
	instructions := "Please choose a recipe to view" +
		"---------------------------------"
	err := selectRecipeLoop(viewRecipe, instructions)
	if err != nil {
		return err
	}

	return nil
}

// Views the recipe chosen by the user, which is indicated by the index number
// passed in args. Prints the recipe name, and all of the recipe ingredients.
func viewRecipe(args []string, menu *climenus.Menu) error {

	chosenRecipeNum, err := strconv.Atoi(strings.TrimSpace(args[0]))

	if err != nil {
		return err
	}

	index := chosenRecipeNum - 1
	recipe, err := getRecipe(index, jsonFileName)
	if err != nil {
		return err
	}

	fmt.Printf("\nRecipe: %s\n", recipe.Name)
	fmt.Println("----------------------------------")

	for _, ingredient := range recipe.Ingredients {
		fmt.Printf("%s: %.2f %s\n", ingredient.Name, ingredient.Quantity, ingredient.Unit)
	}

	fmt.Println("----------------------------------")

	for i, recipeStep := range recipe.Steps {
		stepSplits := getRecipeStepSplits(recipeStep, maxStepLineSize)
		stepNo := i + 1
		fmt.Println("--")
		fmt.Printf("%d: ", stepNo)
		for _, split := range stepSplits {
			fmt.Printf("%s\n", split)
		}
	}

	fmt.Println("--")
	fmt.Println("----------------------------------")

	fmt.Print("\n\n")

	bypassValidator := func(string) (bool, error) { return true, nil }
	input := ""
	for input != "back" {
		input = climenus.UserInput("Enter 'back' to return to previous menu, "+
			"or 'scale X' to scale recipe by X", bypassValidator)
		args := strings.Split(input, " ")
		if args[0] == "scale" {
			scaledRecipeString, err := scaleRecipe(&recipe, args[1])
			if err != nil {
				return err
			}
			fmt.Println(scaledRecipeString)
		}
	}

	return nil
}

// Scales the recipe currently being viewed by a multiplier value indicated
// by user input.
func scaleRecipe(recipe *Recipe, multiplierString string) (string, error) {
	scaledRecipeString := ""

	multiplier, err := strconv.ParseFloat(multiplierString, 32)
	if err != nil {
		return "", err
	}
	for _, ingredient := range recipe.Ingredients {
		scaledRecipeString += fmt.Sprintf("%s: %.2f %s\n", ingredient.Name,
			ingredient.Quantity*float32(multiplier), ingredient.Unit)
	}

	return scaledRecipeString, nil
}

// Splits recipe step into multiple ines of console output based on the maxWidth
// for a line of a recipe step.
func getRecipeStepSplits(step string, maxWidth int) []string {
	// extracts words from the recipe step string, so that splits can be broken
	// such that a word is not broken up in the middle
	words := strings.Split(step, " ")
	splits := make([]string, 0)
	var sb strings.Builder

	currentWidth := 0
	for i, word := range words {
		if len(word)+currentWidth >= maxWidth {
			currentWidth = 0
			splits = append(splits, sb.String())
			sb.Reset()
		}
		sb.WriteString(word)
		currentWidth += len(word)
		if i < len(words)-1 {
			sb.WriteString(" ")
		} else {
			splits = append(splits, sb.String())
		}
	}

	return splits
}
