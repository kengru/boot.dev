package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"internal/pokeapi"
)

type config struct {
	Client   pokeapi.PokeClient
	Pokemons map[string]pokeapi.PokemonStats
	Next     string
	Previous string
}

func main() {
	// creating an instance of scanner
	scnr := bufio.NewScanner(os.Stdin)

	configuration := config{
		Client:   pokeapi.NewPokeClient(time.Second*5, time.Second*5),
		Pokemons: map[string]pokeapi.PokemonStats{},
		Next:     "https://pokeapi.co/api/v2/location-area/",
		Previous: "",
	}

	for {
		// getting input
		fmt.Print("pokedex > ")
		scnr.Scan()
		words := cleanInput(scnr.Text())

		if len(words) == 0 {
			fmt.Println("You need to write a command.")
			fmt.Println("Write 'help' to get a list of commands")
			continue
		}

		command := words[0]
		parameter := ""

		if len(words) > 1 {
			parameter = words[1]
		}

		if cmd, ok := getCommands()[command]; ok {
			err := cmd.callback(&configuration, parameter)
			if err != nil {
				fmt.Println(err)
			}
			continue

		} else {
			fmt.Println("Unknown command")
			fmt.Println("Write 'help' to get a list of commands")
			continue
		}
	}
}
