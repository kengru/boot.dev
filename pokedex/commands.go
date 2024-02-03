package main

import (
	"errors"
	"fmt"
	"math/rand"
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

func commandCatch(conf *config, pokemon string) error {
	if pokemon == "" {
		return errors.New("You should include the name of the pokemon you want to catch!")
	}

	results, err := conf.Client.GetPokemon(pokemon)
	if err != nil {
		return err
	}

	fmt.Printf("Throwing a pokeball at %s...\n", pokemon)

	finalChance := (1.0 / float64(results.BaseExperience)) * 20.0
	luck := rand.Float64()

	if finalChance > luck {
		fmt.Printf("%s was caught!\n", pokemon)
		conf.Pokemons[pokemon] = results
		return nil
	}

	fmt.Printf("%s escaped!\n", pokemon)
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

func commandInspect(conf *config, pokemon string) error {
	if pokemon == "" {
		return errors.New("You should write the name of your pokemon!")
	}

	results, ok := conf.Pokemons[pokemon]
	if !ok {
		return errors.New(fmt.Sprintf("You have not caught a %s", pokemon))
	}

	fmt.Printf("Name: %s\n", results.Name)
	fmt.Printf("Height: %v\n", results.Height)
	fmt.Printf("Stats: \n")
	for _, val := range results.Stats {
		fmt.Printf("\t-%s: %v\n", val.Stat.Name, val.BaseStat)
	}

	return nil
}

func commandPokedex(conf *config, _ string) error {
	if len(conf.Pokemons) < 1 {
		return errors.New("You have not caught any pokemon!")
	}

	fmt.Println("Your pokedex:")
	for k := range conf.Pokemons {
		fmt.Printf("\t-%s\n", k)
	}

	return nil
}

func commandExit(_ *config, _ string) error {
	os.Exit(0)
	return nil
}
