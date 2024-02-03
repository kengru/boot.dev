package main

import (
	"errors"
	"fmt"
	"os"
)

func commandHelp(_ *config, _ string) error {
	fmt.Printf("Welcome to the Pokedex!\n\n")
	fmt.Printf("Here is a list of commands and what they do:\n\n")
	for _, v := range getCommands() {
		fmt.Printf("%v: %v\n", v.name, v.description)
	}
	return nil
}

func commandMap(conf *config, _ string) error {
	if conf.Next == "" {
		return errors.New("There are no more results! Use 'mapb' to go back.")
	}
	results, _ := conf.Client.GetLocationAreas(conf.Next)

	for _, r := range results.Results {
		fmt.Println(r.Name)
	}

	if results.Previous == nil {
		conf.Previous = ""
	} else {
		conf.Previous = *results.Previous
	}
	if results.Next == nil {
		conf.Next = ""
	} else {
		conf.Next = *results.Next
	}

	return nil
}

func commandMapB(conf *config, _ string) error {
	if conf.Previous == "" {
		return errors.New("There are no previous results! Use 'map' to get some.")
	}
	results, _ := conf.Client.GetLocationAreas(conf.Previous)

	for _, r := range results.Results {
		fmt.Println(r.Name)
	}

	if results.Previous == nil {
		conf.Previous = ""
	} else {
		conf.Previous = *results.Previous
	}
	if results.Next == nil {
		conf.Next = ""
	} else {
		conf.Next = *results.Next
	}

	return nil
}

func commandExplore(conf *config, location string) error {
	if location == "" {
		return errors.New("You should write the name of the location to explore!")
	}

	results, err := conf.Client.GetLocationArea(location)
	if err != nil {
		return err
	}

	fmt.Printf("Exploring %s...\n", location)
	fmt.Println("Found Pokemon:")
	for _, val := range results.PokemonEncounters {
		fmt.Printf("- %s\n", val.Pokemon.Name)
	}

	return nil
}

func commandExit(_ *config, _ string) error {
	os.Exit(0)
	return nil
}
