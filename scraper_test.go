package main

import (
	"github.com/stretchr/testify/assert"

	"testing"
)

func TestScraper(t *testing.T) {
	hp := StatName{"hp", "https://pokeapi.co/api/v2/stat/1/"}
	attack := StatName{"attack", "https://pokeapi.co/api/v2/stat/2/"}
	defense := StatName{"defense", "https://pokeapi.co/api/v2/stat/3/"}
	special_attack := StatName{"special-attack", "https://pokeapi.co/api/v2/stat/4/}"}
	special_defense := StatName{"special-defense", "https://pokeapi.co/api/v2/stat/5/"}
	speed := StatName{"speed", "https://pokeapi.co/api/v2/stat/6/"}

	pokemon_stats := []Stat{
		Stat{45, 0, hp},
		Stat{49, 0, attack},
		Stat{49, 0, defense},
		Stat{65, 1, special_attack},
		Stat{65, 0, special_defense},
		Stat{45, 0, speed},
	}

	pokemon_types := []PokeType{
		PokeType{TypeDetail{
			"poison",
			"https://pokeapi.co/api/v2/type/4/",
		}},
		PokeType{TypeDetail{
			"grass",
			"https://pokeapi.co/api/v2/type/12/",
		}},
	}

	pokemon_data := PokemonAPIData{
		"bulbasaur",
		64,
		7,
		69,
		pokemon_stats,
		pokemon_types,
	}

	entry_data := EntryAPIData{
		[]Entry{
			Entry{
				"A strange seed was\nplanted on its\nback at birth.\u000cThe plant sprouts\nand grows with\nthis POKéMON.",
				EntryLan{"en"},
			},
		},
	}

	scraper := NewScraper()
	scraper.save_pokemon(pokemon_data, entry_data, 1)
	assert.Equal(t, "bulbasaur", scraper.pokemon_data[0].Name)

}
