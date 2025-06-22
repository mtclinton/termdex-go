package termdex-go

import (
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"

	"golang.org/x/exp/slices"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type PokeDB struct {
	db *gorm.DB
}

func loadDB(db_name string) (*PokeDB, error) {
	if _, err := os.Stat(db_name); err != nil {
		log.Println("Creating sqlite database")
		file, err := os.Create(db_name) // Create SQLite file
		if err != nil {
			log.Fatal(err.Error())
		}
		file.Close()
		log.Printf("%s created", db_name)
	}
	sqliteDatabase, err := gorm.Open(sqlite.Open(db_name), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic("failed to connect database")
	}
	return &PokeDB{
		db: sqliteDatabase,
	}, nil

}

func (pkd *PokeDB) createTable() {
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
	result := pkd.db.Exec(createPokemonTableSQL) // Prepare SQL Statement
	if result.Error != nil {
		log.Print(result.Error)
	}
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
	result = pkd.db.Exec(createMaxStatsTableSQL) // Prepare SQL Statement
	if result.Error != nil {
		log.Print(result.Error)
	}
	log.Println("Max stats table created")

	createPokemonTypeTableSQL := `CREATE TABLE IF NOT EXISTS pokemon_type (
        "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
        "pokemon_id" integer NOT NULL,
        "type_id" integer NOT NULL
      );` // SQL Statement for Create Table

	log.Println("Create pokemon type table...")
	result = pkd.db.Exec(createPokemonTypeTableSQL) // Prepare SQL Statement
	if result.Error != nil {
		log.Print(result.Error)
	}
	log.Println("Pokemon type table created")

	createTypeNameTableSQL := `CREATE TABLE IF NOT EXISTS type_name (
        "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
        "name" TEXT NOT NULL,
        "url" TEXT NOT NULL
      );` // SQL Statement for Create Table

	log.Println("Create type name table...")
	result = pkd.db.Exec(createTypeNameTableSQL) // Prepare SQL Statement
	if result.Error != nil {
		log.Print(result.Error)
	}
	log.Println("Type name table created")
}

func (pkd PokeDB) getPokemon(search string) (NewPokemon, []string) {
	var pokemon NewPokemon
	var pokemon_types []PokemonType
	var type_names []TypeName
	var names []string
	if _, err := strconv.Atoi(search); err == nil {
		result := pkd.db.Where("pokemon_id = ?", search).First(&pokemon)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			pokemon.Name = "Not Found"
			pokemon.Entry = "Pokemon not found. Try searching again"
		}
	} else {
		result := pkd.db.Where("name = ?", search).First(&pokemon)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			pokemon.Name = "Not Found"
			pokemon.Entry = "Pokemon not found. Try searching again"
		}
	}
	types_result := pkd.db.Where("pokemon_id = ?", pokemon.Pokemon_id).Find(&pokemon_types)
	if !errors.Is(types_result.Error, gorm.ErrRecordNotFound) {
		var pokemon_types_ids []int
		for _, pokemon_type := range pokemon_types {
			r := reflect.ValueOf(pokemon_type)
			f := reflect.Indirect(r).FieldByName("Type_id")
			pokemon_types_ids = append(pokemon_types_ids, int(f.Int()))
		}
		pkd.db.Where("ID in ?", pokemon_types_ids).Find(&type_names)
		for _, type_name := range type_names {
			r := reflect.ValueOf(type_name)
			f := reflect.Indirect(r).FieldByName("Name")
			names = append(names, f.String())
		}
	}

	return pokemon, names
}

func (pkd *PokeDB) getMaxStats() MaxStats {
	var maxStats MaxStats
	result := pkd.db.First(&maxStats)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		panic("No max stats value")
	}
	return maxStats

}

func (pkd *PokeDB) initializePokemon() {
	pkd.createTable()
	var count int64
	pkd.db.Table("pokemon").Count(&count)
	if count == 0 {
		fmt.Println("Initializing pokemon")
		s := NewScraper(pkd)
		s.run()
	}

}

func (pkd *PokeDB) insertPokemon(pokemon_results []NewPokemon) {

	notfound := NewPokemon{
		Pokemon_id: 0,
		Name:       "Not Found",
		Entry:      "Pokemon not found. Please check id or pokemon name",
	}
	pokemon_results = append(pokemon_results, notfound)
	result := pkd.db.Create(&pokemon_results)
	if result.Error != nil {
		log.Print(result.Error)
	}
}

func (pkd *PokeDB) insertMaxStats(max_stats MaxStats) {
	result := pkd.db.Create(&max_stats)
	if result.Error != nil {
		log.Print(result.Error)
	}
}

func (pkd *PokeDB) insertTypeName(type_names []TypeName) {
	result := pkd.db.Create(&type_names)
	if result.Error != nil {
		log.Print(result.Error)
	}
}

func (pkd *PokeDB) insertPokeType(type_poke_tracker []TypePokeTracker) {
	var type_names []TypeName
	type_name_results := pkd.db.Find(&type_names)
	if type_name_results.Error != nil {
		log.Print(type_name_results.Error)
	}

	var poke_types []PokemonType
	for _, tracker := range type_poke_tracker {
		tid := slices.IndexFunc(type_names, func(tn TypeName) bool { return tn.Name == tracker.Name })
		if tid == -1 {
			log.Panic(("Unknown Type"))
		}
		idx := int(tid)
		poke_type := PokemonType{
			Pokemon_id: tracker.Pokemon_id,
			Type_id:    int(type_names[idx].ID),
		}
		poke_types = append(poke_types, poke_type)
	}
	result := pkd.db.Create(&poke_types)
	if result.Error != nil {
		log.Print(result.Error)
	}
}
