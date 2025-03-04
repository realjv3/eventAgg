package main

import (
	"log"
	"net/http"
	"time"

	"github.com/realvjv3/event-agg/interfaces/rest"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	log.Println("Initializing server...")

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(time.Second * 20))

	rest.NewEventHandler(r)

	log.Fatal(http.ListenAndServe(":8080", r))
}
