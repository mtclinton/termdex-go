package main

import (
	"database/sql"
	"fmt"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"image"
	_ "image/png"
	"log"
	"os"
	"strconv"
)

var grid *ui.Grid

type SearchPokemon struct {
	search string
}

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()
	loadDB()
	initializePokemon()

	var termPokemon SearchPokemon
	currentPokemon := getPokemon("1")

	termWidth, termHeight := ui.TerminalDimensions()

	// pokemon := displayPokemon()
	draw := func() {
		reader, err := os.Open("./sprites/" + strconv.Itoa(currentPokemon.Pokemon_id) + ".png")
		if err != nil {
			log.Fatal(err)
		}
		pimage, _, err := image.Decode(reader)
		if err != nil {
			log.Fatal(err)
		}

		img := widgets.NewImage(nil)
		image_width := termWidth / 10 * 7
		img.SetRect(0, 0, int(image_width), termHeight)
		img.Image = pimage

		ui.Render(img)
	}
	draw()

	uiEvents := ui.PollEvents()

	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return

		case "<Enter>":
			currentPokemon = getPokemon(termPokemon.search)
			termPokemon.search = ""
			draw()
		default:
			termPokemon.search = termPokemon.search + e.ID
		}

	}
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

func getPokemon(search string) NewPokemon {
	db, err := gorm.Open(sqlite.Open("pokemon.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	var pokemon NewPokemon
	if _, err := strconv.Atoi(search); err == nil {
		db.Where("pokemon_id = ?", search).First(&pokemon)
	} else {
		db.Where("name = ?", search).First(&pokemon)
	}
	return pokemon

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
