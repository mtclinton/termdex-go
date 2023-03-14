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
import "io/ioutil"
import "github.com/gookit/color"
import "errors"


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




func show_sprite(sprite string, poke_type string) {
    poke_colors := map[string][]int{
        "normal": []int{168, 167, 122},
        "fire": []int{238, 129, 48},
        "water": []int{99, 144, 240},
        "electric": []int{247, 208, 44},
        "grass": []int{122, 199, 76},
        "ice": []int{150, 217, 214},
        "fighting": []int{194, 46, 40},
        "poison": []int{163, 62, 161},
        "ground": []int{226, 191, 101},
        "flying": []int{169, 143, 243},
        "psychic": []int{249, 85, 135},
        "bug": []int{166, 185, 26},
        "rock": []int{182, 161, 54},
        "ghost": []int{115, 87, 151},
        "dragon": []int{111, 53, 252},
        "dark": []int{112, 87, 70},
        "steel": []int{183, 183, 206},
        "fairy": []int{214, 133, 173},
    }
    val, ok := poke_colors[poke_type]
    if ok {
        r, g, b := val[0], val[1], val[2]
        colorSprite := "<fg="+strconv.Itoa(r)+","+strconv.Itoa(g)+","+strconv.Itoa(b)+">"+sprite+"</>"
        color.Println(colorSprite)
    } else {
        log.Fatal("Incorrect type submitted")
    }
}

func show_pokemon(pokemon_id int) (string, error) {
    jsonPokemon, err := os.Open("pokemon.json")
    if err != nil {
        log.Fatal("NewRequest: ", err)
         return "", err
     }
    defer jsonPokemon.Close()
    byteResult, _ := ioutil.ReadAll(jsonPokemon)
    var res map[string]string
    json.Unmarshal([]byte(byteResult), &res)
    val, ok := res[strconv.Itoa(pokemon_id)]
    if ok {
        return val, nil
    } else {
        log.Fatal("Incorrect pokemon id submitted")
        return "", errors.New("Incorrect pokemon id submitted")
    }
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
    pokemon_types :=  []string{}
    for _, ptype := range pokemon.Types {
        pokemon_types = append(pokemon_types, ptype.Type.Name)
    }
    sprite, _ := show_pokemon(pokemon_id)
    show_sprite(sprite, pokemon_types[0])
    fmt.Println(pokemon.Name)
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
