package main

import (
	"errors"
	"fmt"
	"os"

	"internal/pokeapi"
)

func commandHelp(_ *config) error {
	fmt.Printf("Welcome to the Pokedex!\n\n")
	fmt.Printf("Here is a list of commands and what they do:\n\n")
	for _, v := range getCommands() {
		fmt.Printf("%v: %v\n", v.name, v.description)
	}
	return nil
}

func commandMap(conf *config) error {
	if conf.Next == "" {
		return errors.New("There are no more results! Use 'mapb' to go back.")
	}
	n, p := pokeapi.GetLocationAreas(conf.Next)

	conf.Next = n
	conf.Previous = p

	return nil
}

func commandMapB(conf *config) error {
	if conf.Previous == "" {
		return errors.New("There are no previous results! Use 'map' to get some.")
	}
	n, p := pokeapi.GetLocationAreas(conf.Previous)

	conf.Next = n
	conf.Previous = p

	return nil
}

func commandExit(_ *config) error {
	os.Exit(0)
	return nil
}
