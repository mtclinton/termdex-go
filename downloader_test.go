package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
    "log"
)

func router() *gin.Engine {
	router := gin.Default()
	pokemon_json, err := os.Open("pokemon_test.json")
	if err != nil {
		log.Println(err)
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer pokemon_json.Close()
	router.GET("/pokemon/1", func(c *gin.Context) {
		c.JSON(http.StatusOK, pokemon_json)
	})

	return router
}

func makeRequest(method, url string) *httptest.ResponseRecorder {
	request, _ := http.NewRequest(method, url, nil)
	writer := httptest.NewRecorder()
	router().ServeHTTP(writer, request)
	return writer
}

func TestPokemon(t *testing.T) {
	writer := makeRequest("GET", "/pokemon/1")
	assert.Equal(t, 200, writer.Code)
}
