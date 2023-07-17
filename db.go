package main

import (
	"database/sql"
	"errors"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"reflect"
	"strconv"
)

func loadDB() {
	if _, err := os.Stat("pokemon.db"); err != nil {
		log.Println("Creating sqlite-database.db...")
		file, err := os.Create("sqlite-database.db") // Create SQLite file
		if err != nil {
			log.Fatal(err.Error())
		}
		file.Close()
		log.Println("sqlite-database.db created")
	}
	sqliteDatabase, _ := sql.Open("sqlite3", "./pokemon.db") // Open the created SQLite File
	defer sqliteDatabase.Close()                             // Defer Closing the database
	createTable(sqliteDatabase)                              // Create Database Tables

}

func createTable(db *sql.DB) {
	createPokemonTableSQL := `CREATE TABLE IF NOT EXISTS pokemon (
        "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,     
        "pokemon_id" integer NOT NULL,
        "name" TEXT NOT NULL,
        "base_experience" integer NOT NULL,
        "height" integer NOT NULL,
        "weight" integer NOT NULL,
        "hp" integer NOT NULL,
        "attack" integer NOT NULL,
        "defense" integer NOT NULL,
        "special_attack" integer NOT NULL,
        "special_defense" integer NOT NULL,
        "speed" integer NOT NULL,
        "entry" TEXT NOT NULL
      );` // SQL Statement for Create Table

	log.Println("Create pokemon table...")
	statement, err := db.Prepare(createPokemonTableSQL) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec() // Execute SQL Statements
	log.Println("Pokemon table created")

	createMaxStatsTableSQL := `CREATE TABLE IF NOT EXISTS max_stats (
        "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,     
        "hp" integer NOT NULL,
        "attack" integer NOT NULL,
        "defense" integer NOT NULL,
        "special_attack" integer NOT NULL,
        "special_defense" integer NOT NULL,
        "speed" integer NOT NULL
      );` // SQL Statement for Create Table

	log.Println("Create max stats table...")
	stats_statement, err := db.Prepare(createMaxStatsTableSQL) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	stats_statement.Exec() // Execute SQL Statements
	log.Println("Max stats table created")

	createPokemonTypeTableSQL := `CREATE TABLE IF NOT EXISTS pokemon_type (
        "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,     
        "pokemon_id" integer NOT NULL,
        "type_id" integer NOT NULL
      );` // SQL Statement for Create Table

	log.Println("Create pokemon type table...")
	poke_type_statement, err := db.Prepare(createPokemonTypeTableSQL) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	poke_type_statement.Exec() // Execute SQL Statements
	log.Println("Pokemon type table created")

	createTypeNameTableSQL := `CREATE TABLE IF NOT EXISTS type_name (
        "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,     
        "name" TEXT NOT NULL,
        "url" TEXT NOT NULL
      );` // SQL Statement for Create Table

	log.Println("Create type name table...")
	type_name_statement, err := db.Prepare(createTypeNameTableSQL) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	type_name_statement.Exec() // Execute SQL Statements
	log.Println("Type name table created")
}

func getPokemon(search string) (NewPokemon, []string) {
	db, err := gorm.Open(sqlite.Open("pokemon.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic("failed to connect database")
	}
	var pokemon NewPokemon
	var pokemon_types []PokemonType
	var type_names []TypeName
	var names []string
	if _, err := strconv.Atoi(search); err == nil {
		result := db.Where("pokemon_id = ?", search).First(&pokemon)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			db.Where("name = ").First(&pokemon)
		}
	} else {
		result := db.Where("name = ?", search).First(&pokemon)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			db.Where("name = ").First(&pokemon)
		}
	}
	types_result := db.Where("pokemon_id = ?", pokemon.Pokemon_id).Find(&pokemon_types)
	if !errors.Is(types_result.Error, gorm.ErrRecordNotFound) {
		var pokemon_types_ids []int
		for _, pokemon_type := range pokemon_types {
			r := reflect.ValueOf(pokemon_type)
			f := reflect.Indirect(r).FieldByName("Type_id")
			pokemon_types_ids = append(pokemon_types_ids, int(f.Int()))
		}
		db.Where("ID in ?", pokemon_types_ids).Find(&type_names)
		for _, type_name := range type_names {
			r := reflect.ValueOf(type_name)
			f := reflect.Indirect(r).FieldByName("Name")
			names = append(names, f.String())
		}
	}

	return pokemon, names
}

func getMaxStats() MaxStats {
	db, err := gorm.Open(sqlite.Open("pokemon.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic("failed to connect database")
	}
	var maxStats MaxStats
	result := db.First(&maxStats)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		panic("No max stats value")
	}
	return maxStats

}

func initializePokemon() {
	db, err := gorm.Open(sqlite.Open("pokemon.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	var count int64
	db.Table("pokemon").Count(&count)
	if count == 0 {
		fmt.Println("Initializing pokemon")
		s := NewScraper()
		s.run()
	}

}
