package repl

import (
	"fmt"
	"os"
	"os/exec"
)

type cliCommand struct {
	name        string
	description string
	Callback    func() error
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
	}
}
func commandHelp() error {
	commands := CommandsMap()
	fmt.Println("Available commands:")
	for name, command := range commands {
		fmt.Printf("  %s - %s\n", name, command.description)
	}
	return nil
}

func commandExit() error {
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
func clearScreen() error {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
	return nil
}
