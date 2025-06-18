package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/dulshen/goproject/climenus"
)

// Struct used for passing necessary data for selected recipe
// to a climenus.Menu struct, so that this can be accessed by other functions later
type editARecipeMenuData struct {
	Recipe    *Recipe // pointer to the recipe for this menu
	RecipeIdx int     // int indicating the index of this recipe in the list of Recipes in storage
}

// const recipeNameIdx = "recipe name index"
// const ingredientsStartIdx = "ingredients start index"
// const ingredentsEndIdx = "ingredients end index"

// function used to add the Edit Recipe command to the main menu
func registerEditRecipeCommand(menu *climenus.Menu) {
	menu.AddCommand(&climenus.Command{Name: "edit", Description: "Edit a Recipe", Execute: editRecipesLoop})
}

// Main loop for the edit recipes menu, asks the user to select a recipe
// then calls the edit a recipe menu, passing the selected recipe in args
func editRecipesLoop(args []string, menu *climenus.Menu) error {
	instructions := "Please choose a recipe to edit\n" +
		"---------------------------------"

	err := selectRecipeLoop(editRecipe, instructions)
	if err != nil {
		return err
	}

	return nil
}

// Determines which recipe the user selected by parsing the selection in args
// then initializes a menu with options for editing this recipe
// and calls the menu loop for editing this recipe
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

// Initializes edit a recipe menu for the selected recipe
// sets the menu instructions, the column widths and types, and passes the recipe data and index
// to a struct stored in menu.Data, then calls another function to initialize the commands for the menu
// returns the initialized menu for editing this recipe
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

// Initializes (or re-initializes) commands for the edit recipe menu for this recipe
// adds commands for changing the recipe name, changing any ingredients, adding ingredients,
// saving the changes, or going back without saving
func initializeEditRecipeCommands(menu *climenus.Menu, recipe *Recipe) error {
	menu.Commands = nil
	menu.CommandsMap = nil
	menu.AddCommand(&climenus.Command{Name: "", Description: "Recipe Name: " + recipe.Name, Execute: editRecipeName})

	for _, ingredient := range recipe.Ingredients {
		ingredientStr := fmt.Sprintf("%v, %v, %v", ingredient.Name, ingredient.Quantity, ingredient.Unit)
		menu.AddCommand(&climenus.Command{Name: "", Description: ingredientStr, Execute: editRecipeIngredient})
	}

	for _, step := range recipe.Steps {
		menu.AddCommand(&climenus.Command{Name: "", Description: step, Execute: editRecipeStep})
	}

	menu.AddCommand(&climenus.Command{Name: "add", Description: "Add an Ingredient", Execute: addIngredient})
	menu.AddCommand(&climenus.Command{Name: "save", Description: "Save Recipe", Execute: saveChanges})
	menu.AddCommand(&climenus.Command{Name: "back", Execute: climenus.BackFunc})

	return nil
}

// Function used for editing a recipe name. Takes user input for a new name
// then renames the recipe currently being edited
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

// Function used for editing a recipe ingredient. Takes user input for an updated
// ingredient to use for the selected ingredient (indicated by args), and stores this
// in the ingredients list for the recipe being edited currently
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

// Function used when user selects to edit a recipe step.
// Uses the provided arg to get the index of the recipe step to edit
// then prompts the user to give new text to use for this step, and replaces
// the old step in the recipe with this updated step.
func editRecipeStep(args []string, menu *climenus.Menu) error {
	recipe, _, err := extractRecipeData(menu)
	if err != nil {
		return err
	}

	recipeStepIdx, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}
	// ingredients start at # 2
	// so recipeStep starts at 2 + len(ingredients)

	recipeStepIdx = recipeStepIdx - len(recipe.Ingredients) - 2

	prompt := "Provide new data for this recipe step:"
	input := climenus.UserInput(prompt, recipeStepValidator)

	recipe.Steps[recipeStepIdx] = input

	// re-initialize edit a recipe commands in case options changed
	initializeEditRecipeCommands(menu, recipe)
	return nil
}

// Function used for adding a recipe ingredient. Takes user input for a new ingredient to add
// then adds it to the recipe currently being edited
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

// Saves any changes made to the currently selected recipe
// as of now this is done by replacing the recipe that had been selected for editing
// with the newly updated recipe, within the JSON data file
// this can be reworked at a later date when a relational database is added for data storage
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

// Extracts a pointer to a recipe struct, and an int indicating the recipe index from
// the list of recipes, and returns these
func extractRecipeData(menu *climenus.Menu) (*Recipe, int, error) {
	menuData, ok := menu.Data.(*editARecipeMenuData)
	if !ok {
		return nil, -1, errors.New("type assertion failed: edit a recipe menu data is not in correct form")
	}

	recipe := menuData.Recipe
	idx := menuData.RecipeIdx

	return recipe, idx, nil
}
