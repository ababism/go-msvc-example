package core

import (
	"github.com/google/uuid"
	"time"
)

type Dish struct {
	Id          uuid.UUID
	Name        string
	Description string
	Price       float64
	Quantity    int64
	IsAvailable bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
