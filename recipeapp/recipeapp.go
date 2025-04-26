package main

import (
	"fmt"
	"log"

	"github.com/dulshen/goproject/climenus"
)

func main() {
	log.SetPrefix("climenu: ")
	log.SetFlags(0)

	menuData, err := climenus.BuildMenu(mainMenuOptions)

	if err != nil {
		log.Fatal("failed to build menu")
	}

	fmt.Printf(climenus.PrintMenu(menuData))
}
