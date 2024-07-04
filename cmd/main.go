package main

import (
	"bufio"
	"fmt"
	"os"
	"poke-repl/cmd/repl"
	"poke-repl/internal/config"
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
			cmd, err := repl.LookupCommand(command)
			if err != nil {
				fmt.Println(err)
				continue
			}
			err = cmd.Callback(cfg)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
