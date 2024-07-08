package main

import (
	"bufio"
	"fmt"
	"internal/pokecache"
	"os"
	"strings"
	"time"
)

func main() {
	commands := instantiateCommands()
	config := Config{}
	cache := pokecache.NewCache(5 * time.Minute)
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("pokedex > ")
		scanner.Scan()
		text := scanner.Text()
		args := strings.Split(text, " ")
		command, ok := commands[args[0]]
		if !ok {
			fmt.Println("Invalid command, try 'help' for usage details")
		}
		err := command.callback(&config, cache)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		}
		fmt.Println()
	}
}
