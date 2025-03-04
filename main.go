package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/realvjv3/event-agg/domain"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
)

func main() {
	log.Println("Initializing server...")

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	validate := validator.New(validator.WithRequiredStructEnabled())

	r.Post("/track", func(w http.ResponseWriter, r *http.Request) {
		// parse request body
		var event domain.Event
		err := json.NewDecoder(r.Body).Decode(&event)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// validate event
		err = validate.Struct(event)
		if err != nil {
			var validationErrs validator.ValidationErrors
			if errors.As(err, &validationErrs) {
				var invalidFields strings.Builder
				invalidFields.WriteString("Invalid struct fields: ")

				for _, e := range validationErrs {
					invalidFields.WriteString(e.StructField() + "; ")
				}

				http.Error(w, invalidFields.String(), http.StatusBadRequest)
				return
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		// process event
		log.Printf("Received event:%#v\n", event)

		w.Header().Set("Content-type", "application/json")
		err = json.NewEncoder(w).Encode(event)
		return
	})

	log.Fatal(http.ListenAndServe(":8080", r))
}
