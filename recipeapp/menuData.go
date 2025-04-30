package main

const desc = "description"
const menuKey = "menuKey"

const add = "add"
const view = "view"
const edit = "edit"
const del = "del"

// maps menuKeys to their respective description
var mapKeyToDescriptions = map[string]string{
	add:  "Add Recipe",
	view: "View Recipe",
	edit: "Edit Recipe",
	del:  "Delete Recipe",
}

// map containing descriptions of each menu option, as well as a menuKey for each option
var mainMenuOptions = []map[string]string{
	{desc: mapKeyToDescriptions[add], menuKey: add},
	{desc: mapKeyToDescriptions[view], menuKey: view},
	{desc: mapKeyToDescriptions[edit], menuKey: edit},
	{desc: mapKeyToDescriptions[del], menuKey: del},
}
