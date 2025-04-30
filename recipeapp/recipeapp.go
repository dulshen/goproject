package main

import (
	"fmt"
	"log"

	"github.com/dulshen/goproject/climenus"
)

func main() {
	log.SetPrefix("climenu: ")
	log.SetFlags(0)

	menuData, err := climenus.BuildMenu(mainMenuOptions)

	if err != nil {
		log.Fatal("failed to build menu")
	}

	menuString, err := climenus.PrintMenu(menuData)

	if err != nil {
		log.Fatal("failed to print menu")
	}

	fmt.Print(menuString)

	var selection string

	isInput := false

	for !isInput {
		_, err = fmt.Scan(&selection)

		if err == nil {
			selection, err = climenus.MakeSelection(menuData, selection)
		}

		if err != nil {
			fmt.Println("Selection does not exist. Please enter a valid selection.")
		} else {
			isInput = true
		}
	}

	// fmt.Print(selection)

	if selection == "add" {
		PrintAddInstructions()

		// var input string

		// _, err = fmt.Scan(&input)

		// if err == nil {
		// 	fmt.Println("error")
		// }

		ingredient, err := ParseIngredient()

		if err == nil {
			fmt.Printf("ingredient: %v, quantity: %v, unit: %v", ingredient.name, ingredient.quantity, ingredient.unit)
		} else {
			fmt.Print(err)
		}

	}

}
