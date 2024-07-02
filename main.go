package main

import (
	"bufio"
	"fmt"
	"os"
	"poke-repl/repl"
)

const (
	cliName = "Pokedex"
)

func main() {
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
			err = cmd.Callback()
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
