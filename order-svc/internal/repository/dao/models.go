package dao

import (
	"github.com/google/uuid"
	"order-svc/internal/core"
	"time"
)

type DishDAO struct {
	Id          uuid.UUID `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Price       float64   `db:"price"`
	Quantity    int64     `db:"quantity"`
	IsAvailable bool      `db:"is_available"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

func (d *DishDAO) ToDomain() (dish core.Dish) {
	dish.Id = d.Id
	dish.Name = d.Name
	dish.Description = d.Description
	dish.Price = d.Price
	dish.Quantity = d.Quantity
	dish.IsAvailable = d.IsAvailable
	dish.CreatedAt = d.CreatedAt
	dish.UpdatedAt = d.UpdatedAt
	return
}

type OrderDAO struct {
	Id              uuid.UUID `db:"id"`
	UserId          uuid.UUID `db:"user_id"`
	Status          string    `db:"status"`
	SpecialRequests string    `db:"special_requests"`
	ReadyAt         time.Time `db:"ready_at"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
}

func (o *OrderDAO) ToDomain() (order core.Order) {
	order.Id = o.Id
	order.UserId = o.UserId
	order.Status = o.Status
	order.SpecialRequests = o.SpecialRequests
	order.ReadyAt = o.ReadyAt
	order.CreatedAt = o.CreatedAt
	order.UpdatedAt = o.UpdatedAt
	return
}
