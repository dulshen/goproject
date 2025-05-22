package main

import (
	"github.com/dulshen/goproject/climenus"
)

func selectRecipeMenu(executeFunc func([]string) error) (*climenus.Menu, error) {
	var menu climenus.Menu
	recipes, err := readRecipesJSON(jsonFileName)
	if err != nil {
		return nil, err
	}

	for _, recipe := range *recipes {
		menu.AddCommand(&climenus.Command{Description: recipe.Name, Name: "", Execute: executeFunc})
	}

	menu.AddCommand(&climenus.Command{Name: "back", Description: "", Execute: climenus.BackFunc})
	// menu.AddCommand(&climenus.Command{Name: "exit", Description: "Exit Program", Execute: climenus.ExitFunc})

	c1 := climenus.MenuColumn{ColWidth: 5, Label: "#", Type: "string"}
	c2 := climenus.MenuColumn{ColWidth: 1, Label: "", Type: "string"}
	c3 := climenus.MenuColumn{ColWidth: -20, Label: "Recipe Name", Type: "string"}

	menu.Columns = append(menu.Columns, c1, c2, c3)

	return &menu, nil
}
