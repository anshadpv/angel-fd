package model

type DataRequest[T any] struct {
	Data T `json:"data,omitempty"`
}

type APIResponse struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type EmptyJSON struct {
}
