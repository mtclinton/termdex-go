package main

import (
	"golang.org/x/exp/slices"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"strconv"
	"sync"
)

var (
	/// Maximum number of empty recv() from the channel
	MAX_EMPTY_RECEIVES = 10
	/// Sleep duration on empty recv()
	SLEEP_MILLIS = 100
)

type NewPokemon struct {
	ID              uint `gorm:"primarykey"	`
	Pokemon_id      int
	Name            string
	Base_experience int
	Height          int
	Weight          int
	HP              int
	Attack          int
	Defense         int
	Special_attack  int
	Special_defense int
	Speed           int
	Entry           string
}

func (NewPokemon) TableName() string {
	return "pokemon"
}

type PokemonType struct {
	ID         uint `gorm:"primarykey" `
	Pokemon_id int
	Type_id    int
}

func (PokemonType) TableName() string {
	return "pokemon_type"
}

type TypeName struct {
	ID   uint `gorm:"primarykey" `
	Name string
	URL  string
}

func (TypeName) TableName() string {
	return "type_name"
}

type MaxStats struct {
	ID              uint `gorm:"primarykey" `
	HP              int
	Attack          int
	Defense         int
	Special_attack  int
	Special_defense int
	Speed           int
}

type TypePokeTracker struct {
	Pokemon_id int
	Name       string
}

func (MaxStats) TableName() string {
	return "max_stats"
}

type Scraper struct {
	wg           sync.WaitGroup
	mu           sync.Mutex
	balance      int
	downloader   Downloader
	pokemon_data []NewPokemon
	maxstats     MaxStats
	type_names   []TypeName
	tpt          []TypePokeTracker
}

func NewScraper() Scraper {
	return Scraper{
		downloader: NewDownloader(3),
	}
}

func (s *Scraper) save_pokemon(data PokemonAPIData, entry_data EntryAPIData, pid int) {
	var hp, attack, defense, special_attack, special_defense, speed int
	for _, stat := range data.Stats {
		switch stat_name := stat.StatName.Name; stat_name {
		case "hp":
			hp = stat.BaseStat
		case "attack":
			attack = stat.BaseStat
		case "defense":
			defense = stat.BaseStat
		case "special-attack":
			special_attack = stat.BaseStat
		case "special-defense":
			special_defense = stat.BaseStat
		case "speed":
			speed = stat.BaseStat
		}
	}
	entry := entry_data.Entries[0].EntryText
	entry = strings.ReplaceAll(entry, "\n", " ")
	entry = strings.ReplaceAll(entry, "\u000c", " ")

	new_pokemon := NewPokemon{
		Pokemon_id:      pid,
		Name:            data.Name,
		Base_experience: data.BaseExperience,
		Height:          data.Height,
		Weight:          data.Weight,
		HP:              hp,
		Attack:          attack,
		Defense:         defense,
		Special_attack:  special_attack,
		Special_defense: special_defense,
		Speed:           speed,
		Entry:           entry,
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.pokemon_data = append(s.pokemon_data, new_pokemon)
	if hp > s.maxstats.HP {
		s.maxstats.HP = hp
	}
	if attack > s.maxstats.Attack {
		s.maxstats.Attack = attack
	}
	if defense > s.maxstats.Defense {
		s.maxstats.Defense = defense
	}
	if special_attack > s.maxstats.Special_attack {
		s.maxstats.Special_attack = special_attack
	}
	if special_defense > s.maxstats.Special_defense {
		s.maxstats.Special_defense = special_defense
	}
	if speed > s.maxstats.Speed {
		s.maxstats.Speed = speed
	}
	for _, t := range data.Types {
		idx := slices.IndexFunc(s.type_names, func(tn TypeName) bool { return tn.Name == t.TypeDetail.Name })
		if idx == -1 {
			new_type := TypeName{
				Name: t.TypeDetail.Name,
				URL:  t.TypeDetail.URL,
			}
			s.type_names = append(s.type_names, new_type)
		}
		new_tpt := TypePokeTracker{
			Pokemon_id: pid,
			Name:       t.TypeDetail.Name,
		}
		s.tpt = append(s.tpt, new_tpt)
	}

}

func (s *Scraper) handle_url(pid string) {
	url := "https://pokeapi.co/api/v2/pokemon/" + p
	api_data, err := s.downloader.get(url)
	if err != nil {
		log.Print(err)
		return
	}
	entry_url := "https://pokeapi.co/api/v2/pokemon-species/" + p
	entry_data, err := s.downloader.get_entry(entry_url)
	if err != nil {
		log.Print(err)
		return
	}

	s.save_pokemon(api_data, entry_data, pid)
}

func (s *Scraper) pokeGenerator(ch chan string) {
	defer close(ch)
	for i := 1; i <= 151; i++ {
		p := strconv.Itoa(i)
		ch <- p
	}
}
func (s *Scraper) run() {

	queue := make(chan string)

	workers := 8

	s.wg.Add(151)

	go s.pokeGenerator(queue)

	for i := 0; i < workers; i++ {

		go s.pokeHandler(queue, &s.wg)
	}

	s.wg.Wait()
	insertPokemon(s.pokemon_data)
	insertMaxStats(s.maxstats)
	insertTypeName(s.type_names)
	insertPokeType(s.tpt)

}

func (s *Scraper) pokeHandler(ch chan string, wg *sync.WaitGroup) {
	for p := range ch {
		intVar, _ := strconv.Atoi(p)
		s.handle_url("https://pokeapi.co/api/v2/pokemon/", p)
		wg.Done()
	}
}

func insertPokemon(pokemon_results []NewPokemon) {
	db, err := gorm.Open(sqlite.Open("pokemon.db"), &gorm.Config{})
	if err != nil {
		log.Print((err))
	}

	notfound := NewPokemon{
		Pokemon_id: 0,
		Name:       "Not Found",
	}
	pokemon_results = append(pokemon_results, notfound)

	result := db.Create(&pokemon_results)
	if result.Error != nil {
		log.Print((err))
	}
}

func insertMaxStats(max_stats MaxStats) {
	db, err := gorm.Open(sqlite.Open("pokemon.db"), &gorm.Config{})
	if err != nil {
		log.Print((err))
	}

	result := db.Create(&max_stats)
	if result.Error != nil {
		log.Print((err))
	}
}

func insertTypeName(type_names []TypeName) {
	db, err := gorm.Open(sqlite.Open("pokemon.db"), &gorm.Config{})
	if err != nil {
		log.Print((err))
	}

	result := db.Create(&type_names)
	if result.Error != nil {
		log.Print((err))
	}
}

func insertPokeType(type_poke_tracker []TypePokeTracker) {
	db, err := gorm.Open(sqlite.Open("pokemon.db"), &gorm.Config{})
	if err != nil {
		log.Print((err))
	}
	var type_names []TypeName
	type_name_results := db.Find(&type_names)
	if type_name_results.Error != nil {
		log.Print((err))
	}

	var poke_types []PokemonType
	for _, tracker := range type_poke_tracker {
		tid := slices.IndexFunc(type_names, func(tn TypeName) bool { return tn.Name == tracker.Name })
		if tid == -1 {
			log.Panic(("Unknown Type"))
		}
		idx := int(tid)
		poke_type := PokemonType{
			Pokemon_id: tracker.Pokemon_id,
			Type_id:    int(type_names[idx].ID),
		}
		poke_types = append(poke_types, poke_type)
	}
	result := db.Create(&poke_types)
	if result.Error != nil {
		log.Print((err))
	}
}
