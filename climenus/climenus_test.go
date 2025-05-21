package climenus

import (
	"testing"
)

var dummyMenus = []struct {
	description string
	name        string
}{
	{description: "dummy1", name: "dummy1"},
	{description: "dummy2", name: "dummy2"},
	{description: "dummy3", name: "dummy3"},
}

func TestMenuDescriptions(t *testing.T) {
	testCases := []struct {
		name         string
		description  string
		expected     string
		dummyOptions []struct {
			description string
			name        string
		}
	}{
		{
			name:         "test_description_1",
			description:  "descrip1",
			expected:     "descrip1",
			dummyOptions: nil,
		},
		{
			name:         "test_description_2",
			description:  "descrip2",
			expected:     "descrip2",
			dummyOptions: nil,
		},
		{
			name:         "test_description_3",
			description:  "descrip3",
			expected:     "descrip3",
			dummyOptions: dummyMenus,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var menu Menu

			for _, option := range tc.dummyOptions {
				c := Command{
					Name:        option.name,
					Description: option.description,
				}
				menu.AddCommand(&c)
			}

			c := Command{
				Name:        "name1",
				Description: tc.description,
			}
			menu.AddCommand(&c)

			expectedPosition := len(tc.dummyOptions)

			if menu.Commands[expectedPosition].Description != tc.expected {
				t.Errorf(`Expected %v, got %v`, tc.expected, menu.Commands[expectedPosition].Description)
			}
		})
	}
}

func TestMenuOptionNumbers(t *testing.T) {
	testCases := []struct {
		name         string
		expected     int
		dummyOptions []struct {
			description string
			name        string
		}
	}{
		{
			name:         "test_number_1",
			expected:     1,
			dummyOptions: nil,
		},
		{
			name:         "test_number_4",
			expected:     4,
			dummyOptions: dummyMenus,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var menu Menu

			for _, option := range tc.dummyOptions {
				menu.AddCommand(&(Command{
					Name:        option.name,
					Description: option.description,
				}))
			}

			menu.AddCommand(&(Command{
				Name:        "name1",
				Description: "descrip1",
			}))
			expectedPosition := len(tc.dummyOptions)

			if menu.Commands[expectedPosition].OptionNumber != tc.expected {
				t.Errorf(`Expected %v, got %v`, tc.expected, menu.Commands[expectedPosition].OptionNumber)
			}
		})
	}
}

func TestMenuNames(t *testing.T) {
	testCases := []struct {
		name         string
		optionName   string
		expected     string
		dummyOptions []struct {
			description string
			name        string
		}
	}{
		{
			name:         "test_name_1",
			optionName:   "name1",
			expected:     "name1",
			dummyOptions: nil,
		},
		{
			name:         "test_name_2",
			optionName:   "name2",
			expected:     "name2",
			dummyOptions: dummyMenus,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var menu Menu

			for _, option := range tc.dummyOptions {
				menu.AddCommand(&(Command{
					Name:        option.name,
					Description: option.description,
				}))
			}

			menu.AddCommand(&(Command{
				Name:        tc.optionName,
				Description: "descrip",
			}))

			expectedPosition := len(tc.dummyOptions)

			if menu.Commands[expectedPosition].Name != tc.expected {
				t.Errorf(`Expected %v, got %v`, tc.expected, menu.Commands[expectedPosition].Name)
			}
		})
	}
}

func TestCommandByNumber(t *testing.T) {
	testCases := []struct {
		name             string
		numberOfCommands int
		expected         string
		dummyOptions     []struct {
			description string
			name        string
		}
	}{
		{
			name:             "testOption1",
			numberOfCommands: 1,
			expected:         dummyMenus[0].name,
			dummyOptions:     dummyMenus[:1],
		},
		{
			name:             "testOption2",
			numberOfCommands: 2,
			expected:         dummyMenus[1].name,
			dummyOptions:     dummyMenus[:2],
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var menu Menu

			for _, option := range tc.dummyOptions {
				menu.AddCommand(&Command{Name: option.name, Description: option.description})
			}

			c, err := menu.CommandByOptionNumber(tc.numberOfCommands)

			if err != nil {
				t.Errorf("got error %s", err.Error())
			}

			if c.Name != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, c.Name)
			}
		})
	}
}

func TestGetCommand(t *testing.T) {
	testCases := []struct {
		name         string
		commandName  string
		expected     string
		dummyOptions []struct {
			description string
			name        string
		}
	}{
		{
			name:         "testGetCommand1",
			commandName:  "TestCommand1",
			expected:     "TestCommand1",
			dummyOptions: nil,
		},
		{
			name:         "testGetCommand2",
			commandName:  "TestCommand2",
			expected:     "TestCommand2",
			dummyOptions: dummyMenus,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var menu Menu
			for _, option := range tc.dummyOptions {
				menu.AddCommand(&Command{Name: option.name, Description: option.description})
			}

			menu.AddCommand(&Command{Name: tc.commandName, Description: "testDescription"})

			c, err := menu.Command(tc.commandName)

			if err != nil {
				t.Errorf("got error %v", err.Error())
			}

			if c.Name != tc.expected {
				t.Errorf("got %v, expected %v", c.Name, tc.expected)
			}
		})
	}
}
