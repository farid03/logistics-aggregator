package storage

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"logistics-aggregator/src/go/model"
	"os"
	"strings"
)

var Cache = make(map[string]model.User) // перенести кэширование в СУБД
var PG *gorm.DB

func initializeTables(db *gorm.DB) {
	err := db.AutoMigrate(
		&model.User{},
		&model.Specification{},
		&model.Position{},
		&model.Order{},
		&model.Feedback{},
		&model.Car{},
	)
	if err != nil {
		log.Println(err)
	}
}

func ConnectDB() {
	dsn := strings.Join(os.Args[1:], " ")
	log.Printf("DB configurations: %s\n", dsn)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	PG = db
	initializeTables(db)
}
