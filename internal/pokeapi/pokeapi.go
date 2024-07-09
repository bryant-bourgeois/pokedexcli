package pokeapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type LocationData struct {
	Count    int               `json:"count"`
	Next     string            `json:"next"`
	Previous *string           `json:"previous"`
	Results  []LocationResults `json:"results"`
}
type LocationResults struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type LocationDetails struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int   `json:"chance"`
				ConditionValues []any `json:"condition_values"`
				MaxLevel        int   `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

func GetMaps(url string) (LocationData, error) {
	resp, err := http.Get(url)
	if err != nil {
		return LocationData{}, err
	}
	locations := LocationData{}
	decoder := json.NewDecoder(resp.Body)
	defer resp.Body.Close()
	err = decoder.Decode(&locations)
	if err != nil {
		return LocationData{}, err
	}
	return locations, nil
}

func GetLocation(url string) (LocationDetails, error) {
	resp, err := http.Get(url)
	if err != nil {
		return LocationDetails{}, err
	}
	if resp.StatusCode == 404 {
		return LocationDetails{}, errors.New("location not found")
	}
	if resp.StatusCode > 400 {
		fmt.Println("Something wrong with location request")
		return LocationDetails{}, err
	}
	locations := LocationDetails{}
	decoder := json.NewDecoder(resp.Body)
	defer resp.Body.Close()
	err = decoder.Decode(&locations)
	if err != nil {
		return LocationDetails{}, err
	}
	return locations, nil
}
