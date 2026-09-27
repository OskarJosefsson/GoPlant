package main

func main() {
	cfg := &config{
		commands: map[string]cliCommand{
			"help": {
				name:        "help",
				description: "Display a help message",
				callback:    commandHelp,
			},
		},
	}
	startRepl(cfg)
}
