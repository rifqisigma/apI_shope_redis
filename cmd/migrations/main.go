package main

import (
	"api_shope/cmd/database"
	"api_shope/model"
	"log"
)

func main() {

	db, _, err := database.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}

	if err := db.AutoMigrate(&model.User{}, &model.Store{}, &model.Product{}, model.CartItem{}); err != nil {
		log.Fatal(err)
	}

	log.Println("✅ Migrasi selesai dan database siap")
}
