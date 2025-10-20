package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

type cliCommand struct {
	name string
	description string
	callback func() error
}

var cliCommandRegistry map[string]cliCommand = createCommandRegistry()

func createCommandRegistry() map[string]cliCommand {
	cliCommandRegistry := map[string]cliCommand {
		"exit": {
			name: "exit",
			description: "Exit the Pokedex",
			callback: commandExit,
		},
	}

	cliCommandRegistry["help"] = cliCommand {
		name: "help",
		description: "Displays a help message",
		callback: createCommandHelp(cliCommandRegistry),
	}

	return cliCommandRegistry
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)

	return nil
}

func createCommandHelp(cliCommandRegistry map[string]cliCommand) func() error {
	return func() error {
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
			err := command.callback()
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
