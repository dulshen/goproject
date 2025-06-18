package main

import (
	"log"

	"github.com/dulshen/goproject/climenus"
)

const mainMenuInstructions = "----------------------------------------------------------\n" +
	"Please select an option from the menu below:\n" +
	"-----------------------------------------------------------\n\n"

// struct describing a recipe
type Recipe struct {
	Name        string       // name of the recipe
	Ingredients []Ingredient // ingredients list
	Steps       []string     // steps for the recipe
}

// struct describing an Ingredient
type Ingredient struct {
	Name     string  // name of the ingredient
	Quantity float32 // quantity of the ingredient
	Unit     string  // unit of quantity
}

func registerExitCommand(menu *climenus.Menu) {
	menu.AddCommand(&climenus.Command{Name: exit, Description: "Exit Program", Execute: climenus.ExitFunc})
}

// Starts the program
// initializes the json data storage file if needed, then runs the main menu's loop
func main() {

	initializeJSONFile(jsonFileName, jsonDirectoryName, false)

	log.SetPrefix("climenu: ")
	log.SetFlags(0)

	mainMenu := initializeMenu()
	mainMenu.MenuLoop()

}

// Initializes the main menu, setting up the columns and registering the commands
func initializeMenu() *climenus.Menu {
	var menu climenus.Menu

	menu.Instructions = mainMenuInstructions

	optionNumberCol := climenus.MenuColumn{
		ColWidth: optionNumberColWidth,
		Type:     climenus.StringType,
		Label:    optionNumberLabel,
	}
	commandNameCol := climenus.MenuColumn{
		ColWidth: commandNameColWidth,
		Type:     climenus.StringType,
		Label:    commandNameLabel,
	}
	descriptionCol := climenus.MenuColumn{ColWidth: descriptionColWidth,
		Type:  climenus.StringType,
		Label: descriptionLabel,
	}

	menu.Columns = append(menu.Columns, optionNumberCol, commandNameCol, descriptionCol)

	registerAddRecipeCommand(&menu)
	registerViewRecipeCommand(&menu)
	registerEditRecipeCommand(&menu)
	registerDeleteRecipeCommand(&menu)
	registerExitCommand(&menu)

	return &menu
}
