package main

import (
	"fmt"
	"os"
)

func main() {

	apiKey := os.Getenv("PLANT_API_KEY")
	if apiKey == "" {
		fmt.Errorf("PLANT_API_KEY is not set")
	}

	cfg := &config{
		commands: map[string]cliCommand{
			"help": {
				name:        "help",
				description: "Display a help message",
				callback:    commandHelp,
			},
			"search": {
				name:        "search",
				description: "Search for plants",
				callback:    commandSearch,
			},
			"add": {
				name:        "add",
				description: "adding plant to your list of plants",
				callback:    commandAddPlant,
			},
		},
	}
	startRepl(cfg)
}
