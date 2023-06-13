package main

import (
	"database/sql"
	"errors"
	"fmt"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"image"
	_ "image/png"
	"log"
	"os"
	"strconv"
	"strings"
)

var grid *ui.Grid

const values = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

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

	draw := func() {
		var sprite_name string
		if currentPokemon.Pokemon_id == 0 {
			sprite_name = "pokedex"
		} else {
			sprite_name = strconv.Itoa(currentPokemon.Pokemon_id)
		}
		reader, err := os.Open("./sprites/" + sprite_name + ".png")
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
		img.Title = currentPokemon.Name
		img.Image = pimage

		ui.Render(img)

	}
	draw()

	drawInput := func() {
		image_width := int(termWidth / 10 * 7)
		p := widgets.NewParagraph()
		p.Text = termPokemon.search
		p.Title = "Search Pokemon"
		p.SetRect(image_width, 0, termWidth, 3)

		ui.Render(p)

		n := widgets.NewParagraph()
		n.Text = currentPokemon.Name
		n.SetRect(image_width, 5, termWidth, 8)

		ui.Render(n)

        hp := widgets.NewParagraph()
        hp.Text = strconv.Itoa(currentPokemon.HP)
        hp.SetRect(image_width, 10, termWidth, 13)

        ui.Render(hp)

        attack := widgets.NewParagraph()
        attack.Text = strconv.Itoa(currentPokemon.Attack)
        attack.SetRect(image_width, 15, termWidth, 18)

        defense := widgets.NewParagraph()
        defense.Text = strconv.Itoa(currentPokemon.Defense)
        defense.SetRect(image_width, 20, termWidth, 23)

        speed := widgets.NewParagraph()
        speed.Text = strconv.Itoa(currentPokemon.Speed)
        speed.SetRect(image_width, 26, termWidth, 29)

        ui.Render(attack)
        ui.Render(defense)
        ui.Render(speed)

        g0 := widgets.NewGauge()
        g0.Title = "Special Attack"
        g0.SetRect(image_width, 32, termWidth, 42)
        g0.Percent = currentPokemon.Special_attack
        g0.BarColor = ui.ColorRed
        g0.BorderStyle.Fg = ui.ColorWhite
        g0.TitleStyle.Fg = ui.ColorCyan
        ui.Render(g0)
	}

	drawInput()

	uiEvents := ui.PollEvents()

	for {
		e := <-uiEvents
		switch e.ID {
		case "<C-c>":
			return
		case "<Backspace>":
			if len(termPokemon.search) > 0 {
				termPokemon.search = termPokemon.search[:len(termPokemon.search)-1]
				drawInput()
			}

		case "<Enter>":
			currentPokemon = getPokemon(termPokemon.search)
			termPokemon.search = ""
			drawInput()
			draw()
		default:
			if strings.Contains(values, e.ID) {
				termPokemon.search = termPokemon.search + e.ID
			}
			drawInput()
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
        "base_experience" integer NOT NULL,
        "height" integer NOT NULL,
        "weight" integer NOT NULL,
        "hp" integer NOT NULL,
        "attack" integer NOT NULL,
        "defense" integer NOT NULL,
        "special_attack" integer NOT NULL,
        "special_defense" integer NOT NULL,
        "speed" integer NOT NULL
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
	db, err := gorm.Open(sqlite.Open("pokemon.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic("failed to connect database")
	}
	var pokemon NewPokemon
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
