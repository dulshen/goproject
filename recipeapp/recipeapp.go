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

	var selection string

	for selection != "exit" {
		fmt.Print(menuString)
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
			AddRecipeLoop()
		}
	}

}
