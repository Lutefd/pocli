package repl

import (
	"poke-repl/internal/config"
	"testing"
)

func TestCommandsMap(t *testing.T) {
	commands := CommandsMap()
	if len(commands) != 6 {
		t.Errorf("Expected 6 commands, got %d", len(commands))
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
	err := commandHelp(cfg)
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
}
func TestClearScreen(t *testing.T) {
	cfg := &config.Config{}
	err := clearScreen(cfg)
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
}

// Exit command is not tested because it exits the program,
// we could test it by mocking the os.Exit function
// but it's not worth it for this application.

func TestMapCommand(t *testing.T) {
	cfg := &config.Config{}
	err := mapCommand(cfg)
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
}

func TestPageNavigation(t *testing.T) {
	tests := []struct {
		name      string
		setupCfg  func() *config.Config
		testFunc  func(cfg *config.Config) error
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := tt.setupCfg()
			err := tt.testFunc(cfg)
			if tt.expectErr && err == nil {
				t.Errorf("Expected an error, got nil")
			} else if !tt.expectErr && err != nil {
				t.Errorf("Did not expect an error, got %v", err)
			}
		})
	}
}
