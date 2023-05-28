package core

import (
	"github.com/google/uuid"
	"time"
)

const (
	OrderStatusRejected = "cancelled"
	OrderStatusWait     = "pending"
	OrderStatusDone     = "ready"
)

type Order struct {
	Id              uuid.UUID
	UserId          uuid.UUID
	Dishes          []Dish
	Status          string
	SpecialRequests string
	ReadyAt         time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
