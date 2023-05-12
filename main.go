package main

import (
	"database/sql"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
)

func main() {

	loadDB()
	initializePokemon()
	displayPokemon()

}

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
        "large" TEXT NOT NULL,
        "small" TEXT NOT NULL,
        "base_experience" integer NOT NULL,
        "height" integer NOT NULL,
        "weight" integer NOT NULL
      );` // SQL Statement for Create Table

	log.Println("Create pokemon table...")
	statement, err := db.Prepare(createPokemonTableSQL) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec() // Execute SQL Statements
	log.Println("Pokemon table created")
}

func displayPokemon() {
	db, err := gorm.Open(sqlite.Open("pokemon.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	var pokemon NewPokemon
	db.First(&pokemon, 2) // find product with integer primary key
	fmt.Println(pokemon.Small)

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
