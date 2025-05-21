package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/dulshen/goproject/climenus"
)

const maxRecipeNameLength = 50

// const jsonFileName = "../data/recipes.json"

func registerAddRecipeCommand(menu *climenus.Menu) {
	menu.AddCommand(&climenus.Command{
		Name:        "add",
		Description: "Add Recipe",
		Execute:     AddRecipeLoop,
	})
}

func recipeNameValidator(input string) (bool, error) {
	if len(input) > maxRecipeNameLength {
		return false, fmt.Errorf("recipe name must be less than %v characters", maxRecipeNameLength)
	}

	return true, nil
}

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

const saveRecipeInput = "save"
const undoIngredientInput = "undo"

func AddRecipeLoop(args []string) error {

	if len(args) > 0 {
		return errors.New("invalid command")
	}

	prompt := "Enter a name for the recipe:\n----------------------------"
	recipeName := climenus.UserInput(prompt, recipeNameValidator)

	ingredientsList := make([]Ingredient, 0)
	// saveInput := saveRecipeInput
	// ingredientsStrings := climenus.UserInputLoop(prompt, saveInput, ingredientValidator)
	ingredientStrings := getIngredientsInput()
	for _, ingredientString := range ingredientStrings {
		var ingredient Ingredient
		// if ingredientString == undoIngredientInput {
		// 	removeLastIngredient(&ingredientsList)
		// } else {
		var err error
		ingredient, err = parseIngredient(ingredientString)
		if err != nil {
			fmt.Println("ingredientString:" + ingredientString)
			fmt.Println("parseingredient:" + err.Error())
			return err
		}
		// }
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

func yesNoValidator(input string) (bool, error) {
	if !(strings.ToLower(input) == "y" || strings.ToLower(input) == "n") {
		return false, errors.New("must enter either 'Y' or 'N'")
	}

	return true, nil
}

func saveRecipe(recipe *Recipe) error {
	overwrite := false
	err := addRecipe((*recipe), jsonFileName, overwrite)

	if errors.Is(err, errRecipeAlreadyExists) {
		prompt := fmt.Sprintf("A recipe with name %s already exists. Overwrite this recipe? (Y/N)\n", recipe.Name)
		choice := strings.ToLower(climenus.UserInput(prompt, yesNoValidator))
		// fmt.Printf("A recipe with name %s already exists. Overwrite this recipe? (Y/N)\n", recipe.Name)

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
