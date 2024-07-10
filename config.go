package main

import "internal/pokeapi"

type Config struct {
	nextUrl string
	prevUrl string
	Pokedex map[string]pokeapi.PokemonInfo
}
