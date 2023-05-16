package main

import (
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
	SpecialAttack   int
	SpecialDefense  int
	Speed           int
}

func (NewPokemon) TableName() string {
	return "pokemon"
}

type Scraper struct {
	wg           sync.WaitGroup
	mu           sync.Mutex
	balance      int
	downloader   Downloader
	pokemon_data []NewPokemon
}

func NewScraper() Scraper {
	return Scraper{
		downloader: NewDownloader(3),
	}
}

func (s *Scraper) save_pokemon(data PokemonAPIData, pid int) {
	new_pokemon := NewPokemon{
		Pokemon_id:      pid,
		Name:            data.Name,
		Base_experience: data.BaseExperience,
		Height:          data.Height,
		Weight:          data.Weight,
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.pokemon_data = append(s.pokemon_data, new_pokemon)

}

func (s *Scraper) handle_url(url string, pid int) {
	api_data, err := s.downloader.get(url)
	if err != nil {
		log.Print(err)
		return
	}
	s.save_pokemon(api_data, pid)
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

}

func (s *Scraper) pokeHandler(ch chan string, wg *sync.WaitGroup) {
	for p := range ch {
		intVar, _ := strconv.Atoi(p)
		s.handle_url("https://pokeapi.co/api/v2/pokemon/"+p, intVar)
		wg.Done()
	}
}

func insertPokemon(pokemon_results []NewPokemon) {
	db, err := gorm.Open(sqlite.Open("pokemon.db"), &gorm.Config{})
	if err != nil {
		log.Print((err))
	}

	notfound := NewPokemon{
		Pokemon_id:      0,
		Name:            "Not Found",
		Base_experience: -1,
		Height:          -1,
		Weight:          -1,
	}
	pokemon_results = append(pokemon_results, notfound)

	result := db.Create(&pokemon_results)
	if result.Error != nil {
		log.Print((err))
	}
}
