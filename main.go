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
    . "github.com/gizak/termui/v3"
    "strings"
    "golang.org/x/text/cases"
    "golang.org/x/text/language"
)

var grid *ui.Grid

const values = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type SearchPokemon struct {
	search string
}

type Gauge struct {
    Block
    Percent    int
    BarColor   Color
    Label      string
    LabelStyle Style
}

func NewGauge() *Gauge {
    return &Gauge{
        Block:      *NewBlock(),
        BarColor:   Theme.Gauge.Bar,
        LabelStyle: Theme.Gauge.Label,
    }
}

func (self *Gauge) Draw(buf *Buffer) {
    self.Block.Draw(buf)

    label := ""

    // plot bar
    barWidth := int((float64(self.Percent) / 100) * float64(self.Inner.Dx()))
    buf.Fill(
        NewCell(' ', NewStyle(ColorClear, self.BarColor)),
        image.Rect(self.Inner.Min.X, self.Inner.Min.Y, self.Inner.Min.X+barWidth, self.Inner.Max.Y),
    )

    // plot label
    labelXCoordinate := self.Inner.Min.X + (self.Inner.Dx() / 2) - int(float64(len(label))/2)
    labelYCoordinate := self.Inner.Min.Y + ((self.Inner.Dy() - 1) / 2)
    if labelYCoordinate < self.Inner.Max.Y {
        for i, char := range label {
            style := self.LabelStyle
            if labelXCoordinate+i+1 <= self.Inner.Min.X+barWidth {
                style = NewStyle(self.BarColor, ColorClear, ModifierReverse)
            }
            buf.SetCell(NewCell(char, style), image.Pt(labelXCoordinate+i, labelYCoordinate))
        }
    }
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
        // n.Title = "Name"
		n.Text = "["+cases.Title(language.English).String(currentPokemon.Name)+"](fg:blue,mod:bold)"
		n.SetRect(image_width, 5, termWidth, 8)
        n.Border = false

		ui.Render(n)

        height := widgets.NewParagraph()
        // height.Title = "Height"
        height.Text = "Height: "+strconv.Itoa(currentPokemon.Height)
        height.SetRect(image_width, 9, termWidth, 12)
        height.Border = false

        ui.Render(n, height)

        hp := NewGauge()
        hp.Title = "HP"
        hp.SetRect(image_width, termHeight-24, termWidth, termHeight-20)
        hp.Percent = currentPokemon.HP
        hp.BarColor = ui.ColorGreen
        hp.Border = false

        attack := NewGauge()
        attack.Title = "Attack"
        attack.SetRect(image_width, termHeight-20, termWidth, termHeight-16)
        attack.Percent = currentPokemon.Attack
        attack.BarColor = ui.ColorRed
        attack.Border = false

		defense := NewGauge()
        defense.Title = "Defense"
        defense.SetRect(image_width, termHeight-16, termWidth, termHeight-12)
        defense.Percent = currentPokemon.Defense
        defense.BarColor = ui.ColorBlue
        defense.Border = false

        special_attack := NewGauge()
        special_attack.Title = "Special Attack"
        special_attack.SetRect(image_width, termHeight-12, termWidth, termHeight-8)
        special_attack.Percent = currentPokemon.Special_defense
        special_attack.BarColor = ui.ColorMagenta
        special_attack.Border = false

        special_defense := NewGauge()
        special_defense.Title = "Special Defense"
        special_defense.SetRect(image_width, termHeight-8, termWidth, termHeight-4)
        special_defense.Percent = currentPokemon.Special_attack
        special_defense.BarColor = ui.ColorCyan
        special_defense.Border = false

        speed := NewGauge()
        speed.Title = "Speed"
        speed.SetRect(image_width, termHeight-4, termWidth, termHeight)
        speed.Percent = currentPokemon.Speed
        speed.BarColor = ui.ColorYellow
        speed.Border = false

        ui.Render(hp, attack, defense, special_attack, special_defense, speed)
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
