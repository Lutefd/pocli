package pokeapi

import "testing"

func TestAddPokemon(t *testing.T) {
	pokedex := NewPokedex()

	pokemon := PokemonResult{
		Name: "Pikachu",
	}

	pokedex.AddPokemon(pokemon)

	_, ok := pokedex.GetPokemon("Pikachu")
	if !ok {
		t.Errorf("Failed to add pokemon to the pokedex")
	}
}
func TestGetPokemon(t *testing.T) {
	pokedex := NewPokedex()

	pokemon := PokemonResult{
		Name: "Pikachu",
	}

	pokedex.AddPokemon(pokemon)

	retrievedPokemon, ok := pokedex.GetPokemon("Pikachu")
	if !ok {
		t.Errorf("Failed to retrieve pokemon from the pokedex")
	}

	if retrievedPokemon.Name != pokemon.Name {
		t.Errorf("Retrieved pokemon does not match the added pokemon")
	}
}
func TestReleasePokemon(t *testing.T) {
	pokedex := NewPokedex()

	pokemon := PokemonResult{
		Name: "Pikachu",
	}

	pokedex.AddPokemon(pokemon)

	pokedex.ReleasePokemon("Pikachu")

	_, ok := pokedex.GetPokemon("Pikachu")
	if ok {
		t.Errorf("Failed to release pokemon from the pokedex")
	}
}
func TestGetPokemons(t *testing.T) {
	pokedex := NewPokedex()

	pokemon1 := PokemonResult{Name: "Pikachu"}
	pokemon2 := PokemonResult{Name: "Charmander"}
	pokemon3 := PokemonResult{Name: "Bulbasaur"}

	pokedex.AddPokemon(pokemon1)
	pokedex.AddPokemon(pokemon2)
	pokedex.AddPokemon(pokemon3)

	pokemons := pokedex.GetPokemons()

	if len(pokemons) != 3 {
		t.Errorf("Expected 3 pokemons, got %d", len(pokemons))
	}

	pokemonMap := make(map[string]bool)
	for _, p := range pokemons {
		pokemonMap[p.Name] = true
	}

	expectedPokemons := []string{"Pikachu", "Charmander", "Bulbasaur"}

	for _, name := range expectedPokemons {
		if !pokemonMap[name] {
			t.Errorf("Expected pokemon %s not found", name)
		}
	}
	if len(pokemonMap) != len(expectedPokemons) {
		t.Errorf("Expected %d unique pokemons, got %d", len(expectedPokemons), len(pokemonMap))
	}
}
