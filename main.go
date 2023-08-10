package main

import (
	. "github.com/gizak/termui/v3"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"image"
	_ "image/png"
	"log"
	"os"
	"strconv"
	"strings"
)

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
	pokemon_db, _ := loadDB("pokemon.db")
	pokemon_db.initializePokemon()

	var termPokemon SearchPokemon
	currentPokemon, types := pokemon_db.getPokemon("25")

	maxStats := pokemon_db.getMaxStats()

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
		n.Text = "[" + cases.Title(language.English).String(currentPokemon.Name) + "](fg:yellow,mod:bold)"
		n.SetRect(image_width+2, 5, termWidth, 8)
		n.Border = false

		ui.Render(n)

		e := widgets.NewParagraph()
		// n.Title = "Entry"
		e.Text = "[" + cases.Title(language.English).String(currentPokemon.Entry) + "](fg:cyan,mod:bold)"
		e.SetRect(image_width+2, 8, termWidth, 13)
		e.Border = false

		ui.Render(e)

		stats_table := widgets.NewTable()
		stats_table.Rows = [][]string{
			[]string{"HP", "Attack", "Defense", "S Attack", "S Defense", "Speed"},
			[]string{strconv.Itoa(currentPokemon.HP), strconv.Itoa(currentPokemon.Attack), strconv.Itoa(currentPokemon.Defense),
				strconv.Itoa(currentPokemon.Special_attack), strconv.Itoa(currentPokemon.Special_defense), strconv.Itoa(currentPokemon.Speed)},
		}
		stats_table.TextStyle = ui.NewStyle(ui.ColorWhite)
		stats_table.TextAlignment = ui.AlignCenter
		// stats_table.RowSeparator = false
		stats_table.Border = false
		stats_table.SetRect(image_width, 13, termWidth, 18)

		ui.Render(stats_table)

		type_title := widgets.NewParagraph()
		// type_title.Title = "type_title"
		type_title.Text = "[" + "Types" + "](fg:yellow,mod:bold)"
		type_title.SetRect((termWidth-image_width)/2+image_width, 18, termWidth, 21)
		type_title.Border = false

		type_rewrite := widgets.NewParagraph()
		type_rewrite.SetRect(image_width, 21, termWidth, 24)
		type_rewrite.Border = false
		ui.Render(type_rewrite)

		if len(types) == 1 {
			type_data := widgets.NewParagraph()
			// type_title.Title = "type_title"
			type_data.Text = "[" + cases.Title(language.English).String(types[0]) + "](fg:cyan,mod:bold)"
			type_data.SetRect((termWidth-image_width)/2+image_width, 21, termWidth, 24)
			type_data.Border = false
			ui.Render(type_data)
		} else {
			type_data1 := widgets.NewParagraph()
			// type_title.Title = "type_title"
			type_data1.Text = "[" + cases.Title(language.English).String(types[0]) + "](fg:cyan,mod:bold)"
			type_data1.SetRect((termWidth-image_width)/4+image_width, 21, (termWidth-image_width)/2+image_width, 24)
			type_data1.Border = false

			type_data2 := widgets.NewParagraph()
			// type_title.Title = "type_title"
			type_data2.Text = "[" + cases.Title(language.English).String(types[1]) + "](fg:cyan,mod:bold)"
			type_data2.SetRect(((termWidth-image_width)/4*3)+image_width, 21, termWidth, 24)
			type_data2.Border = false

			ui.Render(type_data1, type_data2)

		}

		height := widgets.NewParagraph()
		// height.Title = "Height"
		height.Text = "[" + "Height: " + "](fg:blue)[" + strconv.Itoa(currentPokemon.Height) + "](fg:red,mod:bold)"
		height.SetRect((termWidth-image_width)/4+image_width, 24, (termWidth-image_width)/2+image_width, 27)
		height.Border = false

		weight := widgets.NewParagraph()
		// weight.Title = "Height"
		weight.Text = "[" + "Weight: " + "](fg:red)[" + strconv.Itoa(currentPokemon.Weight) + "](fg:blue,mod:bold)"
		weight.SetRect(((termWidth-image_width)/4*3)+image_width, 24, termWidth, 27)
		weight.Border = false

		ui.Render(n, type_title, height, weight)

		stats_title := widgets.NewParagraph()
		stats_title.Text = "[Stats](fg:green,mod:bold)"
		stats_title.SetRect((termWidth+image_width)/2, termHeight-27, termWidth, termHeight-24)
		stats_title.Border = false
		ui.Render(stats_title)

		hp := NewGauge()
		hp.Title = "HP"
		hp.SetRect(image_width, termHeight-24, termWidth, termHeight-20)
		hp.Percent = currentPokemon.HP * 100 / maxStats.HP
		hp.BarColor = ui.ColorGreen
		hp.Border = false

		attack := NewGauge()
		attack.Title = "Attack"
		attack.SetRect(image_width, termHeight-20, termWidth, termHeight-16)
		attack.Percent = currentPokemon.Attack * 100 / maxStats.Attack
		attack.BarColor = ui.ColorRed
		attack.Border = false

		defense := NewGauge()
		defense.Title = "Defense"
		defense.SetRect(image_width, termHeight-16, termWidth, termHeight-12)
		defense.Percent = currentPokemon.Defense * 100 / maxStats.Defense
		defense.BarColor = ui.ColorBlue
		defense.Border = false

		special_attack := NewGauge()
		special_attack.Title = "Special Attack"
		special_attack.SetRect(image_width, termHeight-12, termWidth, termHeight-8)
		special_attack.Percent = currentPokemon.Special_defense * 100 / maxStats.Special_defense
		special_attack.BarColor = ui.ColorMagenta
		special_attack.Border = false

		special_defense := NewGauge()
		special_defense.Title = "Special Defense"
		special_defense.SetRect(image_width, termHeight-8, termWidth, termHeight-4)
		special_defense.Percent = currentPokemon.Special_attack * 100 / maxStats.Special_attack
		special_defense.BarColor = ui.ColorCyan
		special_defense.Border = false

		speed := NewGauge()
		speed.Title = "Speed"
		speed.SetRect(image_width, termHeight-4, termWidth, termHeight)
		speed.Percent = currentPokemon.Speed * 100 / maxStats.Speed
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
			currentPokemon, types = pokemon_db.getPokemon(termPokemon.search)
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
