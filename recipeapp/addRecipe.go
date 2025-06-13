package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/dulshen/goproject/climenus"
)

// max allowed length for a recipe name
const maxRecipeNameLength = 50

// const jsonFileName = "../data/recipes.json"

// Function used toe register the add recipe command to the main menu
func registerAddRecipeCommand(menu *climenus.Menu) {
	menu.AddCommand(&climenus.Command{
		Name:        "add",
		Description: "Add Recipe",
		Execute:     AddRecipeLoop,
	})
}

// Function used as a validator for recipe name input
// Checks that recipe name is within the allowable length limit
func recipeNameValidator(input string) (bool, error) {
	if len(input) > maxRecipeNameLength {
		return false, fmt.Errorf("recipe name must be less than %v characters", maxRecipeNameLength)
	}

	return true, nil
}

// Function used as a validator for ingredient input
// Checks that the comma delimited list is the correct length for ingredient, quantity
// or ingredient, quantity, unit, and checks that the quantity can be pasrsed as a float
func ingredientValidator(input string) (bool, error) {
	if input == saveRecipeInput || input == undoIngredientInput {
		if input == undoIngredientInput {
			fmt.Println("last ingredient removed.")
		}
		return true, nil
	}

	items := strings.Split(input, ",")

	// check number of args is either 2 or 3
	if len(items) > 3 || len(items) < 2 {
		return false, errors.New("must enter either ingredient, quantity or ingredient, quantity, unit")
	}

	// check that 2nd arg can be converted to a float for quantity value
	_, err := strconv.ParseFloat(strings.TrimSpace(items[1]), 32)
	if err != nil {
		return false, errors.New("ingredient quantity must be a number")
	}

	return true, nil
}

// Parses an Ingredient struct from ingredient text input from the user
// returns the resulting Ingredient struct
func parseIngredient(ingredientString string) (Ingredient, error) {
	fields := strings.Split(ingredientString, ",")
	if !(len(fields) == 2 || len(fields) == 3) {
		return Ingredient{}, errors.New("")
	}
	name := strings.TrimSpace(fields[0])
	quantity, err := strconv.ParseFloat(strings.TrimSpace(fields[1]), 32)
	if err != nil {
		return Ingredient{}, err
	}

	unit := ""
	if len(fields) == 3 {
		unit = strings.TrimSpace(fields[2])
	}
	return Ingredient{Name: name, Quantity: float32(quantity), Unit: unit}, nil
}

// command for saving a recipe
const saveRecipeInput = "save"

// commmand for undoing an added ingredient
const undoIngredientInput = "undo"

// Loop used for adding a new recipe. Prompts the user for a recipe name,
// then has the user add ingredients one at a time, and then saves the new
// recipe to the data file when the user requests to save.
func AddRecipeLoop(args []string, menu *climenus.Menu) error {

	if len(args) > 1 {
		return errors.New("invalid command")
	}

	prompt := "Enter a name for the recipe:\n----------------------------"
	recipeName := climenus.UserInput(prompt, recipeNameValidator)

	ingredientsList := make([]Ingredient, 0)

	ingredientStrings := getIngredientsInput()
	for _, ingredientString := range ingredientStrings {
		var ingredient Ingredient

		var err error
		ingredient, err = parseIngredient(ingredientString)
		if err != nil {
			fmt.Println("ingredientString:" + ingredientString)
			fmt.Println("parseingredient:" + err.Error())
			return err
		}

		ingredientsList = append(ingredientsList, ingredient)
	}

	recipe := Recipe{
		Name:        recipeName,
		Ingredients: ingredientsList,
	}

	err := saveRecipe(&recipe)
	if err != nil {
		return err
	}

	fmt.Printf("Successfully saved recipe: %v.\n", recipe.Name)

	// put in a slight delay before returning to previous menu
	time.Sleep(time.Second * 2)
	return nil

}

// Function used to get user input for ingredients for the recipe.
// Loops through user input ingredients, validates that each user input can be parsed as an ingredient,
// then returns the slice of ingredient strings
func getIngredientsInput() []string {
	prompt := "\nPlease enter recipe ingredients in the following format:"
	prompt += "Ingredient name, ingredient quantity, ingredient unit\n"
	prompt += "(enter 'save' to save recipe once done, enter 'undo' to remove last added ingredient)"

	input := ""
	ingredientStrings := make([]string, 0)
	for input != saveRecipeInput {

		input = climenus.UserInput(prompt, ingredientValidator)
		if input == undoIngredientInput {
			n := len(ingredientStrings)
			if n > 0 {
				ingredientStrings = ingredientStrings[:n-1]
				fmt.Println("removed last ingredient")
			}
		} else if input != saveRecipeInput {
			ingredientStrings = append(ingredientStrings, input)
		}
	}

	return ingredientStrings
}

// Function used for validating input is either Y or N
func yesNoValidator(input string) (bool, error) {
	if !(strings.ToLower(input) == "y" || strings.ToLower(input) == "n") {
		return false, errors.New("must enter either 'Y' or 'N'")
	}

	return true, nil
}

// Saves the recipe that is currently being added to the stored recipe data
func saveRecipe(recipe *Recipe) error {
	overwrite := false
	err := addRecipe((*recipe), jsonFileName, overwrite)

	if errors.Is(err, errRecipeAlreadyExists) {
		prompt := fmt.Sprintf("A recipe with name %s already exists. Overwrite this recipe? (Y/N)\n", recipe.Name)
		choice := strings.ToLower(climenus.UserInput(prompt, yesNoValidator))

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

	return nil
}
