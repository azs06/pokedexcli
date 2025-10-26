package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/azs06/pokedexcli/internal/pokecache"
)

type config struct {
	Url      string
	Next     string
	Previous string
	Cache    *pokecache.Cache
}

type cliCommand struct {
	name        string
	description string
	callback    func(c *config, args ...string) error
}

type Location struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type LocationResponse struct {
	Count     int        `json:"count"`
	Next      string     `json:"next"`
	Previous  string     `json:"previous"`
	Locations []Location `json:"results"`
}

type Pokemon struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type PokemonEncounter struct {
	Pokemon Pokemon `json:"pokemon"`
}

type LocationDetailsResponse struct {
	PokemonEncounters []PokemonEncounter `json:"pokemon_encounters"`
}

type Stat struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}
type StatDetail struct {
	BaseStat int  `json:"base_stat"`
	Stat     Stat `json:"stat"`
}

type Type struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type TypeDetails struct {
	Slot int  `json:"slot"`
	Type Type `json:"type"`
}
type PokemonType struct {
	Name           string        `json:"name"`
	Height         int           `json:"height"`
	Weight         int           `json:"weight"`
	Stats          []StatDetail  `json:"stats"`
	Types          []TypeDetails `json:"types"`
	BaseExperience int           `json:"base_experience"`
}

var apiUrl = "https://pokeapi.co/api/v2/"
var pokeDex = map[string]PokemonType{}

var commands = map[string]cliCommand{
	"exit": {
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    commandExit,
	},
	"help": {
		name:        "help",
		description: "Display available commands",
		callback:    commandHelp,
	},
	"map": {
		name:        "map",
		description: "Display next maps",
		callback:    commandMap,
	},
	"mapb": {
		name:        "map",
		description: "Display previous maps",
		callback:    commandPrevMap,
	},
	"explore": {
		name:        "explore",
		description: "Explore a location",
		callback:    commandExplore,
	},
	"catch": {
		name:        "catch",
		description: "Catch a pokemon",
		callback:    commandCatch,
	},
	"inspect": {
		name:        "inspect",
		description: "Inspect a caught pokemon",
		callback:    commandInspect,
	},
	"pokedex": {
		name:        "pokedex",
		description: "View your pokedex",
		callback:    commandPokedex,
	},
}

func commandPokedex(c *config, args ...string) error {

	fmt.Println("Your Pokedex:")

	for k := range pokeDex {
		fmt.Print(" - ")
		fmt.Println(k)
	}

	return nil
}

func commandCatch(c *config, args ...string) error {
	toCatch := args[0]
	catchPokemon(toCatch, c)
	return nil
}

func cleanInput(text string) []string {
	text = strings.TrimSpace(text) // remove leading/trailing whitespace
	text = strings.ToLower(text)   // normalize case
	words := strings.Fields(text)  // split by any whitespace, ignoring multiples
	return words
}

func catchPokemon(p string, c *config) error {
	printMsg := fmt.Sprintf("Throwing a Pokeball at %s...", p)
	fmt.Println(printMsg)
	response := PokemonType{}
	url := c.Url + "pokemon/" + p

	decodedData, err := fetchData(url, c)
	if err != nil {
		fmt.Println("failed to catch", err)
		return err
	}
	err = json.Unmarshal(decodedData, &response)

	if err != nil {
		fmt.Println(err)
		return err
	}

	baseExperience := response.BaseExperience
	chance := rand.IntN(baseExperience)
	willGotCaught := baseExperience - chance

	if willGotCaught > baseExperience/2 {
		fmt.Println(p + " was caught")
		pokeDex[p] = response
	} else {
		fmt.Println(p + " escaped")
	}
	return nil
}

