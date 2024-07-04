package repl

import (
	"fmt"
	"os"
	"os/exec"
	"poke-repl/internal/api/pokeapi"
	"poke-repl/internal/config"
)

type cliCommand struct {
	name        string
	description string
	url         string
	Callback    func(cfg *config.Config, args []string) error
}

func CommandsMap() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			Callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			Callback:    commandExit,
		},
		"clear": {
			name:        "clear",
			description: "Clear the screen",
			Callback:    clearScreen,
		},
		"next": {
			name:        "next",
			description: "Go to the next page when available",
			Callback:    nextPage,
		},
		"previous": {
			name:        "previous",
			description: "Go to the previous page when available",
			Callback:    previousPage,
		},
		"map": {
			name:        "map",
			description: "Show locations in the pokemon world",
			url:         "https://pokeapi.co/api/v2/location-area?offset=0&limit=20",
			Callback:    mapCommand,
		},
		"explore": {
			name:        "explore",
			description: "Exlore the pokemon world area by area",
			url:         "https://pokeapi.co/api/v2/location-area/",
			Callback:    exploreCommand,
		},
	}
}
func commandHelp(cfg *config.Config, args []string) error {
	commands := CommandsMap()
	fmt.Println("Available commands:")
	for name, command := range commands {
		fmt.Printf("  %s - %s\n", name, command.description)
	}
	cfg.Cmd = "help"
	return nil
}

func commandExit(cfg *config.Config, args []string) error {
	fmt.Println("Bye!")
	os.Exit(1)
	return nil
}
func LookupCommand(name string) (cliCommand, error) {
	command, ok := CommandsMap()[name]
	if !ok {
		return cliCommand{}, fmt.Errorf("command not found")
	}
	return command, nil
}
func clearScreen(cfg *config.Config, args []string) error {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
	cfg.Cmd = "clear"
	return nil
}

func nextPage(cfg *config.Config, args []string) error {
	if cfg.NextUrl == "" {
		return fmt.Errorf("no next page available")
	}
	switch cfg.Cmd {
	case "map":
		cfg.Referrer = "next"
		return mapCommand(cfg, args)

	}
	return nil
}
func previousPage(cfg *config.Config, args []string) error {
	if cfg.PreviousUrl == "" {
		return fmt.Errorf("no previous page available")
	}
	switch cfg.Cmd {
	case "map":
		cfg.Referrer = "previous"
		return mapCommand(cfg, args)

	}
	return nil
}

func mapCommand(cfg *config.Config, args []string) error {
	cfg.Cmd = "map"
	cmd, err := LookupCommand("map")
	if err != nil {
		return err
	}
	defaultUrl := cmd.url
	if cfg.NextUrl != "" && cfg.Referrer == "next" {
		defaultUrl = cfg.NextUrl
	}
	if cfg.PreviousUrl != "" && cfg.Referrer == "previous" {
		defaultUrl = cfg.PreviousUrl
	}
	locations, err := pokeapi.Location.GetLocation(defaultUrl, cfg)
	for _, location := range locations {
		fmt.Printf("- %s\n", location.Name)
	}
	if err != nil {
		return err
	}
	return nil
}

func exploreCommand(cfg *config.Config, args []string) error {
	if len(args) > 1 {
		return fmt.Errorf("only one area can be explored at a time")
	}
	cfg.Cmd = "explore"
	cmd, err := LookupCommand("explore")
	if err != nil {
		return err
	}
	pokemonList, err := pokeapi.Explorer.Explore(cmd.url, args[0], cfg)
	if err != nil {
		return err
	}
	for _, pokemon := range pokemonList {
		fmt.Printf("- %s\n", pokemon)
	}
	return nil
}
