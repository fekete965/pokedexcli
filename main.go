package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/fekete965/pokedexcli/internal"
)

const LOCATION_API_URL string = "https://pokeapi.co/api/v2/location-area"
const POKEMON_API_URL string = "https://pokeapi.co/api/v2/pokemon"

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

type PokemonType struct {
	Slot int `json:"slot"`
	Type NamedResource `json:"type"`
}

type PokemonStat struct {
	BaseStat int `json:"base_stat"`
	Effort int `json:"effort"`
	Stat NamedResource `json:"Stat"`
}

type PokemonInfoResponse struct {
	Id int `json:"id"`
	Name string `json:"name"`
	BaseExperience int `json:"base_experience"`
	Height int `json:"height"`
	IsDefault bool `json:"is_default"`
	Order int `json:"order"`
	Weight int `json:"weight"`
	Types []PokemonType `json:"types"`
	Stats []PokemonStat `json:"stats"`
}

type cliCommand struct {
	name string
	description string
	callback func(cfg *config, args []string) error
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
var pokedexCache map[string]PokemonInfoResponse = make(map[string]PokemonInfoResponse)

func createCommandRegistry() map[string]cliCommand {
	registry := map[string]cliCommand {
		"catch": {
			name: "catch",
			description: "Catch a Pokemon. Usage: catch <pokemon name>",
			callback: commandCatch,
		},
		"exit": {
			name: "exit",
			description: "Exit the Pokedex. Usage: exit",
			callback: commandExit,
		},
		"explore": {
			name: "explore",
			description: "Lists all the Pokémon located in a specific location area. Usage: explore <location name>",
			callback: commandExplore,
		},
		"inspect": {
			name: "inspect",
			description: "Inspect a caught Pokemon. Usage: inspect <pokemon name>",
			callback: commandInspect,
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
		"pokedex": {
			name: "pokedex",
			description: "It displays the list of all caught Pokemon. Usage: pokedex",
			callback: commandPokedex,
		},
	}

	registry["help"] = cliCommand {
		name: "help",
		description: "Displays a help message. Usage: help",
		callback: createCommandHelp(registry),
	}

	return registry
}

func commandExit(cfg *config, args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)

	return nil
}

func createCommandHelp(cliCommandRegistry map[string]cliCommand) func(cfg *config, args []string) error {
	return func(cfg *config, args []string) error {
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

func commandMap(cfg *config, args []string) error {
	apiURL := cfg.Next

	if apiURL == nil {
		fmt.Println("you're on the last page")
		return nil
	}

	result, err := getPokemonLocations(*apiURL)
	if err != nil {
		return err
	}

	mapConfig.Next = result.Next
	mapConfig.Previous = result.Previous
 
	printPokemonLocations(result)

	return nil
}

func commandMapB(cfg *config, args []string) error {
	apiURL := cfg.Previous

	if apiURL == nil {
		fmt.Println("you're on the first page")
		return nil
	}

	result, err := getPokemonLocations(*apiURL)
	if err != nil {
		return err
	}

	mapConfig.Next = result.Next
	mapConfig.Previous = result.Previous
 
	printPokemonLocations(result)

	return nil
}

func commandExplore(cfg *config, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("missing location name. Usage: explore <location name>")
	}

	locationAreaName := args[0]
	
	prompt := fmt.Sprintf("Exploring '%v'...", locationAreaName)
	fmt.Println(prompt)

	data, err := getPokemonLocationAreaDetails(locationAreaName)
	if err != nil {
		return err
	}

	printPokemonLocationAreaDetails(data, locationAreaName)

	return nil
}

func commandCatch(cfg *config, args []string) error {
	if len(args) == 0 || len(strings.TrimSpace(args[0])) == 0 {
		return fmt.Errorf("missing pokemon name. Usage: catch <pokemon name>")
	}

	pokemonName := args[0]
	
	pokemonInfo, err := getPokemonInfo(pokemonName)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("Throwing a Pokeball at %v...", pokemonInfo.Name)
	fmt.Println(msg)

	isSuccessfulCapture := rand.Intn(pokemonInfo.BaseExperience) > pokemonInfo.BaseExperience / 2

	if isSuccessfulCapture {
		pokedexCache[pokemonInfo.Name] = pokemonInfo

		msg := fmt.Sprintf("%v was caught!", pokemonInfo.Name)
		fmt.Println(msg)
		fmt.Println("You may now inspect it with the inspect command.")
	} else {
		msg := fmt.Sprintf("%v escaped!", pokemonInfo.Name)
		fmt.Println(msg)
	}

	return nil
}

func commandInspect(cfg *config, args []string) error {
	if len(args) == 0 || len(strings.TrimSpace(args[0])) == 0 {
		return fmt.Errorf("missing pokemon name. Usage: inspect <pokemon name>")
	}

	pokemonName := args[0]
	pokemonInfo, ok := pokedexCache[pokemonName]
	if !ok {
		return fmt.Errorf("you have not caught that pokemon")
	}
	
	printPokemonInfo(pokemonInfo)

	return nil
}

func commandPokedex(cfg *config, args []string) error {
	fmt.Println("Your Pokedex:")
	
	if len(pokedexCache) == 0 {
		fmt.Println(" - No Pokemon caught yet")
		return nil
	}

	for _, pokemonInfo := range pokedexCache {
		println(" - " + pokemonInfo.Name)
	}	

	return nil
}

func printPokemonInfo(pokemonInfo PokemonInfoResponse) {
	fmt.Println("Name: " + pokemonInfo.Name)

	heightStr := fmt.Sprintf("Height: %v", pokemonInfo.Height)
	fmt.Println(heightStr)

	weightStr := fmt.Sprintf("Weight: %v", pokemonInfo.Weight)
	fmt.Println(weightStr)
	
	fmt.Println("Stats:")
	if len(pokemonInfo.Stats) == 0 {
		fmt.Println("  - No stats found")
	} else {
		for _, stat := range pokemonInfo.Stats {
			statStr := fmt.Sprintf("  - %v: %v", stat.Stat.Name, stat.BaseStat)
			fmt.Println(statStr)
		}
	}

	fmt.Println("Types:")
	if len(pokemonInfo.Types) == 0 {
		fmt.Println("  - No types found")
	} else {
	for _, tpe := range pokemonInfo.Types {
			tpeStr := fmt.Sprintf("  - %v", tpe.Type.Name)
			fmt.Println(tpeStr)
		}
	}
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

func getPokemonInfo(pokemonName string) (PokemonInfoResponse, error) {
	apiURL := POKEMON_API_URL + "/" + pokemonName

	result := PokemonInfoResponse{}

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
			err := command.callback(&mapConfig, input[1:])
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
