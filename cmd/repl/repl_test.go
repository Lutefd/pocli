package repl

import (
	"io"
	"os"
	"poke-repl/internal/api/pokeapi"
	"poke-repl/internal/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommandsMap(t *testing.T) {
	commands := CommandsMap()
	if len(commands) != 10 {
		t.Errorf("Expected 10 commands, got %d", len(commands))
	}
}

func TestLookupCommand(t *testing.T) {
	command, err := LookupCommand("help")
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
	if command.name != "help" {
		t.Errorf("Expected command name to be 'help', got %s", command.name)
	}
}

func TestLookupCommandNotFound(t *testing.T) {
	_, err := LookupCommand("notfound")
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}
func TestCommandHelp(t *testing.T) {
	cfg := &config.Config{}
	args := []string{}
	err := commandHelp(cfg, args)
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
}
func TestClearScreen(t *testing.T) {
	cfg := &config.Config{}
	args := []string{}

	err := clearScreen(cfg, args)
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
}

// Exit command is not tested because it exits the program,
// we could test it by mocking the os.Exit function
// but it's not worth it for this application.

func TestMapCommand(t *testing.T) {
	cfg := &config.Config{}
	args := []string{}
	err := mapCommand(cfg, args)
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
}

func TestPageNavigation(t *testing.T) {
	tests := []struct {
		name      string
		setupCfg  func() *config.Config
		testFunc  func(cfg *config.Config, args []string) error
		expectErr bool
	}{
		{
			name: "PreviousPage with valid URL",
			setupCfg: func() *config.Config {
				return &config.Config{
					PreviousUrl: "https://pokeapi.co/api/v2/location/?offset=40&limit=20",
					Referrer:    "previous",
					Cmd:         "map",
				}
			},
			testFunc:  previousPage,
			expectErr: false,
		},
		{
			name: "PreviousPage with no URL",
			setupCfg: func() *config.Config {
				return &config.Config{
					PreviousUrl: "",
					Referrer:    "previous",
					Cmd:         "map",
				}
			},
			testFunc:  previousPage,
			expectErr: true,
		},
		{
			name: "NextPage with valid URL",
			setupCfg: func() *config.Config {
				return &config.Config{
					NextUrl:  "https://pokeapi.co/api/v2/location/?offset=20&limit=20",
					Referrer: "next",
					Cmd:      "map",
				}
			},
			testFunc:  nextPage,
			expectErr: false,
		},
		{
			name: "NextPage with no URL",
			setupCfg: func() *config.Config {
				return &config.Config{
					NextUrl:  "",
					Referrer: "next",
					Cmd:      "map",
				}
			},
			testFunc:  nextPage,
			expectErr: true,
		},
	}
	args := []string{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := tt.setupCfg()
			err := tt.testFunc(cfg, args)
			if tt.expectErr && err == nil {
				t.Errorf("Expected an error, got nil")
			} else if !tt.expectErr && err != nil {
				t.Errorf("Did not expect an error, got %v", err)
			}
		})
	}
}
func TestInspectCommand(t *testing.T) {
	cfg := &config.Config{}
	args := []string{"Pikachu"}
	pokeDex.AddPokemon(pokeapi.PokemonResult{
		Name:           "Pikachu",
		BaseExperience: 50,
	})
	err := inspectCommand(cfg, args)
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
}

func TestInspectCommandNoPokemon(t *testing.T) {
	cfg := &config.Config{}
	args := []string{"Charmander"}

	err := inspectCommand(cfg, args)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestInspectCommandNoArgs(t *testing.T) {
	cfg := &config.Config{}
	args := []string{}

	err := inspectCommand(cfg, args)
	if err == nil {
		t.Errorf("Expected no error, got %s", err)
	}
}
func TestPokedexCommand(t *testing.T) {
	cfg := &config.Config{}
	args := []string{}

	err := pokedexCommand(cfg, args)
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
}

func TestPokedexCommandWithArgs(t *testing.T) {
	cfg := &config.Config{}
	args := []string{"Pikachu"}

	err := pokedexCommand(cfg, args)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}
func TestExploreCommand(t *testing.T) {
	cfg := &config.Config{}

	tests := []struct {
		name           string
		args           []string
		expectedError  string
		expectedOutput string
	}{
		{
			name:          "no area specified",
			args:          []string{},
			expectedError: "no area specified",
		},
		{
			name:          "more than one area specified",
			args:          []string{"area1", "area2"},
			expectedError: "only one area can be explored at a time",
		},
		{
			name:           "successful exploration",
			args:           []string{"canalave-city-area"},
			expectedOutput: "- finneon\n- gastrodon\n- gyarados\n- lumineon\n- magikarp\n- pelipper\n- shellos\n- staryu\n- tentacool\n- tentacruel\n- wingull\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			err := exploreCommand(cfg, tt.args)

			w.Close()
			out, _ := io.ReadAll(r)
			os.Stdout = oldStdout

			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedOutput, string(out))
			}
		})
	}
}

func TestCatchCommand(t *testing.T) {
	cfg := &config.Config{}

	tests := []struct {
		name           string
		args           []string
		expectedError  string
		expectedOutput string
	}{
		{
			name:          "more than one pokemon specified",
			args:          []string{"pikachu", "bulbasaur"},
			expectedError: "only one pokemon can be caught at a time",
		},
		{
			name:           "pokemon caught successfully",
			args:           []string{"caterpie"},
			expectedOutput: "Throwing a Pokeball at caterpie...\ncaterpie was caught!\n", // Used this to garantee that the pokemon was caught since the test is random and caterpie is one of the most common pokemon
		},
		{
			name:           "pokemon escaped",
			args:           []string{"mew"},
			expectedOutput: "Throwing a Pokeball at mew...\nmew escaped!\n", // Used this to garantee that the pokemon was caught since the test is random and mew is one of the most rare pokemon
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture the output
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			err := catchCommand(cfg, tt.args)

			w.Close()
			out, _ := io.ReadAll(r)
			os.Stdout = oldStdout

			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedOutput, string(out))
			}
		})
	}
}
