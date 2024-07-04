package pokeapi

import (
	"net/http"
	"net/http/httptest"
	"poke-repl/internal/config"
	"reflect"
	"sort"
	"testing"
)

func TestLocationAreaResult_Explore(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		responseBody := `{
			"id": 1,
			"name": "Location Area",
			"game_index": 1,
			"encounter_method_rates": [
				{
					"encounter_method": {
						"name": "Method 1",
						"url": "https://pokeapi.co/api/v2/method/1/"
					},
					"version_details": [
						{
							"rate": 50,
							"version": {
								"name": "Version 1",
								"url": "https://pokeapi.co/api/v2/version/1/"
							}
						},
						{
							"rate": 30,
							"version": {
								"name": "Version 2",
								"url": "https://pokeapi.co/api/v2/version/2/"
							}
						}
					]
				}
			],
			"location": {
				"name": "Location",
				"url": "https://pokeapi.co/api/v2/location/1/"
			},
			"names": [
				{
					"name": "Name 1",
					"language": {
						"name": "Language 1",
						"url": "https://pokeapi.co/api/v2/language/1/"
					}
				},
				{
					"name": "Name 2",
					"language": {
						"name": "Language 2",
						"url": "https://pokeapi.co/api/v2/language/2/"
					}
				}
			],
			"pokemon_encounters": [
				{
					"pokemon": {
						"name": "Pokemon 1",
						"url": "https://pokeapi.co/api/v2/pokemon/1/"
					},
					"version_details": [
						{
							"version": {
								"name": "Version 1",
								"url": "https://pokeapi.co/api/v2/version/1/"
							},
							"max_chance": 80,
							"encounter_details": [
								{
									"min_level": 5,
									"max_level": 10,
									"condition_values": [],
									"chance": 50,
									"method": {
										"name": "Method 1",
										"url": "https://pokeapi.co/api/v2/method/1/"
									}
								}
							]
						},
						{
							"version": {
								"name": "Version 2",
								"url": "https://pokeapi.co/api/v2/version/2/"
							},
							"max_chance": 60,
							"encounter_details": [
								{
									"min_level": 10,
									"max_level": 15,
									"condition_values": [],
									"chance": 30,
									"method": {
										"name": "Method 2",
										"url": "https://pokeapi.co/api/v2/method/2/"
									}
								}
							]
						}
					]
				}
			]
		}`
		_, _ = w.Write([]byte(responseBody))
	}))
	defer server.Close()

	mockConfig := &config.Config{}
	locationAreaResult := &LocationAreaResult{}
	pokemons, err := locationAreaResult.Explore(server.URL, "/area", mockConfig)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expectedPokemons := []string{"Pokemon 1"}
	sort.Strings(expectedPokemons)

	if !reflect.DeepEqual(pokemons, expectedPokemons) {
		t.Errorf("Expected pokemons %v, but got %v", expectedPokemons, pokemons)
	}
}

func TestLocationAreaResult_Explore_DecoderError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("{invalid JSON"))
	}))
	defer server.Close()

	mockConfig := &config.Config{}
	locationAreaResult := &LocationAreaResult{}
	_, err := locationAreaResult.Explore(server.URL, "area", mockConfig)
	if err == nil {
		t.Error("Expected a JSON decoder error, got nil")
	}
}

func TestLocationAreaResult_Explore_CachedData(t *testing.T) {
	mockConfig := &config.Config{}
	locationAreaResult := &LocationAreaResult{}
	pokeCache.Set("area", []byte(`["Cached Pokemon 1", "Cached Pokemon 2"]`))
	pokemons, err := locationAreaResult.Explore("https://pokeapi.co/api/v2/", "area", mockConfig)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	expectedPokemons := []string{"Cached Pokemon 1", "Cached Pokemon 2"}
	sort.Strings(expectedPokemons)
	if !reflect.DeepEqual(pokemons, expectedPokemons) {
		t.Errorf("Expected pokemons %v, but got %v", expectedPokemons, pokemons)
	}
}
