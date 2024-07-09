package main

import (
	"encoding/json"
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
	mapDataCached := pokeapi.LocationData{}
	bytes, ok := c.Get(url)
	if ok {
		json.Unmarshal(bytes, &mapDataCached)
	} else {
		mapData, err := pokeapi.GetMaps(url)
		if err != nil {
			fmt.Printf("Error getting maps from internet: %s\n", err)
			return err
		}
		data, err := json.Marshal(&mapData)
		if err != nil {
			fmt.Printf("An error has occurred marshalling json for cache: %s\n", err)
			return err
		}
		c.Add(url, data)
		cfg.nextUrl = mapData.Next
		if mapData.Previous == nil {
			cfg.prevUrl = ""
		} else {
			cfg.prevUrl = *mapData.Previous
		}
		for _, val := range mapData.Results {
			fmt.Println(val.Name)
		}
		fmt.Printf("Number of cache entries: %d", len(c.Entries))
		return nil
	}
	cfg.nextUrl = mapDataCached.Next
	if mapDataCached.Previous == nil {
		cfg.prevUrl = ""
	} else {
		cfg.prevUrl = *mapDataCached.Previous
	}
	for _, val := range mapDataCached.Results {
		fmt.Println(val.Name)
	}
	fmt.Printf("Number of cache entries: %d", len(c.Entries))
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
