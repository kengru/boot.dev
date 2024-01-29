package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// creating an instance of scanner
	scnr := bufio.NewScanner(os.Stdin)

	for {
		// getting input
		fmt.Print("pokedex > ")
		scnr.Scan()
		words := cleanInput(scnr.Text())

		command := words[0]

		if cmd, ok := getCommands()[command]; ok {
			err := cmd.callback()
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
