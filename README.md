# Pokedex CLI

A command-line interface Pokedex application built in Go that allows you to explore the Pokemon world, catch Pokemon, and manage your collection using the [PokeAPI](https://pokeapi.co/).

## Features

- 🗺️ **Explore Location Areas**: Navigate through different location areas in the Pokemon world
- 🔍 **Discover Pokemon**: Find Pokemon in specific locations
- ⚾ **Catch Pokemon**: Attempt to catch Pokemon with Pokeballs
- 📖 **Pokedex Management**: Keep track of all your caught Pokemon
- 🔎 **Pokemon Inspection**: View detailed stats and information about caught Pokemon
- ⚡ **Caching System**: Fast response times with built-in caching

## Installation

### Prerequisites

- Go 1.23.2 or higher

### Setup

1. Clone the repository:
```bash
git clone https://github.com/fekete965/pokedexcli.git
cd pokedexcli
```

2. Build the application:
```bash
go build -o pokedexcli
```

3. Run the application:
```bash
./pokedexcli
```

## Usage

Once you start the application, you'll see the `Pokedex >` prompt. From here, you can enter various commands to interact with the Pokemon world.

### Example Session

```bash
Pokedex > help
# View all available commands

Pokedex > map
# Display the first 20 location areas

Pokedex > explore canalave-city-area
# Explore a specific location to find Pokemon

Pokedex > catch pikachu
# Try to catch a Pokemon

Pokedex > pokedex
# View all caught Pokemon

Pokedex > inspect pikachu
# View detailed information about a caught Pokemon

Pokedex > exit
# Exit the application
```

## Available Commands

| Command | Description | Usage |
|---------|-------------|-------|
| `help` | Displays a help message with all available commands | `help` |
| `map` | Displays the names of the next 20 location areas | `map` |
| `mapb` | Displays the names of the previous 20 location areas | `mapb` |
| `explore` | Lists all the Pokemon located in a specific location area | `explore <location name>` |
| `catch` | Attempt to catch a Pokemon | `catch <pokemon name>` |
| `inspect` | View detailed information about a caught Pokemon | `inspect <pokemon name>` |
| `pokedex` | Display a list of all caught Pokemon | `pokedex` |
| `exit` | Exit the Pokedex application | `exit` |

## Project Structure

```
pokedexcli/
├── main.go              # Main application and command implementations
├── internal/
│   ├── pokecache.go     # Caching implementation
│   └── pokecache_test.go # Cache tests
├── go.mod               # Go module file
└── README.md            # This file
```

## How It Works

- **Location Navigation**: Use `map` and `mapb` commands to browse through different location areas
- **Pokemon Discovery**: Use `explore` with a location name to see which Pokemon can be found there
- **Catching Mechanics**: The catch success rate is based on the Pokemon's base experience - harder to catch stronger Pokemon!
- **Caching**: API responses are cached for 7 seconds to improve performance and reduce API calls

## Upcoming Features

The following features are planned for future releases:

- [ ] Update the CLI to support the "up" arrow to cycle through previous commands
- [ ] Simulate battles between Pokemon
- [ ] Add more unit tests
- [ ] Refactor your code to organize it better and make it more testable
- [ ] Keep Pokemon in a "party" and allow them to level up
- [ ] Allow for Pokemon that are caught to evolve after a set amount of time
- [ ] Persist a user's Pokedex to disk so they can save progress between sessions
- [ ] Use the PokeAPI to make exploration more interesting. For example, rather than typing the names of areas, maybe you are given choices of areas and just type "left" or "right"
- [ ] Random encounters with wild Pokemon
- [ ] Adding support for different types of balls (Pokeballs, Great Balls, Ultra Balls, etc), which have different chances of catching Pokemon

## Contributing

Contributions are welcome! Feel free to submit issues or pull requests.

## License

This project is for educational purposes.

## Acknowledgments

- [PokeAPI](https://pokeapi.co/) for providing the Pokemon data

