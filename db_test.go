package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	code, err := run(m)
	if err != nil {
		fmt.Println(err)
	}
	os.Exit(code)
}

func run(m *testing.M) (code int, err error) {

	db, err := sql.Open("sqlite3", "file:test.db?cache=shared")
	if err != nil {
		return -1, fmt.Errorf("could not connect to database: %w", err)
	}

	createTable(db)
	defer func() {
		db.Exec(fmt.Sprintf("DELETE FROM pokemon"))
		db.Exec(fmt.Sprintf("DELETE FROM max_stats"))
		db.Exec(fmt.Sprintf("DELETE FROM pokemon_type"))
		db.Exec(fmt.Sprintf("DELETE FROM type_name"))

		db.Close()
	}()

	return m.Run(), nil
}

func TestInsertMaxStats(t *testing.T) {
	max_stats := MaxStats{
		1, 2, 3, 4, 5, 6, 7,
	}

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Print((err))
	}

	result := db.Create(&max_stats)
	if result.Error != nil {
		log.Print(result.Error)
	}

	// using https://github.com/stretchr/testify library for brevity
	require.NoError(t, err)

	var ms MaxStats
	result = db.First(&ms)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		panic("No max stats value")
	}

	assert.Equal(t, 2, ms.HP)
	assert.Equal(t, 3, ms.Attack)
	assert.Equal(t, 4, ms.Defense)
	assert.Equal(t, 5, ms.Special_attack)
	assert.Equal(t, 6, ms.Special_defense)
	assert.Equal(t, 7, ms.Speed)

}

func TestInsertPokemonType(t *testing.T) {
	pokemon_type := PokemonType{
		Pokemon_id: 1,
		Type_id:    1,
	}

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Print((err))
	}

	result := db.Create(&pokemon_type)
	if result.Error != nil {
		log.Print(result.Error)
	}

	// using https://github.com/stretchr/testify library for brevity
	require.NoError(t, err)

	var pt PokemonType
	result = db.First(&pt)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		panic("No type name value")
	}

	assert.Equal(t, 1, pt.Pokemon_id)
	assert.Equal(t, 1, pt.Type_id)
}

func TestInsertTypeName(t *testing.T) {
	type_name := TypeName{
		Name: "grass",
		URL:  "example.com/grass",
	}

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Print((err))
	}

	result := db.Create(&type_name)
	if result.Error != nil {
		log.Print(result.Error)
	}

	// using https://github.com/stretchr/testify library for brevity
	require.NoError(t, err)

	var tn TypeName
	result = db.First(&tn)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		panic("No type name value")
	}

	assert.Equal(t, "grass", tn.Name)
	assert.Equal(t, "example.com/grass", tn.URL)
}
