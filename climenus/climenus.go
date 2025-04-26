package climenus

import (
	"errors"
	"fmt"
	"strings"
)

type MenuOptionData struct {
	OptionNumber int
	Description  string
	MenuKey      string
}

const optionNumberLabel = "#"
const descriptionLabel = "Description"
const menuKeyLabel = "key"
const col1width = 10
const col2width = 10
const col3width = 10

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

func PrintMenu(menuData []MenuOptionData) (string, error) {
	var stringBuilder strings.Builder

	stringBuilder.WriteString("Please select an option from the menu below:\n\n")
	stringBuilder.WriteString(fmt.Sprintf("%*s | %*s | %*s",
		col1width, optionNumberLabel, col2width, descriptionLabel, col3width, menuKeyLabel))
	stringBuilder.WriteString("-------------------------------------------------------------------------")

	for _, optionData := range menuData {
		stringBuilder.WriteString(fmt.Sprintf("%*d | %*s | %*s",
			col1width, optionData.OptionNumber, col2width, optionData.Description, col3width, optionData.MenuKey))
	}

	return stringBuilder.String(), nil

}
