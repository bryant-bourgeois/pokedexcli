package pokeapi

import (
	"encoding/json"
	"net/http"
)

type LocationData struct {
	Count    int       `json:"count"`
	Next     string    `json:"next"`
	Previous *string       `json:"previous"`
	Results  []LocationResults `json:"results"`
}
type LocationResults struct {
	Name string `json:"name"`
	URL  string `json:"url"`
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
