package climenus

import (
	"testing"
)

var dummyMenus = []struct {
	description string
	menuKey     string
}{
	{description: "dummy1", menuKey: "dummy1"},
	{description: "dummy2", menuKey: "dummy2"},
	{description: "dummy3", menuKey: "dummy3"},
}

func TestBuildMenuDescription(t *testing.T) {
	testCases := []struct {
		name         string
		description  string
		expected     string
		dummyOptions []struct {
			description string
			menuKey     string
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
			var menuOptions []map[string]string

			for _, option := range tc.dummyOptions {
				menuOptions = append(menuOptions, map[string]string{"description": option.description, "menuKey": option.menuKey})
			}

			menuOptions = append(menuOptions, map[string]string{"description": tc.description, "menuKey": "key1"})

			results, err := BuildMenu(menuOptions)

			expectedPosition := len(tc.dummyOptions)

			if results[expectedPosition].Description != tc.expected {
				t.Errorf(`Expected %v, got %v`, tc.expected, results[expectedPosition].Description)
			}

			if err != nil {
				t.Errorf(`got error: %v`, err)
			}

		})
	}
}

func TestBuildMenuNumbers(t *testing.T) {
	testCases := []struct {
		name         string
		expected     int
		dummyOptions []struct {
			description string
			menuKey     string
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
			var menuOptions []map[string]string

			for _, option := range tc.dummyOptions {
				menuOptions = append(menuOptions, map[string]string{"description": option.description, "menuKey": option.menuKey})
			}

			menuOptions = append(menuOptions, map[string]string{"description": "descrip1", "menuKey": "key1"})

			results, err := BuildMenu(menuOptions)

			expectedPosition := len(tc.dummyOptions)

			if results[expectedPosition].OptionNumber != tc.expected {
				t.Errorf(`Expected %v, got %v`, tc.expected, results[expectedPosition].OptionNumber)
			}

			if err != nil {
				t.Errorf(`got error: %v`, err)
			}

		})
	}
}

func TestBuildMenuKeys(t *testing.T) {
	testCases := []struct {
		name         string
		menuKey      string
		expected     string
		dummyOptions []struct {
			description string
			menuKey     string
		}
	}{
		{
			name:         "test_menukey_1",
			menuKey:      "key1",
			expected:     "key1",
			dummyOptions: nil,
		},
		{
			name:         "test_menukey_2",
			menuKey:      "key2",
			expected:     "key2",
			dummyOptions: dummyMenus,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var menuOptions []map[string]string

			for _, option := range tc.dummyOptions {
				menuOptions = append(menuOptions, map[string]string{"description": option.description, "menuKey": option.menuKey})
			}

			menuOptions = append(menuOptions, map[string]string{"description": "descrip", "menuKey": tc.menuKey})

			results, err := BuildMenu(menuOptions)

			expectedPosition := len(tc.dummyOptions)

			if results[expectedPosition].MenuKey != tc.expected {
				t.Errorf(`Expected %v, got %v`, tc.expected, results[expectedPosition].MenuKey)
			}

			if err != nil {
				t.Errorf(`got error: %v`, err)
			}

		})
	}
}
