package main

import (
	"github.com/dulshen/goproject/climenus"
)

// need to be able to access the menu if deleting a recipe, as this requires
// reinitializing the menu after recipe is removed
// could consider refactoring so that Exec functions pass a pointer to the current menu instead
var selectRecipeMenu *climenus.Menu

func selectRecipeLoop(executeFunc func([]string) error, instructions string) (*climenus.Menu, error) {
	var menu climenus.Menu
	selectRecipeMenu = &menu
	recipes, err := readRecipesJSON(jsonFileName)
	if err != nil {
		return nil, err
		// return err
	}

	InitializeSelectRecipeCommands(&menu, recipes, executeFunc)

	// for _, recipe := range *recipes {
	// 	menu.AddCommand(&climenus.Command{Description: recipe.Name, Name: "", Execute: executeFunc})
	// }

	// menu.AddCommand(&climenus.Command{Name: "back", Description: "", Execute: climenus.BackFunc})
	// menu.AddCommand(&climenus.Command{Name: "exit", Description: "Exit Program", Execute: climenus.ExitFunc})

	c1 := climenus.MenuColumn{ColWidth: 5, Label: "#", Type: "string"}
	c2 := climenus.MenuColumn{ColWidth: 1, Label: "", Type: "string"}
	c3 := climenus.MenuColumn{ColWidth: -20, Label: "Recipe Name", Type: "string"}

	menu.Columns = append(menu.Columns, c1, c2, c3)

	// if err != nil {
	// 	return err
	// }

	menu.Instructions = instructions

	// err = menu.ShowMenu()
	// if err != nil {
	// 	return nil, err
	// 	// return err
	// }

	// fmt.Println(&menu)

	err = menu.MenuLoop()
	if err != nil {
		return nil, err
		// return err
	}
	// return nil

	return &menu, nil
	// return nil
}

func InitializeSelectRecipeCommands(menu *climenus.Menu, recipes *[]Recipe, executeFunc func([]string) error) error {
	menu.Commands = []*climenus.Command{}
	menu.CommandsMap = map[string]*climenus.Command{}

	for _, recipe := range *recipes {
		menu.AddCommand(&climenus.Command{Description: recipe.Name, Name: "", Execute: executeFunc})
	}

	menu.AddCommand(&climenus.Command{Name: "back", Description: "", Execute: climenus.BackFunc})

	return nil

}
