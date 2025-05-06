package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const maxRecipeNameLength = 50

// const jsonFileName = "../data/recipes.json"

// Prints the instructions for the Add Recipe option
func printAddInstructions() {
	fmt.Println("\nPlease enter recipe ingredients in the following format:")
	fmt.Println("Ingredient name, ingredient quantity, ingredient unit")
}

func printSaveUndoReminder() {
	fmt.Println("(Enter 'save' when done adding ingredients. Enter 'undo' to remove last ingredient entered.)")
	fmt.Println("--------------------------------------------------------------------------------------------")
}

func getRecipeName() (string, error) {
	fmt.Println("Enter a name for the recipe:")
	fmt.Println("----------------------------")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	name := scanner.Text()

	if len(name) > maxRecipeNameLength {
		errorStr := fmt.Sprintf("max length for recipe name is %d characters", maxRecipeNameLength)
		return "", errors.New(errorStr)
	}

	return name, nil
}

// Gets user input and parses an ingredient from it
// Returns an Ingredient, and an error if ingredient could not be parsed
func ParseIngredient(ingredientsList *[]Ingredient) (Ingredient, string, error) {

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()

	if input == "save" || input == "undo" {
		return Ingredient{}, input, nil
	}

	// if err != nil {
	// 	fmt.Println(err)
	// 	return Ingredient{}, errors.New("ingredient format error (use format: Ingredient name, ingredient quantity, ingredient unit)")
	// }

	items := strings.Split(input, ",")

	if len(items) > 3 || len(items) < 2 {
		return Ingredient{}, "", errors.New("must enter either ingredient, quantity or ingredient, quantity, unit")
	}

	name := strings.TrimSpace(items[0])
	quantity, err := strconv.Atoi(strings.TrimSpace(items[1]))
	var unit string
	if len(items) == 3 {
		unit = strings.TrimSpace(items[2])
	} else {
		unit = ""
	}

	if err != nil {
		return Ingredient{}, "", errors.New("ingredient quantity was not an integer")
	}

	ingredientData := Ingredient{
		Name:     name,
		Quantity: quantity,
		Unit:     unit,
	}

	*ingredientsList = append(*ingredientsList, ingredientData)

	return ingredientData, "", nil
}

func removeLastIngredient(ingredientList *[]Ingredient) error {
	if len(*ingredientList) == 0 {
		fmt.Println("undo failed: no ingredients to remove")
		return errors.New("no ingredients to remove")
	}

	ingredientRemoved := (*ingredientList)[len(*ingredientList)-1]
	*ingredientList = (*ingredientList)[:len(*ingredientList)-1]

	fmt.Printf("Removed ingredient: %v, quantity: %v, unit: %v\n",
		ingredientRemoved.Name, ingredientRemoved.Quantity, ingredientRemoved.Unit)

	return nil
}

func AddRecipeLoop() {

	// get user input for naming the recipe
	recipeName := ""
	var err error

	for recipeName == "" {
		recipeName, err = getRecipeName()
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	recipe := Recipe{
		Name: recipeName,
	}

	ingredientsInputLoop(&recipe)

	// // initialize variables for the list of Ingredients
	// // and for current line of user input
	// var input string
	// var ingredientList []Ingredient

	// // write the new recipe to JSON file
	// recipe := Recipe{
	// 	Name:        recipeName,
	// 	Ingredients: ingredientList,
	// }
	err = saveRecipe(&recipe)

	if err == nil {
		fmt.Println("Successfully saved recipe.")
	} else {
		fmt.Println(err.Error())
	}

}

func ingredientsInputLoop(recipe *Recipe) {
	// print instructions for adding ingredients
	printAddInstructions()

	var input string

	// loop while user adds ingredients until input == "save"
	for input != "save" {
		printSaveUndoReminder()

		input = ingredientInput()
		ingredient, err := parseIngredient(input)

		if err != nil {
			fmt.Println(err.Error())
		} else {
			recipe.Ingredients = append(recipe.Ingredients, ingredient)
			fmt.Printf("Added ingredient: %v, quantity: %v, unit: %v\n", ingredient.Name, ingredient.Quantity, ingredient.Unit)
		}

		// var ingredient Ingredient
		// ingredient, input, err := ParseIngredient(&ingredientList)

		// if err == nil {
		// 	if input != "save" && input != "undo" {
		// 		fmt.Printf("Added ingredient: %v, quantity: %v, unit: %v\n", ingredient.Name, ingredient.Quantity, ingredient.Unit)
		// 	} else if input == "undo" {
		// 		removeLastIngredient(&ingredientList)
		// 	}

		// } else {
		// 	fmt.Println(err)
		// }
	}
}

func ingredientInput() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()

	return input

	// if input == "save" || input == "undo" {
	// 	return Ingredient{}, input, nil
	// }

	// ingredient, err := parseIngredient(input)

	// if err != nil {
	// 	return Ingredient{}, "", err
	// }

	// return ingredient, "", nil
}

func parseIngredient(ingredientString string) (Ingredient, error) {
	items := strings.Split(ingredientString, ",")

	if len(items) > 3 || len(items) < 2 {
		return Ingredient{}, errors.New("must enter either ingredient, quantity or ingredient, quantity, unit")
	}

	name := strings.TrimSpace(items[0])
	quantity, err := strconv.Atoi(strings.TrimSpace(items[1]))
	if err != nil {
		return Ingredient{}, errors.New("ingredient quantity must be an integer")
	}

	unit := ""
	if len(items) == 3 {
		unit = strings.TrimSpace(items[2])
	}

	ingredientData := Ingredient{
		Name:     name,
		Quantity: quantity,
		Unit:     unit,
	}

	// *ingredientsList = append(*ingredientsList, ingredientData)

	return ingredientData, nil
}

// func addIngredient(recipe *Recipe, ingredientStr string) error {

// 	ingredient, err := ParseIngredient(ingredientStr)
// 	if err != nil {
// 		return err
// 	}

// 	recipe.Ingredients = append(recipe.Ingredients, ingredient)

// 	return nil
// }

func saveRecipe(recipe *Recipe) error {
	overwrite := false
	err := addRecipe((*recipe), jsonFileName, overwrite)

	if errors.Is(err, errRecipeAlreadyExists) {
		fmt.Printf("A recipe with name %s already exists. Overwrite this recipe? (Y/N)\n", recipe.Name)
		choice := ""
		for choice != "n" && choice != "y" {
			fmt.Scan(&choice)
			choice = strings.ToLower(choice)
			if choice == "n" {
				return errors.New("aborted creating new recipe due to conflicting recipe name")
			} else if choice == "y" {
				overwrite = true
				err = addRecipe((*recipe), jsonFileName, overwrite)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
