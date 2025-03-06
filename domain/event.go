package domain

type Event struct {
	Event      string     `json:"event" validate:"required"`
	Properties Properties `json:"properties" validate:"required"`
}

type Properties struct {
	DistinctID string         `json:"distinct_id" validate:"required"`
	Token      string         `json:"token" validate:"required"`
	Time       int64          `json:"time" validate:"required"`
	DeviceID   string         `json:"device_id,omitempty"`
	UserProps  map[string]any `json:"user_props,omitempty"`
}
