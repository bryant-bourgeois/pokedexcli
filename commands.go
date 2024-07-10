package main

import (
	"encoding/json"
	"fmt"
	"internal/pokeapi"
	"internal/pokecache"
	"math/rand"
	"os"
	"time"
)

type Command struct {
	name        string
	description string
	callback    func(*Config, *pokecache.Cache, []string) error
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

		"explore": {
			name:        "explore",
			description: "Find pokemon in a map locaiton. `explore <MAP_LOCATION>`",
			callback:    commandExplore,
		},

		"catch": {
			name:        "catch",
			description: "Attempt to catch pokemon. `catch <POKEMON>`",
			callback:    commandCatch,
		},

		"inspect": {
			name:        "inspect",
			description: "Inspect pokemon in your pokedex. `inspect <POKEMON>`",
			callback:    commandInspect,
		},
	}
}

func commandExit(cfg *Config, c *pokecache.Cache, args []string) error {
	os.Exit(0)
	return nil
}

func commandHelp(cfg *Config, c *pokecache.Cache, args []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	commands := instantiateCommands()
	for _, val := range commands {
		fmt.Printf("%s: %s\n", val.name, val.description)
	}
	return nil
}

func commandMap(cfg *Config, c *pokecache.Cache, args []string) error {
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
	return nil
}

func commandMapb(cfg *Config, c *pokecache.Cache, args []string) error {
	if cfg.prevUrl == "" {
		fmt.Println("No previous maps available...")
		return nil
	}
	mapDataCached := pokeapi.LocationData{}
	bytes, ok := c.Get(cfg.prevUrl)
	if ok {
		json.Unmarshal(bytes, &mapDataCached)
	} else {
		mapData, err := pokeapi.GetMaps(cfg.prevUrl)
		if err != nil {
			return err
		}
		data, err := json.Marshal(&mapData)
		if err != nil {
			fmt.Printf("An error has occurred marshalling json for cache: %s\n", err)
			return err
		}
		c.Add(cfg.prevUrl, data)
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
	cfg.nextUrl = mapDataCached.Next
	if mapDataCached.Previous == nil {
		cfg.prevUrl = ""
	} else {
		cfg.prevUrl = *mapDataCached.Previous
	}
	for _, val := range mapDataCached.Results {
		fmt.Println(val.Name)
	}
	return nil
}

func commandExplore(cfg *Config, c *pokecache.Cache, args []string) error {
	if len(args) < 1 {
		fmt.Println("No city to explore. See 'help' command.")
		return nil
	}
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", args[0])
	locDataCached := pokeapi.LocationDetails{}
	bytes, ok := c.Get(url)
	if ok {
		json.Unmarshal(bytes, &locDataCached)
	} else {
		locationData, err := pokeapi.GetLocation(url)
		if err != nil {
			return err
		}
		data, err := json.Marshal(&locationData)
		if err != nil {
			return err
		}
		c.Add(url, data)
		for _, val := range locationData.PokemonEncounters {
			fmt.Printf("- %s\n", val.Pokemon.Name)
		}
		return nil
	}
	for _, val := range locDataCached.PokemonEncounters {
		fmt.Printf("- %s\n", val.Pokemon.Name)
	}
	return nil
}

func commandCatch(cfg *Config, c *pokecache.Cache, args []string) error {
	if len(args) < 1 {
		fmt.Println("No pokemon to catch. See 'help' command.")
		return nil
	}
	_, ok := cfg.Pokedex[args[0]]
	if ok {
		fmt.Println("Pokemon has already been caught")
		return nil
	}
	fmt.Printf("Throwing a Pokeball at %s...\n", args[0])
	time.Sleep(1 * time.Second)
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s/", args[0])
	pokeDataCached := pokeapi.PokemonInfo{}
	bytes, ok := c.Get(url)
	if ok {
		json.Unmarshal(bytes, &pokeDataCached)
	} else {
		pokemonData, err := pokeapi.GetPokemon(url)
		if err != nil {
			return err
		}
		data, err := json.Marshal(&pokemonData)
		if err != nil {
			return err
		}
		c.Add(url, data)
		caught := AttemptCatch(pokemonData.BaseExperience)
		if caught {
			cfg.Pokedex[pokemonData.Name] = pokemonData
			fmt.Printf("%s was caught!\n", pokemonData.Name)
		} else {
			fmt.Printf("%s has escaped!\n", pokemonData.Name)
		}
		return nil
	}
	caught := AttemptCatch(pokeDataCached.BaseExperience)
	if caught {
		cfg.Pokedex[pokeDataCached.Name] = pokeDataCached
		fmt.Printf("%s was caught!\n", pokeDataCached.Name)
	} else {
		fmt.Printf("%s has escaped!\n", pokeDataCached.Name)
	}
	return nil
}

func AttemptCatch(bxp int) bool {
	comp := rand.Intn(370)
	return comp > bxp
}

func commandInspect(cfg *Config, c *pokecache.Cache, args []string) error {
	if len(args) < 1 {
		fmt.Println("No pokemon to catch. See 'help' command.")
		return nil
	}
	pokemon, ok := cfg.Pokedex[args[0]]
	if !ok {
		fmt.Println("you have not caught that pokemon")
		return nil
	}
	fmt.Printf("Name: %s\nHeight: %d\nWeight: %d\n", pokemon.Name, pokemon.Height, pokemon.Weight)
	fmt.Println("Stats:")
	for _, val := range pokemon.Stats {
		fmt.Printf("  -%s: %v\n", val.Stat.Name, val.BaseStat)
	}
	fmt.Println("Types:")
	for _, val := range pokemon.Types {
		fmt.Printf("  - %s\n", val.Type.Name)
	}
	return nil
}
