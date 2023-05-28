package repository

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"order-svc/internal/core"
)

type Order interface {
	Create(d core.Order) (core.Order, error)
	Get(id uuid.UUID) (core.Order, error)
}

type Dish interface {
	GetAll(limit, offset int64) ([]core.Dish, error)
	Create(d core.Dish) (core.Dish, error)
	Get(id uuid.UUID) (core.Dish, error)
	Update(d core.Dish) (core.Dish, error)
	Delete(id uuid.UUID) (core.Dish, error)
}

type Repository struct {
	Order
	Dish
}

func New(db *sqlx.DB) *Repository {
	return &Repository{
		Order: NewOrderRepository(db),
		Dish:  NewDishRepository(db),
	}
}
