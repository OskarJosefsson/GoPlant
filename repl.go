package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"
)

type config struct {
	commands    map[string]cliCommand
	store       *plantStore
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

func commandSearch(cfg *config, args []string) error {
	if len(args) == 0 {
		fmt.Println("Please provide a search term.")
		return nil
	}

	search := strings.Join(args, " ")
	fmt.Println("Searching:", search)
	searchResult, err := SearchPlants(search)

	if err != nil {
		return err
	}

	for _, s := range searchResult.Data {
		println(s.Name)
	}

	return nil
}

func commandAddPlant(cfg *config, args []string) error {
	if len(args) == 0 {
		fmt.Println("Please provide a search term.")
		return nil
	}

	add := strings.Join(args, " ")
	fmt.Println("Adding:", add)
	searchResult, err := SearchPlants(add)
	if err != nil {
		return err
	}
	if len(searchResult.Data) == 0 {
		return fmt.Errorf("no plants found for %q", add)
	}

	GetResult, err := GetFullPlant(searchResult.Data[0].Slug)

	if err != nil {
		return err
	}
	if err := cfg.store.AddPlant(context.Background(), GetResult); err != nil {
		return err
	}

	fmt.Printf("ID: %d\nCommon name: %s\nScientific name: %s\nFamily: %s\nGenus: %s\nStatus: %s\nYear: %d\nImage: %s\n",
		GetResult.ID,
		GetResult.CommonName,
		GetResult.ScientificName,
		GetResult.Family.Name,
		GetResult.Genus.Name,
		GetResult.Status,
		GetResult.Year,
		GetResult.ImageURL,
	)

	return nil
}

func commandListPlants(cfg *config, args []string) error {
	plants, err := cfg.store.ReadAllPlants(context.Background())
	if err != nil {
		return err
	}
	if len(plants) == 0 {
		fmt.Println("No plants have been added yet.")
		return nil
	}

	for _, plant := range plants {
		fmt.Printf("ID: %d\nCommon name: %s\nScientific name: %s\nFamily: %s\nGenus: %s\nStatus: %s\nYear: %d\nImage: %s\n",
			plant.ID,
			plant.CommonName,
			plant.ScientificName,
			plant.Family.Name,
			plant.Genus.Name,
			plant.Status,
			plant.Year,
			plant.ImageURL,
		)
	}

	return nil
}

func commandListWeatherStations(cfg *config, args []string) error {
	stations, err := listStations()
	if err != nil {
		return err
	}
	if len(args) > 0 {
		name := strings.Join(args, " ")
		for _, station := range stations.Station {
			if !strings.EqualFold(station.Name, name) {
				continue
			}

			weather, err := latestStationWeather(station.ID)
			if err != nil {
				return err
			}
			fmt.Printf("%s (station %d)\n", station.Name, station.ID)
			fmt.Printf("Temperature: %s °C at %s\n",
				weather.Temperature.Value,
				weather.Temperature.Time().Format(time.RFC3339))
			rainDate := weather.Rainfall.Ref
			if rainDate == "" {
				rainDate = weather.Rainfall.Time().Format("2006-01-02")
			}
			fmt.Printf("Daily rainfall: %s mm for %s\n", weather.Rainfall.Value, rainDate)
			return nil
		}
		return fmt.Errorf("no active station reporting both hourly temperature and daily rainfall named %q", name)
	}
	if len(stations.Station) == 0 {
		fmt.Println("No active stations report both hourly temperature and daily rainfall.")
		return nil
	}

	fmt.Printf("%-8s %-32s %10s %11s\n", "ID", "Station", "Latitude", "Longitude")
	for _, station := range stations.Station {
		fmt.Printf("%-8d %-32s %10.4f %11.4f\n",
			station.ID, station.Name, station.Latitude, station.Longitude)
	}
	return nil
}
