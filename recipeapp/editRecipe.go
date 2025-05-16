package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"time"

	"github.com/dulshen/goproject/climenus"
)

const recipeNameIdx = "recipe name index"
const ingredientsStartIdx = "ingredients start index"
const ingredentsEndIdx = "ingredients end index"

func editRecipesMenu(filename string) error {
	log.SetPrefix("editRecipeMenu: ")

	recipes, err := readRecipesJSON(filename)

	if err != nil {
		log.Fatal(err.Error())
	}

	var selection string
	for selection != "back" {
		selection, err = selectRecipeMenuNew(recipes)
		if err != nil {
			fmt.Println(err.Error())
		}
		if selection != "back" {
			err = editOneRecipeMenu(recipes, selection)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}
	return nil
}

func editOneRecipeMenu(recipes *[]Recipe, selection string) error {
	selectionInt, err := strconv.Atoi(selection)
	if err != nil {
		return err
	}

	// switch from 1 indexed to 0 indexed
	index := selectionInt - 1

	recipe := (*recipes)[index]

	var option string
	for option != "back" && option != "save" {
		menuData, indicesMap, err := buildEditOptionsMenu(recipe)
		fmt.Println(indicesMap)
		if err != nil {
			return err
		}
		editMenuString, err := climenus.PrintMenu((*menuData))
		if err != nil {
			return err
		}
		fmt.Print(editMenuString)
		_, err = fmt.Scan(&option)
		if err != nil {
			return err
		}

		option, err := climenus.MakeSelection(*menuData, option)
		if err != nil {
			fmt.Println("invalid selection")
		} else {
			if option == "add" {
				addIngredient(&recipe)
			} else if option == "rename" {
				err = renameRecipe(&recipe)
				for err != nil {
					fmt.Println(err.Error())
					err = renameRecipe(&recipe)
				}
			} else if option == "remove" {
				removeIngredientMenu(&recipe)
			} else if option == "save" {
				err := saveChanges(recipe, index)
				if err != nil {
					fmt.Println(err.Error())
				} else {
					fmt.Println("Successfully saved changes to recipe.")
					time.Sleep(2 * time.Second)
				}
			}
		}
	}

	return nil
}

func removeIngredientMenu(recipe *Recipe) error {
	var removeIngMenuOptions []map[string]string

	removeIngMenuOptions = addIngredientMenuOptions(removeIngMenuOptions, *recipe, 1)

	// add one more option to menu for going back
	removeIngMenuOptions = append(removeIngMenuOptions, map[string]string{
		"description": "Return to main menu",
		"menuKey":     "back",
	})

	menuOptionData, err := climenus.BuildMenu(removeIngMenuOptions)
	if err != nil {
		return err
	}

	menuString, err := climenus.PrintMenu(menuOptionData)
	if err != nil {
		return err
	}

	fmt.Println(menuString)

	var choice string
	_, err = fmt.Scan(&choice)
	if err != nil {
		return err
	}

	err = removeIngredient(choice, recipe)
	if err != nil {
		return err
	}
	return nil
}

func removeIngredient(choice string, recipe *Recipe) error {
	choiceInt, err := strconv.Atoi(choice)
	index := choiceInt - 1
	if err != nil {
		return err
	}

	recipe.Ingredients = slices.Delete(recipe.Ingredients, index, index+1)
	return nil
}

func renameRecipe(recipe *Recipe) error {
	fmt.Println("Please provide an updated name for the recipe:")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	newName := scanner.Text()

	if len(newName) > maxRecipeNameLength {
		return fmt.Errorf("name must be %d characters or less", maxRecipeNameLength)
	}

	recipe.Name = newName

	return nil
}

func buildEditOptionsMenu(recipe Recipe) (*[]climenus.MenuOptionData, map[string]int, error) {
	indicesMap := make(map[string]int)
	var editMenuOptions []map[string]string

	// count := 0

	// add the recipe name as an option
	option := map[string]string{
		"description": recipe.Name,
		"menuKey":     "rename",
	}
	editMenuOptions = append(editMenuOptions, option)
	// count++

	count := 1
	indicesMap[recipeNameIdx] = count
	indicesMap[ingredientsStartIdx] = count + 1

	editMenuOptions = addIngredientMenuOptions(editMenuOptions, recipe, indicesMap[ingredientsStartIdx])

	indicesMap[ingredentsEndIdx] = count

	// add an option to add new ingredient
	editMenuOptions = append(editMenuOptions, map[string]string{
		"description": "add ingredient",
		"menuKey":     "add",
	})
	// add an option to remove an ingredient
	editMenuOptions = append(editMenuOptions, map[string]string{
		"description": "remove ingredient",
		"menuKey":     "remove",
	})
	// TODO: add recipe steps as options once those are added
	// add one more option to menu for going back
	editMenuOptions = append(editMenuOptions, map[string]string{
		"description": "Return to main menu",
		"menuKey":     "back",
	})
	editMenuOptions = append(editMenuOptions, map[string]string{
		"description": "Save changes to recipe",
		"menuKey":     "save",
	})

	editMenuData, err := climenus.BuildMenu(editMenuOptions)
	if err != nil {
		return nil, nil, err
	}

	return &editMenuData, indicesMap, nil
}

func addIngredient(recipe *Recipe) error {

	fmt.Println("Enter an ingredient to add")
	input := ingredientInput()
	ingredient, err := parseIngredient(input)

	for err != nil {
		fmt.Println(err.Error())
		input = ingredientInput()
		ingredient, err = parseIngredient(input)
	}

	recipe.Ingredients = append(recipe.Ingredients, ingredient)

	return nil
}

func addIngredientMenuOptions(ingredientMenuOptions []map[string]string, recipe Recipe, ingStartIdx int) []map[string]string {

	for i, ingredient := range recipe.Ingredients {
		option := map[string]string{
			"description": fmt.Sprintf("%s %d %s", ingredient.Name, ingredient.Quantity, ingredient.Unit),
			"menuKey":     strconv.Itoa(i + ingStartIdx), // no support for menuKey for view recipe at the moment, just use number
		}
		ingredientMenuOptions = append(ingredientMenuOptions, option)
	}

	return ingredientMenuOptions

}

func saveChanges(recipe Recipe, index int) error {
	err := replaceRecipe(recipe, jsonFileName, index)

	if err != nil {
		return err
	}

	return nil
}
