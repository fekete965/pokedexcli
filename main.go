package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Pokedex > ")
	
	for scanner.Scan() {
		cleanedInput := cleanInput(scanner.Text())
		fmt.Println("Your command was: " + cleanedInput[0])
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
