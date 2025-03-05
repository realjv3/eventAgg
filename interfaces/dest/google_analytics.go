package dest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/realjv3/event-agg/domain"

	"github.com/spf13/viper"
)

const url = "https://www.google-analytics.com/mp/collect"

type ga4Payload struct {
	ClientID string     `json:"client_id"`
	Events   []ga4Event `json:"events"`
}

type ga4Event struct {
	Name   string         `json:"name"`
	Params map[string]any `json:"params,omitempty"`
}

func SendGoogleAnalytics(event *domain.Event) error {
	measurementId := viper.GetString("GOOGLE_ANALYTICS_MEASUREMENTS_ID")
	apiSecret := viper.GetString("GOOGLE_ANALYTICS_API_KEY")

	if measurementId == "" || apiSecret == "" {
		return errors.New("failed to load Google Analytics secrets")
	}

	payload, err := json.Marshal(ga4Payload{
		ClientID: "github.com/realjv3/event-agg",
		Events: []ga4Event{
			{
				Name: event.Event,
				Params: map[string]any{
					"distinct_id": event.Properties.DistinctID,
					"token":       event.Properties.Token,
					"time":        event.Properties.Time,
					"device_id":   event.Properties.DeviceID,
					"user_props":  event.Properties.UserProps,
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to marshal Google Analytics event payload: %v", err)
	}

	resp, err := http.Post(
		fmt.Sprintf("%s?measurement_id=%s&api_secret=%s", url, measurementId, apiSecret),
		"application/json",
		bytes.NewReader(payload),
	)
	if err != nil {
		return fmt.Errorf("error sending Google Analytics event ga4Payload.: %v", err)
	}

	if resp.StatusCode != 204 {
		return fmt.Errorf("error status code from Google Analytics : %d", resp.StatusCode)
	} else {
		log.Println("Event sent to Google Analytics ok")
	}

	return nil
}
