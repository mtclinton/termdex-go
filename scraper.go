package main

import (
   "fmt"
   "sync"
   "time"
   "log"
   "strconv"
)

var (
    /// Maximum number of empty recv() from the channel
    MAX_EMPTY_RECEIVES := 10;
    /// Sleep duration on empty recv()
    SLEEP_MILLIS := 100;
)

type Scraper struct {
    wg sync.WaitGroup
    mu      sync.Mutex
    balance int
    downloader Downloader
    pokemon_data []NewPokemon
}


func NewScraper() Scraper {
    return Scraper{
        wg: sync.WaitGroup,
        mu: sync.Mutex
        downloader: NewDownloader(3),
        pokemon_data: []string,
    }
}

func (s Scraper) save_pokemon(data downloader.PokemonAPIData, pid int) {
    l_path := "sprites/large/" + data.Name
    s_path := "sprites/small/" + data.name
    l_data, err := os.ReadFile(l_path)
    if err != nil {
        log.Print(err)
        return
    }
    s_data, err := os.ReadFile(s_path)
    if err != nil {
        log.Print(err)
        return
    }
    new_pokemon := NewPokemon {
            pokemon_id: pid,
            name: data.Name,
            large: l_data,
            small: s_data,
            base_experience: data.BaseExperience,
            height: data.Height,
            weight: data.Weight,
    }
    mu.Lock()
    defer mu.Unlock()
    s.pokemon_data = append(s.pokemon_data, new_pokemon)
}

func (s Scraper) handle_url(url string pid int) {
    api_data, err := s.downloader.get(url)
    if err != nil {
        log.Print(err)
        return
    }
    s.save_pokemon(api_data, pid)
}

func (s Scraper) pokeGenerator(ch chan string) {
    defer close(ch)
    for for i := 1; i <= 151; i++ {
        p := strconv.Itoa(i)
        ch <- p
    }
}
func (s Scraper) run() {

    queue := make(chan string)

    workers := 3

    s.wg.Add(151)

    go s.pokeGenerator(queue, tasks)

    for i := 0; i < workers; i++ {
      
       go s.pokeHandler(queue, &s.wg)
    }

    wg.Wait()

}

func (s Scraper) pokeHandler(ch chan string, wg *sync.WaitGroup) {
   for p := range ch {
        intVar, _ := strconv.Atoi(p)
        s.handle_url("https://pokeapi.co/api/v2/pokemon/"+p, intVar)
        wg.Done()
   }
}
