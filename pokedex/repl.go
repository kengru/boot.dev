package main

import "strings"

func cleanInput(text string) []string {
	output := strings.ToLower(text)
	words := strings.Fields(output)
	return words
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config, string) error
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"catch": {
			name:        "catch",
			description: "Attempts to catch a pokemon, the chance to catch it depends on the pokemon base experience.\n\nUsage: catch <pokemon_name>",
			callback:    commandCatch,
		},
		"explore": {
			name:        "explore",
			description: "Shows the available pokemon in a given area.\n\nUsage: explore <location_area>",
			callback:    commandExplore,
		},
		"inspect": {
			name:        "inspect",
			description: "Shows the stats of pokemon you already own.\n\nUsage: inspect <pokemon_name>",
			callback:    commandInspect,
		},
		"map": {
			name:        "map",
			description: "Displays a list of 20 location areas within the world",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays a list of 20 location areas within the world",
			callback:    commandMapB,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Shows a list of the pokemon you have caught!",
			callback:    commandPokedex,
		},

		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the pokedex",
			callback:    commandExit,
		},
	}
}
