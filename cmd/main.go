package main

import (
	"log"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/karambo3a/wbtech_test_task/internal/cache"
	"github.com/karambo3a/wbtech_test_task/internal/consumer"
	"github.com/karambo3a/wbtech_test_task/internal/handler"
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
	service := service.NewService(repository, consumer, cache)
	log.Println("service created")
	defer service.CloseConsumer()

	handler := handler.NewHandler(service)
	log.Println("handler created")

	if err := http.ListenAndServe(":8081", handler.InitRouts()); err != nil {
		panic(err)
	}
}
