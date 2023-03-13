package main

import "fmt"
import "termdex-go/pokeball"
import "os"
import "bufio"
import "strconv"
import "log"
import "strings"
import "net/http"
import "encoding/json"

type TypeName struct {
    Name            string `json:"name"`
}

type PokeType struct {
    Type            TypeName `json:"type"`
}

type PokeData struct{
    Name              string `json:"name"`
    Types             []PokeType `json:"types"`
}

func searchPokemon(pokemon_id int) {
    url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%d", pokemon_id)

	// Build the request
	req, err := http.NewRequest("GET", url, nil)

    if err != nil {
		log.Fatal("NewRequest: ", err)
		return
	}

    client := &http.Client{}

    resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return
	}

    defer resp.Body.Close()

    var pokemon PokeData

    if err := json.NewDecoder(resp.Body).Decode(&pokemon); err != nil {
		log.Println(err)
	}
    fmt.Println(pokemon.Name)
    for _, ptype := range pokemon.Types {
        fmt.Println(ptype.Type.Name)
    }
}

func main() {
	pokeball.ShowPokeball()
	fmt.Println("Welcome to TermDex")
	fmt.Println("Input a pokemon ID")
	reader := bufio.NewReader(os.Stdin)
	userInput, _ := reader.ReadString('\n')
    userInput = strings.ReplaceAll(userInput, "\n", "")
    pokemon_id, err := strconv.Atoi(userInput)
    if err != nil {
        log.Fatal(err)
    }
    searchPokemon(pokemon_id)
}
