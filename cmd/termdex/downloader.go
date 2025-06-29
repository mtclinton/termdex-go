package termdex

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

type PokemonAPIData struct {
	Name           string     `json:"name"`
	BaseExperience int        `json:"base_experience"`
	Height         int        `json:"height"`
	Weight         int        `json:"weight"`
	Stats          []Stat     `json:"stats"`
	Types          []PokeType `json:"types"`
}

type EntryAPIData struct {
	Entries []Entry `json:"flavor_text_entries"`
}

type Entry struct {
	EntryText string   `json:"flavor_text"`
	EntryLan  EntryLan `json:"language"`
}

type EntryLan struct {
	Name string
}

type Stat struct {
	BaseStat int      `json:"base_stat"`
	Effort   int      `json:"effort"`
	StatName StatName `json:"stat"`
}

type StatName struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type PokeType struct {
	TypeDetail TypeDetail `json:"type"`
}

type TypeDetail struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// /A Downloader to download web content
type Downloader struct {
	client http.Client
	tries  int
}

func NewDownloader(tries int) Downloader {
	return Downloader{
		client: http.Client{},
		tries:  tries,
	}
}

func (d Downloader) make_request(url string) (PokemonAPIData, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return PokemonAPIData{}, err
	}
	res, getErr := d.client.Do(req)
	if getErr != nil {
		return PokemonAPIData{}, getErr
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return PokemonAPIData{}, readErr
	}

	apiData := PokemonAPIData{}
	jsonErr := json.Unmarshal(body, &apiData)
	if jsonErr != nil {
		return PokemonAPIData{}, jsonErr
	}
	return apiData, nil
}

func (d Downloader) get(url string) (PokemonAPIData, error) {
	for i := 0; i < d.tries; i++ {
		apiData, err := d.make_request(url)
		if err == nil {
			return apiData, nil
		}
		if i+1 == d.tries {
			return PokemonAPIData{}, err
		}
	}
	return PokemonAPIData{}, errors.New("Something went wrong downloading")
}

func (d Downloader) make_entry_request(url string) (EntryAPIData, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return EntryAPIData{}, err
	}
	res, getErr := d.client.Do(req)
	if getErr != nil {
		return EntryAPIData{}, getErr
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return EntryAPIData{}, readErr
	}

	apiData := EntryAPIData{}
	jsonErr := json.Unmarshal(body, &apiData)
	if jsonErr != nil {
		return EntryAPIData{}, jsonErr
	}
	return apiData, nil
}

func (d Downloader) get_entry(url string) (EntryAPIData, error) {
	for i := 0; i < d.tries; i++ {
		apiData, err := d.make_entry_request(url)
		if err == nil {
			return apiData, nil
		}
		if i+1 == d.tries {
			return EntryAPIData{}, err
		}
	}
	return EntryAPIData{}, errors.New("Something went wrong downloading entry")
}
