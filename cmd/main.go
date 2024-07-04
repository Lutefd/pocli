package main

import (
	"bufio"
	"fmt"
	"os"
	"poke-repl/cmd/repl"
	"poke-repl/internal/config"
	"strings"
)

const (
	cliName = "Pokedex"
)

func main() {
	cfg := &config.Config{}

	for {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Printf("%s > ", cliName)
		if scanner.Scan() {
			command := scanner.Text()
			commandArgs := strings.Split(command, " ")
			cmd, err := repl.LookupCommand(commandArgs[0])
			if err != nil {
				fmt.Println(err)
				continue
			}
			err = cmd.Callback(cfg, commandArgs[1:])
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
