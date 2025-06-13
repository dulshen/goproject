package main

import (
	"github.com/dulshen/goproject/climenus"
)

// The main loop for selecting a recipe, used by the view recipe, edit recipe, and delete recipe functions.
// Prints a list of recipes for the user to select from, then calls the appropriate function (view, edit, delete)
// as indicated by the executeFunc argument, with the selected recipe index as an argument
func selectRecipeLoop(executeFunc func([]string, *climenus.Menu) error, instructions string) error {
	var menu climenus.Menu
	recipes, err := readRecipesJSON(jsonFileName)
	if err != nil {
		return err
	}

	InitializeSelectRecipeCommands(&menu, recipes, executeFunc)

	c1 := climenus.MenuColumn{ColWidth: 5, Label: "#", Type: "string"}
	c2 := climenus.MenuColumn{ColWidth: 1, Label: "", Type: "string"}
	c3 := climenus.MenuColumn{ColWidth: -20, Label: "Recipe Name", Type: "string"}

	menu.Columns = append(menu.Columns, c1, c2, c3)

	menu.Instructions = instructions

	err = menu.MenuLoop()
	if err != nil {
		return err
	}

	return nil

}

// Initializes the select recipe menu with commands for each recipe in the recipe data
func InitializeSelectRecipeCommands(
	menu *climenus.Menu, recipes *[]Recipe, executeFunc func([]string, *climenus.Menu) error,
) error {
	menu.Commands = []*climenus.Command{}
	menu.CommandsMap = map[string]*climenus.Command{}

	for _, recipe := range *recipes {
		menu.AddCommand(&climenus.Command{Description: recipe.Name, Name: "", Execute: executeFunc})
	}

	menu.AddCommand(&climenus.Command{Name: "back", Description: "", Execute: climenus.BackFunc})

	return nil

}
