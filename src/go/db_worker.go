package swagger

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"logistics-aggregator/src/go/model"
	"os"
	"strings"
)

var pg *gorm.DB

func initializeTables(db *gorm.DB) {
	errors := make([]error, 6)
	errors[0] = db.AutoMigrate(&model.User{})
	errors[1] = db.AutoMigrate(&model.Specification{})
	errors[2] = db.AutoMigrate(&model.Position{})
	errors[3] = db.AutoMigrate(&model.Order{})
	errors[4] = db.AutoMigrate(&model.Feedback{})
	errors[5] = db.AutoMigrate(&model.Car{})
	for _, err := range errors {
		if err != nil {
			log.Println(err)
		}
	}
}

func ConnectDB() {
	dsn := strings.Join(os.Args[1:], " ")
	log.Printf("DB configurations: %s\n", dsn)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	pg = db
	initializeTables(db)
}
