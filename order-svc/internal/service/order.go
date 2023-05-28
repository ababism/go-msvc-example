package service

import (
	"errors"
	"github.com/google/uuid"
	"order-svc/internal/core"
	"order-svc/internal/repository"
)

func NewOrderService(dr repository.Order) *OrderService {
	return &OrderService{r: dr}
}

type OrderService struct {
	r repository.Order
}

func (os OrderService) Create(o core.Order) (core.Order, error) {
	if o.Dishes == nil || len(o.Dishes) <= 0 {
		return core.Order{}, errors.New("order should contain at least one dish")
	}
	return os.r.Create(o)
}

func (os OrderService) Get(id uuid.UUID) (core.Order, error) {
	return os.r.Get(id)
}
