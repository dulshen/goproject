package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/dulshen/goproject/climenus"
)

func registerViewRecipeCommand(menu *climenus.Menu) {
	c := climenus.Command{Name: view, Description: "View a Recipe", Execute: viewRecipeLoop}
	menu.AddCommand(&c)
}

func viewRecipeLoop(args []string) error {
	menu, err := selectRecipeMenu(viewRecipe)
	menu.Instructions = "Please choose a recipe to view" +
		"---------------------------------"
	if err != nil {
		return err
	}

	err = menu.ShowMenu()
	if err != nil {
		return err
	}

	err = menu.MenuLoop()
	if err != nil {
		return err
	}
	return nil
}

func viewRecipe(args []string) error {

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
		fmt.Printf("%s: %f %s\n", ingredient.Name, ingredient.Quantity, ingredient.Unit)
	}

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

func scaleRecipe(recipe *Recipe, multiplierString string) (string, error) {
	scaledRecipeString := ""

	multiplier, err := strconv.ParseFloat(multiplierString, 32)
	if err != nil {
		return "", err
	}
	for _, ingredient := range recipe.Ingredients {
		scaledRecipeString += fmt.Sprintf("%s: %f %s\n", ingredient.Name,
			ingredient.Quantity*float32(multiplier), ingredient.Unit)
	}

	return scaledRecipeString, nil
}
