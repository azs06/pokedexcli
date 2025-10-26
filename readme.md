# pokedexcli

A command line pokedex application created using Go and PokeAPI. As part of the boot.dev GoLang course.

## Running the program

To run the program, use the following command:

```bash
go run .
```

Alternatively, build the program first:

```bash
go build -o pokedexcli
```

Then run it with:

```bash
./pokedexcli
```

## Available Commands

- exit: Exit the application.
- help: Display available commands.
- map : Show available areas to explore.
- mapb: Show previous areas explored.
- explore [area]: Explore a specified area to find Pokémon.
- catch [pokemon]: Attempt to catch a specified Pokémon.
- pokedex: Display all caught Pokémon.

## Improvement Options

- [ ] Update the CLI to support the "up" arrow to cycle through previous commands
- [ ] Simulate battles between pokemon
- [ ] Add more unit tests
- [ ] Refactor your code to organize it better and make it more testable
- [ ] Keep pokemon in a "party" and allow them to level up
- [ ] Allow for pokemon that are caught to evolve after a set amount of time
- [ ] Persist a user's Pokedex to disk so they can save progress between sessions
- [ ] Use the PokeAPI to make exploration more interesting. For example, rather than typing the names of areas, maybe you are given choices of areas and just type "left" or "right"
- [ ] Random encounters with wild pokemon
- [ ] Adding support for different types of balls (Pokeballs, Great Balls, Ultra Balls, etc), which have different chances of catching pokemon
å