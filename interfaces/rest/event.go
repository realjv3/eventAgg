package rest

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/realvjv3/event-agg/domain"
	"github.com/realvjv3/event-agg/util"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type EventHandler struct {
	validate *validator.Validate
}

func NewEventHandler(r chi.Router) *EventHandler {
	v := validator.New(validator.WithRequiredStructEnabled())

	h := &EventHandler{
		validate: v,
	}

	h.registerRoutes(r)

	return h
}

func (h *EventHandler) registerRoutes(r chi.Router) {
	r.With(eventCtx).Post("/track", h.TrackEvent)
}

func (h *EventHandler) TrackEvent(w http.ResponseWriter, r *http.Request) {
	// get event from middleware context
	ctx := r.Context()
	event, ok := ctx.Value("event").(*domain.Event)
	if !ok {
		http.Error(w, "missing or malformed event payload", http.StatusBadRequest)
		return
	}

	// validate event
	err := h.validate.Struct(*event)
	if err != nil {
		var validationErrs validator.ValidationErrors
		if errors.As(err, &validationErrs) {
			var invalidFields strings.Builder
			invalidFields.WriteString("Invalid request fields: ")

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
}

func eventCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// parse request body
		var event domain.Event
		err := json.NewDecoder(r.Body).Decode(&event)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// obfuscate personally identifying information
		util.Obfuscate(&event)

		ctx := context.WithValue(r.Context(), "event", &event)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
