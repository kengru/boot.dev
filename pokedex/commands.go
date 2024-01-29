package main

import (
	"fmt"
	"os"
)

func commandHelp() error {
	fmt.Printf("Welcome to the Pokedex!\n\n")
	for _, v := range getCommands() {
		fmt.Printf("%v: %v\n", v.name, v.description)
	}
	return nil
}

func commandExit() error {
	os.Exit(0)
	return nil
}
