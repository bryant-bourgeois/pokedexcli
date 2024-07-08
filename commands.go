package main

import (
	"fmt"
	"internal/pokeapi"
	"internal/pokecache"
	"os"
)

type Command struct {
	name        string
	description string
	callback    func(*Config, *pokecache.Cache) error
}

func instantiateCommands() map[string]Command {
	return map[string]Command{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},

		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},

		"map": {
			name:        "map",
			description: "Print the next 20 map locations",
			callback:    commandMap,
		},

		"mapb": {
			name:        "mapb",
			description: "Print the previous 20 map locations",
			callback:    commandMapb,
		},
	}
}

func commandExit(cfg *Config, c *pokecache.Cache) error {
	os.Exit(0)
	return nil
}

func commandHelp(cfg *Config, c *pokecache.Cache) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	commands := instantiateCommands()
	for _, val := range commands {
		fmt.Printf("%s: %s\n", val.name, val.description)
	}
	return nil
}

func commandMap(cfg *Config, c *pokecache.Cache) error {
	url := ""
	if cfg.nextUrl == "" {
		url = "https://pokeapi.co/api/v2/location-area"
	} else {
		url = cfg.nextUrl
	}
	mapData, err := pokeapi.GetMaps(url)
	if err != nil {
		return err
	}
	cfg.nextUrl = mapData.Next
	if mapData.Previous == nil {
		cfg.prevUrl = ""
	} else {
		cfg.prevUrl = *mapData.Previous
	}
	for _, val := range mapData.Results {
		fmt.Println(val.Name)
	}
	return nil
}

func commandMapb(cfg *Config, c *pokecache.Cache) error {
	if cfg.prevUrl == "" {
		fmt.Println("No previous maps available...")
		return nil
	}
	mapData, err := pokeapi.GetMaps(cfg.prevUrl)
	if err != nil {
		return err
	}
	cfg.nextUrl = mapData.Next
	if mapData.Previous == nil {
		cfg.prevUrl = ""
	} else {
		cfg.prevUrl = *mapData.Previous
	}
	for _, val := range mapData.Results {
		fmt.Println(val.Name)
	}
	return nil
}
