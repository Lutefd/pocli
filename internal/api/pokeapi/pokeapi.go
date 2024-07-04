package pokeapi

import (
	"poke-repl/internal/cache"
	"time"
)

var pokeCache = cache.NewCache(time.Minute * 5)