func fetchData(url string, c *config) ([]byte, error) {
	if strings.TrimSpace(url) == "" {
		return []byte{}, errors.New("Invalid input")
	}

	decodedData, ok := c.Cache.Get(url)
	if ok {
		return decodedData, nil
	}

	res, err := http.Get(url)
	if err != nil {
		return []byte{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("failed to fetch data: %s", res.Status)
	}

	decodedData, err = io.ReadAll(res.Body)
	c.Cache.Add(url, decodedData)

	if err != nil {
		return []byte{}, err
	}
	return decodedData, nil
}

func commandExit(c *config, args ...string) error {
	fmt.Print("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(c *config, args ...string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("help: Displays a help message")
	fmt.Println("exit: Exit the Pokedex")
	return nil
}

func fetchLocationDetails(url string, c *config) (LocationDetailsResponse, error) {
	response := LocationDetailsResponse{}
	decodedData, err := fetchData(url, c)

	if err != nil {
		return response, err
	}

	err = json.Unmarshal(decodedData, &response)

	if err != nil {
		return response, err
	}

	return response, nil
}

func commandExplore(c *config, args ...string) error {
	area := args[0]
	response, err := fetchLocationDetails(c.Url+"location-area/"+area, c)
	pokemonEncounters := response.PokemonEncounters
	if err != nil {
		return err
	}
	if len(pokemonEncounters) > 0 {
		for _, pokemonEncounter := range pokemonEncounters {
			fmt.Println(pokemonEncounter.Pokemon.Name)
		}
	}
	return nil
}

func commandMap(c *config, args ...string) error {
	locations := []Location{}
	response := LocationResponse{}
	mapUrl := fmt.Sprintf("%s/location-area", c.Url)
	if c.Next != "" {
		mapUrl = c.Next
	}
	response, err := fetchLocations(mapUrl, c)

	if err != nil {
		return err
	}

	locations = response.Locations
	c.Next = response.Next
	c.Previous = response.Previous

	for _, location := range locations {
		fmt.Println(location.Name)
	}

	return nil
}

func fetchLocations(url string, c *config) (LocationResponse, error) {
	response := LocationResponse{}
	decodedData, err := fetchData(url, c)

	if err != nil {
		return response, err
	}

	err = json.Unmarshal(decodedData, &response)

	if err != nil {
		return response, err
	}

	return response, nil
}

func commandPrevMap(c *config, args ...string) error {
	locations := []Location{}
	response := LocationResponse{}
	mapUrl := ""
	if c.Previous == "" {
		fmt.Println("you're on the first page")
		return nil
	} else {
		mapUrl = c.Previous
	}
	response, err := fetchLocations(mapUrl, c)

	if err != nil {
		return err
	}

	locations = response.Locations
	c.Next = response.Next
	c.Previous = response.Previous

	for _, location := range locations {
		fmt.Println(location.Name)
	}

	return nil
}

func commandInspect(c *config, args ...string) error {
	pokemonName := args[0]
	pokemon, exists := pokeDex[pokemonName]
	if !exists {
		fmt.Println("You haven't caught", pokemonName)
		return nil
	}

	fmt.Printf("Details of %s:\n", pokemonName)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Printf("Base Experience: %d\n", pokemon.BaseExperience)

	fmt.Println("Types:")
	for _, t := range pokemon.Types {
		fmt.Printf("- %s (Slot %d)\n", t.Type.Name, t.Slot)
	}

	fmt.Println("Stats:")
	for _, s := range pokemon.Stats {
		fmt.Printf("- %s: %d\n", s.Stat.Name, s.BaseStat)
	}

	return nil
}

func main() {
	cache := pokecache.NewCache(5 * time.Minute)

	scanner := bufio.NewScanner(os.Stdin)
	apiConfig := config{
		Url:      apiUrl,
		Next:     "",
		Previous: "",
		Cache:    cache,
	}

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		text := scanner.Text()
		words := cleanInput(text)
		if len(words) == 0 {
			continue
		}
		command := words[0]

		if cmd, ok := commands[command]; ok {
			err := cmd.callback(&apiConfig, words[1:]...)
			if err != nil {
				fmt.Println("Error:", err)
			}
		} else {
			fmt.Println("Unknown command:", command)
		}
	}

}
