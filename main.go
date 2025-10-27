package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/fekete965/pokedexcli/internal"
)

const LOCATION_API_URL string = "https://pokeapi.co/api/v2/location-area/"

type NamedResource struct {
	Name string `json:"name"`
	Url string `json:"url"`
}

type PokemonLocationAreaResponse struct {
	Count int `json:"count"`
	Next *string `json:"next"`
	Previous *string `json:"previous"`
	Results []NamedResource `json:"results"`
}

type PokemonEncounter struct {
	Pokemon NamedResource `json:"pokemon"`
}

type PokemonLocationAreaDetailsResponse struct {
	GameIndex int `json:"game_index"`
	Location NamedResource `json:"location"`
	Name string `json:"name"`
	PokemonEncounters []PokemonEncounter `json:"pokemon_encounters"`
}

type cliCommand struct {
	name string
	description string
	callback func(cfg *config) error
}

type config struct {
	Next *string
	Previous *string
}

func constToPtr(s string) *string { return &s }

var mapConfig config = config{
	Next: constToPtr(LOCATION_API_URL),
	Previous: nil,
}
var cliCommandRegistry map[string]cliCommand = createCommandRegistry()
var pokemonLocationAreaCache = internal.NewCache(time.Second * 7)
var pokemonLocationAreaDetailCache = internal.NewCache(time.Second * 7)

func createCommandRegistry() map[string]cliCommand {
	registry := map[string]cliCommand {
		"exit": {
			name: "exit",
			description: "Exit the Pokedex. Usage: exit",
			callback: commandExit,
		},
		"map": {
			name: "map",
			description: "It displays the names of the next 20 location areas in the Pokemon world. Usage: map",
			callback: commandMap,
		},
		"mapb": {
			name: "mapb",
			description: "It displays the names of previous 20 location areas in the Pokemon world. Usage: mapb",
			callback: commandMapB,
		},
	}

	registry["help"] = cliCommand {
		name: "help",
		description: "Displays a help message. Usage: help",
		callback: createCommandHelp(registry),
	}

	return registry
}

func commandExit(cfg *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)

	return nil
}

func createCommandHelp(cliCommandRegistry map[string]cliCommand) func(cfg *config) error {
	return func(cfg *config) error {
		fmt.Println("")

		fmt.Println("\nWelcome to the Pokedex!")
		fmt.Println("Usage:")
		fmt.Println("")
		
		sortedKeys := make([]string, 0, len(cliCommandRegistry))
		for key := range cliCommandRegistry {
			sortedKeys = append(sortedKeys, key)
		}
		sort.Strings(sortedKeys)

		for _, sortedKey := range sortedKeys {
			command := cliCommandRegistry[sortedKey]
			fmt.Printf("%s: %s\n", command.name, command.description)
		}

		fmt.Println("")

		return nil
	}
}

func commandMap(cfg *config) error {
	apiURL := cfg.Next

	if apiURL == nil {
		fmt.Println("you're on the last page")
		return nil
	}

	result, err := getPokemonLocation(*apiURL)
	if err != nil {
		return err
	}

	mapConfig.Next = result.Next
	mapConfig.Previous = result.Previous
 
	printPokemonLocations(result)

	return nil
}

func commandMapB(cfg *config) error {
	apiURL := cfg.Previous

	if apiURL == nil {
		fmt.Println("you're on the first page")
		return nil
	}

	result, err := getPokemonLocation(*apiURL)
	if err != nil {
		return err
	}

	mapConfig.Next = result.Next
	mapConfig.Previous = result.Previous
 
	printPokemonLocations(result)

	return nil
}


func printPokemonLocations(data PokemonLocationAreaResponse) {
	for _, location := range data.Results {
		fmt.Println(location.Name)
	}
}

func printPokemonLocationAreaDetails(data PokemonLocationAreaDetailsResponse, locationAreaName string) {
	if len(data.PokemonEncounters) == 0 {
		fmt.Println("No Pokemon found in " + locationAreaName)
		
		return 
	}
	
	fmt.Println("Found Pokemon:")
	for _, encounter := range data.PokemonEncounters {
		fmt.Println(" - " +encounter.Pokemon.Name)
	}
}

func getPokemonLocations(apiURL string) (PokemonLocationAreaResponse, error) {
	result := PokemonLocationAreaResponse{}

	if cacheData, hasCache := pokemonLocationAreaCache.Get(apiURL); hasCache {
		err := json.Unmarshal(cacheData, &result)
		if err != nil {
			return result, err
		}
		
		return result, nil
	}


	res, err := http.Get(apiURL)
	if err != nil {
		return result, err
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return result, err
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		err := fmt.Errorf("response failed with status code: %d and\nbody: %s", res.StatusCode, data)
		return result, err
	}

	
	err = json.Unmarshal(data, &result)
	if err != nil {
		return result, err
	}

	pokemonLocationAreaCache.Add(apiURL, data)

	return result, nil
}

func getPokemonLocationAreaDetails(locationAreaName string) (PokemonLocationAreaDetailsResponse, error) {
	apiURL := LOCATION_API_URL + "/" + locationAreaName

	result := PokemonLocationAreaDetailsResponse{}

	if cacheData, hasCache := pokemonLocationAreaDetailCache.Get(apiURL); hasCache {
		err := json.Unmarshal(cacheData, &result)
		if err != nil {
			return result, err
		}

		return result, nil
	}

	res, err := http.Get(apiURL)
	if err != nil {
		return result, err
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return result, err
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		err := fmt.Errorf("response failed with status code: %d and\nbody: %s", res.StatusCode, data)
		return result, err
	}

	err = json.Unmarshal(data, &result)
	if err != nil {
		return result, err
	}

	pokemonLocationAreaDetailCache.Add(apiURL, data)

	return result, nil
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Pokedex > ")

	for scanner.Scan() {
		input := cleanInput(scanner.Text())
		
		isInputEmpty := len(input) == 0
		if isInputEmpty {
			continue
		}

		commandWord := input[0]
		command, isCommandExists := cliCommandRegistry[commandWord]

		if !isCommandExists {
			fmt.Println("Unknown command")
		} else {
			err := command.callback(&mapConfig)
			if err != nil {
				fmt.Println(err)
			}
		}

		fmt.Print("Pokedex > ")
	}
}

func cleanInput(text string) []string {
	words := strings.Fields(text)

	for i, word := range words {
		words[i] = strings.ToLower(word)
	}

	return words
}
