package pokeapi

import "sync"

type pokedex map[string]PokemonResult

type PokemonCaught struct {
	Dex pokedex
	mu  sync.RWMutex
}

func NewPokedex() *PokemonCaught {
	return &PokemonCaught{
		Dex: make(pokedex),
	}
}

func (p *PokemonCaught) AddPokemon(pokemon PokemonResult) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Dex[pokemon.Name] = pokemon
}

func (p *PokemonCaught) GetPokemon(name string) (PokemonResult, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	pokemon, ok := p.Dex[name]
	return pokemon, ok
}

func (p *PokemonCaught) GetPokemons() []PokemonResult {
	p.mu.RLock()
	defer p.mu.RUnlock()
	pokemons := make([]PokemonResult, 0, len(p.Dex))
	for _, pokemon := range p.Dex {
		pokemons = append(pokemons, pokemon)
	}
	return pokemons
}

func (p *PokemonCaught) ReleasePokemon(name string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.Dex, name)
}
