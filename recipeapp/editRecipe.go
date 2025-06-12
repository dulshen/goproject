package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/dulshen/goproject/climenus"
)

type editARecipeMenuData struct {
	Recipe    *Recipe
	RecipeIdx int
}

// const recipeNameIdx = "recipe name index"
// const ingredientsStartIdx = "ingredients start index"
// const ingredentsEndIdx = "ingredients end index"

func registerEditRecipeCommand(menu *climenus.Menu) {
	menu.AddCommand(&climenus.Command{Name: "edit", Description: "Edit a Recipe", Execute: editRecipesLoop})
}

func editRecipesLoop(args []string, menu *climenus.Menu) error {
	instructions := "Please choose a recipe to edit\n" +
		"---------------------------------"

	_, err := selectRecipeLoop(editRecipe, instructions)
	if err != nil {
		return err
	}

	return nil
}

func editRecipe(args []string, menu *climenus.Menu) error {
	chosenRecipeNum, err := strconv.Atoi(strings.TrimSpace(args[0]))
	if err != nil {
		return err
	}

	recipes, err := readRecipesJSON(jsonFileName)
	if err != nil {
		return err
	}

	index := chosenRecipeNum - 1
	recipe := &((*recipes)[index])

	editThisRecipeMenu := initializeEditARecipeMenu(recipe, index)

	err = editThisRecipeMenu.MenuLoop()
	if err != nil {
		return err
	}

	// on exiting this menu need to re-read recipes to see if any chnages were saved
	recipes, err = readRecipesJSON(jsonFileName)
	if err != nil {
		return err
	}
	// re-initialize the parent menu before it is re-displayed, in case of changes
	err = InitializeSelectRecipeCommands(menu, recipes, editRecipe)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil

}

func initializeEditARecipeMenu(recipe *Recipe, index int) *climenus.Menu {
	var menu climenus.Menu

	menu.Instructions = "Choose an item from the recipe to edit:"
	c1 := climenus.MenuColumn{ColWidth: 5, Type: climenus.StringType, Label: optionNumberLabel}
	c2 := climenus.MenuColumn{ColWidth: -5, Type: climenus.StringType, Label: commandNameLabel}
	c3 := climenus.MenuColumn{ColWidth: -40, Type: climenus.StringType, Label: descriptionLabel}
	menu.Columns = append(menu.Columns, c1, c2, c3)

	menu.Data = &(editARecipeMenuData{Recipe: recipe, RecipeIdx: index})
	initializeEditRecipeCommands(&menu, recipe)

	return &menu
}

func initializeEditRecipeCommands(menu *climenus.Menu, recipe *Recipe) error {
	menu.Commands = nil
	menu.CommandsMap = nil
	menu.AddCommand(&climenus.Command{Name: "", Description: "Recipe Name: " + recipe.Name, Execute: editRecipeName})

	for _, ingredient := range recipe.Ingredients {
		ingredientStr := fmt.Sprintf("%v, %v, %v", ingredient.Name, ingredient.Quantity, ingredient.Unit)
		menu.AddCommand(&climenus.Command{Name: "", Description: ingredientStr, Execute: editRecipeIngredient})
	}

	menu.AddCommand(&climenus.Command{Name: "add", Description: "Add an Ingredient", Execute: addIngredient})
	menu.AddCommand(&climenus.Command{Name: "save", Description: "Save Recipe", Execute: saveChanges})
	menu.AddCommand(&climenus.Command{Name: "back", Execute: climenus.BackFunc})

	return nil
}

func editRecipeName(args []string, menu *climenus.Menu) error {
	recipe, _, err := extractRecipeData(menu)
	if err != nil {
		return err
	}

	prompt := "Provide a new name for this recipe:"
	input := climenus.UserInput(prompt, recipeNameValidator)
	recipe.Name = input

	// re-initialize edit a recipe commands in case options changed
	initializeEditRecipeCommands(menu, recipe)

	return nil
}

func editRecipeIngredient(args []string, menu *climenus.Menu) error {
	// args[0] will be the option number chosen
	// idx 1 is recipe name, so recipe ingredients start at idx 2 so adjust by 2 to get 0-indexed
	ingredientIdx, err := strconv.Atoi(args[0])
	ingredientIdx -= 2

	if err != nil {
		return err
	}

	recipe, _, err := extractRecipeData(menu)
	if err != nil {
		return err
	}

	prompt := "Provide new data for this ingredient (in the form ingredient name, quantity, unit(optional)):"
	input := climenus.UserInput(prompt, ingredientValidator)

	ingredient, err := parseIngredient(input)
	if err != nil {
		return err
	}

	recipe.Ingredients[ingredientIdx] = ingredient

	// re-initialize edit a recipe commands in case options changed
	initializeEditRecipeCommands(menu, recipe)

	return nil
}

func addIngredient(args []string, menu *climenus.Menu) error {
	recipe, _, err := extractRecipeData(menu)
	if err != nil {
		return err
	}

	prompt := "\nPlease enter recipe ingredients in the following format:"
	prompt += "Ingredient name, ingredient quantity, ingredient unit\n"
	input := climenus.UserInput(prompt, ingredientValidator)
	ingredient, err := parseIngredient(input)
	if err != nil {
		return err
	}

	recipe.Ingredients = append(recipe.Ingredients, ingredient)

	// re-initialize edit a recipe commands in case options changed
	initializeEditRecipeCommands(menu, recipe)

	return nil
}

func saveChanges(args []string, menu *climenus.Menu) error {
	recipe, recipeIdx, err := extractRecipeData(menu)
	if err != nil {
		return err
	}

	err = replaceRecipe(*recipe, jsonFileName, recipeIdx)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Printf("Successfully saved changes to %s\n", recipe.Name)
	time.Sleep(1 * time.Second)

	return nil
}

func extractRecipeData(menu *climenus.Menu) (*Recipe, int, error) {
	menuData, ok := menu.Data.(*editARecipeMenuData)
	if !ok {
		return nil, -1, errors.New("type assertion failed: edit a recipe menu data is not in correct form")
	}

	recipe := menuData.Recipe
	idx := menuData.RecipeIdx

	return recipe, idx, nil
}
