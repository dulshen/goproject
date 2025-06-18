package climenus

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

const StringType = "string"
const IntType = "int"
const FloatType = "float"

const ExitProgram = "exit program command issued"
const BackCommand = "back command issued"

const optionNumberColIdx = 0
const nameColIdx = 1
const descriptionColIdx = 2

// struct representing a CLI Menu
type Menu struct {
	Commands    []*Command          // slice of commands for the menu (for printing in order)
	CommandsMap map[string]*Command // map of commands for the menu (for lookup of command by name)
	Columns     []MenuColumn        // slice with specifiers for the columns of this menu
	// ColWidths []int // slice containing column widths for each column of the menu
	// Labels []string // slice containing column labels for each column of the menu
	Instructions string      // instructions to print when menu is reached
	Data         interface{} // field for storing additional data that may need to be accessed by commands
}

// add a new command to the menu, adds the command to the list of commands
// maps the command's name to the command, and assigns a option number to the command
func (menu *Menu) AddCommand(command *Command) {
	if menu.CommandsMap == nil {
		menu.CommandsMap = make(map[string]*Command, 1)
	}
	// add the command to the menu command list and map, then update its option number
	menu.Commands = append(menu.Commands, command)
	menu.CommandsMap[command.Name] = command
	command.OptionNumber = len(menu.Commands)
}

// Prints the menu and shows the options for its commands
func (menu *Menu) ShowMenu() error {
	fmt.Println("\n\n" + menu.Instructions)
	// if menu just presents instructions then can return here
	if len(menu.Columns) == 0 {
		return nil
	}

	totalWidth := 0
	for _, col := range menu.Columns {
		formatString, _ := col.typeFormatString()
		fmt.Printf(formatString, col.ColWidth, col.Label)
		totalWidth += int(math.Abs(float64(col.ColWidth)))
	}
	fmt.Print("\n")

	for i := 0; i < totalWidth+5; i++ {
		fmt.Print("-")
	}
	fmt.Print("\n")
	for _, command := range menu.Commands {
		menu.renderCommand(command)
	}

	return nil
}

// Renders the text representing a command within the menu,
// handles formatting of command data into columns using format strings
// and splits lines that are too long for the column width into multiple rows
// then prints the processed text to the console.
// Returns the resulting format strings and args used for Printf (for testing purposes)
func (menu *Menu) renderCommand(command *Command) ([]string, [][]interface{}) {

	// get splits for the command, breaking up column context into rows by column width
	splits, height := getSplits(&menu.Columns, command)

	// pad shorter splits with empty strings to match max height column
	padWithEmptyStrings(&splits, height)

	// put together the args for the format string for each column in a []interface{}
	// (as a pair of column width, column contents for each column)
	// so that these can be used to supply the column widths and content to Printf
	fstringArgs := combineFstringArgs(&menu.Columns, &splits, height)

	// get the format strings for each column to use for Printf
	formatStrings, err := getFormatStrings(&menu.Columns)
	if err != nil {
		return nil, nil
	}

	// use the formatStrings and the args to render the text with Printf
	for row := range height {
		fmt.Printf(strings.Join(formatStrings, "")+"\n", fstringArgs[row]...)
	}

	return formatStrings, fstringArgs
}

