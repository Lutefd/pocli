package pokeapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"poke-repl/internal/config"
)

var Location LocationResult

type LocationList []LocationInfo
type LocationInfo struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type LocationResult struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous any    `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func (l *LocationResult) GetLocation(url string, cfg *config.Config) (LocationList, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("error fetching %s: %s", url, res.Status)
	}
	var result LocationResult
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	cfg.NextUrl = result.Next
	if prevUrl, ok := result.Previous.(string); ok {
		cfg.PreviousUrl = prevUrl
	}
	var locations LocationList
	for _, item := range result.Results {
		locations = append(locations, LocationInfo{Name: item.Name, URL: item.URL})
	}

	return locations, nil
}
