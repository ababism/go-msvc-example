package service

import (
	"github.com/google/uuid"
	"order-svc/internal/core"
	"order-svc/internal/repository"
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

type Service struct {
	Order
	Dish
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Order: NewOrderService(repos.Order),
		Dish:  NewDishService(repos.Dish),
	}
}
