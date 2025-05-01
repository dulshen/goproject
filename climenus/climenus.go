package climenus

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// struct representing a Menu option and its respective data
type MenuOptionData struct {
	OptionNumber int    // the option number for this option
	Description  string // description of the option
	MenuKey      string // menu key to enter to select option
}

const optionNumberLabel = "#"
const descriptionLabel = "Description"
const menuKeyLabel = "key"
const col1width = 2
const col2width = 20
const col3width = 5

// Builds a menu in the form of a slice of MenuOptionData, from a slice of maps of strings to strings
// which contains the description and menuKey of each option
func BuildMenu(menuOptions []map[string]string) ([]MenuOptionData, error) {

	var menuData []MenuOptionData

	if len(menuOptions) <= 0 {
		return nil, errors.New("no menu options were provided")
	}

	for i, option := range menuOptions {
		description := option["description"]
		menuKey := option["menuKey"]

		if description == "" || menuKey == "" {
			return nil, errors.New("description and menuKey cannot be empty")
		}

		optionData := MenuOptionData{OptionNumber: i + 1, Description: description, MenuKey: menuKey}
		menuData = append(menuData, optionData)
	}

	return menuData, nil
}

// Prints the menu represented by menuData
// menuData is a slice of MenuOptionData
func PrintMenu(menuData []MenuOptionData) (string, error) {

	if len(menuData) <= 0 {
		return "", errors.New("menuData cannot be empty")
	}

	var stringBuilder strings.Builder

	stringBuilder.WriteString("Please select an option from the menu below:\n\n")
	stringBuilder.WriteString(fmt.Sprintf("%*s | %*s | %*s\n",
		col1width, optionNumberLabel, -col2width, descriptionLabel, -col3width, menuKeyLabel))

	totalWidth := col1width + col2width + col3width + 10

	for range totalWidth {
		stringBuilder.WriteString("-")
	}
	stringBuilder.WriteString("\n")

	for _, optionData := range menuData {
		stringBuilder.WriteString(fmt.Sprintf("%*d   %*s   %*s\n",
			col1width, optionData.OptionNumber, -col2width, optionData.Description, -col3width, optionData.MenuKey))
	}

	return stringBuilder.String(), nil

}

// Takes user input selection and checks if it matches a valid selection from the MenuOptionData
func MakeSelection(menuData []MenuOptionData, selection string) (string, error) {
	// if menus were longer would be good to initialize with a map for O(1) lookup
	// for shorter menus probably fine to just search the slice
	for _, option := range menuData {
		optionNum, intErr := strconv.Atoi(selection)
		isNumeric := intErr == nil
		if (!isNumeric && option.MenuKey == selection) ||
			(isNumeric && option.OptionNumber == optionNum) {
			return option.MenuKey, nil
		}
	}

	// if got here then selection wasn't found
	return "", errors.New("selected option does not exist")
}
