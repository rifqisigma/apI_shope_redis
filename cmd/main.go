package main

import (
	"api_shope/cmd/database"
	"api_shope/cmd/routes"
	"api_shope/internal/handler"
	"api_shope/internal/repository"
	"api_shope/internal/usecase"
	"api_shope/internal/worker"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file")
	}

	db, rdb, err := database.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}

	//auth
	authRepo := repository.NewAuthRepo(db, rdb)
	authUsecase := usecase.NewAuthUsecase(authRepo)
	authHandler := handler.NewAuthHandler(authUsecase)

	//shop
	shopRepo := repository.NewShopRepo(db, rdb)
	shopUsecase := usecase.NewShopUsecase(shopRepo)
	shopHandler := handler.NewShopHandler(shopUsecase)

	r := routes.SetupRoutes(authHandler, shopHandler)

	go func() {
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}

		fmt.Println("Server running on http://localhost:" + port)
		log.Fatal(http.ListenAndServe(":"+port, r))
	}()

	//worker queue redis
	w := worker.NewWorker(db, rdb)
	w.StartFlushWorker(10 * time.Second)
	log.Println("Worker started...")

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	<-stopChan
	log.Println("Stopping worker...")
	defer w.StopFlushWorker()

}
