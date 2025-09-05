package main

import (
	"log"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/karambo3a/wbtech_test_task/internal/cache"
	"github.com/karambo3a/wbtech_test_task/internal/consumer"
	"github.com/karambo3a/wbtech_test_task/internal/handlers"
	"github.com/karambo3a/wbtech_test_task/internal/repository"
	"github.com/karambo3a/wbtech_test_task/internal/service"
)

func main() {
	log.SetFlags(log.Lshortfile)

	db, err := repository.NewPostgresDB()
	if err != nil {
		log.Fatalf("failed connect to db: %s", err.Error())
	}

	defer db.Close()
	log.Println("connected to data base")

	repository := repository.NewRepository(db)
	log.Println("repository created")

	consumer := consumer.NewConsumer()
	cache := cache.NewRedisCache(50)
	service := service.NewService(repository, consumer, cache, int64(100))
	log.Println("service created")
	defer service.CloseConsumer()

	handler := handlers.NewHandler(service)
	log.Println("handler created")

	if err := http.ListenAndServe(":"+os.Getenv("SERVER_PORT"), handler.InitRouts()); err != nil {
		panic(err)
	}
}
