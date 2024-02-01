package main

import (
	"bufio"
	"fmt"
	"os"
)

type config struct {
	Next     string
	Previous string
}

func main() {
	// creating an instance of scanner
	scnr := bufio.NewScanner(os.Stdin)
	configuration := config{
		Next:     "https://pokeapi.co/api/v2/location-area/",
		Previous: "",
	}

	for {
		// getting input
		fmt.Print("pokedex > ")
		scnr.Scan()
		words := cleanInput(scnr.Text())

		command := words[0]

		if cmd, ok := getCommands()[command]; ok {
			err := cmd.callback(&configuration)
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
