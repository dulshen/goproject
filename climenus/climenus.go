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
		// first add option number, name, description
		formatStrings := []string{"%*d ", "%*s ", "%*s "}
		args := []interface{}{menu.Columns[0].ColWidth, command.OptionNumber,
			menu.Columns[1].ColWidth, command.Name,
			menu.Columns[2].ColWidth, command.Description}
		for i, additionalColumn := range command.AdditionalColumns {
			colIdx := i + 3
			formatString, err := menu.Columns[colIdx].typeFormatString()

			if err != nil {
				return err
			}

			formatStrings = append(formatStrings, formatString)
			args = append(args, menu.Columns[colIdx].ColWidth)
			args = append(args, additionalColumn)
		}
		formatStrings = append(formatStrings, "\n")

		fmt.Printf(strings.Join(formatStrings, ""), args...)

	}

	return nil
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
