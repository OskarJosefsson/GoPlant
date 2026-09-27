package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type config struct {
	commands    map[string]cliCommand
	NextUrl     string
	PreviousUrl string
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config, []string) error
}

func startRepl(cfg *config) {

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Welcome to PlantWorld!")

	for {
		fmt.Print("Plant > ")
		if !scanner.Scan() {
			break
		}
		input := scanner.Text()
		args := cleanInput(input)
		if len(args) == 0 {
			continue
		}
		cmdName := args[0]
		cmd, exists := cfg.commands[cmdName]
		if !exists {
			fmt.Printf("Unknown command: %s\n", cmdName)
			continue
		}
		err := cmd.callback(cfg, args[1:])
		if err != nil {
			fmt.Printf("Error executing command: %v\n", err)
		}
	}

}

func cleanInput(input string) []string {

	input = strings.ToLower(input)
	input = strings.TrimSpace(input)
	if input == "" {
		return []string{""}
	}
	return strings.Fields(input)
}

func registerCommand(cfg *config, name string, description string, callback func(*config, []string) error) {
	cfg.commands[name] = cliCommand{
		name:        name,
		description: description,
		callback:    callback,
	}
}

func commandHelp(cfg *config, args []string) error {
	fmt.Println("Available commands:")
	for _, cmd := range cfg.commands {
		fmt.Printf("  %s: %s\n", cmd.name, cmd.description)
	}
	return nil
}
