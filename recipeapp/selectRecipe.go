package main

import (
	"fmt"
	"strconv"

	"github.com/dulshen/goproject/climenus"
)

func selectRecipeMenuNew(recipes *[]Recipe) (string, error) {
	viewMenuData, err := buildSelectRecipeMenu(recipes)
	if err != nil {
		return "", err
	}

	menuString, err := climenus.PrintMenu(*viewMenuData)
	if err != nil {
		return "", err
	}

	var selection string
	fmt.Print(menuString)
	isInput := false
	// loop until valid user input
	for !isInput {
		_, err = fmt.Scan(&selection)
		// check if valid selection
		if err == nil {
			selection, err = climenus.MakeSelection(*viewMenuData, selection)
		}
		// if err for either scanning or MakeSelection, then wasn't a valid selection
		if err != nil {
			fmt.Println("Selection does not exist. Please enter a valid selection.")
		} else {
			isInput = true
		}
	}
	return selection, nil
}

func buildSelectRecipeMenu(recipes *[]Recipe) (*[]climenus.MenuOptionData, error) {
	var viewMenuOptions []map[string]string

	for i, recipe := range *recipes {
		option := map[string]string{
			"description": recipe.Name,
			"menuKey":     strconv.Itoa(i + 1), // no support for menuKey for view recipe at the moment, just use number
		}
		viewMenuOptions = append(viewMenuOptions, option)
	}
	// add one more option to menu for going back
	viewMenuOptions = append(viewMenuOptions, map[string]string{
		"description": "Return to main menu",
		"menuKey":     "back",
	})

	viewMenuData, err := climenus.BuildMenu(viewMenuOptions)
	if err != nil {
		return nil, err
	}
	return &viewMenuData, nil
}
