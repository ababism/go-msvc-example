package http_models

import (
	"github.com/google/uuid"
	"time"
)

type OrderResponse struct {
	Id              uuid.UUID      `json:"id"`
	UserId          uuid.UUID      `json:"user_id"`
	Dishes          []DishResponse `json:"dishes"`
	Status          string         `json:"status"`
	SpecialRequests string         `json:"special_requests"`
	ReadyAt         time.Time      `json:"ready_at"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}
type DishResponse struct {
	Id          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Quantity    int64     `json:"quantity"`
	IsAvailable bool      `json:"is_available"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
