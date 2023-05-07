package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type PokemonAPIData struct {
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
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