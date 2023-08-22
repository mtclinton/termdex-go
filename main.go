package main

import (
	"image"
	_ "image/png"
	"log"
	"os"
	"strconv"
	"strings"

	. "github.com/gizak/termui/v3"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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
	ui.Theme.Block.Border = ui.NewStyle(ui.Color(ui.ColorRed))
	ui.Theme.Block.Title = ui.NewStyle(ui.Color(ui.ColorRed))
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
		img.Title = cases.Title(language.English).String(currentPokemon.Name)
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
		n.Title = cases.Title(language.English).String(currentPokemon.Name)
		n.SetRect(image_width, 5, termWidth, termHeight)

		ui.Render(n)

		e := widgets.NewParagraph()
		e.Text = "[" + cases.Title(language.English).String(currentPokemon.Entry) + "](fg:yellow,mod:bold)"
		e.SetRect(image_width+2, 6, termWidth-1, 11)
		e.Border = false

		ui.Render(e)

		height := widgets.NewParagraph()
		height.Title = "Height"
		height.SetRect(image_width+2, 13, (termWidth-image_width)/2+image_width-2, 19)

		height_text := widgets.NewParagraph()
		height_text.Text = "[" + strconv.Itoa(currentPokemon.Height) + "](fg:yellow,mod:bold)"
		height_text.SetRect(((termWidth-image_width)/4)+image_width, 14, (termWidth-image_width)/2+image_width-3, 17)
		height_text.Border = false

		weight := widgets.NewParagraph()
		weight.Title = "Weight"
		weight.SetRect((termWidth-image_width)/2+image_width+2, 13, termWidth-2, 19)

		weight_text := widgets.NewParagraph()
		weight_text.Text = "[" + strconv.Itoa(currentPokemon.Weight) + "](fg:yellow,mod:bold)"
		weight_text.SetRect(((termWidth-image_width)/4*3)+image_width, 14, termWidth-3, 17)
		weight_text.Border = false

		ui.Render(height, height_text, weight, weight_text)

		if len(types) == 1 {
			type_border := widgets.NewParagraph()
			type_border.Title = "Type"
			type_border.SetRect((termWidth-image_width)/4+image_width, 21, (termWidth-image_width)*3/4+image_width, 26)

			type_data := widgets.NewParagraph()
			type_data.Text = "[" + cases.Title(language.English).String(types[0]) + "](fg:yellow,mod:bold)"
			type_data.SetRect((termWidth-image_width)/4+image_width+6, 22, (termWidth-image_width)*3/4+image_width-2, 25)
			type_data.Border = false
			ui.Render(type_border, type_data)
		} else {
			type_border1 := widgets.NewParagraph()
			type_border1.Title = "Type"
			type_border1.SetRect(image_width+2, 21, (termWidth-image_width)/2+image_width-2, 26)

			type_data1 := widgets.NewParagraph()
			type_data1.Text = "[" + cases.Title(language.English).String(types[0]) + "](fg:yellow,mod:bold)"
			type_data1.SetRect(((termWidth-image_width)/4)+image_width, 22, (termWidth-image_width)/2+image_width-3, 25)
			type_data1.Border = false

			type_border2 := widgets.NewParagraph()
			type_border2.Title = "Type"
			type_border2.SetRect((termWidth-image_width)/2+image_width+2, 21, termWidth-2, 26)

			type_data2 := widgets.NewParagraph()
			type_data2.Text = "[" + cases.Title(language.English).String(types[1]) + "](fg:yellow,mod:bold)"
			type_data2.SetRect(((termWidth-image_width)/4*3)+image_width, 22, termWidth-3, 25)
			type_data2.Border = false

			ui.Render(type_border1, type_border2, type_data1, type_data2)

		}

		hp := NewGauge()
		hp.Title = "HP"
		hp.SetRect(image_width+2, termHeight-26, termWidth, termHeight-22)
		hp.Percent = currentPokemon.HP * 100 / maxStats.HP
		hp.BarColor = ui.ColorYellow

		attack := NewGauge()
		attack.Title = "Attack"
		attack.SetRect(image_width+2, termHeight-22, termWidth, termHeight-18)
		attack.Percent = currentPokemon.Attack * 100 / maxStats.Attack
		attack.BarColor = ui.ColorYellow

		defense := NewGauge()
		defense.Title = "Defense"
		defense.SetRect(image_width+2, termHeight-18, termWidth, termHeight-14)
		defense.Percent = currentPokemon.Defense * 100 / maxStats.Defense
		defense.BarColor = ui.ColorYellow

		special_attack := NewGauge()
		special_attack.Title = "Special Attack"
		special_attack.SetRect(image_width+2, termHeight-14, termWidth, termHeight-10)
		special_attack.Percent = currentPokemon.Special_defense * 100 / maxStats.Special_defense
		special_attack.BarColor = ui.ColorYellow

		special_defense := NewGauge()
		special_defense.Title = "Special Defense"
		special_defense.SetRect(image_width+2, termHeight-10, termWidth, termHeight-6)
		special_defense.Percent = currentPokemon.Special_attack * 100 / maxStats.Special_attack
		special_defense.BarColor = ui.ColorYellow

		speed := NewGauge()
		speed.Title = "Speed"
		speed.SetRect(image_width+2, termHeight-6, termWidth, termHeight-2)
		speed.Percent = currentPokemon.Speed * 100 / maxStats.Speed
		speed.BarColor = ui.ColorYellow

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
