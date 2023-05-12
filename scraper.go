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
	Large           string
	Small           string
	Base_experience int
	Height          int
	Weight          int
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
	large_sprite, err := s.downloader.get_sprite("https://raw.githubusercontent.com/mtclinton/pokemon-sprites/master/large/" + data.Name)
	if err != nil {
		log.Print(err)
		return
	}
	small_sprite, err := s.downloader.get_sprite("https://raw.githubusercontent.com/mtclinton/pokemon-sprites/master/small/" + data.Name)
	if err != nil {
		log.Print(err)
		return
	}
	new_pokemon := NewPokemon{
		Pokemon_id:      pid,
		Name:            data.Name,
		Large:           string(large_sprite),
		Small:           string(small_sprite),
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
	for i := 1; i <= 15; i++ {
		p := strconv.Itoa(i)
		ch <- p
	}
}
func (s *Scraper) run() {

	queue := make(chan string)

	workers := 8

	s.wg.Add(15)

	go s.pokeGenerator(queue)

	for i := 0; i < workers; i++ {

		go s.pokeHandler(queue, &s.wg)
	}

	s.wg.Wait()
	insertPokemon(s.pokemon_data[0])

}

func (s *Scraper) pokeHandler(ch chan string, wg *sync.WaitGroup) {
	for p := range ch {
		intVar, _ := strconv.Atoi(p)
		s.handle_url("https://pokeapi.co/api/v2/pokemon/"+p, intVar)
		wg.Done()
	}
}

func insertPokemon(p NewPokemon) {
	db, err := gorm.Open(sqlite.Open("pokemon.db"), &gorm.Config{})
	if err != nil {
		log.Print((err))
	}
	log.Print(p.Name, p.Pokemon_id)
	result := db.Create(&p)
	if result.Error != nil {
		log.Print((err))
	}
}
