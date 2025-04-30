package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// struct describing an Ingredient
type Ingredient struct {
	name     string // name of the ingredient
	quantity int    // quantity of the ingredient
	unit     string // unit of quantity
}

// Prints the instructions for the Add Recipe option
func PrintAddInstructions() {
	fmt.Println("Please enter recipe ingredients in the following format:")
	fmt.Println("Ingredient name, ingredient quantity, ingredient unit")
}

// Gets user input and parses an ingredient from it
// Returns an Ingredient, and an error if ingredient could not be parsed
func ParseIngredient() (Ingredient, error) {

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()

	// if err != nil {
	// 	fmt.Println(err)
	// 	return Ingredient{}, errors.New("ingredient format error (use format: Ingredient name, ingredient quantity, ingredient unit)")
	// }

	items := strings.Split(input, ",")

	if len(items) > 3 || len(items) < 2 {
		return Ingredient{}, errors.New("must enter either ingredient, quantity or ingredient, quantity, unit")
	}

	name := strings.TrimSpace(items[0])
	quantity, err := strconv.Atoi(strings.TrimSpace(items[1]))
	var unit string
	if len(items) == 3 {
		unit = strings.TrimSpace(items[2])
	} else {
		unit = ""
	}

	if err != nil {
		return Ingredient{}, errors.New("ingredient quantity was not an integer")
	}

	ingredientData := Ingredient{
		name:     name,
		quantity: quantity,
		unit:     unit,
	}

	return ingredientData, nil
}
