package repl

import (
	"testing"
)

func TestCommandsMap(t *testing.T) {
	commands := CommandsMap()
	if len(commands) != 3 {
		t.Errorf("Expected 3 commands, got %d", len(commands))
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

func TestClearScreen(t *testing.T) {
	err := clearScreen()
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
}

func TestCommandHelp(t *testing.T) {
	err := commandHelp()
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
}

// Exit command is not tested because it exits the program,
// we could test it by mocking the os.Exit function
// but it's not worth it for this application.
