package termdex-go

import (
	"errors"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestMain(m *testing.M) {
	code, err := run(m)
	if err != nil {
		fmt.Println(err)
	}
	os.Exit(code)
}

func run(m *testing.M) (code int, err error) {

	poke_db, _ := loadDB("test.db")
	poke_db.createTable()
	defer func() {
		poke_db.db.Exec("DELETE FROM pokemon")
		poke_db.db.Exec("DELETE FROM max_stats")
		poke_db.db.Exec("DELETE FROM pokemon_type")
		poke_db.db.Exec("DELETE FROM type_name")
		os.Remove("test.db")

	}()

	return m.Run(), nil
}

func TestInsertPokemon(t *testing.T) {
	pokemon := NewPokemon{
		ID:              1,
		Pokemon_id:      1,
		Name:            "bulbasaur",
		Base_experience: 2,
		Height:          3,
		Weight:          4,
		HP:              5,
		Attack:          6,
		Defense:         7,
		Special_attack:  8,
		Special_defense: 9,
		Speed:           10,
		Entry:           "tst entry",
	}
	var pokemon_data []NewPokemon
	pokemon_data = append(pokemon_data, pokemon)

	pkd, _ := loadDB("test.db")
	pkd.insertPokemon(pokemon_data)

	p, _ := pkd.getPokemon("bulbasaur")

	assert.Equal(t, 1, p.Pokemon_id)
	assert.Equal(t, "bulbasaur", p.Name)
}

func TestInsertMaxStats(t *testing.T) {
	max_stats := MaxStats{
		1, 2, 3, 4, 5, 6, 7,
	}

	pkd, _ := loadDB("test.db")
	pkd.insertMaxStats(max_stats)

	ms := pkd.getMaxStats()

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
