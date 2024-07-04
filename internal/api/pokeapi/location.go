package pokeapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"poke-repl/internal/cache"
	"poke-repl/internal/config"
	"time"
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

var pokeCache = cache.NewCache(time.Minute * 5)

func (l *LocationResult) GetLocation(url string, cfg *config.Config) (LocationList, error) {
	if cached, ok := pokeCache.Get(url); ok {
		var cachedLocations LocationList
		err := json.Unmarshal(cached, &cachedLocations)
		if err != nil {
			return nil, fmt.Errorf("error deserializing cached data: %w", err)
		}
		return cachedLocations, nil
	}
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
	locationsData, err := json.Marshal(locations)
	if err != nil {
		return nil, err
	}
	pokeCache.Set(url, locationsData)
	return locations, nil
}
