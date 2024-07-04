package pokeapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"poke-repl/internal/config"
	"sort"
)

var Explorer LocationAreaResult

type LocationAreaResult struct {
	ID                   int    `json:"id"`
	Name                 string `json:"name"`
	GameIndex            int    `json:"game_index"`
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
	Location struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Names []struct {
		Name     string `json:"name"`
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
			MaxChance        int `json:"max_chance"`
			EncounterDetails []struct {
				MinLevel        int   `json:"min_level"`
				MaxLevel        int   `json:"max_level"`
				ConditionValues []any `json:"condition_values"`
				Chance          int   `json:"chance"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
			} `json:"encounter_details"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

func (l *LocationAreaResult) Explore(url string, area string, cfg *config.Config) ([]string, error) {
	if cached, ok := pokeCache.Get(area); ok {
		var cachedPokemons []string
		err := json.Unmarshal(cached, &cachedPokemons)
		if err != nil {
			return nil, fmt.Errorf("error deserializing cached data: %w", err)
		}
		return cachedPokemons, nil
	}
	res, err := http.Get(url + area)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("error fetching %s: %s", url, res.Status)
	}
	var result LocationAreaResult
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	var pokemons []string
	for _, pokemon := range result.PokemonEncounters {
		pokemons = append(pokemons, pokemon.Pokemon.Name)
	}
	sort.Strings(pokemons)
	pokemonsData, err := json.Marshal(pokemons)
	if err != nil {
		return nil, err
	}
	pokeCache.Set(area, pokemonsData)
	return pokemons, nil
}
