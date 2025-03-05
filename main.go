package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/time/rate"

	"github.com/realjv3/event-agg/interfaces/rest"
	"github.com/realjv3/event-agg/services/events"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/spf13/viper"
)

var limiter *rate.Limiter

func init() {
	limiter = rate.NewLimiter(10, 20)
}

func main() {
	log.Println("Initializing config...")

	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Println("Initializing event queue...")

	queue := events.NewQueue()
	go queue.Process()

	log.Println("Initializing server...")

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(time.Second * 20))
	r.Use(rateLimiter)

	rest.NewEventHandler(r, queue)

	log.Fatal(http.ListenAndServe(":8080", r))
}

func rateLimiter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := limiter.Wait(r.Context())
		if err != nil {
			fmt.Println("Rate limit exceeded")
		}

		next.ServeHTTP(w, r)
	})
}
