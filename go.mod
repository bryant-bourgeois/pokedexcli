module github.com/bryant-bourgeois/pokedexcli

go 1.22.0

require internal/pokeapi v1.0.0

replace internal/pokeapi => ./internal/pokeapi

require internal/pokecache v1.0.0

replace internal/pokecache => ./internal/pokecache
