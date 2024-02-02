package main

import (
	"errors"
	"fmt"
	"os"
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

func commandMapB(conf *config) error {
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

func commandExit(_ *config) error {
	os.Exit(0)
	return nil
}