// Takes the contents of the command data to be printed,
// and splits it into multiple rows if the data for that column is longer than
// the column width. Returns a slice of slices of strings representing these splits
// as well as an int representing the height in rows needed for this command.
func getSplits(columns *[]MenuColumn, command *Command) ([][]string, int) {
	// make an empty list of lists of strings for each column
	// after the contents of that columnn are split into rows
	height := 0
	n := len(*columns)
	// formatStrings := make([][]string, n)

	optionNumber := strconv.Itoa(command.OptionNumber)
	contentsOfDefaultColumns := map[int]string{
		optionNumberColIdx: optionNumber,
		nameColIdx:         command.Name,
		descriptionColIdx:  command.Description,
	}

	// TODO: add support for AdditionalColumns

	splits := make([][]string, n)
	// loop through each column and generate a list of strings for that column
	// split into separate rows of max length colWidth[i]
	for i := 0; i < n; i++ {
		// for width get abs value since formatting can use negatives to indicate left justify
		width := int(math.Abs(float64((*columns)[i].ColWidth)))
		s := contentsOfDefaultColumns[i]
		// break up the string into words to build splits one word at a time
		// (prevents breaking up a word in the middle when splitting lines)
		words := strings.Split(strings.TrimSpace(s), " ")

		// empty slice to use for splits for this column
		splits[i] = make([]string, 0)

		currentWidth := 0
		var sb strings.Builder
		for wordNum, word := range words {
			// if the word would exceed allowable width, then add current sb.String() as a split
			// and start a new one
			if len(word)+currentWidth > width {
				splits[i] = append(splits[i], sb.String()[:len(sb.String())-1])
				sb.Reset()
				currentWidth = 0
			}
			// after checking if split is needed, write current word to sb
			sb.WriteString(word)
			currentWidth += len(word)
			// if this isn't the last word, add a space
			if wordNum < len(words)-1 {
				sb.WriteString(" ")
				currentWidth += 1
			}

			// if it is the last word, then just add sb so far as a new split
			if wordNum == len(words)-1 {
				splits[i] = append(splits[i], sb.String())
			}
		}

		// height for the whole row should be the height of the max len split for the row
		height = max(height, len(splits[i]))
	}

	return splits, height
}

// Gets format strings (e.g. "%*s ", "%*d ") to use in the Printf call for rendering.
// Computes these by using the column type for each column of the menu.
// Returns the list of format strings.
func getFormatStrings(columns *[]MenuColumn) ([]string, error) {
	n := len(*columns)
	formatStrings := make([]string, n)

	for i := 0; i < n; i++ {
		typeString, err := (*columns)[i].typeFormatString()
		if err != nil {
			return nil, err
		}
		formatStrings[i] = typeString
	}

	return formatStrings, nil
}

// Pads splits in place with empty strings, based on the height in rows needed
// for this command (i.e. the height of the largest column for this command).
func padWithEmptyStrings(splits *[][]string, height int) {
	n := len(*splits)
	// For any columns that have fewer splits than the other columns
	// need to add empty strings at the end
	for i := 0; i < n; i++ {
		for len((*splits)[i]) < height {
			(*splits)[i] = append((*splits)[i], "")
		}
	}
}

// Combines the column widths for each column of menu with the splits containing the data
// to display in the column, by putting these into a slice of slices of interface{}
// so that each index in the slice of slices can be used with the ... operator
// to function as the additional arguments in the Printf call for rendering.
func combineFstringArgs(columns *[]MenuColumn, splits *[][]string, height int) [][]interface{} {
	n := len(*columns)
	// put together the args for the format string for each column in a []interface{}
	// so that these can be used to supply the column widths and content
	fstringArgs := make([][]interface{}, height)
	// loop by row first instead of by column
	for i := 0; i < height; i++ {
		fstringArgs[i] = make([]interface{}, 0)
		for j := 0; j < n; j++ {
			fstringArgs[i] = append(fstringArgs[i], (*columns)[j].ColWidth, (*splits)[j][i])
		}
	}

	return fstringArgs
}

// looks up a command by the option number, and returns the command if valid option
func (menu *Menu) CommandByOptionNumber(optionNumber int) (*Command, error) {
	// slice is 0 indexed, but optionNumber is 1 indexed so adjust
	index := optionNumber - 1

	if index >= len(menu.Commands) || index < 0 {
		return nil, errors.New("not a valid option number")
	}

	return menu.Commands[index], nil
}

// look up a command from the option number or name
// note if there is only one possible command for this menu, then
// just return that command
func (menu *Menu) Command(commandString string) (*Command, error) {
	// if len(menu.Commands) == 1 {
	// 	return menu.Commands[0], nil
	// }

	optionNumber, err := strconv.Atoi(commandString)
	// if it was an int use function to get command from option number
	if err == nil {
		return menu.CommandByOptionNumber(optionNumber)
	}

	// otherwise get command from the map
	command, isValidCommand := menu.CommandsMap[commandString]
	if !isValidCommand {
		return nil, errors.New("invalid command")
	}

	return command, nil

}

