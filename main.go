package main

import (
	"log"
	"net/http"
	"time"

	"github.com/realjv3/event-agg/interfaces/rest"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/spf13/viper"
)

func main() {
	log.Println("Initializing config...")

	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Println("Initializing server...")

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(time.Second * 20))

	rest.NewEventHandler(r)

	log.Fatal(http.ListenAndServe(":8080", r))
}
