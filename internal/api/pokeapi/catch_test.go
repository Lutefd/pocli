package pokeapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"poke-repl/internal/config"
	"testing"
)

func TestCatchPokemon(t *testing.T) {
	cfg := &config.Config{}

	pokemon := "pikachu"
	expectedURL := "https://pokeapi.co/api/v2/pokemon/" + pokemon

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.String() != expectedURL {
			t.Errorf("expected URL %s, got %s", expectedURL, r.URL.String())
		}

		response := PokemonResult{
			ID:   25,
			Name: "pikachu",
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()
	p := &PokemonResult{}

	result, err := p.CatchPokemon("https://pokeapi.co/api/v2/pokemon/", pokemon, cfg)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if result.Name != "pikachu" {
		t.Errorf("expected name %s, got %s", "pikachu", result.Name)
	}
}
func TestCatchPokemon_CachedData(t *testing.T) {
	cfg := &config.Config{}

	// Pre-populate the cache
	pokemonName := "bulbasaur"
	cachedPokemon := PokemonResult{
		ID:   1,
		Name: "bulbasaur",
	}
	cachedData, err := json.Marshal(cachedPokemon)
	if err != nil {
		t.Fatalf("Failed to marshal cached data: %v", err)
	}
	pokeCache.Set(pokemonName, cachedData)

	// Mock server to ensure no HTTP request is made
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("Expected no HTTP request")
	}))
	defer server.Close()

	// Test CatchPokemon with cached data
	p := &PokemonResult{}
	result, err := p.CatchPokemon(server.URL, pokemonName, cfg)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Verify the result matches the cached data
	if result.Name != cachedPokemon.Name {
		t.Errorf("Expected name %s, got %s", cachedPokemon.Name, result.Name)
	}
}