// a default validator to use for menu commands
// checks if the provided input matches a valid command name,
// or if it matches a valid command number.
func (menu *Menu) commandValidator(commandString string) (bool, error) {
	if commandString == "back" {
		return true, nil
	}
	// strings := strings.Split(commandString, " ")
	// if len(strings) > 1 {
	// 	return false, errors.New("this command does not support additional arguments")
	// }
	// command := strings[0]
	optionNumber, err := strconv.Atoi(commandString)
	if err != nil {
		// see if input s is in CommandsMap for non-numeric
		_, isValid := menu.CommandsMap[commandString]
		if !isValid {
			return false, errors.New("not a valid command")
		}
		return true, nil
	}
	// for numeric check if index is valid
	isValid := (optionNumber > 0 && optionNumber <= len(menu.Commands))
	if !isValid {
		return false, errors.New("optionNumber is outside valid range")
	}
	return true, nil
}

func BackFunc(args []string, menu *Menu) error {
	return errors.New(BackCommand)
}

func ExitFunc(args []string, menu *Menu) error {
	return errors.New(ExitProgram)
}

// main loop for a CLI menu, takes user input until a valid command
// is issued or user elects to go back or exit the program
func (menu *Menu) MenuLoop() error {

	commandString := ""

	for commandString != "back" && commandString != "exit" {
		menu.ShowMenu()
		// prompts := []string{""}
		// validators := []func(string, []string) (bool, error){menu.commandValidator}
		input := UserInput("", menu.commandValidator)
		// inputStrings := strings.Split(input, " ")
		args := strings.Split(input, " ")
		// commandString = inputStrings[0]
		commandString = args[0]
		command, err := menu.Command(commandString)
		// args := inputStrings[1:]

		if err != nil {
			fmt.Println(err.Error())
		} else {
			err := command.Execute(args, menu)
			if err != nil && err.Error() == ExitProgram {
				return err
			} else if err != nil && err.Error() == BackCommand {
				return nil
			} else if err != nil {
				// fmt.Println("debug1")
				fmt.Println(err.Error())
			}
		}
	}

	// else command was "back", so just return
	return nil
}

// Specifier for a MenuColumn with header label, width, and type
type MenuColumn struct {
	ColWidth int    // print width for this column
	Type     string // type of this column (string, int, float)
	Label    string // label to print for this column header
}

// returns the format string to use for the column type
func (c *MenuColumn) typeFormatString() (string, error) {
	fmtStr := ""
	if c.Type == "string" {
		fmtStr = "%*s "
	} else if c.Type == "int" {
		fmtStr = "%*d "
	} else if c.Type == "float" {
		fmtStr = "%*f "
	} else {
		return "", errors.New("invalid column type")
	}

	return fmtStr, nil
}

type Command struct {
	// option number for this command, can be used to select this option
	// from menu. Value automatically set by AddCommand
	OptionNumber int
	// name of the command, can be used to select this command from menu
	Name string
	// longer description of the command
	Description string
	// additional columns to print in the menu for this command, if needed
	AdditionalColumns []string
	// function to execute when this command is issued
	Execute func(args []string, menu *Menu) error
	// SubMenu that should be displayed when this command is issued
	SubMenu *Menu
}

func UserInput(prompt string, validator func(string) (bool, error)) string {

	scanner := bufio.NewScanner(os.Stdin)

	isValid := false
	err := error(nil)
	input := ""
	// inputSlice := make([]string, 0)
	for !isValid {
		fmt.Println(prompt)
		scanner.Scan()
		input = strings.TrimSpace(scanner.Text())
		// inputSlice = strings.Split(input, " ")
		// commandString := inputSlice[0]
		// args := inputSlice[1:]
		isValid, err = validator(input)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	return input
}

func UserInputLoop(prompt string, exitLoop string, validator func(string) (bool, error)) []string {
	input := ""
	inputStrings := make([]string, 0)

	for input != exitLoop {
		input = UserInput(prompt, validator)
		if input != exitLoop {
			inputStrings = append(inputStrings, input)
		}
	}

	return inputStrings
}
